package core

import (
	"testing"
)

func TestNewGRBSolver(t *testing.T) {
	solver := NewGurobiSolver()
	if solver == nil {
		t.Errorf("Error creating Gurobi Solver!")
		return
	}
}
