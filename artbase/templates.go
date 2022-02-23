package artbase

import (
	"html/template"
	"bytes"
	"fmt"
)



//////
// default built-in templates
//   use artbase.Templates to overwrite / change



const templateHome = `
<h1>{{ len . }} Collections</h1>
<ul>
    {{range .}}
        <li>
				  <a href="/{{.Name}}">/{{.Name}}</a>
					({{ .Width}}×{{ .Height}})

          {{if .Background}}
					   - incl. backgrounds
          {{end}}

					<span style="font-size: 80%">
					  <a href="{{ .Url }}">(Download .png)</a>
					</span>
				</li>
    {{end}}
</ul>


<hr>

<p style="font-size: 80%">
  New to artbase? Find out more at
	<a href="https://github.com/pixelartexchange/artbase">/artbase »</a>
</p>
`


const templateIndex = `
<p style="font-size: 80%">
  <a href="/">« Collections</a>
</p>

<h1>{{ .Name}}   ({{ .Width}}×{{ .Height}}) Collection</h1>

<p>
  To get images, use /{{ .Name}}/<em>:id</em>
</p>

<p>
	Example:
	<a href="/{{ .Name}}/0">/{{ .Name}}/0</a>,
	<a href="/{{ .Name}}/1">/{{ .Name}}/1</a>,
	<a href="/{{ .Name}}/2">/{{ .Name}}/2</a>, ...
</p>

<p>
	Note: The default image size is ({{ .Width}}×{{ .Height}}).
	 Use the z (zoom) parameter to upsize.
</p>
<p>
	 Try 2x:
	 <a href="/{{ .Name}}/0?z=2">/{{ .Name}}/0?z=2</a>,
	 <a href="/{{ .Name}}/1?z=2">/{{ .Name}}/1?z=2</a>,
	 <a href="/{{ .Name}}/2?z=2">/{{ .Name}}/2?z=2</a>, ...
</p>

<p>
	 Try 8x:
	 <a href="/{{ .Name}}/0?z=8">/{{ .Name}}/0?z=8</a>,
	 <a href="/{{ .Name}}/1?z=8">/{{ .Name}}/1?z=8</a>,
	 <a href="/{{ .Name}}/2?z=8">/{{ .Name}}/2?z=8</a>, ...

	 And so on.
</p>



<hr>

<p style="font-size: 80%">
  New to artbase? Find out more at
	<a href="https://github.com/pixelartexchange/artbase">/artbase »</a>
</p>

`



///
// todo/check if template.Template is an interface - if yes, no need for pointer!!
//
//  note: template.Template is a struct with a pointer to parse.Tree
//             (NOT an interface), thus, use *template (pointer)

var Templates = make( map[string]*template.Template )





func init() {
	fmt.Println( "  [artbase.init] compileTemplates" )
  compileTemplates()
}

func compileTemplates() {
  Templates["home"]  = template.Must( template.New("home").Parse( templateHome ))
	Templates["index"] = template.Must( template.New("index").Parse( templateIndex ))
}



func RenderHome( data []Collection ) []byte {
	buf := new( bytes.Buffer )
	Templates["home"].Execute( buf, data )
	return buf.Bytes()
}


func RenderCollection( data *Collection ) []byte {
	buf := new( bytes.Buffer )
	Templates["index"].Execute( buf, data )
	return buf.Bytes()
}

