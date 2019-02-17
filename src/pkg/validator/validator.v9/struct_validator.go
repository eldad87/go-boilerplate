package validator

import (
	"context"
	"fmt"
	baseValidator "github.com/eldad87/go-boilerplate/src/pkg/validator"
	v9 "gopkg.in/go-playground/validator.v9"
)

func NewStructVallidator(validator *v9.Validate) *StructValidator {
	return &StructValidator{validator: validator}
}

type StructValidator struct {
	validator *v9.Validate
}

func (vs *StructValidator) Struct(s interface{}) error {
	err := vs.validator.Struct(s)
	if _, ok := err.(*v9.InvalidValidationError); ok {
		return err
	}

	return vs.toStructViolation(err)
}

func (vs *StructValidator) StructCtx(ctx context.Context, s interface{}) error {
	err := vs.validator.StructCtx(ctx, s)
	if _, ok := err.(*v9.InvalidValidationError); ok {
		return err
	}

	return vs.toStructViolation(err)
}

func (vs *StructValidator) toStructViolation(errs error) error {
	if errs == nil {
		return nil
	}

	sv := baseValidator.StructViolation{Description: errs.Error()}

	for _, err := range errs.(v9.ValidationErrors) {
		sv.AddViolation(err.StructField(), fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", err.Namespace(), err.Field(), err.Tag()))
	}

	return &sv
}
