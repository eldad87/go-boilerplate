package app

import (
	"context"
	"time"
)

type Visit struct {
	ID        int       `json:"id" validate:"gte=0"`
	FirstName string    `json:"first_name" validate:"required,gte=2,lte=254"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VisitService interface {
	Get(c context.Context, id *int) (*Visit, error)
	Set(c context.Context, v *Visit) (*Visit, error)
	//Validate(c context.Context, v *Visit) ([]*validator.FieldViolation, error)
}
