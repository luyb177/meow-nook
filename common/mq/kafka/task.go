package kafka

// Task represents a unit of work to be dispatched via Kafka.
// Type returns the dot-separated task type (e.g. "user.send_email").
// Biz returns the business identifier for the task (e.g. a user ID).
// Payload returns the serialised task data.
type Task interface {
	Type() string
	Biz() string
	Payload() ([]byte, error)
}
