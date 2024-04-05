package generator

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/eebor/sweetquery/querygen/model"
)

const header = `
// Copyright 2024 eebor. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by "querygen"; DO NOT EDIT.
`

const queryPackage = "github.com/eebor/sweetquery/query"

func Generate(pkg *model.Package, out io.Writer) error {
	g := builderGenerator{}

	for _, task := range pkg.Tasks {
		err := g.ProcessTask(&task)
		if err != nil {
			return fmt.Errorf("\npackage: %s\ntarget: %s\n\n%s", pkg.Source.Name, task.TypeSpec.Name.Name, err.Error())
		}
	}

	builderDecls := g.GetBuilders()

	astOut := &ast.File{
		Name: pkg.Source.Name,
	}

	astOut.Decls = append(astOut.Decls, builderDecls...)

	astutil.AddNamedImport(token.NewFileSet(), astOut, "query", queryPackage)

	out.Write([]byte(header))

	err := printer.Fprint(out, token.NewFileSet(), astOut)
	if err != nil {
		return fmt.Errorf("package: %s\n\n%s", pkg.Source.Name, err.Error())
	}

	return nil
}
