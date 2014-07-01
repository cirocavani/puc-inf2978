package engine

import (
	"fmt"
	"parallax/core"
	"parallax/fct"
	"parallax/graph"
)

type FirstEdges struct {
	graphs fct.GraphLoader
	factor int
}

func NewFirstEdges(g fct.GraphLoader, factor int) core.BidEngine {
	return &FirstEdges{g, factor}
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
		source, sink, price := n.bid(e)
		pack.Bid(source, sink, price)
	}
	return pack
}

func (n *FirstEdges) bid(e *graph.Edge) (int, int, float64) {
	source := e.I.Data.(fct.VertexData).Id
	sink := e.J.Data.(fct.VertexData).Id
	price := e.Data.(fct.EdgeData).VCost
	factor := float64(n.factor)
	return source, sink, factor * price
}

func (n *FirstEdges) Update(f *core.Flow) {
}
