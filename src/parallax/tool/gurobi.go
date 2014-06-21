package tool

/*
#include "gurobi_c.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type GRBEnv struct {
	log *C.char
	env *C.struct__GRBenv
}

func NewGRBEnv(log string) (*GRBEnv, error) {
	clog := C.CString(log)
	var env *C.struct__GRBenv = nil
	result := C.GRBloadenv(&env, clog)
	if int(result) != 0 {
		C.free(unsafe.Pointer(clog))
		return nil, fmt.Errorf("%d", int(result))
	}
	return &GRBEnv{clog, env}, nil
}

func (env *GRBEnv) Dispose() {
	C.GRBfreeenv(env.env)
	C.free(unsafe.Pointer(env.log))
	env.env = nil
	env.log = nil
}

func (env *GRBEnv) ErrorMessage() string {
	if m := C.GRBgeterrormsg(env.env); m != nil {
		return C.GoString(m)
	} else {
		return "empty"
	}
}

type GRBModel struct {
	name  *C.char
	model *C.struct__GRBmodel

	env     *GRBEnv
	vars    map[string]*GRBVar
	constrs map[string]*GRBConstr
}

type GRBVar struct {
	name *C.char

	model        *GRBModel
	varType      int
	obj          float64
	lower, upper float64
}

type GRBConstr struct {
	name *C.char

	model *GRBModel
}

func NewGRBModel(env *GRBEnv, name string) (*GRBModel, error) {
	cname := C.CString(name)
	var model *C.struct__GRBmodel = nil
	result := C.GRBnewmodel(env.env, &model, cname, C.int(0), nil, nil, nil, nil, nil)
	if int(result) != 0 {
		C.free(unsafe.Pointer(cname))
		return nil, fmt.Errorf("%d: %s", int(result), env.ErrorMessage())
	}
	return &GRBModel{
		cname,
		model,
		env,
		make(map[string]*GRBVar),
		make(map[string]*GRBConstr),
	}, nil
}

func (model *GRBModel) Dispose() {
	model.env = nil

	C.GRBfreemodel(model.model)
	C.free(unsafe.Pointer(model.name))
	for _, v := range model.vars {
		v.dispose()
	}
	for _, c := range model.constrs {
		c.dispose()
	}
	model.vars = nil
	model.constrs = nil
	model.model = nil
	model.name = nil
}

func (v *GRBVar) dispose() {
	C.free(unsafe.Pointer(v.name))
	v.name = nil
	v.model = nil
}

func (c *GRBConstr) dispose() {
	C.free(unsafe.Pointer(c.name))
	c.name = nil
	c.model = nil
}

const (
	GRB_INFINITY  float64 = 1e100
	GRB_UNDEFINED         = 1e101
	GRB_MAXINT            = 2000000000
)

func (v *GRBVar) Index() (int, error) {
	index := C.int(-1)
	result := C.GRBgetvarbyname(v.model.model, v.name, &index)
	if int(result) == 0 {
		return int(index), nil
	} else {
		return -1, fmt.Errorf("%d: %s", int(result), v.model.env.ErrorMessage())
	}
}

func (v *GRBVar) Value() (float64, error) {
	i, err := v.Index()
	if err != nil {
		return 0., err
	}
	ATTR := C.CString("X")
	defer C.free(unsafe.Pointer(ATTR))
	index := C.int(i)
	value := C.double(0.)
	result := C.GRBgetdblattrelement(v.model.model, ATTR, index, &value)
	if int(result) == 0 {
		return float64(value), nil
	} else {
		return 0., fmt.Errorf("%d: %s", int(result), v.model.env.ErrorMessage())
	}
}

func (model *GRBModel) addVar(name string, t int, obj, lower, upper float64) *GRBVar {
	cname := C.CString(name)
	ctype := C.char(t)
	cobj := C.double(obj)
	clower := C.double(lower)
	cupper := C.double(upper)
	result := C.GRBaddvar(model.model, 0, nil, nil, cobj, clower, cupper, ctype, cname)
	if int(result) != 0 {
		C.free(unsafe.Pointer(cname))
		return nil
	}
	v := &GRBVar{cname, model, t, obj, lower, upper}
	model.vars[name] = v
	return v
}

func (model *GRBModel) AddContVar(name string, obj, lower, upper float64) *GRBVar {
	return model.addVar(name, 'C', obj, lower, upper)
}

func (model *GRBModel) AddSemiContVar(name string, obj, lower, upper float64) *GRBVar {
	return model.addVar(name, 'S', obj, lower, upper)
}

func (model *GRBModel) AddIntVar(name string, lower, upper, obj float64) *GRBVar {
	return model.addVar(name, 'I', obj, lower, upper)
}

func (model *GRBModel) AddSemiIntVar(name string, obj, lower, upper float64) *GRBVar {
	return model.addVar(name, 'N', obj, lower, upper)
}

func (model *GRBModel) AddBinaryVar(name string, obj, lower, upper float64) *GRBVar {
	return model.addVar(name, 'B', obj, lower, upper)
}

type ConstrExpr map[*GRBVar]float64

type ConstrOp int

const (
	GRB_LESS_EQUAL    ConstrOp = '<'
	GRB_GREATER_EQUAL          = '>'
	GRB_EQUAL                  = '='
)

func (model *GRBModel) AddConstr(name string, expr ConstrExpr, op ConstrOp, value float64) *GRBConstr {
	cnumv := C.int(len(expr))
	cindex := make([]C.int, len(expr))
	ccoef := make([]C.double, len(expr))
	k := 0
	for v, coef := range expr {
		i, err := v.Index()
		if err != nil {
			return nil
		}
		cindex[k] = C.int(i)
		ccoef[k] = C.double(coef)
		k++
	}
	cop := C.char(op)
	cvalue := C.double(value)
	cname := C.CString(name)
	result := C.GRBaddconstr(model.model, cnumv, &cindex[0], &ccoef[0], cop, cvalue, cname)
	if int(result) != 0 {
		C.free(unsafe.Pointer(cname))
		return nil
	}
	c := &GRBConstr{cname, model}
	model.constrs[name] = c
	return c
}

func (model *GRBModel) GetIntAttr(attr string) (int, error) {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.int(-1)
	result := C.GRBgetintattr(model.model, ATTR, &VALUE)
	if int(result) == 0 {
		return int(VALUE), nil
	} else {
		return -1, fmt.Errorf("%d: %s", int(result), model.env.ErrorMessage())
	}
}

func (model *GRBModel) SetIntAttr(attr string, value int) error {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.int(value)
	result := C.GRBsetintattr(model.model, ATTR, VALUE)
	if int(result) == 0 {
		return nil
	} else {
		return fmt.Errorf("%d: %s", int(result), model.env.ErrorMessage())
	}
}

func (model *GRBModel) GetDoubleAttr(attr string) (float64, error) {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.double(0.)
	result := C.GRBgetdblattr(model.model, ATTR, &VALUE)
	if int(result) == 0 {
		return float64(VALUE), nil
	} else {
		return 0., fmt.Errorf("%d: %s", int(result), model.env.ErrorMessage())
	}
}

func (model *GRBModel) SetDoubleAttr(attr string, value float64) error {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.double(value)
	result := C.GRBsetdblattr(model.model, ATTR, VALUE)
	if int(result) == 0 {
		return nil
	} else {
		return fmt.Errorf("%d: %s", int(result), model.env.ErrorMessage())
	}
}

func (model *GRBModel) SetMaximize() {
	model.SetIntAttr("ModelSense", -1 /* Maximize */)
}

func (model *GRBModel) SetMinimize() {
	model.SetIntAttr("ModelSense", 1 /* Minimize */)
}

func (model *GRBModel) Update() error {
	result := C.GRBupdatemodel(model.model)
	if int(result) == 0 {
		return nil
	} else {
		return fmt.Errorf("%d: %s", int(result), model.env.ErrorMessage())
	}
}

func (model *GRBModel) Optimize() error {
	result := C.GRBoptimize(model.model)
	if int(result) == 0 {
		return nil
	} else {
		return fmt.Errorf("%d: %s", int(result), model.env.ErrorMessage())
	}
}

func (model *GRBModel) Optimal() (bool, error) {
	v, err := model.GetIntAttr("Status")
	if err != nil {
		return false, err
	}
	return v == 2, nil // GRB_OPTIMAL 2
}

func (model *GRBModel) ObjectiveValue() (float64, error) {
	v, err := model.GetDoubleAttr("ObjVal")
	if err != nil {
		return 0., err
	}
	return v, nil
}
