package validator

import "context"

type StructValidator interface {
	Struct(interface{}) error
	StructCtx(context.Context, interface{}) error
}
