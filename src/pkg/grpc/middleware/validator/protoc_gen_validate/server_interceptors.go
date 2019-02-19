// Based on: https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/validator/validator.go

package protoc_gen_validate

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	grpcErrors "github.com/eldad87/go-boilerplate/src/pkg/grpc/error"
)

type validator interface {
	Validate() error
}

type FieldError interface {
	Field() string
	Reason() string
}

// UnaryServerInterceptor returns a new unary server interceptor that validates incoming messages.
//
// Invalid messages will be rejected with `InvalidArgument` before reaching any userspace handlers.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if res := ErrorHandler(req); res != nil {
			return nil, res
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that validates incoming messages.
//
// The stage at which invalid messages will be rejected with `InvalidArgument` varies based on the
// type of the RPC. For `ServerStream` (1:m) requests, it will happen before reaching any userspace
// handlers. For `ClientStream` (n:1) or `BidiStream` (n:m) RPCs, the messages will be rejected on
// calls to `stream.Recv()`.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &recvWrapper{stream}
		return handler(srv, wrapper)
	}
}

type recvWrapper struct {
	grpc.ServerStream
}

func (s *recvWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if res := ErrorHandler(m); res != nil {
		return res
	}
	return nil
}

func ErrorHandler(msg interface{}) error {
	if v, ok := msg.(validator); ok {
		if err := v.Validate(); err != nil {
			br := grpcErrors.NewBadRequest()
			if v, ok := err.(FieldError); ok {
				br.AddViolation(v.Field(), v.Reason())
			}
			return br.GetStatusError(codes.InvalidArgument, err.Error())
		}
	}

	return nil
}
