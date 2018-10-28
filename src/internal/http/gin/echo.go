package gin

import (
	"github.com/RichardKnop/machinery/v1"
	gintracing "github.com/eldad87/go-boilerplate/src/pgk/bose/go-gin-opentracing"
	"github.com/eldad87/go-boilerplate/src/pgk/machinery/v1/tasks"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func NewEcho(logger *log.Logger, server *machinery.Server) *Echo {
	return &Echo{logger: logger, server: server}
}

type Echo struct {
	logger *log.Logger
	server *machinery.Server
}

func (e *Echo) Repeat(c *gin.Context) {
	sig, err := tasks.NewSignature("repeat", "hello")
	if err != nil {
		e.logger.Error(err.Error())
		c.JSON(500, gin.H{})
		return
	}

	// Set Req's span as parent, pass it along to Machinery -> Worker
	ctx := opentracing.ContextWithSpan(context.Background(), gintracing.GetSpan(c))
	res, err := e.server.SendTaskWithContext(ctx, sig)
	if err != nil {
		e.logger.Error(err.Error())
		c.JSON(501, gin.H{})
		return
	}

	code := 200
	if res.GetState().IsFailure() {
		code = 401
	} else if res.GetState().IsSuccess() {
		code = 201
	}

	msg := "Hello"
	if res.GetState().IsCompleted() {
		msg = "Done"
	}

	// Expect to see 201/Done
	c.JSON(code, gin.H{
		"message": msg,
		"results": res.GetState().Results,
	})
}
