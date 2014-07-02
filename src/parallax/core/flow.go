package core

import (
	"fmt"
	"parallax/fct"
)

type FlowEngine struct {
	graph  *fct.Graph
	solver Solver
}

func NewFlowEngine(graph *fct.Graph, solver Solver) *FlowEngine {
	return &FlowEngine{graph, solver}
}

func (n *FlowEngine) ComputeFlow(bids map[string]*BidPack) (*Flow, error) {
	_g, bidMap := BidGraph(n.graph, bids)
	flow, err := n.solver.ComputeFlow(_g)
	if err != nil {
		return nil, err
	}
	return BidFlow(flow, bidMap), nil
}

type Solver interface {
	ComputeFlow(g *fct.Graph) ([]*EdgeFlow, error)
}

type EdgeFlow struct {
	Source, Sink int
	Amount       float64
}

func (e *EdgeFlow) String() string {
	return fmt.Sprintf("(%d)-[%.2f]->(%d)", e.Source, e.Amount, e.Sink)
}

type EdgeBid struct {
	source, sink int
	owners       []string
	price        float64
	count        int
}

func BidGraph(g *fct.Graph, bids map[string]*BidPack) (*fct.Graph, map[string]*EdgeBid) {
	result := fct.NewGraph()

	for _, v := range g.Sources {
		_v := v.Data.(*fct.VertexData)
		result.SourceSize(_v.Id, _v.Size)
	}

	for _, v := range g.Sinks {
		_v := v.Data.(*fct.VertexData)
		result.SinkSize(_v.Id, _v.Size)
	}

	for _, e := range g.Edges {
		source := e.I.Data.(*fct.VertexData)
		sink := e.J.Data.(*fct.VertexData)
		_e := e.Data.(*fct.EdgeData)
		vcost := 20 * _e.VCost
		fcost := _e.FCost
		result.NewEdge(source.Id, sink.Id, vcost, fcost)
	}

	bidMap := make(map[string]*EdgeBid)
	for owner, pack := range bids {
		for _, bid := range pack.bids {
			e, key := result.Edge(bid.source, bid.sink)
			if e == nil {
				continue
			}
			ebid, found := bidMap[key]
			if !found {
				ebid = &EdgeBid{
					bid.source,
					bid.sink,
					make([]string, 0),
					0.,
					0,
				}
				bidMap[key] = ebid
			}
			_e := e.Data.(*fct.EdgeData)
			if bid.price < _e.VCost {
				ebid.owners = []string{owner}
				ebid.price = bid.price
				_e.VCost = bid.price
			} else if bid.price-_e.VCost < 0.001 {
				ebid.owners = append(ebid.owners, owner)
				ebid.price = bid.price
			} // bid.price > _e.VCost + 0.001
			ebid.count++
		}
	}

	return result, bidMap
}

func BidFlow(edges []*EdgeFlow, bids map[string]*EdgeBid) *Flow {
	result := make([]*Stream, 0)
	for _, e := range edges {
		key := fct.EdgeKey(e.Source, e.Sink)
		bid, found := bids[key]
		if !found {
			continue
		}
		n := float64(len(bid.owners))
		for _, owner := range bid.owners {
			s := &Stream{
				e.Source,
				e.Sink,
				e.Amount / n,
				owner,
				bid.price,
				bid.count,
			}
			result = append(result, s)
		}
	}
	return &Flow{result}
}
