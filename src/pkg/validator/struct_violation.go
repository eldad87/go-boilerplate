package validator

type StructViolation struct {
	Description    string
	FieldViolation []*FieldViolation
}

func (sv *StructViolation) AddViolation(field string, description string) {
	v := FieldViolation{
		Field:       field,
		Description: description,
	}

	sv.FieldViolation = append(sv.FieldViolation, &v)
}

func (sv *StructViolation) Error() string {
	return sv.Description
}
