package model

import "go/ast"

type Package struct {
	Source *ast.File
	Tasks  []GenTask
}

type GenTask struct {
	Struct   *ast.StructType
	TypeSpec *ast.TypeSpec
}
