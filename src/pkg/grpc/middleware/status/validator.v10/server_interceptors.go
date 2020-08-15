package v9_validation_error

import (
	"context"
	"fmt"
	grpcErrors "github.com/eldad87/go-boilerplate/src/pkg/grpc/error"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// UnaryServerInterceptor returns a new unary server interceptor that transform v10 validation errors to gRPC status code.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, ErrorHandler(err)
		}

		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that transform v10 validation errors to gRPC status code.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, stream)
		return ErrorHandler(err)
	}
}

func ErrorHandler(error error) error {
	if vErrors, ok := error.(validator.ValidationErrors); ok {
		br := grpcErrors.NewBadRequest()

		for _, err := range vErrors {
			br.AddViolation(err.StructField(), fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", err.Namespace(), err.Field(), err.Tag()))
		}

		return br.GetStatusError(codes.InvalidArgument, vErrors.Error())
	}

	return error
}
