package generator

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/eebor/sweetquery/querygen/generator/gentempl"
)

type operationInterface interface {
	Build() (*bytes.Buffer, error)
	GetValue() string
	AddBuildValuePrefix(string)
	GetValuePrefix() string
	AddBuildValueSufix(string)
	GetValueSufix() string
	GetKey() string
	AddBuildKeyPrefix(string)
	GetKeyPrefix() string
	AddBuildKeySufix(string)
	GetKeySufix() string
}

type defaultOperation struct {
	Key      string
	Value    string
	v_prefix string
	v_sufix  string
	k_prefix string
	k_sufix  string
}

func (o *defaultOperation) Build() (*bytes.Buffer, error) {
	return nil, errors.New("defaultOperation does not implement build function")
}

func (o *defaultOperation) GetValue() string {
	return o.Value
}

func (o *defaultOperation) AddBuildValuePrefix(prefix string) {
	o.v_prefix = prefix + o.v_prefix
}

func (o *defaultOperation) GetValuePrefix() string {
	return o.v_prefix
}

func (o *defaultOperation) AddBuildValueSufix(sufix string) {
	o.v_sufix += sufix
}

func (o *defaultOperation) GetValueSufix() string {
	return o.v_sufix
}

func (o *defaultOperation) GetKey() string {
	return o.Key
}

func (o *defaultOperation) AddBuildKeyPrefix(prefix string) {
	o.k_prefix = prefix + o.k_prefix
}

func (o *defaultOperation) GetKeyPrefix() string {
	return o.k_prefix
}

func (o *defaultOperation) AddBuildKeySufix(sufix string) {
	o.k_sufix += sufix
}

func (o *defaultOperation) GetKeySufix() string {
	return o.k_sufix
}

var queryTypeRelataion = map[string]string{
	"string":  "String",
	"int":     "Int",
	"int8":    "Int",
	"int16":   "Int",
	"int32":   "Int",
	"int64":   "Int",
	"uint":    "Uint",
	"bool":    "Bool",
	"float":   "Float",
	"float64": "Float",
}

var queryTypeConvertion = map[string]string{
	"int":   "int64",
	"int8":  "int64",
	"int16": "int64",
	"int32": "int64",
}

type queryWriteOperation struct {
	*defaultOperation
	Type       string
	CustomType bool
}

func (o *queryWriteOperation) Build() (*bytes.Buffer, error) {
	typ, ok := queryTypeRelataion[o.Type]
	if !ok {
		return nil, fmt.Errorf("QueryWrite is not support %s", o.Type)
	}

	val := o.v_prefix + o.Value + o.v_sufix

	if o.CustomType {
		conv, ok := queryTypeConvertion[o.Type]
		if ok {
			val = conv + "(" + val + ")"
		} else {
			val = o.Type + "(" + val + ")"
		}

	}

	key := `"` + o.k_prefix + o.Key + o.k_sufix + `"`

	params := struct {
		Key   string
		Value string
		Type  string
	}{
		Key:   key,
		Value: val,
		Type:  typ,
	}

	var buf bytes.Buffer

	err := gentempl.OpWriteTempl.Execute(&buf, &params)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

type structOperation struct {
	*defaultOperation
	BuildFuncName string
}

func (o *structOperation) Build() (*bytes.Buffer, error) {
	val := o.Value + o.v_sufix

	if len(o.v_prefix) > 0 && o.v_prefix[0] == '*' {
		val = o.v_prefix[1:] + val
	} else {
		val = "&" + o.v_prefix + val
	}

	params := struct {
		BuildFuncName string
		Value         string
	}{
		BuildFuncName: o.Key,
		Value:         val,
	}

	var buf bytes.Buffer

	err := gentempl.OpStructTempl.Execute(&buf, &params)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

type pointerCondOperation struct {
	operationInterface
}

func (o *pointerCondOperation) Build() (*bytes.Buffer, error) {
	prefix := o.GetValuePrefix()

	o.AddBuildValuePrefix("*")

	op, err := o.operationInterface.Build()
	if err != nil {
		return nil, err
	}

	params := struct {
		Value     string
		Operation string
	}{
		Value:     prefix + o.GetValue(),
		Operation: op.String(),
	}

	var buf bytes.Buffer

	err = gentempl.OpPointerCondTempl.Execute(&buf, &params)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

type arrayOperation struct {
	operationInterface
}

func (o *arrayOperation) Build() (*bytes.Buffer, error) {
	prefix := o.GetValuePrefix()
	sufix := o.GetValueSufix()

	if len(prefix) > 0 && prefix[0] == '*' {
		o.AddBuildValuePrefix("(")
		o.AddBuildValueSufix(")")
	}

	o.AddBuildValueSufix("[i]")
	o.AddBuildKeySufix("[]")

	op, err := o.operationInterface.Build()
	if err != nil {
		return nil, err
	}

	val := prefix + o.GetValue() + sufix

	params := struct {
		Value     string
		Operation string
	}{
		Value:     val,
		Operation: op.String(),
	}

	var buf bytes.Buffer

	err = gentempl.OpArrayTempl.Execute(&buf, &params)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

type mapOperation struct {
	operationInterface
	CustomKeyType bool
}

func (o *mapOperation) Build() (*bytes.Buffer, error) {
	prefix := o.GetValuePrefix()

	if len(prefix) > 0 && prefix[0] == '*' {
		o.AddBuildValuePrefix("(")
		o.AddBuildValueSufix(")")
	}

	o.AddBuildKeySufix("[\" + ")

	if o.CustomKeyType {
		o.AddBuildKeySufix("string(key)")
	} else {
		o.AddBuildKeySufix("key")
	}

	o.AddBuildKeySufix(" + \"]")
	o.AddBuildValueSufix("[key]")

	op, err := o.operationInterface.Build()
	if err != nil {
		return nil, err
	}

	val := prefix + o.GetValue()

	params := struct {
		Value     string
		Operation string
	}{
		Value:     val,
		Operation: op.String(),
	}

	var buf bytes.Buffer

	err = gentempl.OpMapTempl.Execute(&buf, &params)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
