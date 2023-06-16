package stamets

import (
	"fmt"
	"time"
)

type BaseMetrics[T any] struct {
	// Time it took to perform task
	time.Duration
	error

	Payload T
}

// Metrics is an interface implemented by all metrics structures.
type Metrics interface {
	fmt.Stringer
	error

	UnpackAny() (any, error)
	Ok() bool
}

// Unpack extracts the payload wrapped in the metrics,
// and an error, if one was produced.
func (m BaseMetrics[T]) Unpack() (T, error) {
	return m.Payload, m.error
}

// UnpackAny extracts the payload wrapped in the metrics,
// and an error, if one was produced. UnpackAny allows
// the implementation of the Metrics interface.
func (m BaseMetrics[T]) UnpackAny() (any, error) {
	return m.Payload, m.error
}

// Ok checks whether the desired function failed to execute.
func (m BaseMetrics[T]) Ok() bool {
	return m.error == nil
}

func (m BaseMetrics[T]) String() string {
	return fmt.Sprintf(`
BASE METRICS:
- Duration: %fs
`,
		m.Duration.Seconds())
}
