package utils

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

const HtmlErr = `
<html>
  <body>
    <h2>An error occurrred!</h2>
    <p>
    <img height="50px" width="120px" src="http://i.imgur.com/DwFKI0J.png"/>
    	{{.}}
    </p>
  </body>
</html>
`

func GetFileTpl(path string) (*template.Template, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t, err := template.New("Single").Parse(string(f))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func ExecTpl(w http.ResponseWriter, tpl *template.Template, data interface{}) {
	err := tpl.Execute(w, data)
	if err != nil {
		w.Write([]byte("Error executing template: " + err.Error()))
	}
	return
}
