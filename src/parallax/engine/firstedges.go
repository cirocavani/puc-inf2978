package engine

import (
	"fmt"
	"parallax/core"
	"parallax/fct"
	"parallax/graph"
)

type FirstEdges struct {
	*graphEngine
	factor float64
}

func NewFirstEdges(g fct.GraphLoader, factor float64) core.BidEngine {
	return &FirstEdges{newGraphEngine(g), factor}
}

func (n *FirstEdges) ComputeBid(m *core.Match) *core.BidPack {
	n.setup(m.InstanceName)
	if n.current == nil {
		fmt.Println("Instance not found:", m.InstanceName)
		return core.EmptyBidPack()
	}
	pack := core.NewBidPack(m.NumberOfEdges)
	for i := 0; i < m.NumberOfEdges; i++ {
		e := n.current.Edges[i]
		source, sink, price := n.bid(e)
		pack.Bid(source, sink, price)
	}
	return pack
}

func (n *FirstEdges) bid(e *graph.Edge) (int, int, float64) {
	source := e.I.Data.(*fct.VertexData).Id
	sink := e.J.Data.(*fct.VertexData).Id
	price := e.Data.(*fct.EdgeData).VCost
	factor := n.factor
	return source, sink, factor * price
}
