package machinery

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/afex/hystrix-go/hystrix"
	reHystrix "github.com/eldad87/go-boilerplate/src/pkg/concurrency/hystrix"
	"time"

	pkgMachineryTasks "github.com/eldad87/go-boilerplate/src/pkg/machinery/v1/tasks"
	"github.com/eldad87/go-boilerplate/src/pkg/task/producer"
	"github.com/eldad87/go-boilerplate/src/pkg/task/request"
	pkgResult "github.com/eldad87/go-boilerplate/src/pkg/task/result"
)

func NewProducer(server *machinery.Server, hyConf hystrix.CommandConfig, reHyConf reHystrix.RetryConfig) producer.Producer {
	md5HashInBytes := md5.Sum([]byte(server.GetConfig().Broker))
	commandName := hex.EncodeToString(md5HashInBytes[:])

	hystrix.ConfigureCommand(commandName, hyConf)
	return &Producer{server: server, commandName: commandName, reConf: reHyConf}
}

type Producer struct {
	server      *machinery.Server
	commandName string
	reConf      reHystrix.RetryConfig
}

func (p *Producer) NewRequest(name string, options map[string]interface{}, args ...interface{}) (producer.Request, error) {
	return request.NewRequest(name, options, args...)
}

// options[sleepDuration] = int
func (p *Producer) Produce(req producer.Request, options map[string]interface{}) (producer.AsyncResponse, error) {
	sig, err := p.requestToSignature(req)
	if err != nil {
		return nil, err
	}

	sleepDuration := time.Duration(0)
	if options != nil {
		if val, ok := options["sleepDuration"]; ok {
			if i, ok := val.(int); ok {
				sleepDuration = time.Duration(i)
			}
		}
	}

	asyncRes, err := reHystrix.ReHystrixWithRes(
		func() (interface{}, error) {
			if ctx, ok := options["ctx"]; ok {
				if ctx, ok := ctx.(context.Context); ok {
					// ToDo:
					// Make NewRequest() accept options["ctx"], append it to producer.Request (add new attribute Headers map[string]interface{}
					// if task already have a context with Tracing, link them as FollowsFrom / ChildOf, for additional info:
					// https://github.com/opentracing/specification/blob/master/specification.md
					return p.server.SendTaskWithContext(ctx, sig)
				}
			}
			return p.server.SendTask(sig)
		},
		p.commandName, p.reConf)
	if err != nil {
		return nil, err
	}

	if asyncRes, ok := asyncRes.(*result.AsyncResult); ok {
		return NewAsyncResponse(asyncRes, nil, req, sleepDuration), nil
	} else {
		return nil, producer.ErrInvalidResponse(errors.New("Invalid Response"))
	}
}

// options[sleepDuration] = int
func (p *Producer) ProduceWithContext(ctx context.Context, req producer.Request, options map[string]interface{}) (producer.AsyncResponse, error) {
	sig, err := p.requestToSignature(req)
	if err != nil {
		return nil, err
	}

	sleepDuration := time.Duration(0)
	if options != nil {
		if val, ok := options["sleepDuration"]; ok {
			if i, ok := val.(int); ok {
				sleepDuration = time.Duration(i)
			}
		}
	}

	asyncRes, err := reHystrix.ReHystrixWithRes(
		func() (interface{}, error) { return p.server.SendTaskWithContext(ctx, sig) },
		p.commandName, p.reConf)
	if err != nil {
		return nil, err
	}

	if asyncRes, ok := asyncRes.(*result.AsyncResult); ok {
		return NewAsyncResponse(asyncRes, nil, req, sleepDuration), nil
	} else {
		return nil, producer.ErrInvalidResponse(errors.New("Invalid Response"))
	}
}

func (p *Producer) requestToSignature(req producer.Request) (*tasks.Signature, error) {
	sig, err := pkgMachineryTasks.NewSignature(req.Name(), *req.Args()...)

	if err != nil {
		return nil, err
	}

	if req.ETA() != nil {
		sig.ETA = req.ETA()
	}

	if req.RetryCount() != 0 {
		sig.RetryCount = req.RetryCount()
	}

	if req.RetryDelay() > 0 {
		sig.RetryCount = int(req.RetryDelay() / time.Second)
	}
	if req.RoutingKey() != "" {
		sig.RoutingKey = req.RoutingKey()
	}

	if req.OnError() != nil {
		onError := make([]*tasks.Signature, len(req.OnError()))

		for i, r := range req.OnError() {
			if fn, err := p.requestToSignature(r); err != nil {
				return nil, err
			} else {
				onError[i] = fn
			}
		}

		sig.OnError = onError
	}

	if req.OnSuccess() != nil {
		onSuccess := make([]*tasks.Signature, len(req.OnSuccess()))

		for i, r := range req.OnSuccess() {
			if fn, err := p.requestToSignature(r); err != nil {
				return nil, err
			} else {
				onSuccess[i] = fn
			}
		}

		sig.OnSuccess = onSuccess
	}

	return sig, nil
}

func NewAsyncResponse(asyncResult *result.AsyncResult, taskState *tasks.TaskState, req producer.Request, sleepDuration time.Duration) *AsyncResponse {
	ar := &AsyncResponse{asyncResult: asyncResult, taskState: taskState, req: req, sleepDuration: sleepDuration}
	ar.Sync()
	return ar
}

type AsyncResponse struct {
	asyncResult   *result.AsyncResult
	taskState     *tasks.TaskState
	req           producer.Request
	sleepDuration time.Duration
}

func (ar *AsyncResponse) UUID() string {
	return ar.asyncResult.Signature.UUID
}

func (ar *AsyncResponse) Request() producer.Request {
	return ar.req
}

func (ar *AsyncResponse) Sync() error {
	ar.taskState = ar.asyncResult.GetState()
	return nil
}

func (ar *AsyncResponse) Status() string {
	switch ar.taskState.State {
	case tasks.StatePending:
		return producer.ProduceStatusInit
	case tasks.StateReceived:
		return producer.ProduceStatusQueued
	case tasks.StateStarted:
		return producer.ProduceStatusInProgress
	case tasks.StateRetry:
		return producer.ProduceStatusRetry
	case tasks.StateSuccess:
		return producer.ProduceStatusSuccess
	case tasks.StateFailure:
		return producer.ProduceStatusFailure
	default:
		return producer.ProduceStatusUnknown
	}
}

// If in final status (success/failure)
func (ar *AsyncResponse) IsCompleted() bool {
	return ar.IsSuccess() || ar.IsFailure()
}

// Check if status is success
func (ar *AsyncResponse) IsSuccess() bool {
	return ar.Status() == producer.ProduceStatusSuccess
}

// Check if status is failure
func (ar *AsyncResponse) IsFailure() bool {
	return ar.Status() == producer.ProduceStatusFailure
}

func (ar *AsyncResponse) Error() (error, error) {
	if ar.taskState.Error == "" {
		return nil, nil
	}

	return errors.New(ar.taskState.Error), nil
}

func (ar *AsyncResponse) Subscribe(timeoutDuration time.Duration) error {
	timeout := time.NewTimer(timeoutDuration)

	for {
		select {
		case <-timeout.C:
			return producer.ErrTimeoutReached(errors.New("Timeout reached"))
		default:
			results, err := ar.asyncResult.Touch()

			if results == nil && err == nil {
				time.Sleep(ar.sleepDuration)
			} else {
				return err
			}
		}
	}
}

func (ar *AsyncResponse) Result() ([]producer.Result, error) { //ErrNotSuccessful, ErrCannotSyncState
	if !ar.taskState.IsSuccess() {
		return nil, producer.ErrNotSuccessful(errors.New("Not Successful response"))
	}

	res := ar.taskState.Results

	results := make([]producer.Result, len(res))
	for i, r := range res {
		results[i] = pkgResult.NewResult(r.Type, r.Value)
	}

	return results, nil
}
