package null_v4

import (
	"database/sql/driver"
	"github.com/go-playground/validator/v10"
	"gopkg.in/guregu/null.v4"
	"reflect"
)

// Check if ID already presents && Set
func IsUpdate(fl validator.FieldLevel) bool {
	p := fl.Parent()
	f := reflect.Indirect(p).FieldByName("ID")
	return f.IsValid() && f.Uint() > 0
}

// Check if ID already presents && Set
func IsNew(fl validator.FieldLevel) bool {
	p := fl.Parent()
	f := reflect.Indirect(p).FieldByName("ID")
	return f.IsValid() && f.Uint() == 0
}

func RegisterIsUpdate(v *validator.Validate) error {
	return v.RegisterValidation("isUpdate", IsUpdate)
}

func RegisterIsNew(v *validator.Validate) error {
	return v.RegisterValidation("isNew", IsNew)
}

func RegisterSQLNullValuer(v *validator.Validate) {
	v.RegisterCustomTypeFunc(validateValuer, null.String{}, null.Int{}, null.Bool{}, null.Float{}, null.Time{})
}

func validateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}
	return nil
}
