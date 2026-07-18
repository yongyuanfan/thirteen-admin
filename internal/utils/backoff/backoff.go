package backoff

import (
	"context"
	"time"
)

// Sleep 按指数退避等待下一次重试。
func Sleep(ctx context.Context, baseDelay, maxDelay time.Duration, attempt int) error {
	delay := baseDelay
	if delay <= 0 {
		delay = time.Millisecond
	}
	if maxDelay > 0 && delay > maxDelay {
		delay = maxDelay
	}

	for i := 1; i < attempt; i++ {
		delay *= 2
		if maxDelay > 0 && delay >= maxDelay {
			delay = maxDelay
			break
		}
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
