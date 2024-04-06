package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/inspector"

	"github.com/eebor/sweetquery/querygen/model"
)

const (
	taskName = "querygen:query"
)

func parseSource(sourcePath string) ([]model.GenTask, *ast.File, error) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, sourcePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing source file %v: %v", sourcePath, err)
	}

	return findTasks(file), file, nil
}

func findTasks(file *ast.File) []model.GenTask {
	tasks := []model.GenTask{}

	i := inspector.New([]*ast.File{file})
	//Подготовим фильтр для этого инспектора
	iFilter := []ast.Node{
		//Нас интересуют декларации
		&ast.GenDecl{},
	}

	i.Nodes(iFilter, func(n ast.Node, push bool) (proceed bool) {
		genDecl := n.(*ast.GenDecl)

		if genDecl.Doc == nil {
			return false
		}

		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return false
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}

		for _, comment := range genDecl.Doc.List {
			if strings.Contains(comment.Text, taskName) {
				tasks = append(tasks, model.GenTask{
					Struct:   structType,
					TypeSpec: typeSpec,
				})
			}
		}

		return false
	})

	return tasks
}
