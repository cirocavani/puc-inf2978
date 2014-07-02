package engine

import (
	"fmt"
	"parallax/core"
	"parallax/fct"
	"sort"
)

type GurobiEdges struct {
	*graphEngine
	factor float64
	solver core.Solver
}

func NewGurobiEdges(g fct.GraphLoader, factor float64) core.BidEngine {
	solver := core.NewGurobiSolver()
	return &GurobiEdges{
		newGraphEngine(g),
		factor,
		solver,
	}
}

func (n *GurobiEdges) ComputeBid(m *core.Match) *core.BidPack {
	n.setup(m.InstanceName)
	if n.current == nil {
		fmt.Println("Instance not found:", m.InstanceName)
		return core.EmptyBidPack()
	}
	r, err := n.solver.ComputeFlow(n.current)
	if err != nil {
		fmt.Println("Error computing flow:", err)
		return core.EmptyBidPack()
	}
	sort.Sort(core.FlowSort(r))

	pack := core.NewBidPack(m.NumberOfEdges)
	factor := n.factor
	for i, k := len(r)-1, 0; k < m.NumberOfEdges && i > -1; i, k = i-1, k+1 {
		ef := r[i]
		e, _ := n.current.Edge(ef.Source, ef.Sink)
		price := e.Data.(*fct.EdgeData).VCost
		pack.Bid(ef.Source, ef.Sink, price*factor)
	}
	return pack
}
