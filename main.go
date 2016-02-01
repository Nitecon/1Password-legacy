package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/user"

	"github.com/chzyer/logex"
	"github.com/julienschmidt/httprouter"
)

var (
	passwordLoc = "/home/whatting/GoogleDrive/1Password.agilekeychain"
)

func handler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	basedir, _ := getUserData()
	dataFileServer := http.FileServer(http.Dir(basedir))
	if p.ByName("filepath") != "/" {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=60")
		dataFileServer.ServeHTTP(w, r)
		return
	}
	f, err := ioutil.ReadFile(basedir + "/1Password.html")
	if err != nil {
		fmt.Fprint(w, errorMessage)
		return
	}
	fmt.Fprint(w, string(f))

	return
}

func getRandPort() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func getUserData() (basedir string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	logex.Println(usr.HomeDir)
	basedir = passwordLoc
	return
}

func main() {
	randPort := getRandPort()

	router := httprouter.New()
	router.GET("/*filepath", handler)
	fmt.Printf("http://localhost:%d", randPort)
	http.ListenAndServe(fmt.Sprintf(":%d", randPort), router)

}

var errorMessage = `
<html>
  <body>
    <h2>An error occurrred!</h2>
    <p>
    <img height="50px" width="120px" src="http://i.imgur.com/DwFKI0J.png"/>
    	Could not load home directory for user to get preferences.
    </p>

  </body>
</html>
`
