package result

func NewResult(typ string, val interface{}) *Result {
	return &Result{typ, val}
}

type Result struct {
	Type  string
	Value interface{}
}
