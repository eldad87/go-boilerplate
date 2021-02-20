package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
)

const RequestMethod = "request_method"
const RequestPath = "request_path"
const RequesURI = "request_uri"
const RequesId = "request_id"

func Opentracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		route := mux.CurrentRoute(r)

		span, ctx := opentracing.StartSpanFromContext(ctx, "http")
		defer span.Finish()

		span.SetTag(RequestMethod, r.Method)
		span.SetTag(RequesURI, r.RequestURI)

		path, err := route.GetPathTemplate()
		if err == nil && path != "" {
			span.SetTag(RequestPath, path)
		}

		req_id, exists := GetReqIdFromContext(ctx)
		if exists {
			span.SetTag(RequesId, req_id)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
