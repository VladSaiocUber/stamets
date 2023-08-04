package stamets

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/pointer"
)

func TestGetPTAResultsFromReader(t *testing.T) {
	require.Empty(t, UnparsePTAResultsFromReader(strings.NewReader("")))

	expectedMetrics := []PTAMetrics{{
		BaseMetrics: BaseMetrics[*pointer.Result]{
			Duration: time.Millisecond * 500,
		},
		Queries:             100,
		IndirectQueries:     10,
		PointsToSetSizeP50:  1,
		PointsToSetSizeP90:  2,
		PointsToSetSizeP99:  3,
		PointsToSetSizeMax:  4,
		PointsToSetSizeMode: 5,
	}, {
		BaseMetrics: BaseMetrics[*pointer.Result]{
			Duration: time.Second,
		},
		Queries:             200,
		IndirectQueries:     20,
		PointsToSetSizeP50:  6,
		PointsToSetSizeP90:  7,
		PointsToSetSizeP99:  8,
		PointsToSetSizeMax:  9,
		PointsToSetSizeMode: 10,
	}}

	resultMetrics := UnparsePTAResultsFromReader(strings.NewReader(`
	PTA METRICS
	- Duration: 0.5
	- Number of PTA queries: 100
	- Number of indirect PTA queries: 10
	- P50 points-to set size: 1
	- P90 points-to set size: 2
	- P99 points-to set size: 3
	- Max points-to set size: 4
	- Most common points-to set size: 5
		`))

	compare := func(i int) {
		require.Equal(t, expectedMetrics[i].Duration, resultMetrics[i].Duration)
		require.Equal(t, expectedMetrics[i].PointsToSetSizeP50, resultMetrics[i].PointsToSetSizeP50)
		require.Equal(t, expectedMetrics[i].PointsToSetSizeP90, resultMetrics[i].PointsToSetSizeP90)
		require.Equal(t, expectedMetrics[i].PointsToSetSizeP99, resultMetrics[i].PointsToSetSizeP99)
		require.Equal(t, expectedMetrics[i].PointsToSetSizeMax, resultMetrics[i].PointsToSetSizeMax)
		require.Equal(t, expectedMetrics[i].PointsToSetSizeMode, resultMetrics[i].PointsToSetSizeMode)
		require.NotNil(t, resultMetrics[i].Payload)
		require.Equal(t, expectedMetrics[i].Queries, resultMetrics[i].Queries)
		require.Equal(t, expectedMetrics[i].IndirectQueries, resultMetrics[i].IndirectQueries)
	}

	require.Len(t, resultMetrics, 1)
	compare(0)

	resultMetrics = UnparsePTAResultsFromReader(strings.NewReader(`
	PTA METRICS
	- Duration: 0.6
	PTA METRICS
	- Duration: 0.5
	- Number of PTA queries: 100
	- Number of indirect PTA queries: 10
	- P50 points-to set size: 1
	- P90 points-to set size: 2
	- P99 points-to set size: 3
	- Max points-to set size: 4
	- Most common points-to set size: 5
		`))

	require.Len(t, resultMetrics, 1)
	compare(0)

	resultMetrics = UnparsePTAResultsFromReader(strings.NewReader(`
	PTA METRICS
	- Duration: 0.5
	- Number of PTA queries: 100
	- Number of indirect PTA queries: 10
	- P50 points-to set size: 1
	- P90 points-to set size: 2
	- P99 points-to set size: 3
	- Max points-to set size: 4
	- Most common points-to set size: 5

	PTA METRICS
	- Duration: 1
	- Number of PTA queries: 200
	- Number of indirect PTA queries: 20
	- P50 points-to set size: 6
	- P90 points-to set size: 7
	- P99 points-to set size: 8
	- Max points-to set size: 9
	- Most common points-to set size: 10
		`))

	require.Len(t, resultMetrics, 2)
	compare(0)
	compare(1)
}

func TestGetCallGraphResultsFromReader(t *testing.T) {
	require.Empty(t, UnparseCallGraphMetricsFromReader(strings.NewReader("")))

	expectedMetrics := []CallGraphMetrics{{
		Functions:     10,
		InDegreeP50:   1,
		InDegreeP90:   2,
		InDegreeP99:   3,
		InDegreeMax:   4,
		InDegreeMode:  5,
		OutDegreeP50:  6,
		OutDegreeP90:  7,
		OutDegreeP99:  8,
		OutDegreeMax:  9,
		OutDegreeMode: 10,
	}, {
		Functions:     20,
		InDegreeP50:   11,
		InDegreeP90:   12,
		InDegreeP99:   13,
		InDegreeMax:   14,
		InDegreeMode:  15,
		OutDegreeP50:  16,
		OutDegreeP90:  17,
		OutDegreeP99:  18,
		OutDegreeMax:  19,
		OutDegreeMode: 20,
	}}

	resultMetrics := UnparseCallGraphMetricsFromReader(strings.NewReader(`
	CALL GRAPH METRICS
	- Number of functions: 10
	Call site out-degree metrics:
	- P50: 6
	- P90: 7
	- P99: 8
	- Max: 9
	- Most common out-degree: 10
	Callee in-degree metrics:
	- P50: 1
	- P90: 2
	- P99: 3
	- Max: 4
	- Most common in-degree: 5
		`))

	compare := func(i int) {
		require.Equal(t, expectedMetrics[i].InDegreeP50, resultMetrics[i].InDegreeP50)
		require.Equal(t, expectedMetrics[i].InDegreeP90, resultMetrics[i].InDegreeP90)
		require.Equal(t, expectedMetrics[i].InDegreeP99, resultMetrics[i].InDegreeP99)
		require.Equal(t, expectedMetrics[i].InDegreeMax, resultMetrics[i].InDegreeMax)
		require.Equal(t, expectedMetrics[i].InDegreeMode, resultMetrics[i].InDegreeMode)
		require.Equal(t, expectedMetrics[i].OutDegreeP50, resultMetrics[i].OutDegreeP50)
		require.Equal(t, expectedMetrics[i].OutDegreeP90, resultMetrics[i].OutDegreeP90)
		require.Equal(t, expectedMetrics[i].OutDegreeP99, resultMetrics[i].OutDegreeP99)
		require.Equal(t, expectedMetrics[i].OutDegreeMax, resultMetrics[i].OutDegreeMax)
		require.Equal(t, expectedMetrics[i].OutDegreeMode, resultMetrics[i].OutDegreeMode)
		require.NotNil(t, resultMetrics[i].Payload)
		require.Equal(t, expectedMetrics[i].Functions, resultMetrics[i].Functions)
	}

	require.Len(t, resultMetrics, 1)
	compare(0)

	resultMetrics = UnparseCallGraphMetricsFromReader(strings.NewReader(`
	CALL GRAPH METRICS
	- Number of functions: 20
	Call site out-degree metrics:
	- P50: 60
	- P90: 70
	- P99: 80
	- Max: 90
	CALL GRAPH METRICS
	- Number of functions: 10
	Call site out-degree metrics:
	- P50: 6
	- P90: 7
	- P99: 8
	- Max: 9
	- Most common out-degree: 10
	Callee in-degree metrics:
	- P50: 1
	- P90: 2
	- P99: 3
	- Max: 4
	- Most common in-degree: 5
		`))

	require.Len(t, resultMetrics, 1)
	compare(0)

	resultMetrics = UnparseCallGraphMetricsFromReader(strings.NewReader(`
	CALL GRAPH METRICS
	- Number of functions: 10
	Call site out-degree metrics:
	- P50: 6
	- P90: 7
	- P99: 8
	- Max: 9
	- Most common out-degree: 10
	Callee in-degree metrics:
	- P50: 1
	- P90: 2
	- P99: 3
	- Max: 4
	- Most common in-degree: 5

	CALL GRAPH METRICS
	- Number of functions: 20
	Call site out-degree metrics:
	- P50: 16
	- P90: 17
	- P99: 18
	- Max: 19
	- Most common out-degree: 20
	Callee in-degree metrics:
	- P50: 11
	- P90: 12
	- P99: 13
	- Max: 14
	- Most common in-degree: 15
		`))

	require.Len(t, resultMetrics, 2)
	compare(0)
	compare(1)
}
