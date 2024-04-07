package gentempl

import "text/template"

var BuilderTempl = template.Must(template.New("").Parse(`
package main

func Build{{ .QueryName }}(obj *{{ .QueryNamePrefix }}{{ .QueryName }}, q *query.Query) {
{{ .Operations }}
}`))
