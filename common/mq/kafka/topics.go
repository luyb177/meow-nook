package kafka

import (
	"fmt"

	"github.com/google/uuid"
)

type Topics struct {
	Pending string
	Retry   string
	DLQ     string
}

func BuildTopics(serviceName string) Topics {
	return Topics{
		Pending: fmt.Sprintf("%s.task.pending", serviceName),
		Retry:   fmt.Sprintf("%s.task.retry", serviceName),
		DLQ:     fmt.Sprintf("%s.task.dlq", serviceName),
	}
}

// BuildTaskID constructs a TaskID in the format "{type}:{biz}:{uuid}".
func BuildTaskID(taskType, biz string) string {
	return fmt.Sprintf("%s:%s:%s", taskType, biz, uuid.New().String())
}
