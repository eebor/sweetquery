package gentempl

import "text/template"

var BuilderTempl = template.Must(template.New("").Parse(`
package main

func Build{{ .QueryName }}(obj *{{ .QueryNamePrefix }}{{ .QueryName }}) *query.Query {
	q := query.NewQuery()	
{{ .Operations }}
	return q
}`))
