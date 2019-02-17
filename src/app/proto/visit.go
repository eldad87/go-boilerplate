package pb

import (
	"context"
	"github.com/eldad87/go-boilerplate/src/app"
	grpcErrors "github.com/eldad87/go-boilerplate/src/pkg/grpc/error"
	"github.com/eldad87/go-boilerplate/src/pkg/validator"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
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

	gVis, err := vs.VisitService.Set(c, aVis)
	if err != nil {
		// TODO: Improve the way we convert and return
		br := grpcErrors.NewBadRequest()
		if errs, ok := err.(*validator.StructViolation); ok {
			for _, err := range errs.FieldViolation {
				br.AddViolation(err.Field, err.Description)
			}

			return nil, br.GetStatusError(codes.InvalidArgument, err.Error())
		}

		return nil, br.GetStatusError(codes.Unknown, err.Error())
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
