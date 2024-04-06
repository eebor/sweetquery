package model

import "go/ast"

type Package struct {
	Name               string
	Imports            map[string]string
	Tasks              []GenTask
	SrcPkgHandlePrefix string
}

type GenTask struct {
	Struct   *ast.StructType
	TypeSpec *ast.TypeSpec
}
