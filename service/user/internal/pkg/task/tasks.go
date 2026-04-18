// Package task defines Kafka task types for the user service.
// TaskID format: {type}:{biz}:{uuid}  (built by kafka.BuildTaskID).
package task

import "encoding/json"

// Task type constants (dot-separated, routed by Envelope.Type).
const (
	TypeSendEmail = "user.send_email"
)

// SendEmailTask is a task that sends a transactional email to a user.
type SendEmailTask struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (t *SendEmailTask) Type() string { return TypeSendEmail }
func (t *SendEmailTask) Biz() string  { return t.To }
func (t *SendEmailTask) Payload() ([]byte, error) {
	return json.Marshal(t)
}
