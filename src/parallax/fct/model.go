package fct

import (
	"fmt"
	"parallax/graph"
)

// FCT data model

type VertexData struct {
	Id   int
	Size float64
}

func (v *VertexData) String() string {
	return fmt.Sprintf("[%d, %.2f]", v.Id, v.Size)
}

type EdgeData struct {
	VCost, FCost float64
}

func (e *EdgeData) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", e.VCost, e.FCost)
}

func EdgeKey(source, sink int) string {
	return fmt.Sprint(source, ":", sink)
}

type Graph struct {
	*graph.Graph
	Sources, Sinks map[int]*graph.Vertex
	EdgeMap        map[string]*graph.Edge
}

func (g *Graph) String() string {
	return fmt.Sprintf("Sources %d, Sinks %d, Edges %d", g.SourceOrder(), g.SinkOrder(), g.Graph.Size())
}

func NewGraph() *Graph {
	return &Graph{
		graph.New(),
		make(map[int]*graph.Vertex),
		make(map[int]*graph.Vertex),
		make(map[string]*graph.Edge),
	}
}

func (g *Graph) Clone() *Graph {
	result := NewGraph()

	for _, v := range g.Sources {
		_v := v.Data.(*VertexData)
		result.SourceSize(_v.Id, _v.Size)
	}

	for _, v := range g.Sinks {
		_v := v.Data.(*VertexData)
		result.SinkSize(_v.Id, _v.Size)
	}

	for _, e := range g.Edges {
		source := e.I.Data.(*VertexData)
		sink := e.J.Data.(*VertexData)
		_e := e.Data.(*EdgeData)
		vcost := _e.VCost
		fcost := _e.FCost
		result.NewEdge(source.Id, sink.Id, vcost, fcost)
	}

	return result
}

func (g *Graph) SourceOrder() int {
	return len(g.Sources)
}

func (g *Graph) SinkOrder() int {
	return len(g.Sinks)
}

func (g *Graph) v(m map[int]*graph.Vertex, id int) *graph.Vertex {
	if v, found := m[id]; found {
		return v
	}
	v := g.Vertex()
	v.Data = &VertexData{id, .0}
	m[id] = v
	return v
}

func (g *Graph) NewEdge(source, sink int, v, f float64) *graph.Edge {
	vsource := g.v(g.Sources, source)
	vsink := g.v(g.Sinks, sink)
	e := g.Connect(vsource).To(vsink)
	e.Data = &EdgeData{v, f}
	key := fmt.Sprint(source, ":", sink)
	g.EdgeMap[key] = e
	return e
}

func (g *Graph) Edge(source, sink int) (*graph.Edge, string) {
	key := EdgeKey(source, sink)
	if e, found := g.EdgeMap[key]; found {
		return e, key
	}
	return nil, key
}

func (g *Graph) EdgeCost(source, sink int, v float64) (*graph.Edge, string) {
	key := EdgeKey(source, sink)
	e, found := g.EdgeMap[key]
	if !found {
		return nil, key
	}
	_e := e.Data.(*EdgeData)
	_e.VCost = v
	return e, key
}

func (g *Graph) SourceSize(id int, s float64) *graph.Vertex {
	v := g.v(g.Sources, id)
	v.Data.(*VertexData).Size = s
	return v
}

func (g *Graph) SinkSize(id int, s float64) *graph.Vertex {
	v := g.v(g.Sinks, id)
	v.Data.(*VertexData).Size = s
	return v
}
