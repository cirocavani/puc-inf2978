package core

import (
	"errors"
	"fmt"
	"parallax/fct"
	"parallax/graph"
	"parallax/gurobi"
)

func NewGurobiFlowEngine(graph *fct.Graph) *FlowEngine {
	return &FlowEngine{graph, NewGurobiSolver()}
}

func NewGurobiSolver() *GurobiSolver {
	return &GurobiSolver{}
}

type GurobiSolver struct {
}

func (*GurobiSolver) ComputeFlow(g *fct.Graph) ([]*EdgeFlow, error) {
	env, err := grb.NewEnv("gurobi_solver.log")
	if err != nil {
		return nil, err
	}
	defer env.Dispose()
	model, err := grb.NewModel(env, "TransportModel")
	if err != nil {
		return nil, err
	}
	defer model.Dispose()

	// minimize n(i,j) * v(i,j)
	// 0 <= n(i,j) <= min{si,sj}
	// each i sum(i) n(i,j) = si
	// each j sum(j) n(i,j) = sj

	edges := make(map[*graph.Edge]*grb.Var)
	for _, e := range g.Edges {
		name, obj, upper := edge(e)
		edges[e] = model.AddContVar(name, obj, 0., upper)
	}

	model.SetMinimize()
	model.Update()

	expr := func(_edges []*graph.Edge) grb.ConstrExpr {
		expr := make(grb.ConstrExpr)
		for _, e := range _edges {
			evar := edges[e]
			expr[evar] = 1.
		}
		return expr
	}

	for _, v := range g.Sources {
		name, size := vertex(v)
		expr := expr(v.EdgeOut)
		model.AddConstr(name, expr, grb.EQUAL, size)
	}

	for _, v := range g.Sinks {
		name, size := vertex(v)
		expr := expr(v.EdgeIn)
		model.AddConstr(name, expr, grb.EQUAL, size)
	}

	model.Optimize()

	opt, err := model.Optimal()
	if err != nil {
		return nil, err
	}
	if !opt {
		return nil, errors.New("Model is not optimal!")
	}
	obj, err := model.ObjectiveValue()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Optimal Objective: %f\n", obj)

	result := make([]*EdgeFlow, 0, len(edges))
	for e, v := range edges {
		m, err := v.Value()
		if err != nil {
			return nil, err
		}
		if m < 0.01 {
			continue
		}
		result = append(result, flow(e, m))
	}
	return result, nil
}

func edge(e *graph.Edge) (string, float64, float64) {
	source := e.I.Data.(*fct.VertexData)
	sink := e.J.Data.(*fct.VertexData)
	name := fmt.Sprint(source.Id, ":", sink.Id)

	obj := e.Data.(*fct.EdgeData).VCost

	upper := source.Size
	if upper > sink.Size {
		upper = sink.Size
	}

	return name, obj, upper
}

func vertex(v *graph.Vertex) (string, float64) {
	_v := v.Data.(*fct.VertexData)
	return fmt.Sprint("Vertex ", _v.Id), _v.Size
}

func flow(e *graph.Edge, amount float64) *EdgeFlow {
	source := e.I.Data.(*fct.VertexData)
	sink := e.J.Data.(*fct.VertexData)
	return &EdgeFlow{
		source.Id,
		sink.Id,
		amount,
	}
}
