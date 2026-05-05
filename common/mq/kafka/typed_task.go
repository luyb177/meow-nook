package kafka

import (
	"encoding/json"
	"errors"
)

type TypedTask[T any] struct {
	TaskType string
	BizKey   string
	Data     T
}

func NewTypedTask[T any](taskType, biz string, data T) TypedTask[T] {
	return TypedTask[T]{TaskType: taskType, BizKey: biz, Data: data}
}

func (t TypedTask[T]) Type() string {
	return t.TaskType
}
func (t TypedTask[T]) Biz() string {
	return t.BizKey
}

func (t TypedTask[T]) Payload() ([]byte, error) {
	if t.TaskType == "" {
		return nil, errors.New("kafka task type is empty")
	}
	return json.Marshal(t.Data)
}
