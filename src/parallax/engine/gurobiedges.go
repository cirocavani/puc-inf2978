package engine

import (
	"fmt"
	"parallax/core"
	"parallax/fct"
	"sort"
)

type GurobiEdges struct {
	graphs fct.GraphLoader
	factor int
	solver core.Solver
}

func NewGurobiEdges(g fct.GraphLoader, factor int) core.BidEngine {
	s := core.NewGurobiSolver()
	return &GurobiEdges{g, factor, s}
}

func (n *GurobiEdges) ComputeBid(m *core.Match) *core.BidPack {
	g := n.graphs.Instance(m.InstanceName)
	if g == nil {
		fmt.Println("Instance not found:", m.InstanceName)
		return core.EmptyBidPack()
	}
	r, err := n.solver.ComputeFlow(g)
	if err != nil {
		fmt.Println("Error computing flow:", err)
		return core.EmptyBidPack()
	}
	sort.Sort(flowSort(r))

	pack := core.NewBidPack(m.NumberOfEdges)
	factor := float64(n.factor)
	for i, k := len(r)-1, 0; k < m.NumberOfEdges && i > -1; i, k = i-1, k+1 {
		ef := r[i]
		e, _ := g.Edge(ef.Source, ef.Sink)
		price := e.Data.(*fct.EdgeData).VCost
		pack.Bid(ef.Source, ef.Sink, price*factor)
	}
	return pack
}

type flowSort []*core.EdgeFlow

func (v flowSort) Len() int           { return len(v) }
func (v flowSort) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v flowSort) Less(i, j int) bool { return v[i].Amount < v[j].Amount }

func (n *GurobiEdges) Update(f *core.Flow) {
}
