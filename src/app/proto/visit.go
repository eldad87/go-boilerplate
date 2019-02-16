package pb

import (
	"context"
	"fmt"
	"github.com/eldad87/go-boilerplate/src/app"
	grpcErrors "github.com/eldad87/go-boilerplate/src/pkg/grpc/error"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"gopkg.in/go-playground/validator.v9"
)

type VisitService struct {
	VisitService app.VisitService
}

func (vs *VisitService) Get(c context.Context, id *ID) (*VisitResponse, error) {
	i := uint(id.GetId())
	v, err := vs.VisitService.Get(c, &i)
	if err != nil {
		return nil, err
	}

	return vs.visitToProto(v)
}

// Update/Create a device
func (vs *VisitService) Set(c context.Context, v *VisitRequest) (*VisitResponse, error) {
	aVis, err := vs.protoToVisit(v)
	if err != nil {
		return nil, err
	}

	if errs := vs.VisitService.Validate(c, aVis); errs != nil {
		br := grpcErrors.NewBadRequest()
		for _, err := range errs.(validator.ValidationErrors) {
			// TODO: use FieldViolation
			br.AddViolation(err.StructField(), fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", err.Namespace(), err.Field(), err.Tag()))
		}

		return nil, br.GetStatusError(codes.InvalidArgument, errs.Error())
	}

	gVis, err := vs.VisitService.Set(c, aVis)
	if err != nil {
		return nil, err
	}

	return vs.visitToProto(gVis)
}

func (vs *VisitService) visitToProto(visit *app.Visit) (*VisitResponse, error) {
	created, err := ptypes.TimestampProto(visit.CreatedAt)
	if err != nil {
		return nil, err
	}

	updated, err := ptypes.TimestampProto(visit.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &VisitResponse{Id: uint32(visit.ID), FirstName: visit.FirstName, LastName: visit.LastName, CreatedAt: created, UpdatedAt: updated}, nil
}

func (vs *VisitService) protoToVisit(visit *VisitRequest) (*app.Visit, error) {
	return &app.Visit{
		ID:        uint(visit.Id),
		FirstName: visit.FirstName,
		LastName:  visit.LastName,
	}, nil
}
