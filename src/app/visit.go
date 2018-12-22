package app

import "time"

type Visit struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type VisitService interface {
	Get(id *int) (*Visit, error)
	Set(v *Visit) (*Visit, error)
}

type VisitServiceDemo struct {
}

func (vsd *VisitServiceDemo) Get(id *int) (*Visit, error) {
	return &Visit{*id, time.Now()}, nil
}

func (vsd *VisitServiceDemo) Set(v *Visit) (*Visit, error) {
	return v, nil
}
