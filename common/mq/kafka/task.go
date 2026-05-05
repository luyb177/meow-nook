package kafka

type Task interface {
	Type() string
	Biz() string
	Payload() ([]byte, error)
}
