package graph

import (
	"fmt"
)

// FCTP data model

type Vertex struct {
	Id    int
	Size  float64
	Edges []*Edge
}

func (v *Vertex) String() string {
	return fmt.Sprintf("(%d, [%.2f,%d])", v.Id, v.Size, len(v.Edges))
}

type Edge struct {
	I, J         *Vertex
	VCost, FCost float64
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d)-[%.2f,%.2f]->(%d)", e.I.Id, e.VCost, e.FCost, e.J.Id)
}

type Graph struct {
	Sources, Sinks map[int]*Vertex
	Edges          []*Edge
}

func (g *Graph) String() string {
	return fmt.Sprintf("Sources %d, Sinks %d, Edges %d", len(g.Sources), len(g.Sinks), len(g.Edges))
}

func NewGraph() *Graph {
	return &Graph{
		Sources: make(map[int]*Vertex),
		Sinks:   make(map[int]*Vertex),
		Edges:   make([]*Edge, 0),
	}
}

func (g *Graph) v(m map[int]*Vertex, id int) *Vertex {
	v, found := m[id]
	if !found {
		v = &Vertex{id, .0, make([]*Edge, 0)}
		m[id] = v
	}
	return v
}

func (g *Graph) NewEdge(i, j int, v, f float64) *Edge {
	vi := g.v(g.Sources, i)
	vj := g.v(g.Sinks, j)
	e := &Edge{vi, vj, v, f}
	vi.Edges = append(vi.Edges, e)
	vj.Edges = append(vj.Edges, e)
	g.Edges = append(g.Edges, e)
	return e
}

func (g *Graph) SourceSize(i int, s float64) *Vertex {
	vi := g.v(g.Sources, i)
	vi.Size = s
	return vi
}

func (g *Graph) SinkSize(j int, s float64) *Vertex {
	vj := g.v(g.Sinks, j)
	vj.Size = s
	return vj
}
