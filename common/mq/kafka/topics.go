package kafka

import "fmt"

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
