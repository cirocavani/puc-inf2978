package engine

import (
	"parallax/core"
	"parallax/fct"
)

const (
	BID_RANDOM_EDGES string = "RandomEdges"
	BID_FIRST_EDGES         = "FirstEdges"
	BID_GUROBI_EDGES string = "GurobiEdges"
)

func New(name string, graphs fct.GraphLoader, factor float64) core.BidEngine {
	switch name {
	case BID_RANDOM_EDGES:
		return NewRandomEdges(graphs, factor)
	case BID_FIRST_EDGES:
		return NewFirstEdges(graphs, factor)
	case BID_GUROBI_EDGES:
		return NewGurobiEdges(graphs, factor)
	default:
		return nil
	}
}
