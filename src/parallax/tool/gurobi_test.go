package tool

import (
	"testing"
)

func TestGRBEnv(t *testing.T) {
	env, err := NewGRBEnv("test.log")
	if err != nil {
		t.Errorf("Error creating Gurobi Env: ", err)
		return
	}
	env.Dispose()
}

func TestGRBModel(t *testing.T) {
	env, _ := NewGRBEnv("test.log")
	defer env.Dispose()
	model, err := NewGRBModel(env, "Model Test")
	if err != nil {
		t.Errorf("Error creating Gurobi Model: ", err)
		return
	}
	model.Dispose()
}

func TestGRBVar(t *testing.T) {
	env, _ := NewGRBEnv("test.log")
	defer env.Dispose()
	model, _ := NewGRBModel(env, "Model Test")
	defer model.Dispose()
	v_cont := model.AddContVar("v_cont", 1., 0., 10.)
	if v_cont == nil {
		t.Errorf("Error creating Gurobi Variable: Continuous!")
	}
	v_semicont := model.AddSemiContVar("v_semicont", 1., 0., 10.)
	if v_semicont == nil {
		t.Errorf("Error creating Gurobi Variable: Semi Continuous!")
	}
	v_int := model.AddIntVar("v_int", 1., 0., 10.)
	if v_int == nil {
		t.Errorf("Error creating Gurobi Variable: Integer!")
	}
	v_semiint := model.AddSemiIntVar("v_semiint", 1., 0., 10.)
	if v_semiint == nil {
		t.Errorf("Error creating Gurobi Variable: Semi Integer!")
	}
	v_binary := model.AddBinaryVar("v_binary", 1., 0., 10.)
	if v_binary == nil {
		t.Errorf("Error creating Gurobi Variable: Binary!")
	}
}

func TestGRBConstr(t *testing.T) {
	env, _ := NewGRBEnv("test.log")
	defer env.Dispose()
	model, _ := NewGRBModel(env, "Model Test")
	defer model.Dispose()

	/* maximize: x + y + 2 z */
	x := model.AddBinaryVar("x", 1., 0., 1.)
	y := model.AddBinaryVar("y", 1., 0., 1.)
	z := model.AddBinaryVar("z", 2., 0., 1.)

	model.SetMaximize()
	model.Update()

	/* First constraint: x + 2 y + 3 z <= 4 */
	c1 := model.AddConstr("1", ConstrExpr{x: 1., y: 2., z: 3.}, GRB_LESS_EQUAL, 4.)
	if c1 == nil {
		t.Errorf("Error creating Gurobi Constraint: x + 2 y + 3 z <= 4!")
	}
	/* Second constraint: x + y >= 1 */
	c2 := model.AddConstr("2", ConstrExpr{x: 1., y: 1.}, GRB_GREATER_EQUAL, 1.)
	if c2 == nil {
		t.Errorf("Error creating Gurobi Constraint: x + y >= 1!")
	}

	//model.Optimize()
}

func TestGRBModelUpdate(t *testing.T) {
	env, _ := NewGRBEnv("test.log")
	defer env.Dispose()
	model, _ := NewGRBModel(env, "Model Test")
	defer model.Dispose()
	model.AddContVar("v_cont", 1., 0., 10.)
	err := model.Update()
	if err != nil {
		t.Errorf("Error updating Gurobi Model: ", err)
	}
}

func TestGRBModelOptimize(t *testing.T) {
	env, _ := NewGRBEnv("test.log")
	defer env.Dispose()
	model, _ := NewGRBModel(env, "Model Test")
	defer model.Dispose()
	model.AddContVar("v_cont", 1., 0., 10.)
	model.Update()
	err := model.Optimize()
	if err != nil {
		t.Errorf("Error optimizing Gurobi Model: ", err)
	}
}

func TestQuickStart(t *testing.T) {
	env, _ := NewGRBEnv("test.log")
	defer env.Dispose()
	model, _ := NewGRBModel(env, "Model Test")
	defer model.Dispose()

	/* maximize: x + y + 2 z */
	x := model.AddBinaryVar("x", 1., 0., 1.)
	y := model.AddBinaryVar("y", 1., 0., 1.)
	z := model.AddBinaryVar("z", 2., 0., 1.)

	model.SetMaximize()
	model.Update()

	/* First constraint: x + 2 y + 3 z <= 4 */
	model.AddConstr("1", ConstrExpr{x: 1., y: 2., z: 3.}, GRB_LESS_EQUAL, 4.)
	/* Second constraint: x + y >= 1 */
	model.AddConstr("2", ConstrExpr{x: 1., y: 1.}, GRB_GREATER_EQUAL, 1.)

	model.Optimize()

	opt, err := model.Optimal()
	if err != nil {
		t.Errorf("Error reading Optimal Status: ", err)
		return
	}
	if !opt {
		t.Errorf("Model is not optimal!")
		return
	}
	obj, err := model.ObjectiveValue()
	if err != nil {
		t.Errorf("Error reading Optimal Objective: ", err)
		return
	}
	t.Logf("Optimal Objective: %f\n", obj)
	vx, err := x.Value()
	if err != nil {
		t.Errorf("Error reading X: ", err)
		return
	}
	t.Logf("x: %f\n", vx)
	vy, err := y.Value()
	if err != nil {
		t.Errorf("Error reading Y: ", err)
		return
	}
	t.Logf("y: %f\n", vy)
	vz, err := z.Value()
	if err != nil {
		t.Errorf("Error reading Z: ", err)
		return
	}
	t.Logf("z: %f\n", vz)
}
