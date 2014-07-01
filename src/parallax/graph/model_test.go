package graph

import (
	"testing"
)

func TestNew(t *testing.T) {
	g := New()
	if g == nil {
		t.Fatal("Error creating Graph: nil")
	}
	if n := g.Order(); n != 0 {
		t.Error("Graph is not empty (order):", n)
	}
	if m := g.Size(); m != 0 {
		t.Error("Graph is not empty (size):", m)
	}
}

func TestNewVertex(t *testing.T) {
	g := New()
	v := g.Vertex()
	if v == nil {
		t.Fatal("Error creating Vertex: nil")
	}
	if d := v.Degree(); d != 0 {
		t.Error("Vertex is not empty (degree):", d)
	}
	if d := v.InDegree(); d != 0 {
		t.Error("Vertex is not empty (indegree):", d)
	}
	if d := v.OutDegree(); d != 0 {
		t.Error("Vertex is not empty (outdegree):", d)
	}
	if n := g.Order(); n != 1 {
		t.Error("Graph missing vertex (1):", n)
	}
}

func TestNewUndirectedEdge(t *testing.T) {
	g := New()
	v1 := g.Vertex()
	v2 := g.Vertex()
	e := g.Connect(v1).With(v2)
	if e == nil {
		t.Fatal("Error creating Edge: nil")
	}
	if o := e.Other(v1); o != v2 {
		t.Error("Edge is not bind to v2 from v1")
	}
	if o := e.Other(v2); o != v1 {
		t.Error("Edge is not bind to v1 from v2")
	}
	if d := v1.Degree(); d != 1 {
		t.Error("Vertex v1 missing degree:", d)
	}
	if d := v1.InDegree(); d != 1 {
		t.Error("Vertex v1 missing indegree:", d)
	}
	if d := v1.OutDegree(); d != 1 {
		t.Error("Vertex v1 missing outdegree:", d)
	}
	if d := v2.Degree(); d != 1 {
		t.Error("Vertex v2 missing degree:", d)
	}
	if d := v2.InDegree(); d != 1 {
		t.Error("Vertex v2 missing indegree:", d)
	}
	if d := v2.OutDegree(); d != 1 {
		t.Error("Vertex v2 missing outdegree:", d)
	}
	if n := g.Order(); n != 2 {
		t.Error("Graph missing vertex (2):", n)
	}
	if n := g.Size(); n != 1 {
		t.Error("Graph missing edge (1):", n)
	}
}

func TestNewDirectedEdge(t *testing.T) {
	g := New()
	v1 := g.Vertex()
	v2 := g.Vertex()
	e := g.Connect(v1).To(v2)
	if e == nil {
		t.Fatal("Error creating Edge: nil")
	}
	if o := e.Other(v1); o != v2 {
		t.Error("Edge is not bind to v2 from v1")
	}
	if o := e.Other(v2); o != v1 {
		t.Error("Edge is not bind to v1 from v2")
	}
	if d := v1.Degree(); d != 1 {
		t.Error("Vertex v1 missing degree:", d)
	}
	if d := v1.InDegree(); d != 0 {
		t.Error("Vertex v1 exceeding indegree:", d)
	}
	if d := v1.OutDegree(); d != 1 {
		t.Error("Vertex v1 missing outdegree:", d)
	}
	if d := v2.Degree(); d != 1 {
		t.Error("Vertex v2 missing degree:", d)
	}
	if d := v2.InDegree(); d != 1 {
		t.Error("Vertex v2 missing indegree:", d)
	}
	if d := v2.OutDegree(); d != 0 {
		t.Error("Vertex v2 exceeding outdegree:", d)
	}
	if n := g.Order(); n != 2 {
		t.Error("Graph missing vertex (2):", n)
	}
	if n := g.Size(); n != 1 {
		t.Error("Graph missing edge (1):", n)
	}
}
