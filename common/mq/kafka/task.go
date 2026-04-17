package kafka

type Task interface {
	ID() string
	Payload() ([]byte, error)
}
