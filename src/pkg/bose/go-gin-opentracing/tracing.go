package ginopentracing

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func GetSpanTraceID() string {
	return "tracing-context"
}

func GetSpan(ctx *gin.Context) (span opentracing.Span) {
	if bspan, ok := ctx.Get(GetSpanTraceID()); !ok {
		return nil
	} else if cspan, ok := bspan.(opentracing.Span); !ok {
		return nil
	} else {
		return cspan
	}
}
