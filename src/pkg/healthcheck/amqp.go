package healthcheck

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

func AMQPCheck(dsn string, timeout time.Duration) func() error {
	return func() error {
		conn, err := amqp.Dial(dsn)
		if err != nil {
			return fmt.Errorf("AMQP cannot dial: %v", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			return fmt.Errorf("AMQP cannot get a channel: %v", err)
		}
		defer ch.Close()

		return nil
	}
}
