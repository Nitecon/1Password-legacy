package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"

	"os"

	"encoding/json"

	"log"

	"github.com/Nitecon/1Password/utils"
	"github.com/julienschmidt/httprouter"
)

func show1PData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cfg := utils.GetConfig()
	dfs := http.FileServer(http.Dir(cfg.MainLocation))
	fp := p.ByName("filepath")
	if fp == "" {
		fp = "/"
	}
	if fp != "/" {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=60")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		dfs.ServeHTTP(w, r)
		return
	}
	f, err := ioutil.ReadFile(filepath.ToSlash(cfg.MainLocation + "/1Password.html"))
	if err != nil {
		log.Printf("Error reading 1password file: %s", err.Error())
		fmt.Fprint(w, utils.HTMLErr)
		return
	}
	fmt.Fprint(w, string(f))
	return
}

func showSetConfData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
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

func getRandPort() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func getConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	cfg := utils.GetConfig()
	d, err := json.Marshal(cfg)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(d)
	return
}

func validateVaultPath(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	appData := utils.AppData{}
	appData.StatusCode = 200
	if err != nil {
		appData.IsError = true
		appData.StatusCode = 400
		appData.Error = "Bad Request Body Received"
		appData.WriteRestResponse(w)
		return
	}
	defer r.Body.Close()
	var d utils.Configuration
	err = json.Unmarshal(body, &d)
	if err != nil {
		appData.IsError = true
		appData.StatusCode = 500
		appData.Error = "Invalid Data Encoding"
		appData.WriteRestResponse(w)
		return
	}
	log.Printf("Vault Location: %s", d.MainLocation)

	vaultPath := d.MainLocation + "/1Password.html"
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		appData.IsError = true
		appData.StatusCode = 404
		appData.Error = "Vault not found at location"
		appData.WriteRestResponse(w)
		return
	}
	d.InitialConfig = true
	err = utils.UpdateConfig(d)
	if err != nil {
		appData.IsError = true
		appData.StatusCode = 500
		appData.Error = err.Error()
		appData.WriteRestResponse(w)
	}
	utils.SetConfig()

	appData.Content = "Success"
	appData.WriteRestResponse(w)
	return
}

func main() {
	err := utils.SetConfig()
	if err != nil {
		panic("Could not set base config: \n" + err.Error())
	}
	randPort := getRandPort()
	router := httprouter.New()
	cfg := utils.GetConfig()
	if !cfg.InitialConfig {
		router.GET("/", showSetConfData)
		router.GET("/rest/config", getConfig)
		router.POST("/rest/config/validateVault", validateVaultPath)
		fileServer := http.FileServer(http.Dir("static"))
		router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Header().Set("Vary", "Accept-Encoding")
			w.Header().Set("Cache-Control", "public, max-age=60")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			r.URL.Path = p.ByName("filepath")
			fileServer.ServeHTTP(w, r)
		})
	} else {
		router.GET("/*filepath", show1PData)

	}
	fmt.Printf("http://localhost:%d", randPort)
	http.ListenAndServe(fmt.Sprintf(":%d", randPort), router)
}
