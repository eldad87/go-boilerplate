package tasks

import (
	"fmt"
	machineryTasks "github.com/RichardKnop/machinery/v1/tasks"
)

func NewSignature(name string, args ...interface{}) (*machineryTasks.Signature, error) {
	sig := &machineryTasks.Signature{
		Name: name,
		Args: make([]machineryTasks.Arg, len(args)),
	}

	for i, arg := range args {
		sig.Args[i] = machineryTasks.Arg{
			Type:  fmt.Sprintf("%T", arg),
			Value: arg,
		}
	}

	return sig, nil
}
