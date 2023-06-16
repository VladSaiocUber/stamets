package stamets

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
)

// PTAMetrics aggregates important metrics about the Points-To
// analysis.
type PTAMetrics struct {
	BaseMetrics[*pointer.Result]

	// Aliasing metrics
	PointsToSetSizeMax  int
	PointsToSetSizeP50  int
	PointsToSetSizeP90  int
	PointsToSetSizeP99  int
	PointsToSetSizeMode int
}

func (m PTAMetrics) String() string {
	return fmt.Sprintf(`
PTA METRICS
- Duration: %fs
- Number of PTA queries: %d
- Number of indirect PTA queries: %d
- P50 points-to set size: %d
- P90 points-to set size: %d
- P99 points-to set size: %d
- Max points-to set size: %d
- Most common points-to set size: %d
`,
		m.Duration.Seconds(),
		len(m.Payload.Queries),
		len(m.Payload.IndirectQueries),
		m.PointsToSetSizeP50,
		m.PointsToSetSizeP90,
		m.PointsToSetSizeP99,
		m.PointsToSetSizeMax,
		m.PointsToSetSizeMode,
	)
}

// Analyze runs the points-to analysis with the given configuration,
// collecting metrics i.e., duration and information about the call graph.
func Analyze(config *pointer.Config) PTAMetrics {
	start := time.Now()

	res, err := pointer.Analyze(config)
	if err != nil {
		return PTAMetrics{
			BaseMetrics: BaseMetrics[*pointer.Result]{
				error: err,
			},
		}
	}

	m := PTAMetrics{
		BaseMetrics: BaseMetrics[*pointer.Result]{
			Duration: time.Since(start),
			Payload:  res,
		},
	}

	m = m.PointsToSetMetrics()

	return m
}

func callgraphMetrics(res *pointer.Result) (m PTAMetrics) {
	if res == nil {
		return
	}

	return
}

func visitCallgraph(cg *callgraph.Graph, f func(n *callgraph.Node)) {
	if cg == nil || f == nil {
		return
	}

	visited := make(map[*callgraph.Node]struct{})

	var visit func(*callgraph.Node)
	visit = func(n *callgraph.Node) {
		if _, ok := visited[n]; ok {
			return
		}
		visited[n] = struct{}{}

		f(n)
		for _, e := range n.Out {
			visit(e.Callee)
		}
	}

	visit(cg.Root)
}

// PointsToSetMetrics computes metrics about the sizes of points-to sets.
func (m PTAMetrics) PointsToSetMetrics() PTAMetrics {
	// Max size points-to set.
	ptSizes := make([]int, 0, len(m.Payload.Queries))
	for _, pt := range m.Payload.Queries {
		ptSizes = append(ptSizes, len(pt.PointsTo().Labels()))

		if m.PointsToSetSizeMax < len(pt.PointsTo().Labels()) {
			m.PointsToSetSizeMax = len(pt.PointsTo().Labels())
		}
	}

	slices.Sort(ptSizes)
	m.PointsToSetSizeMode = mode(ptSizes)
	m.PointsToSetSizeP50 = p50(ptSizes)
	m.PointsToSetSizeP90 = p90(ptSizes)
	m.PointsToSetSizeP99 = p99(ptSizes)

	return m
}
