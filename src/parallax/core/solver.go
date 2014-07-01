package core

import (
	"fmt"
	"parallax/fct"
)

type NamedBid map[string]*BidPack

type EdgeFlow struct {
	Source, Sink int
	Amount       float64
}

func (e *EdgeFlow) String() string {
	return fmt.Sprintf("(%d)-[%.2f]->(%d)", e.Source, e.Amount, e.Sink)
}

type Solver interface {
	ComputeFlow(g *fct.Graph) ([]*EdgeFlow, error)
}
