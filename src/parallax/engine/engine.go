package engine

import (
	"fmt"
	"parallax/core"
	"parallax/fct"
)

type graphEngine struct {
	graphs  fct.GraphLoader
	data    map[string]*fct.Graph
	current *fct.Graph
}

func newGraphEngine(g fct.GraphLoader) *graphEngine {
	return &graphEngine{
		g,
		make(map[string]*fct.Graph),
		nil,
	}
}

func (n *graphEngine) setup(name string) {
	if g, found := n.data[name]; found {
		n.current = g
		return
	}
	n.current = nil
	g := n.graphs.Instance(name)
	if g == nil {
		return
	}
	n.current = g.Clone()
	n.data[name] = n.current
}

func (n *graphEngine) Update(f *core.Flow) {
	if n.current == nil {
		fmt.Println("Instance not found")
		return
	}
	for _, s := range f.Streams {
		n.current.EdgeCost(s.Source, s.Sink, s.Price)
	}
}
