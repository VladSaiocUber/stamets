package stamets

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
