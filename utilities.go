package stamets

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Series is a sequence of orderable values.
type Series[T constraints.Ordered] []T

// P50 returns the smallest value in the 50'th percentile of the series.
func (s Series[T]) P50() T {
	return p50(s)
}

// P90 returns the smallest value in the 90'th percentile of the series.
func (s Series[T]) P90() T {
	return p90(s)
}

// P99 returns the smallest value in the 99'th percentile of the series.
func (s Series[T]) P99() T {
	return p99(s)
}

// Max returns the largest value in the series.
func (s Series[T]) Max() T {
	if len(s) == 0 {
		var t T
		return t
	}
	return s[len(s)-1]
}

// Mode returns the most common value in the series.
func (s Series[T]) Mode() T {
	return mode(s)
}

// Order returns a sorted copy of the original series, but where
// all elements have been ordered.
func (s Series[T]) Order() []T {
	s2 := make(Series[T], 0, len(s))
	for _, o := range s {
		s2 = append(s2, o)
	}
	slices.Sort[T](s2)
	return s2
}

// MakeSeries creates an ordered series from a data list, given a transformation
// function over data in the list.
func MakeSeries[T any, U constraints.Ordered](get func(T) U, ts ...T) Series[U] {
	s := make(Series[U], 0, len(ts))
	for _, t := range ts {
		s = append(s, get(t))
	}
	return s.Order()
}

func p50[T any](list []T) T {
	if len(list) == 0 {
		var t T
		return t
	}

	return list[len(list)/2]
}

func p90[T any](list []T) T {
	if len(list) == 0 {
		var t T
		return t
	}

	return list[len(list)*9/10]
}

func p99[T any](list []T) T {
	if len(list) == 0 {
		var t T
		return t
	}

	return list[len(list)*99/100]
}

func mode[T comparable](list []T) T {
	cardinality := make(map[T]int)

	for _, l := range list {
		cardinality[l]++
	}

	modeCount := 0
	var m T
	for x, count := range cardinality {
		if modeCount < count {
			m = x
			modeCount = count
		}
	}

	return m
}
