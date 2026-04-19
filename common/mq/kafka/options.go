package kafka

import "time"

type DispatchOptions struct {
	MaxRetry    *int
	BaseBackoff *time.Duration
}

type DispatchOption func(*DispatchOptions)

func WithMaxRetry(n int) DispatchOption {
	return func(o *DispatchOptions) { o.MaxRetry = &n }
}

func WithBaseBackoff(d time.Duration) DispatchOption {
	return func(o *DispatchOptions) { o.BaseBackoff = &d }
}
