package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	log "github.com/sirupsen/logrus"
	"time"
)

// CommandConfig is used to tune circuit settings at runtime
type RetryConfig struct {
	Attempts int
	Delay    time.Duration
}

func ReHystrix(fn func() error, breakerName string, reConf RetryConfig) chan error {
	return hystrix.Go(breakerName, func() error {
		r := retrier.New(retrier.ConstantBackoff(reConf.Attempts, reConf.Delay*time.Millisecond), nil)
		err := r.Run(func() error {
			return fn()
		})
		return err
	}, func(err error) error {
		circuit, _, _ := hystrix.GetCircuit(breakerName)
		log.Error("In fallback function for breaker %v, Circuit state is: %v, error: %v",
			breakerName, circuit.IsOpen(), err.Error())
		return err
	})
}

func ReHystrixWithRes(fn func() (interface{}, error), breakerName string, reConf RetryConfig) (interface{}, error) {
	output := make(chan interface{})

	errors := ReHystrix(func() error {
		res, err := fn()
		if res != nil {
			output <- res
		}
		return err
	}, breakerName, reConf)

	select {
	case out := <-output:
		log.Debug("Call in breaker %v successful", breakerName)
		return out, nil
	case err := <-errors:
		return false, err
	}
}
