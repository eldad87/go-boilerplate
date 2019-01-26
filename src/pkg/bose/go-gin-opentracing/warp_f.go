package ginopentracing

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func WarpF(f http.HandlerFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		span := GetSpan(c)
		if span != nil {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"gRPC-Gateway-ServeHTTP",
				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
				ext.RPCServerOption(span.Context()),
				grpcGatewayTag,
			)
			ext.HTTPMethod.Set(serverSpan, c.Request.Method)
			ext.HTTPUrl.Set(serverSpan, c.Request.URL.String())

			defer serverSpan.Finish()
		} else {
			// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/docs/_docs/customizingyourgateway.md
			parentSpanContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(c.Request.Header))
			if err == nil || err == opentracing.ErrSpanContextNotFound {
				serverSpan := opentracing.GlobalTracer().StartSpan(
					"ServeHTTP",
					// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
					ext.RPCServerOption(parentSpanContext),
					grpcGatewayTag,
				)
				c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), serverSpan))
				ext.HTTPMethod.Set(serverSpan, c.Request.Method)
				ext.HTTPUrl.Set(serverSpan, c.Request.URL.String())
				defer serverSpan.Finish()
			}
		}

		f(c.Writer, c.Request)
	}
}
