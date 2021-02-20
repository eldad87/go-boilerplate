package middleware

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
)

const requestId = contextKey("request_id")

func ContextReqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqId := r.Header.Get("X-Request-Id")
		if reqId == "" {
			u, err := uuid.NewV4()
			if err == nil {
				reqId = u.String()
			}
		}
		// return the request id for future tracking
		w.Header().Set("X-Request-Id", reqId)

		ctx := r.Context()
		ctx = context.WithValue(ctx, requestId, reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetReqIdFromContext(ctx context.Context) (reqId string, exists bool) {
	reqId, exists = ctx.Value(requestId).(string)
	return
}
