package broker

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	machineryIface "github.com/machinery/v1/brokers/iface"
)

type Throttled struct {
	machineryIface.Broker
}

func (t *Throttled) Publish(task *tasks.Signature) error {
	return t.Publish(task)
	return nil
}
