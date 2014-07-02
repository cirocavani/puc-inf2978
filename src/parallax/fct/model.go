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

func (g *Graph) Edge(source, sink int) *graph.Edge {
	key := fmt.Sprint(source, ":", sink)
	if e, found := g.EdgeMap[key]; found {
		return e
	}
	return nil
}
