package stamets

import (
	"time"
)

// BaseMetrics provides a foundation for other metrics to build upon.
// Basic information include the duration to perform a task, whether
// the task produced an error, and the result produced by the task.
type BaseMetrics[T any] struct {
	// Time it took to perform task
	time.Duration
	// Metrics produced error
	err error

	Payload T
}

type Metrics interface {
	UnpackAny() (any, error)
	Ok() bool
}

// Unpack extracts the payload wrapped in the metrics,
// and an error, if one was produced.
func (m BaseMetrics[T]) Unpack() (T, error) {
	return m.Payload, m.err
}

// UnpackAny extracts the payload wrapped in the metrics,
// and an error, if one was produced. UnpackAny allows
// the implementation of the Metrics interface.
func (m BaseMetrics[T]) UnpackAny() (any, error) {
	return m.Payload, m.err
}

// Ok checks whether the desired function failed to execute.
func (m BaseMetrics[T]) Ok() bool {
	return m.err == nil
}
