package utils

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

// HTMLErr is the standard http error page template for the application.
const HTMLErr = `
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

// GetFileTpl parse a file and returns an executable html template or an error.
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

// ExecTpl is a utility function to execute the template that may be removed shortly in the future.
func ExecTpl(w http.ResponseWriter, tpl *template.Template, data interface{}) {
	err := tpl.Execute(w, data)
	if err != nil {
		w.Write([]byte("Error executing template: " + err.Error()))
	}
	return
}
