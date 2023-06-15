package stamets

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestP50(t *testing.T) {
	require.Equal(t, 0, p50([]int{}))
	require.Equal(t, 1, p50([]int{1}))
	require.Equal(t, 2, p50([]int{1, 2}))
	require.Equal(t, 1, p50([]int{1, 1, 2}))
}

func TestP90(t *testing.T) {
	require.Equal(t, 0, p90([]int{}))
	require.Equal(t, 1, p90([]int{1}))
	require.Equal(t, 2, p90([]int{1, 2}))
	require.Equal(t, 2, p90([]int{1, 1, 2}))

	var values []int
	for i := 1; i <= 10; i++ {
		values = append(values, i)
	}
	require.Equal(t, 10, p90(values))

	values = []int{}
	for i := 1; i <= 100; i++ {
		values = append(values, i)
	}
	require.Equal(t, 91, p90(values))

	values = []int{}
	for i := 1; i <= 1000; i++ {
		values = append(values, i)
	}
	require.Equal(t, 901, p90(values))
}

func TestP99(t *testing.T) {
	require.Equal(t, 0, p99([]int{}))
	require.Equal(t, 1, p99([]int{1}))
	require.Equal(t, 2, p99([]int{1, 2}))
	require.Equal(t, 2, p99([]int{1, 1, 2}))

	var values []int
	for i := 1; i <= 10; i++ {
		values = append(values, i)
	}
	require.Equal(t, 10, p99(values))

	values = []int{}
	for i := 1; i <= 100; i++ {
		values = append(values, i)
	}
	require.Equal(t, 100, p99(values))

	values = []int{}
	for i := 1; i <= 1000; i++ {
		values = append(values, i)
	}
	require.Equal(t, 991, p99(values))
}

func TestMode(t *testing.T) {
	require.Equal(t, 0, mode([]int{}))
	require.Equal(t, 1, mode([]int{0, 0, 1, 1, 1}))
	require.Equal(t, 1, mode([]int{0, 1, 1, 1, 2}))
	require.Equal(t, 0, mode([]int{0, 0, 0, 0, 1, 2}))
}
