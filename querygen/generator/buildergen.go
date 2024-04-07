package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"

	"github.com/eebor/sweetquery/querygen/generator/gentempl"
	"github.com/eebor/sweetquery/querygen/model"
)

type queryBuilderParams struct {
	QueryName       string
	QueryNamePrefix string
	Operations      string
}

type builderGenerator struct {
	buidlers          []ast.Decl
	buildersNameCache map[string]bool
}

func (g *builderGenerator) ProcessTask(task *model.GenTask, prefix string) error {
	operations := make([]operationInterface, 0)

	if g.buildersNameCache != nil {
		_, isProcessed := g.buildersNameCache[task.TypeSpec.Name.Name]
		if isProcessed {
			return nil
		}
	} else {
		g.buildersNameCache = make(map[string]bool)
	}

	for _, field := range task.Struct.Fields.List {
		tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

		key := tag.Get("query")
		if key == "" {
			continue
		}
		if prefix != "" {
			key = prefix + "[" + key + "]"
		}

		t := NewUniType(field.Type)

		var op operationInterface
		for _, name := range field.Names {
			t.Reset()
			op = t.GetOpertion(key, "obj."+name.Name)
			if op == nil {
				return fmt.Errorf("type of %s is not supported", field.Names[0].Name)
			}

			operations = append(operations, op)
		}

		for _, task := range t.GetTasks() {
			g.ProcessTask(&task, key)
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

	buildBuf := bytes.Buffer{}

	gentempl.BuilderTempl.Execute(&buildBuf, params)

	build, err := parser.ParseFile(token.NewFileSet(), "", buildBuf.Bytes(), 0)
	if err != nil {
		return err
	}

	g.buidlers = append(g.buidlers, build.Decls[0])

	g.buildersNameCache[task.TypeSpec.Name.Name] = true

	return nil
}

func (g *builderGenerator) GetBuilders() []ast.Decl {
	return g.buidlers
}
