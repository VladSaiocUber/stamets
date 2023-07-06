package stamets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeout(t *testing.T) {
	m, ok := TaskWIthTimeout(time.Second, func() BaseMetrics[string] {
		return BaseMetrics[string]{
			Payload: "I made it",
		}
	})
	require.True(t, ok)
	msg, err := m.Unpack()
	require.NoError(t, err)
	require.Equal(t, msg, "I made it")

	m, ok = TaskWIthTimeout(100*time.Millisecond, func() BaseMetrics[string] {
		select {
		case <-time.After(3 * time.Second):
			t.FailNow()
		}

		return BaseMetrics[string]{
			Payload: "I made it",
		}
	})
	require.Zero(t, m)
	require.False(t, ok)
}
