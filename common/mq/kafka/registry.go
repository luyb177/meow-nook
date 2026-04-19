package kafka

import (
	"context"
	"encoding/json"
	"fmt"
)

type TaskHandler interface {
	Handle(ctx context.Context, env *Envelope) error
}

type HandlerFunc func(ctx context.Context, env *Envelope) error

func (f HandlerFunc) Handle(ctx context.Context, env *Envelope) error {
	return f(ctx, env)
}

type Registry struct {
	m map[string]HandlerFunc
}

func NewRegistry() *Registry {
	return &Registry{
		m: make(map[string]HandlerFunc),
	}
}

func (r *Registry) Register(taskType string, fn HandlerFunc) {
	if taskType == "" {
		panic("taskType is empty")
	}
	if fn == nil {
		panic("handler func is nil")
	}
	if _, ok := r.m[taskType]; ok {
		panic(fmt.Sprintf("duplicate handler for type=%s", taskType))
	}
	r.m[taskType] = fn
}

func (r *Registry) All() map[string]HandlerFunc {
	return r.m
}

func Decode[T any](env *Envelope, out *T) error {
	return json.Unmarshal(env.Data, out)
}
