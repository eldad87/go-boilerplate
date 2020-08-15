package sqlmw

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/ngrok/sqlmw"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"reflect"
	"strings"
)

// Inspired by https://github.com/luna-duclos/instrumentedsql
type Interceptor struct {
	sqlmw.NullInterceptor
	Tracer opentracing.Tracer
}

// Connection interceptors
func (in Interceptor) ConnBeginTx(ctx context.Context, tx driver.ConnBeginTx, opts driver.TxOptions) (t driver.Tx, err error) {
	span, _ := in.startSpanFromContext(ctx, "ConnBeginTx")

	defer func() {
		if err != nil {
			span = in.spanError(span, err)
		}
		span.Finish()
	}()

	return tx.BeginTx(ctx, opts)
}

// Rows interceptors
func (in Interceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) (err error) {
	span, _ := in.startSpanFromContext(ctx, "RowsNext")

	defer func() {
		if err != io.EOF && err != nil {
			span = in.spanError(span, err)
		}
		span.Finish()
	}()

	return rows.Next(dest)
}

// Stmt interceptors
func (in Interceptor) StmtClose(ctx context.Context, stmt driver.Stmt) (err error) {
	span, _ := in.startSpanFromContext(ctx, "StmtClose")

	defer func() {
		if err != nil {
			span = in.spanError(span, err)
		}
		span.Finish()
	}()

	return stmt.Close()
}

func (in Interceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, query string, args []driver.NamedValue) (rows driver.Result, err error) {
	span, ctx := in.startSpanFromContext(ctx, "StmtExecContext")
	span.SetTag("query", query)
	span.SetTag("args", formatArgs(args))

	defer func() {
		if err != nil {
			span = in.spanError(span, err)
		}
		span.Finish()
	}()

	return stmt.ExecContext(ctx, args)
}

func (in Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	span, ctx := in.startSpanFromContext(ctx, "StmtQueryContext")
	span.SetTag("query", query)
	span.SetTag("args", formatArgs(args))

	defer func() {
		if err != nil {
			span = in.spanError(span, err)
		}
		span.Finish()
	}()

	return stmt.QueryContext(ctx, args)
}

// Tx interceptors
func (in Interceptor) TxCommit(ctx context.Context, tx driver.Tx) (err error) {
	span, _ := in.startSpanFromContext(ctx, "TxCommit")

	defer func() {
		if err != nil {
			span = in.spanError(span, err)
		}
		span.Finish()
	}()

	return tx.Commit()
}

func (in Interceptor) TxRollback(ctx context.Context, tx driver.Tx) (err error) {
	span, _ := in.startSpanFromContext(ctx, "TxRollback")

	defer func() {
		if err != nil {
			span = in.spanError(span, err)
		}
		span.Finish()
	}()

	return tx.Rollback()
}

func (in Interceptor) spanError(span opentracing.Span, err error) opentracing.Span {
	ext.Error.Set(span, true)
	span.LogFields(
		log.String("event", "error"),
		log.String("message", err.Error()),
	)

	return span
}

func (in Interceptor) startSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, in.Tracer, operationName)
	span.SetTag("component", "database/sql")

	return span, ctx
}

func formatArgs(args interface{}) string {
	argsVal := reflect.ValueOf(args)
	if argsVal.Kind() != reflect.Slice {
		return "<unknown>"
	}

	strArgs := make([]string, 0, argsVal.Len())
	for i := 0; i < argsVal.Len(); i++ {
		strArgs = append(strArgs, formatArg(argsVal.Index(i).Interface()))
	}

	return fmt.Sprintf("{%s}", strings.Join(strArgs, ", "))
}

func formatArg(arg interface{}) string {
	strArg := ""
	switch arg := arg.(type) {
	case []uint8:
		strArg = fmt.Sprintf("[%T len:%d]", arg, len(arg))
	case string:
		strArg = fmt.Sprintf("[%T %q]", arg, arg)
	case driver.NamedValue:
		if arg.Name != "" {
			strArg = fmt.Sprintf("[%T %s=%v]", arg.Value, arg.Name, formatArg(arg.Value))
		} else {
			strArg = formatArg(arg.Value)
		}
	default:
		strArg = fmt.Sprintf("[%T %v]", arg, arg)
	}

	return strArg
}
