package gentempl

import "text/template"

var BuilderTempl = template.Must(template.New("").Parse(`
package main

func Build{{ .QueryName }}(req *{{ .QueryNamePrefix }}{{ .QueryName }}) []byte {
	q := query.NewQuery()	
{{ .Operations }}
	return q.Bytes()
}`))
