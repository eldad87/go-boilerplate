package consumer

type Consumer interface {
	RegisterTask(name string, taskFunc interface{}) error
}
