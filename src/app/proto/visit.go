package pb

import (
	"context"
	"github.com/eldad87/go-boilerplate/src/app"
	grpcErrors "github.com/eldad87/go-boilerplate/src/pkg/grpc/error"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VisitService struct {
	VisitService app.VisitService
}

func (vs *VisitService) Get(c context.Context, id *ID) (*Visit, error) {
	i := int(id.GetId())

	// Validation
	if i != 1 {
		br := grpcErrors.NewBadRequest()
		br.AddViolation("Id", "Id must be 1")

		// we set the status here
		st := status.New(codes.InvalidArgument, "Invalid Argument")
		if det, err := st.WithDetails(br.GetDetails()); err != nil {
			return nil, st.Err()
		} else {
			return nil, det.Err()
		}
	}

	v, err := vs.VisitService.Get(c, &i)
	if err != nil {
		return nil, err
	}

	return vs.visitToProto(v)
}

// Update/Create a device
func (vs *VisitService) Set(c context.Context, v *Visit) (*Visit, error) {
	aVis, err := vs.protoToVisit(v)
	if err != nil {
		return nil, err
	}

	gVis, err := vs.VisitService.Set(c, aVis)
	if err != nil {
		return nil, err
	}

	return vs.visitToProto(gVis)
}

func (vs *VisitService) visitToProto(visit *app.Visit) (*Visit, error) {
	t, err := ptypes.TimestampProto(visit.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &Visit{Id: int32(visit.ID), CreatedAt: t}, nil
}

func (vs *VisitService) protoToVisit(visit *Visit) (*app.Visit, error) {
	t, err := ptypes.Timestamp(visit.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &app.Visit{
		ID:        int(visit.Id),
		CreatedAt: t,
	}, nil
}
