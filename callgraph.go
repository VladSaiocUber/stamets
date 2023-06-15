package stamets

import (
	"fmt"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

type CallGraphMetrics struct {
	BaseMetrics[*callgraph.Graph]

	// Call graph out degree metrics
	OutDegreeMax  int
	OutDegreeP50  int
	OutDegreeP90  int
	OutDegreeP99  int
	OutDegreeMode int

	// Call graph in degree metrics
	InDegreeMax  int
	InDegreeP50  int
	InDegreeP90  int
	InDegreeP99  int
	InDegreeMode int
}

func (m CallGraphMetrics) String() string {
	return fmt.Sprintf(`
CALL GRAPH METRICS
- Number of functions: %d
Call site out-degree metrics:
	- P50: %d
	- P90: %d
	- P99: %d
	- Max: %d
	- Most common out-degree: %d
Callee in-degree metrics:
	- P50: %d
	- P90: %d
	- P99: %d
	- Max: %d
	- Most common in-degree: %d
`,
		m.NumberOfFunctions(),
		m.OutDegreeP50,
		m.OutDegreeP90,
		m.OutDegreeP99,
		m.OutDegreeMax,
		m.OutDegreeMode,
		m.InDegreeP50,
		m.InDegreeP90,
		m.InDegreeP99,
		m.InDegreeMax,
		m.InDegreeMode,
	)
}

func GetCallGraphMetrics(cg *callgraph.Graph) CallGraphMetrics {
	m := CallGraphMetrics{
		BaseMetrics: BaseMetrics[*callgraph.Graph]{
			Payload: cg,
		},
	}

	m = m.CallGraphInDegreeMetrics()
	m = m.CallGraphOutDegreeMetrics()
	return m
}

// CallGraphOutDegreeMetrics computes out-degree metrics on the call graph. The out degree
// is computed per-call site. Every function without outgoing calls contributes with a 0
// to the statistics.
func (m CallGraphMetrics) CallGraphOutDegreeMetrics() CallGraphMetrics {
	if m.Payload == nil {
		return m
	}

	res := m.Payload

	// Out-degree mode
	outDegrees := make([]int, 0, len(res.Nodes))
	// Maximum out-degree
	visitCallgraph(res, func(n *callgraph.Node) {
		outs := make(map[ssa.CallInstruction]int)

		if len(n.Out) == 0 {
			outDegrees = append(outDegrees, 0)
		}

		for _, e := range n.Out {
			outs[e.Site]++
		}

		for _, count := range outs {
			if m.OutDegreeMax < count {
				m.OutDegreeMax = count
			}
			outDegrees = append(outDegrees, count)
		}
	})

	slices.Sort(outDegrees)
	m.OutDegreeP50 = p50(outDegrees)
	m.OutDegreeP90 = p90(outDegrees)
	m.OutDegreeP99 = p99(outDegrees)
	m.OutDegreeMode = mode(outDegrees)

	return m
}

// CallGraphOutdegreeMetrics computes in-degree metrics on the call graph.
func (m CallGraphMetrics) CallGraphInDegreeMetrics() CallGraphMetrics {
	if m.Payload == nil {
		return m
	}

	res := m.Payload

	// In-degree mode
	cardinality := make(map[int]int)
	inDegrees := make([]int, 0, len(m.Payload.Nodes))
	visitCallgraph(res, func(n *callgraph.Node) {
		cardinality[len(n.In)]++
		inDegrees = append(inDegrees, len(n.In))
		if m.InDegreeMax < len(n.In) {
			m.InDegreeMax = len(n.In)
		}
	})

	slices.Sort(inDegrees)
	m.InDegreeP50 = p50(inDegrees)
	m.InDegreeP90 = p90(inDegrees)
	m.InDegreeP99 = p99(inDegrees)
	m.InDegreeMode = mode(inDegrees)

	return m
}

// NumberOfFunctions produces the number of functions in the call-graph produced
// by the Points-To analysis.
func (m CallGraphMetrics) NumberOfFunctions() int {
	if m.Payload == nil {
		return 0
	}

	return len(m.Payload.Nodes)
}
