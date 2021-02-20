package error

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type errorBody struct {
	Error           string                                  `json:”error"`
	FieldViolations []*errdetails.BadRequest_FieldViolation `json:”fieldViolation”`
}

func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType("application/json"))
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))

	eb := errorBody{Error: grpc.ErrorDesc(err)}
	st := status.Convert(err)
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range t.GetFieldViolations() {
				eb.FieldViolations = append(eb.FieldViolations, violation)
			}
		}
	}

	jErr := json.NewEncoder(w).Encode(eb)

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
