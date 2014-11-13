package server

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path"
	"github.com/eris-ltd/deCerver-interfaces/core"
)

var deCerver core.DeCerver

func Init(dc core.DeCerver) {
	deCerver = dc
}

func handleDecerverGET(w http.ResponseWriter, r *http.Request){
	fmt.Println("[martini] GET deCerver config")
	cfg := deCerver.GetConfig()

	bts, err := json.Marshal(cfg)

	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bts))
}

func handleDecerverPOST(w http.ResponseWriter, r *http.Request){
	fmt.Println("[martini] POST deCerver config")
	contentType := r.Header.Get("Content-Type")

	idx := strings.Index(contentType, ";")
	if idx != -1 {
		contentType = contentType[:idx]
	}
	ct := strings.ToLower(contentType)

	if ct != "application/json" {
		writeError(w, 415, "unrecognized Content-Type: " + contentType)
		return
	}

	bts, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	cfg := &core.DCConfig{}
	fmt.Print(string(bts))
	err = json.Unmarshal(bts,cfg)

	if err != nil {
		writeError(w,422,err.Error())
		return
	}

	deCerver.WriteConfig(cfg)

	w.WriteHeader(204)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

// Modules
func handleModuleGET(w http.ResponseWriter, r *http.Request){
	url := r.URL.String()
	mName := path.Base(url)
	if mName == "." || mName == "/" {
		writeError(w, 404, "Malformed URL")
		return
	}

	fio := deCerver.GetPaths()

	pt := fio.Modules() + "/" + mName
	fmt.Printf("[martini] GET %s config\n", mName)

	bts, err := fio.ReadFile(pt,"config")

	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bts))
}

func handleModulePOST(w http.ResponseWriter, r *http.Request){

	contentType := r.Header.Get("Content-Type")
	idx := strings.Index(contentType, ";")
	if idx != -1 {
		contentType = contentType[:idx]
	}
	ct := strings.ToLower(contentType)

	if ct != "application/json" {
		writeError(w, 415, "unrecognized Content-Type: " + contentType)
		return
	}

	url := r.URL.String()
	mName := path.Base(url)
	if mName == "." || mName == "/" {
		writeError(w, 404, "Malformed URL")
		return
	}
	fmt.Printf("[martini] POST %s config\n", mName)

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	fio := deCerver.GetPaths()
	pt := fio.Modules() + "/" + mName

	var tmp_int interface{}
	err = json.Unmarshal(bts,&tmp_int)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	bts, err = json.MarshalIndent(tmp_int, "", "    ")
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	fio.WriteFile(pt,"config",bts)
	w.WriteHeader(204)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, msg)
}
