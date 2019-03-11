package grpc

import (
	"context"
	"github.com/eldad87/go-boilerplate/src/app"
	pb "github.com/eldad87/go-boilerplate/src/transport/grpc/proto"
	"github.com/golang/protobuf/ptypes"
)

type VisitTransport struct {
	VisitService app.VisitService
}

func (vs *VisitTransport) Get(c context.Context, id *pb.ID) (*pb.VisitResponse, error) {
	i := uint(id.GetId())
	v, err := vs.VisitService.Get(c, &i)
	if err != nil {
		return nil, err
	}

	return vs.visitToProto(v)
}

// Update/Create a device
func (vs *VisitTransport) Set(c context.Context, v *pb.VisitRequest) (*pb.VisitResponse, error) {
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

func (vs *VisitTransport) visitToProto(visit *app.Visit) (*pb.VisitResponse, error) {
	created, err := ptypes.TimestampProto(visit.CreatedAt)
	if err != nil {
		return nil, err
	}

	updated, err := ptypes.TimestampProto(visit.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &pb.VisitResponse{Id: uint32(visit.ID), FirstName: visit.FirstName, LastName: visit.LastName, CreatedAt: created, UpdatedAt: updated}, nil
}

func (vs *VisitTransport) protoToVisit(visit *pb.VisitRequest) (*app.Visit, error) {
	return &app.Visit{
		ID:        uint(visit.Id),
		FirstName: visit.FirstName,
		LastName:  visit.LastName,
	}, nil
}
