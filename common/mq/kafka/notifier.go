package kafka

import "context"

type Notifier interface {
	Notify(ctx context.Context, env *Envelope) error
}
