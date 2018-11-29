package healthcheck

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func RedisCheck(dsn string, timeout time.Duration) func() error {
	return func() error {
		pool := &redis.Pool{
			MaxIdle:     1,
			IdleTimeout: 10 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.DialURL(dsn) },
		}

		conn := pool.Get()
		defer conn.Close()

		data, err := redis.DoWithTimeout(conn, timeout, "PING")
		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("PING command returned an empty response")
		}

		if data != "PONG" {
			return fmt.Errorf("PING command returned an unexpected response")
		}

		return nil
	}
}
