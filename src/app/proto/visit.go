package pb

import (
	"context"
	"github.com/eldad87/go-boilerplate/src/app"
	"github.com/golang/protobuf/ptypes"
)

type VisitService struct {
	VisitService app.VisitService
}

func (vs *VisitService) Get(c context.Context, id *ID) (*Visit, error) {
	i := int(id.GetId())
	v, err := vs.VisitService.Get(&i)
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

	gVis, err := vs.VisitService.Set(aVis)
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
