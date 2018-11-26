package request

import (
	"github.com/eldad87/go-boilerplate/src/pkg/task/producer"
	"github.com/mitchellh/mapstructure"
	"time"
)

func NewRequest(name string, options map[string]interface{}, args ...interface{}) (producer.Request, error) {
	req := &Request{name: name, args: &args}
	if options != nil {
		mapstructure.Decode(options, req)
	}

	return req, nil
}

type Request struct {
	uuid       string
	name       string
	routingKey string
	eta        *time.Time
	args       *[]interface{}
	retryCount int
	retryDelay time.Duration
	onSuccess  []producer.Request
	onError    []producer.Request
}

func (t *Request) Name() string {
	return t.name
}

func (t *Request) SetName(name string) {
	t.name = name
}

func (t *Request) Args() *[]interface{} {
	return t.args
}
func (t *Request) SetArgs(args *[]interface{}) error {
	t.args = args
	return nil
}

func (t *Request) SetETA(time *time.Time) error {
	t.eta = time
	return nil
}
func (t *Request) ETA() *time.Time {
	return t.eta
}

func (t *Request) SetRetryCount(count int) error {
	t.retryCount = count
	return nil
}
func (t *Request) RetryCount() int {
	return t.retryCount
}

func (t *Request) SetRetryDelay(delay time.Duration) error {
	t.retryDelay = delay
	return nil
}
func (t *Request) RetryDelay() time.Duration {
	return t.retryDelay
}

func (t *Request) SetRoutingKey(routingKey string) error {
	t.routingKey = routingKey
	return nil
}
func (t *Request) RoutingKey() string {
	return t.routingKey
}

func (t *Request) SetOnError(requests []producer.Request) error {
	t.onError = requests
	return nil
}

func (t *Request) OnError() []producer.Request {
	return t.onError
}

func (t *Request) SetOnSuccess(requests []producer.Request) error {
	t.onSuccess = requests
	return nil
}

func (t *Request) OnSuccess() []producer.Request {
	return t.onSuccess
}
