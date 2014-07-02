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
	flow, err := n.solver.ComputeFlow(n.current)
	if err != nil {
		fmt.Println("Error computing flow:", err)
		return core.EmptyBidPack()
	}
	//sort.Sort(core.FlowSort(flow))
	sort.Sort(NewProfitSort(n.current, flow))
	pack := core.NewBidPack(m.NumberOfEdges)
	factor := n.factor
	for i, k := len(flow)-1, 0; k < m.NumberOfEdges && i > -1; i, k = i-1, k+1 {
		ef := flow[i]
		e, _ := n.current.Edge(ef.Source, ef.Sink)
		price := e.Data.(*fct.EdgeData).VCost
		pack.Bid(ef.Source, ef.Sink, price*factor)
	}
	return pack
}

type ProfitSort struct {
	profit []float64
	flow   []*core.EdgeFlow
}

func NewProfitSort(g *fct.Graph, flow []*core.EdgeFlow) *ProfitSort {
	profit := make([]float64, len(flow))
	for i, ef := range flow {
		e, _ := g.Edge(ef.Source, ef.Sink)
		v := e.Data.(*fct.EdgeData).VCost
		profit[i] = ef.Amount * v
	}
	return &ProfitSort{profit, flow}
}

func (v *ProfitSort) Len() int           { return len(v.flow) }
func (v *ProfitSort) Less(i, j int) bool { return v.profit[i] < v.profit[j] }

func (v *ProfitSort) Swap(i, j int) {
	v.profit[i], v.profit[j] = v.profit[j], v.profit[i]
	v.flow[i], v.flow[j] = v.flow[j], v.flow[i]
}
