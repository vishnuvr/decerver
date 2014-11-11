package server



import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/eris-ltd/deCerver-interfaces/core"
)

var deCerver core.DeCerver

func Init(dc core.DeCerver) {
	deCerver = dc
}

func handleDecerverGET(w http.ResponseWriter, r *http.Request){
	fmt.Println("Decerver GET")
	cfg := deCerver.GetConfig()
	
	bts, _ := json.Marshal(cfg)
	
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, bts)
}

func handleDecerverPOST(w http.ResponseWriter, r *http.Request){
	fmt.Println("Decerver POST")
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
		fmt.Println(err.Error())
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
	
	path := r.URL.String()
	if path[len(path) - 1] != '/' {
		path = path[:len(path) - 1]
	}
	
	segs := strings.Split(path,"/")
	if segs == nil || len(segs) == 0 {
		writeError(w, 404, "Malformed URL")
		return
	}
	
	mName := segs[len(segs) - 1]
	fmt.Println(mName)
	
	cfg := deCerver.GetConfig()
	
	bts, _ := json.Marshal(cfg)
	
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, bts)
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
	
	bts, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	cfg := &core.DCConfig{}
	
	err = json.Unmarshal(bts,cfg)
	
	if err != nil {
		writeError(w,422,err.Error())
	}
	
	deCerver.WriteConfig(cfg)
	
	w.WriteHeader(204)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") 
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, msg)
}
