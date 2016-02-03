package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"os"

	"github.com/Nitecon/1Password/utils"
	"github.com/julienschmidt/httprouter"
)

func show1PData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cfg := utils.GetConfig()
	dataFileServer := http.FileServer(http.Dir(cfg.MainLocation))
	if p.ByName("filepath") != "/" {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=60")
		dataFileServer.ServeHTTP(w, r)
		return
	}
	f, err := ioutil.ReadFile(fmt.Sprintf("%s%q%s", cfg.MainLocation, os.PathSeparator, "1Password.html"))
	if err != nil {
		fmt.Fprint(w, utils.HtmlErr)
		return
	}
	fmt.Fprint(w, string(f))
	return
}

func showSetConfData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpl, err := utils.GetFileTpl("tpl/main_config.gohtml")
	if err != nil {
		w.Write([]byte("A critical error occurred, exiting: " + err.Error()))
		os.Exit(1)
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		w.Write([]byte("Could not execute template: " + err.Error()))
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cfg := utils.GetConfig()
	if !cfg.InitialConfig {
		showSetConfData(w, r, p)
		return
	}

	show1PData(w, r, p)
	return
}

func getRandPort() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func main() {
	err := utils.SetConfig()
	if err != nil {
		panic("Could not set base config: \n" + err.Error())
	}
	randPort := getRandPort()
	router := httprouter.New()
	router.GET("/*filepath", handler)
	fmt.Printf("http://localhost:%d", randPort)
	http.ListenAndServe(fmt.Sprintf(":%d", randPort), router)
}
