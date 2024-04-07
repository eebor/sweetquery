package generator

import "go/ast"

type UniType struct {
	typ ast.Expr
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
	}

	return nil
}

func (t *UniType) identCase(key string, value string) operationInterface {
	id := t.typ.(*ast.Ident)
	if id.Obj != nil {
		t.typ = t.unwrapCustomType(id)
		return t.GetOpertion(key, value)
	}

	return &queryWriteOperation{
		defaultOperation: &defaultOperation{
			Key:   key,
			Value: value,
		},
		Type: id.Name,
	}
}

func (t *UniType) pointerCase(key string, value string) operationInterface {
	t.typ = t.typ.(*ast.StarExpr).X
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
	op := t.GetOpertion(key, value)
	if op == nil {
		return nil
	}

	_, isArray := op.(*arrayOperation)
	if isArray {
		return nil
	}

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

func (t *UniType) unwrapCustomType(id *ast.Ident) ast.Expr {
	wt := id.Obj.Decl.(*ast.TypeSpec).Type
	idwt, isIdent := wt.(*ast.Ident)
	if isIdent && idwt.Obj != nil {
		return t.unwrapCustomType(idwt)
	}

	return wt
}
