package server

import (
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/core"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type DecerverAPIServer struct {
	decerver core.DeCerver
}

func NewDecerverAPIServer(dc core.DeCerver) *DecerverAPIServer {
	return &DecerverAPIServer{dc}
}

func (das *DecerverAPIServer) handleDecerverGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[martini] GET decerver config")
	cfg := das.decerver.GetConfig()

	bts, err := json.Marshal(cfg)

	if err != nil {
		das.writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bts))
}

func (das *DecerverAPIServer) handleDecerverPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[martini] POST decerver config")
	contentType := r.Header.Get("Content-Type")

	idx := strings.Index(contentType, ";")
	if idx != -1 {
		contentType = contentType[:idx]
	}
	ct := strings.ToLower(contentType)

	if ct != "application/json" {
		das.writeError(w, 415, "unrecognized Content-Type: "+contentType)
		return
	}

	bts, err := ioutil.ReadAll(r.Body)

	if err != nil {
		das.writeError(w, 400, err.Error())
		return
	}
	cfg := &core.DCConfig{}
	fmt.Print(string(bts))
	err = json.Unmarshal(bts, cfg)

	if err != nil {
		das.writeError(w, 422, err.Error())
		return
	}

	das.decerver.WriteConfig(cfg)

	w.WriteHeader(204)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

// Modules
func (das *DecerverAPIServer) handleModuleGET(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	mName := path.Base(url)
	if mName == "." || mName == "/" {
		das.writeError(w, 404, "Malformed URL")
		return
	}

	fio := das.decerver.GetPaths()

	pt := fio.Modules() + "/" + mName
	fmt.Printf("[martini] GET %s config\n", mName)

	bts, err := fio.ReadFile(pt, "config")

	if err != nil {
		das.writeError(w, 500, err.Error())
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bts))
}

func (das *DecerverAPIServer) handleModulePOST(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	idx := strings.Index(contentType, ";")
	if idx != -1 {
		contentType = contentType[:idx]
	}
	ct := strings.ToLower(contentType)

	if ct != "application/json" {
		das.writeError(w, 415, "unrecognized Content-Type: "+contentType)
		return
	}

	url := r.URL.String()
	mName := path.Base(url)
	if mName == "." || mName == "/" {
		das.writeError(w, 404, "Malformed URL")
		return
	}
	fmt.Printf("[martini] POST %s config\n", mName)

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		das.writeError(w, 400, err.Error())
		return
	}

	fio := das.decerver.GetPaths()
	pt := fio.Modules() + "/" + mName

	var tmp_int interface{}
	err = json.Unmarshal(bts, &tmp_int)
	if err != nil {
		das.writeError(w, 400, err.Error())
		return
	}

	bts, err = json.MarshalIndent(tmp_int, "", "    ")
	if err != nil {
		das.writeError(w, 400, err.Error())
		return
	}

	fio.WriteFile(pt, "config", bts)
	w.WriteHeader(204)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (das *DecerverAPIServer) writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, msg)
}
