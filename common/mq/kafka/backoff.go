package kafka

import "time"

func ExponentialBackoff(base time.Duration, retry int) time.Duration {
	if retry <= 1 {
		return base
	}
	return base * time.Duration(1<<(retry-1))
}
