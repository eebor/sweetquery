package generator

import (
	"go/ast"

	"github.com/eebor/sweetquery/querygen/model"
)

type UniType struct {
	typ          ast.Expr
	originalType ast.Expr
	baseTyp      ast.Expr
	tasks        []model.GenTask
}

func NewUniType(typ ast.Expr) *UniType {
	return &UniType{
		typ:          typ,
		originalType: typ,
		baseTyp:      typ,
		tasks:        make([]model.GenTask, 0),
	}
}

func (t *UniType) Reset() {
	t.typ = t.originalType
	t.baseTyp = t.originalType
	t.tasks = make([]model.GenTask, 0)
}

func (t *UniType) GetTasks() []model.GenTask {
	return t.tasks
}

func (t *UniType) GetOpertion(key string, value string) operationInterface {
	switch t.typ.(type) {
	case *ast.Ident:
		return t.identCase(key, value)
	case *ast.StarExpr:
		return t.pointerCase(key, value)
	case *ast.ArrayType:
		return t.arrayCase(key, value)
	case *ast.MapType:
		return t.mapCase(key, value)
	case *ast.StructType:
		return t.structCase(key, value)
	}

	return nil
}

func (t *UniType) identCase(key string, value string) operationInterface {
	id := t.typ.(*ast.Ident)
	if id.Obj != nil {
		t.typ = t.unwrapCustomType(id)
		return t.GetOpertion(key, value)
	}

	isCustomType := id != t.baseTyp

	return &queryWriteOperation{
		defaultOperation: &defaultOperation{
			Key:   key,
			Value: value,
		},
		Type:       id.Name,
		CustomType: isCustomType,
	}
}

func (t *UniType) pointerCase(key string, value string) operationInterface {
	t.typ = t.typ.(*ast.StarExpr).X

	t.baseTyp = t.typ

	op := t.GetOpertion(key, value)
	if op == nil {
		return nil
	}

	return &pointerCondOperation{
		operationInterface: op,
	}
}

func (t *UniType) arrayCase(key string, value string) operationInterface {
	t.typ = t.typ.(*ast.ArrayType).Elt

	t.baseTyp = t.typ

	op := t.GetOpertion(key, value)
	if op == nil {
		return nil
	}

	_, isArray := op.(*arrayOperation)
	if isArray {
		return nil
	}

	t.baseTyp = t.typ

	return &arrayOperation{
		operationInterface: op,
	}
}

func (t *UniType) mapCase(key string, value string) operationInterface {
	mt := t.typ.(*ast.MapType)

	keyid, keyIsIdent := mt.Key.(*ast.Ident)
	if !keyIsIdent {
		return nil
	}
	if keyid.Obj != nil {
		nt := t.unwrapCustomType(keyid)
		keyid, keyIsIdent = nt.(*ast.Ident)
		if !keyIsIdent {
			return nil
		}
	}

	if keyid.Name != "string" {
		return nil
	}

	t.typ = mt.Value
	op := t.GetOpertion(key, value)
	if op == nil {
		return nil
	}

	_, isMap := op.(*mapOperation)
	if isMap {
		return nil
	}

	return &mapOperation{
		operationInterface: op,
	}
}

func (t *UniType) structCase(key string, value string) operationInterface {
	st := t.typ.(*ast.StructType)

	task := model.GenTask{
		Struct: st,
	}

	stname := "Param_" + key

	baseid, baseIsIdent := t.baseTyp.(*ast.Ident)
	if baseIsIdent {
		stname = baseid.Name
		task.TypeSpec = &ast.TypeSpec{
			Name: baseid,
		}
	} else {
		task.TypeSpec = &ast.TypeSpec{
			Name: &ast.Ident{
				Name: stname,
			},
		}
	}

	t.tasks = append(t.tasks, task)

	buildFuncName := "Build" + stname

	return &structOperation{
		defaultOperation: &defaultOperation{
			Key:   buildFuncName,
			Value: value,
		},
	}
}

func (t *UniType) unwrapCustomType(id *ast.Ident) ast.Expr {
	wt := id.Obj.Decl.(*ast.TypeSpec).Type
	idwt, isIdent := wt.(*ast.Ident)
	if isIdent && idwt.Obj != nil {
		return t.unwrapCustomType(idwt)
	}

	return wt
}
