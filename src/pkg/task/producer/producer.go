package producer

import (
	"context"
	"time"
)

const (
	ProduceStatusInit       = "INIT" // Initial state of a Request
	ProduceStatusQueued     = "QUEUED"
	ProduceStatusInProgress = "INPROGRESS"
	ProduceStatusRetry      = "RETRY"
	ProduceStatusUnknown    = "UNKNOWN"
	ProduceStatusSuccess    = "SUCCESS" // Final status
	ProduceStatusFailure    = "FAILURE" // Final status
)

type ErrInvalidResponse error
type ErrInvalidRequest error
type ErrTimeoutReached error
type ErrNotSuccessful error
type ErrCannotSyncState error

type Producer interface {
	NewRequest(name string, options map[string]interface{}, args ...interface{}) (Request, error)
	Produce(req Request, options map[string]interface{}) (AsyncResponse, error)
	ProduceWithContext(ctx context.Context, req Request, options map[string]interface{}) (AsyncResponse, error)
}

type Request interface {
	Name() string
	SetName(name string)

	Args() *[]interface{}
	SetArgs(args *[]interface{}) error

	SetETA(time *time.Time) error
	ETA() *time.Time
	SetRetryCount(count int) error
	RetryCount() int

	RetryDelay() time.Duration
	SetRetryDelay(time.Duration) error

	SetRoutingKey(routingKey string) error
	RoutingKey() string
	SetOnError(requests []Request) error // ErrInvalidRequest
	OnError() []Request
	SetOnSuccess(reqs []Request) error // ErrInvalidRequest
	OnSuccess() []Request
}

type AsyncResponse interface {
	UUID() string
	Status() string
	Result() ([]Result, error) //ErrNotSuccessful, ErrCannotSyncState
	Error() (error, error)
	Subscribe(timeout time.Duration) error // Subscriber for status update. ErrWaitTimeout, ErrCannotSyncState
	Sync() error                           // Manually check and sync latest State
	Request() Request
	IsCompleted() bool // If in final status (success/failure)
	IsSuccess() bool   // Check if status is success
	IsFailure() bool   // Check if status is failure
}

type Result interface{}
