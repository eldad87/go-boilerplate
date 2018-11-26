package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

// https://console.bluemix.net/docs/go/fault_tolerance.html#fault-tolerance
func HystrixHandler(command string, config hystrix.CommandConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		hystrix.ConfigureCommand(command, config)

		hystrix.Do(command, func() error {
			c.Next()
			return nil
		}, func(err error) error {
			//Handle failures
			return err
		})
	}
}
