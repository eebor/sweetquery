package gentempl

import "text/template"

var OpWriteTempl = template.Must(template.New("").Parse(`
q.Write{{ .Type }}({{ .Key }}, {{ .Value }})
`))

var OpPointerCondTempl = template.Must(template.New("").Parse(`
if {{ .Value }} != nil {
	{{ .Operation }}
}
`))

var OpArrayTempl = template.Must(template.New("").Parse(`
for i := 0; i < len({{ .Value }}); i++ {
	{{ .Operation }}
}
`))

var OpMapTempl = template.Must(template.New("").Parse(`
for key := range {{ .Value }} {
	{{ .Operation }}
}
`))
