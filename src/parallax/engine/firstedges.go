package engine

import (
	"fmt"
	"parallax/core"
	"parallax/fct"
)

type FirstEdges struct {
	graphs *fct.GraphLoader
}

func NewFirstEdges(g *fct.GraphLoader) core.Engine {
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
		e := g.Graph.Edges[i]
		source := e.I.Data.(fct.VertexData).Id
		sink := e.J.Data.(fct.VertexData).Id
		price := e.Data.(fct.EdgeData).VCost
		pack.Bid(source, sink, price)
	}
	return pack
}

func (n *FirstEdges) Update(f *core.Flow) {
}
