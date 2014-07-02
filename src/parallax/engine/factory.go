package engine

import (
	"parallax/core"
	"parallax/fct"
)

const (
	BID_RANDOM_EDGES string = "RandomEdges"
	BID_FIRST_EDGES         = "FirstEdges"
)

func New(name string, graphs fct.GraphLoader, factor int) core.BidEngine {
	switch name {
	case BID_RANDOM_EDGES:
		return NewRandomEdges(graphs, factor)
	case BID_FIRST_EDGES:
		return NewFirstEdges(graphs, factor)
	default:
		return nil
	}
}
