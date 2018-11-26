package gin

import (
	gintracing "github.com/eldad87/go-boilerplate/src/pkg/bose/go-gin-opentracing"
	"github.com/eldad87/go-boilerplate/src/pkg/task/producer"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func NewEcho(logger *log.Logger, producer producer.Producer) *Echo {
	return &Echo{logger: logger, producer: producer}
}

type Echo struct {
	logger   *log.Logger
	producer producer.Producer
}

func (e *Echo) Repeat(c *gin.Context) {
	// Set Req's span as parent, pass it along to Machinery -> Worker
	ctx := opentracing.ContextWithSpan(context.Background(), gintracing.GetSpan(c))
	options := map[string]interface{}{} // Yeah, I can use nil - but its an example. This is how you can inject options
	req, err := e.producer.NewRequest("repeat", options, "hello")
	if err != nil {
		e.logger.Error(err.Error())
		c.JSON(500, gin.H{})
		return
	}

	res, err := e.producer.Produce(req, map[string]interface{}{"ctx": ctx})
	if err != nil {
		e.logger.Error(err.Error())
		c.JSON(500, gin.H{})
		return
	}

	code := 500
	if res.IsFailure() {
		code = 401
	} else if res.IsSuccess() {
		code = 201
	}

	msg := "Hello"
	if res.IsCompleted() {
		msg = "Done"
	}

	resArr, err := res.Result()

	// Expect to see 201/Done
	c.JSON(code, gin.H{
		"message": msg,
		"results": resArr,
		"error":   err,
	})
}
