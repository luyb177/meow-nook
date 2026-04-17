package kafka

type Envelope struct {
	TaskID    string `json:"task_id"`
	Retry     int    `json:"retry"`
	MaxRetry  int    `json:"max_retry"`
	CreatedAt int64  `json:"created_at"`

	// retry mover 使用：下一次允许运行的时间戳（秒）
	NextRunAt int64 `json:"next_run_at"`

	// 业务 payload
	Data []byte `json:"data"`

	// DLQ 信息（可选）
	DLQReason string `json:"dlq_reason,omitempty"`
	DLQStage  string `json:"dlq_stage,omitempty"`
	Error     string `json:"error,omitempty"`
}
