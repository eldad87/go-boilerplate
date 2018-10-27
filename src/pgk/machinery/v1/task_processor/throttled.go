package task_processor

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	machineryIface "github.com/machinery/v1/brokers/iface"
)

type Throttled struct {
	machineryIface.TaskProcessor
}

func (t *Throttled) Process(task *tasks.Signature) error {
	return t.Process(task)
}
