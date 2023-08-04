package stamets

import (
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
)

// mockInterfaceMaps creates a map of size n, where the keys are pointers wrapped in an interface.
// The map keys are merely representative and are only meant to inflate the map size.
func mockInterfaceMaps[T comparable, U any](fresh func() T, n int) map[T]U {
	qs := make(map[T]U)
	for i := 0; i < n; i++ {
		var u U
		qs[fresh()] = u
	}
	return qs
}

// mockPointerMaps creates a map of size n, where the keys are pointers.
// The map keys are merely representative and are only meant to inflate the map size.
func mockPointerMaps[T any, U any](n int) map[*T]U {
	return mockInterfaceMaps[*T, U](func() *T {
		return new(T)
	}, n)
}

func getRowValue(head, l string) string {
	split := strings.Split(l, head)
	return strings.TrimSpace(split[len(split)-1])
}

// UnparsePTAResultsFromReader reads the whole content of a reader and then
// unparses it line by line. Any reconstructed PTAMetrics values are aggregated
// and then returned in a slice.
func UnparsePTAResultsFromReader(r io.Reader) []PTAMetrics {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil
	}

	// Relevant rows
	const (
		TITLE    = "PTA METRICS"
		DURATION = "- Duration:"
		QUERIES  = "- Number of PTA queries:"
		IQUERIES = "- Number of indirect PTA queries:"
		P50      = "- P50 points-to set size:"
		P90      = "- P90 points-to set size:"
		P99      = "- P99 points-to set size:"
		MAX      = "- Max points-to set size:"
		MODE     = "- Most common points-to set size:"
	)

	content := string(bs)
	results := make([]PTAMetrics, 0, 1)

	fresh := func() PTAMetrics {
		return PTAMetrics{
			BaseMetrics: BaseMetrics[*pointer.Result]{
				Payload: new(pointer.Result),
			},
		}
	}
	current, unparsing := fresh(), false
	flush := func() {
		results = append(results, current)
		current = fresh()
		unparsing = false
	}

	for _, l := range strings.Split(content, "\n") {
		l = strings.TrimSpace(l)
		if !unparsing && l == TITLE {
			unparsing = true
		} else if unparsing {
			switch {
			case strings.HasPrefix(l, TITLE):
				current = fresh()
			case strings.HasPrefix(l, DURATION):
				if t, err := time.ParseDuration(getRowValue(DURATION, l) + "s"); err == nil {
					current.Duration = t
				}
			case strings.HasPrefix(l, QUERIES):
				if v, err := strconv.Atoi(getRowValue(QUERIES, l)); err == nil {
					current.Queries = v
				}
			case strings.HasPrefix(l, IQUERIES):
				if v, err := strconv.Atoi(getRowValue(IQUERIES, l)); err == nil {
					current.IndirectQueries = v
				}
			case strings.HasPrefix(l, P50):
				if v, err := strconv.Atoi(getRowValue(P50, l)); err == nil {
					current.PointsToSetSizeP50 = v
				}
			case strings.HasPrefix(l, P90):
				if v, err := strconv.Atoi(getRowValue(P90, l)); err == nil {
					current.PointsToSetSizeP90 = v
				}
			case strings.HasPrefix(l, P99):
				if v, err := strconv.Atoi(getRowValue(P99, l)); err == nil {
					current.PointsToSetSizeP99 = v
				}
			case strings.HasPrefix(l, MAX):
				if v, err := strconv.Atoi(getRowValue(MAX, l)); err == nil {
					current.PointsToSetSizeMax = v
				}
			case strings.HasPrefix(l, MODE):
				if v, err := strconv.Atoi(getRowValue(MODE, l)); err == nil {
					current.PointsToSetSizeMode = v
				}
				flush()
			}
		}
	}

	return results
}

// UnparseCallGraphMetricsFromReader reads the whole content of a reader and then
// unparses it line by line. Any reconstructed CallGraphMetrics values are aggregated
// and then returned in a slice.
func UnparseCallGraphMetricsFromReader(r io.Reader) []CallGraphMetrics {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil
	}

	// Relevant rows
	const (
		TITLE     = "CALL GRAPH METRICS"
		FUNCTIONS = "- Number of functions:"
		OUT       = "Call site out-degree metrics:"
		IN        = "Callee in-degree metrics:"
		P50       = "- P50:"
		P90       = "- P90:"
		P99       = "- P99:"
		MAX       = "- Max:"
		OUTMODE   = "- Most common out-degree:"
		INMODE    = "- Most common in-degree:"
	)

	content := string(bs)
	results := make([]CallGraphMetrics, 0, 1)

	fresh := func() CallGraphMetrics {
		return CallGraphMetrics{
			BaseMetrics: BaseMetrics[*callgraph.Graph]{
				Payload: new(callgraph.Graph),
			},
		}
	}
	current, unparsing, in, out := fresh(), false, false, false
	flush := func() {
		results = append(results, current)
		current = fresh()
		unparsing = false
		in = false
		out = false
	}

	for _, l := range strings.Split(content, "\n") {
		l = strings.TrimSpace(l)
		if !unparsing && l == TITLE {
			unparsing = true
		} else if unparsing {
			switch {
			case strings.HasPrefix(l, TITLE):
				current = fresh()
				in = false
				out = false
			case strings.HasPrefix(l, FUNCTIONS):
				if v, err := strconv.Atoi(getRowValue(FUNCTIONS, l)); err == nil {
					current.Functions = v
				}
			case strings.HasPrefix(l, OUT):
				out = true
				in = false
			case strings.HasPrefix(l, IN):
				in = true
				out = false
			case strings.HasPrefix(l, P50):
				if v, err := strconv.Atoi(getRowValue(P50, l)); err == nil {
					if out && !in {
						current.OutDegreeP50 = v
					} else if in && !out {
						current.InDegreeP50 = v
					}
				}
			case strings.HasPrefix(l, P90):
				if v, err := strconv.Atoi(getRowValue(P90, l)); err == nil {
					if out && !in {
						current.OutDegreeP90 = v
					} else if in && !out {
						current.InDegreeP90 = v
					}
				}
			case strings.HasPrefix(l, P99):
				if v, err := strconv.Atoi(getRowValue(P99, l)); err == nil {
					if out && !in {
						current.OutDegreeP99 = v
					} else if in && !out {
						current.InDegreeP99 = v
					}
				}
			case strings.HasPrefix(l, MAX):
				if v, err := strconv.Atoi(getRowValue(MAX, l)); err == nil {
					if out && !in {
						current.OutDegreeMax = v
					} else if in && !out {
						current.InDegreeMax = v
					}
				}
			case strings.HasPrefix(l, OUTMODE):
				if v, err := strconv.Atoi(getRowValue(OUTMODE, l)); err == nil {
					current.OutDegreeMode = v
				}
			case strings.HasPrefix(l, INMODE):
				if v, err := strconv.Atoi(getRowValue(INMODE, l)); err == nil {
					current.InDegreeMode = v
				}
				flush()
			}
		}
	}

	return results
}
