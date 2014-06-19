package engine

import (
	"fmt"
	"parallax/core"
	"parallax/graph"
)

type FirstEdges struct {
	graphs *graph.GraphLoader
}

func NewFirstEdges(g *graph.GraphLoader) core.Engine {
	return &FirstEdges{g}
}

func (n *FirstEdges) ComputeBid(m *core.Match) *core.BidPack {
	g := n.graphs.Instance(m.InstanceName)
	if g == nil {
		fmt.Println("Instance not found:", m.InstanceName)
		return core.EmptyBidPack()
	}
	pack := core.NewBidPack(m.NumberOfEdges)
	for i := 0; i < m.NumberOfEdges; i++ {
		e := g.Edges[i]
		pack.Bid(e.I.Id, e.J.Id, e.VCost)
	}
	return pack
}

func (n *FirstEdges) Update(f *core.Flow) {
}
