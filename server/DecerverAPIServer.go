package server

import (
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/decerver/interfaces/core"
	"github.com/eris-ltd/decerver/interfaces/dapps"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type SwitchName struct {
	Name string `json:"name"`
}

type DecerverAPIServer struct {
	decerver core.DeCerver
	dappreg  dapps.DappRegistry
}

func NewDecerverAPIServer(dc core.DeCerver, dr dapps.DappRegistry) *DecerverAPIServer {
	return &DecerverAPIServer{dc, dr}
}

func (das *DecerverAPIServer) handleReadyGET(w http.ResponseWriter, r *http.Request) {
	logger.Println("GET decerver ready")

	if !das.decerver.IsStarted() {
		das.writeError(w, 400, "decerver not started")
	}

	dapplist := das.dappreg.GetDappList()
	bts, err := json.Marshal(dapplist)

	if err != nil {
		das.writeError(w, 500, err.Error())
		return
	}
	jsn := string(bts)
	// DEBUG
	logger.Println("Dapplist:\n" + jsn)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, jsn)
}

func (das *DecerverAPIServer) handleDecerverGET(w http.ResponseWriter, r *http.Request) {
	logger.Println("GET decerver config")
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
	logger.Println("POST decerver config")
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

	fio := das.decerver.GetFileIO()

	pt := fio.Modules() + "/" + mName
	logger.Printf("GET %s config\n", mName)

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
	logger.Printf(" POST %s config\n", mName)

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		das.writeError(w, 400, err.Error())
		return
	}

	fio := das.decerver.GetFileIO()
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

func (das *DecerverAPIServer) handleDappSwitch(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	mName := path.Base(url)
	fmt.Println("Url: " + url)
	if mName == "." || mName == "/" || mName == "" {
		das.writeError(w, 404, "Malformed URL")
		return
	}
	logger.Println("Switching to dapp: ", mName)
	err := das.dappreg.LoadDapp(mName)

	if err != nil {
		das.writeError(w, 400, err.Error())
		return
	}

	w.WriteHeader(200)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Whatever...
	fmt.Fprint(w, "success")
}

func (das *DecerverAPIServer) handleFoF(w http.ResponseWriter, r *http.Request) {
	das.writeError(w, 400, "The route not open (the dapp is not in focus).")
}

func (das *DecerverAPIServer) writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, msg)
}
