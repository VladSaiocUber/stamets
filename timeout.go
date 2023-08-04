package stamets

import (
	"context"
	"time"
)

// TaskWithTimeout performs a task with collectible metrics in the alloted time limit.
// It returns the Metrics-wrapped result and 'true' if the operation was completed before the timeout,
// or the zero Metrics value and 'false' otherwise.
func TaskWithTimeout[T Metrics](timeout time.Duration, f func() T) (T, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ch := make(chan T, 1)

	go func() {
		ch <- f()
	}()

	select {
	case res := <-ch:
		return res, true
	case <-ctx.Done():
		var x T
		return x, false
	}
}
