package kafka

import "time"

type Envelope struct {
	TaskID    string `json:"task_id"`
	Retry     int    `json:"retry"`
	MaxRetry  int    `json:"max_retry"`
	CreatedAt int64  `json:"created_at"`

	// 下一次允许运行的时间戳（秒或毫秒都行，统一即可）
	NextRunAt int64 `json:"next_run_at"`

	// 业务 payload
	Data []byte `json:"data"`

	// 可选：死信信息
	DLQReason string `json:"dlq_reason,omitempty"`
	DLQStage  string `json:"dlq_stage,omitempty"`
	Error     string `json:"error,omitempty"`
}

func (e *Envelope) NextRunTime() time.Time {
	return time.Unix(e.NextRunAt, 0)
}
