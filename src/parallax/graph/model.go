package graph

import (
	"fmt"
)

type Vertex struct {
	Graph   *Graph
	Edges   []*Edge
	EdgeOut []*Edge
	EdgeIn  []*Edge
	Data    interface{}
}

func (v *Vertex) Degree() int {
	return len(v.Edges) + len(v.EdgeOut) + len(v.EdgeIn)
}

func (v *Vertex) OutDegree() int {
	return len(v.Edges) + len(v.EdgeOut)
}

func (v *Vertex) InDegree() int {
	return len(v.Edges) + len(v.EdgeIn)
}

type Edge struct {
	Graph *Graph
	I, J  *Vertex
	Data  interface{}
}

func (e *Edge) Other(v *Vertex) *Vertex {
	switch v {
	case e.I:
		return e.J
	case e.J:
		return e.I
	default:
		return nil
	}
}

type Graph struct {
	Vertices []*Vertex
	Edges    []*Edge
}

func New() *Graph {
	return &Graph{
		make([]*Vertex, 0),
		make([]*Edge, 0),
	}
}

func (g *Graph) String() string {
	return fmt.Sprintf("G(V, E) = [%d, %d]", g.Order(), g.Size())
}

func (g *Graph) Order() int {
	return len(g.Vertices)
}

func (g *Graph) Size() int {
	return len(g.Edges)
}

func (g *Graph) Vertex() *Vertex {
	v := &Vertex{
		g,
		make([]*Edge, 0),
		make([]*Edge, 0),
		make([]*Edge, 0),
		nil,
	}
	g.Vertices = append(g.Vertices, v)
	return v
}

func (g *Graph) edge(vi, vj *Vertex) *Edge {
	e := &Edge{g, vi, vj, nil}
	g.Edges = append(g.Edges, e)
	return e
}

type EdgeBuilder struct {
	g *Graph
	v *Vertex
}

func (b *EdgeBuilder) With(other *Vertex) *Edge {
	// undirected
	one := b.v
	e := b.g.edge(one, other)
	one.Edges = append(one.Edges, e)
	other.Edges = append(other.Edges, e)
	return e
}

func (b *EdgeBuilder) To(head *Vertex) *Edge {
	// directed
	tail := b.v
	e := b.g.edge(tail, head)
	head.EdgeIn = append(head.EdgeIn, e)
	tail.EdgeOut = append(tail.EdgeOut, e)
	return e
}

func (g *Graph) Connect(v *Vertex) *EdgeBuilder {
	return &EdgeBuilder{g, v}
}
