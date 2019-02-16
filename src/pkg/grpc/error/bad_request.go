package error

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewBadRequest() *BadRequest {
	return &BadRequest{badRequest: &errdetails.BadRequest{}}
}

type BadRequest struct {
	badRequest *errdetails.BadRequest
}

func (br *BadRequest) AddViolation(field string, description string) {
	v := &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: description,
	}

	br.badRequest.FieldViolations = append(br.badRequest.FieldViolations, v)
}

func (br *BadRequest) GetDetails() *errdetails.BadRequest {
	return br.badRequest
}

func (br *BadRequest) GetStatusError(c codes.Code, msg string) error {
	st := status.New(c, msg)
	if det, err := st.WithDetails(br.GetDetails()); err != nil {
		return st.Err()
	} else {
		return det.Err()
	}
}
