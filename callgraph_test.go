package stamets

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

// Construct a call graph to use for testing
//
// Call graph:
//
//	                .--->[5]
//	 .------>[4]    |
//	 |        |     c--->[6]
//	 |        V     |
//	[0]----->[1]-->[2]-->[3]
//	          Λ     |
//	          |_____|
//
// Topology:
// 0: Out(0.a) = 1, Out(0.b) = 1, In(0) = 0
// 1: Out(1) = 1, In(1) = 3
// 2: Out(2.a) = 1, Out(2.b) = 1, Out(2.c) = 2, In(2) = 1
// 3: In(3) = 1
// 4: Out(4) = 1, In(4) = 1
// 5: In(5) = 1
// 6: In(6) = 1
//
// Out-degress: [0, 0, 0, 1, 1, 1, 1, 1, 1, 2]
// In-degrees: [0, 1, 1, 1, 1, 1, 3]
func makeCallgraph(t *testing.T) (*callgraph.Graph, map[int]*callgraph.Node) {
	nodes := make(map[int]*callgraph.Node)
	makeNode := func(id int) *callgraph.Node {
		n := &callgraph.Node{ID: id, Func: &ssa.Function{}}
		_, ok := nodes[id]
		require.False(t, ok)
		nodes[id] = n
		return n
	}
	addEdge := func(from, to int, site ssa.CallInstruction) {
		n1, ok := nodes[from]
		require.True(t, ok)
		n2, ok := nodes[to]
		require.True(t, ok)

		callgraph.AddEdge(n1, site, n2)
	}
	root := makeNode(0)
	makeNode(1)
	makeNode(2)
	makeNode(3)
	makeNode(4)
	makeNode(5)
	makeNode(6)

	addEdge(0, 1, &ssa.Call{})
	addEdge(0, 4, &ssa.Call{})
	addEdge(4, 1, &ssa.Call{})
	addEdge(1, 2, &ssa.Call{})
	addEdge(2, 1, &ssa.Call{})
	addEdge(2, 3, &ssa.Call{})
	impreciseCallsite := &ssa.Call{}
	addEdge(2, 5, impreciseCallsite)
	addEdge(2, 6, impreciseCallsite)

	funs := make(map[*ssa.Function]*callgraph.Node)
	for _, n := range nodes {
		funs[n.Func] = n
	}
	return &callgraph.Graph{
		Root:  root,
		Nodes: funs,
	}, nodes
}

func TestVisitCallgraph(t *testing.T) {
	// No panic
	visitCallgraph(nil, func(n *callgraph.Node) {})

	// No panic
	visitCallgraph(&callgraph.Graph{}, nil)

	cg, nodes := makeCallgraph(t)

	visited := make(map[int]struct{})
	visitCallgraph(cg, func(n *callgraph.Node) {
		existingNode, ok := nodes[n.ID]
		require.True(t, ok)
		require.Equal(t, existingNode, n)

		require.NotContains(t, visited, n.ID)

		visited[n.ID] = struct{}{}
	})

	for _, n := range nodes {
		require.Contains(t, visited, n.ID)
	}

	for n := range visited {
		require.Contains(t, nodes, n)
	}
}

func TestCallGraphOutDegreeMetrics(t *testing.T) {
	cg, _ := makeCallgraph(t)

	m := CallGraphMetrics{
		BaseMetrics: BaseMetrics[*callgraph.Graph]{
			Payload: cg,
		},
	}

	require.Equal(t, 7, m.NumberOfFunctions())

	m = m.CallGraphOutDegreeMetrics()

	require.Equal(t, 2, m.OutDegreeMax)
	require.Equal(t, 1, m.OutDegreeP50)
	require.Equal(t, 2, m.OutDegreeP90)
	require.Equal(t, 2, m.OutDegreeP99)
}

func TestCallGraphInDegreeMetrics(t *testing.T) {
	cg, _ := makeCallgraph(t)

	m := CallGraphMetrics{
		BaseMetrics: BaseMetrics[*callgraph.Graph]{
			Payload: cg,
		},
	}

	require.Equal(t, 7, m.NumberOfFunctions())

	m = m.CallGraphInDegreeMetrics()

	require.Equal(t, 3, m.InDegreeMax)
	require.Equal(t, 1, m.InDegreeP50)
	require.Equal(t, 3, m.InDegreeP90)
	require.Equal(t, 3, m.InDegreeP99)
	require.Equal(t, 1, m.InDegreeMode)
}
