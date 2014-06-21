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

const (
	GRB_INFINITY  float64 = 1e100
	GRB_UNDEFINED         = 1e101
	GRB_MAXINT            = 2000000000
)

type ConstrExpr map[*GRBVar]float64

type ConstrOp int

const (
	GRB_LESS_EQUAL    ConstrOp = '<'
	GRB_GREATER_EQUAL          = '>'
	GRB_EQUAL                  = '='
)

type GRBEnv struct {
	log *C.char
	env *C.struct__GRBenv
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
	expr  ConstrExpr
	op    ConstrOp
	value float64
}

func NewGRBEnv(log string) (*GRBEnv, error) {
	clog := C.CString(log)
	var env *C.struct__GRBenv = nil
	result := int(C.GRBloadenv(&env, clog))
	if result != 0 {
		C.free(unsafe.Pointer(clog))
		return nil, fmt.Errorf("%d", result)
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
	}
	return "empty"
}

func (env *GRBEnv) error(code int) error {
	return fmt.Errorf("%d: %s", code, env.ErrorMessage())
}

func NewGRBModel(env *GRBEnv, name string) (*GRBModel, error) {
	cname := C.CString(name)
	cnumv := C.int(0)
	var model *C.struct__GRBmodel = nil
	result := int(C.GRBnewmodel(env.env, &model, cname, cnumv, nil, nil, nil, nil, nil))
	if result != 0 {
		C.free(unsafe.Pointer(cname))
		return nil, env.error(result)
	}
	return &GRBModel{
		cname,
		model,
		env,
		make(map[string]*GRBVar),
		make(map[string]*GRBConstr),
	}, nil
}

func (m *GRBModel) Dispose() {
	m.env = nil

	C.GRBfreemodel(m.model)
	C.free(unsafe.Pointer(m.name))
	for _, v := range m.vars {
		v.dispose()
	}
	for _, c := range m.constrs {
		c.dispose()
	}
	m.vars = nil
	m.constrs = nil
	m.model = nil
	m.name = nil
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

func (v *GRBVar) Index() (int, error) {
	index := C.int(-1)
	result := int(C.GRBgetvarbyname(v.model.model, v.name, &index))
	if result != 0 {
		return -1, v.model.env.error(result)
	}
	return int(index), nil
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
	result := int(C.GRBgetdblattrelement(v.model.model, ATTR, index, &value))
	if result != 0 {
		return 0., v.model.env.error(result)
	}
	return float64(value), nil
}

func (m *GRBModel) addVar(name string, t int, obj, lower, upper float64) *GRBVar {
	cname := C.CString(name)
	ctype := C.char(t)
	cobj := C.double(obj)
	clower := C.double(lower)
	cupper := C.double(upper)
	result := int(C.GRBaddvar(m.model, 0, nil, nil, cobj, clower, cupper, ctype, cname))
	if result != 0 {
		C.free(unsafe.Pointer(cname))
		return nil
	}
	v := &GRBVar{cname, m, t, obj, lower, upper}
	m.vars[name] = v
	return v
}

func (m *GRBModel) AddContVar(name string, obj, lower, upper float64) *GRBVar {
	return m.addVar(name, 'C', obj, lower, upper)
}

func (m *GRBModel) AddSemiContVar(name string, obj, lower, upper float64) *GRBVar {
	return m.addVar(name, 'S', obj, lower, upper)
}

func (m *GRBModel) AddIntVar(name string, lower, upper, obj float64) *GRBVar {
	return m.addVar(name, 'I', obj, lower, upper)
}

func (m *GRBModel) AddSemiIntVar(name string, obj, lower, upper float64) *GRBVar {
	return m.addVar(name, 'N', obj, lower, upper)
}

func (m *GRBModel) AddBinaryVar(name string, obj, lower, upper float64) *GRBVar {
	return m.addVar(name, 'B', obj, lower, upper)
}

func (m *GRBModel) AddConstr(name string, expr ConstrExpr, op ConstrOp, value float64) *GRBConstr {
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
	result := int(C.GRBaddconstr(m.model, cnumv, &cindex[0], &ccoef[0], cop, cvalue, cname))
	if result != 0 {
		C.free(unsafe.Pointer(cname))
		return nil
	}
	c := &GRBConstr{cname, m, expr, op, value}
	m.constrs[name] = c
	return c
}

func (m *GRBModel) GetIntAttr(attr string) (int, error) {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.int(-1)
	result := int(C.GRBgetintattr(m.model, ATTR, &VALUE))
	if result != 0 {
		return -1, m.env.error(result)
	}
	return int(VALUE), nil
}

func (m *GRBModel) SetIntAttr(attr string, value int) error {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.int(value)
	result := int(C.GRBsetintattr(m.model, ATTR, VALUE))
	if result != 0 {
		return m.env.error(result)
	}
	return nil
}

func (m *GRBModel) GetDoubleAttr(attr string) (float64, error) {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.double(0.)
	result := int(C.GRBgetdblattr(m.model, ATTR, &VALUE))
	if result != 0 {
		return 0., m.env.error(result)
	}
	return float64(VALUE), nil
}

func (m *GRBModel) SetDoubleAttr(attr string, value float64) error {
	ATTR := C.CString(attr)
	defer C.free(unsafe.Pointer(ATTR))
	VALUE := C.double(value)
	result := int(C.GRBsetdblattr(m.model, ATTR, VALUE))
	if result != 0 {
		return m.env.error(result)
	}
	return nil
}

func (m *GRBModel) SetMaximize() {
	m.SetIntAttr("ModelSense", -1 /* Maximize */)
}

func (m *GRBModel) SetMinimize() {
	m.SetIntAttr("ModelSense", 1 /* Minimize */)
}

func (m *GRBModel) Update() error {
	result := int(C.GRBupdatemodel(m.model))
	if result != 0 {
		return m.env.error(result)
	}
	return nil
}

func (m *GRBModel) Optimize() error {
	result := int(C.GRBoptimize(m.model))
	if result != 0 {
		return m.env.error(result)
	}
	return nil
}

func (m *GRBModel) Optimal() (bool, error) {
	v, err := m.GetIntAttr("Status")
	if err != nil {
		return false, err
	}
	return v == 2, nil // GRB_OPTIMAL 2
}

func (m *GRBModel) ObjectiveValue() (float64, error) {
	return m.GetDoubleAttr("ObjVal")
}
