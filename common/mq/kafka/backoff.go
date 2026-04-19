package kafka

import "time"

// ExponentialBackoff  base * 2^(retry-1)，retry 从 1 开始更符合直觉
func ExponentialBackoff(base time.Duration, retry int) time.Duration {
	if retry <= 1 {
		return base
	}
	return base * time.Duration(1<<(retry-1))
}
