package generator

import (
	"bytes"
	"text/template"
)

var opWriteTempl = template.Must(template.New("").Parse(`
q.Write{{ .Type }}({{ .Key }}, {{ .Value }})
`))

var opWriteCondTempl = template.Must(template.New("").Parse(`
if {{ .Value }} != nil {
	q.Write{{ .Type }}({{ .Key }}, *{{ .Value }})
}
`))

var opTypeRelataion = map[string]string{
	"string": "String",
	"int":    "Int",
	"uint":   "Uint",
	"bool":   "Bool",
	"float":  "Float",
}

type writeOperation struct {
	Type      string
	Key       string
	Value     string
	CheckNull bool
}

func (o *writeOperation) Build() (*bytes.Buffer, error) {
	var buf bytes.Buffer

	var err error
	if o.CheckNull {
		err = opWriteCondTempl.Execute(&buf, o)
	} else {
		err = opWriteTempl.Execute(&buf, o)
	}

	if err != nil {
		return nil, err
	}

	return &buf, nil
}
