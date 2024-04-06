package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"text/template"

	"github.com/eebor/sweetquery/querygen/model"
)

var builderTempl = template.Must(template.New("").Parse(`
package main

func Build{{ .QueryName }}(req *{{ .QueryNamePrefix }}{{ .QueryName }}) []byte {
	q := query.NewQuery()	
{{ .Operations }}
	return q.Bytes()
}`))

type queryBuilderParams struct {
	QueryName       string
	QueryNamePrefix string
	Operations      string
}

type builderGenerator struct {
	buidlers []ast.Decl
}

func (g *builderGenerator) ProcessTask(task *model.GenTask, prefix string) error {
	operations := make([]writeOperation, len(task.Struct.Fields.List))

	for i, field := range task.Struct.Fields.List {
		tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

		key := tag.Get("query")
		if key == "" {
			continue
		}

		t, _ := field.Type.(*ast.Ident)

		pointer, isPointer := field.Type.(*ast.StarExpr)
		if isPointer {
			t = pointer.X.(*ast.Ident)
		}

		if t == nil {
			return fmt.Errorf("type of %s is not supported", field.Names[0].Name)
		}

		operations[i] = writeOperation{
			Type:      opTypeRelataion[t.Name],
			Key:       strconv.Quote(key),
			Value:     "req." + field.Names[0].Name,
			CheckNull: isPointer,
		}
	}

	opsBuf := bytes.Buffer{}
	for i := 0; i < len(operations); i++ {
		op, err := operations[i].Build()
		if err != nil {
			return err
		}

		op.WriteTo(&opsBuf)
	}

	params := queryBuilderParams{
		QueryName:  task.TypeSpec.Name.Name,
		Operations: opsBuf.String(),
	}

	if prefix != "" {
		params.QueryNamePrefix = prefix + "."
	}

	buildBuf := bytes.Buffer{}

	builderTempl.Execute(&buildBuf, params)

	build, err := parser.ParseFile(token.NewFileSet(), "", buildBuf.Bytes(), 0)
	if err != nil {
		return err
	}

	g.buidlers = append(g.buidlers, build.Decls[0])

	return nil
}

func (g *builderGenerator) GetBuilders() []ast.Decl {
	return g.buidlers
}
