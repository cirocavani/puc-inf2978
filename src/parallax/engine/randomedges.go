package engine

import (
	"fmt"
	"math/rand"
	"parallax/core"
	"parallax/fct"
	"parallax/graph"
	"time"
)

type RandomEdges struct {
	*graphEngine
	maxFactor  int
	irnd, frnd *rand.Rand
}

func NewRandomEdges(g fct.GraphLoader, factor float64) core.BidEngine {
	t := time.Now().UnixNano()
	src := rand.NewSource(t)
	irnd := rand.New(src)
	frnd := rand.New(src)
	return &RandomEdges{newGraphEngine(g), int(factor), irnd, frnd}
}

func (n *RandomEdges) ComputeBid(m *core.Match) *core.BidPack {
	n.setup(m.InstanceName)
	if n.current == nil {
		fmt.Println("Instance not found:", m.InstanceName)
		return core.EmptyBidPack()
	}

	pack := core.NewBidPack(m.NumberOfEdges)
	emax := n.current.Size()
	eidx := make(map[int]bool)
	for k := 0; k < m.NumberOfEdges; {
		i := n.irnd.Intn(emax)
		if eidx[i] {
			continue
		}
		eidx[i] = true
		k++
		e := n.current.Edges[i]
		source, sink, price := n.bid(e)
		pack.Bid(source, sink, price)
	}
	return pack
}

func (n *RandomEdges) bid(e *graph.Edge) (int, int, float64) {
	source := e.I.Data.(*fct.VertexData).Id
	sink := e.J.Data.(*fct.VertexData).Id
	price := e.Data.(*fct.EdgeData).VCost
	factor := 1. + float64(n.frnd.Intn(n.maxFactor))
	return source, sink, factor * price
}
