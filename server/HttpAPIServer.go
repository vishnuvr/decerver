package server

import (
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/core"
	"net/http"
	"path"
)

type HttpResp struct {
	Status int				  `json:"status"`
	Header map[string]string  `json:"header"`
	Body   string             `json:"body"`
}

type HttpAPIServer struct {
	ate core.RuntimeManager
}

func NewHttpAPIServer(rm core.RuntimeManager) *HttpAPIServer {
	return &HttpAPIServer{rm}
}

// This is our basic http receiver that takes the request and passes it into the js runtime.
func (has *HttpAPIServer) handleHttp(w http.ResponseWriter, r *http.Request) {

	u := r.URL
	p := u.Path
	caller := path.Base(p)
	
	rt := has.ate.GetRuntime(caller)
	// TODO Update this. It's basically how we check if dapp is ready now.
	if rt == nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, "Dapp not in focus")
		return
	}
	
	reqJson, _ := json.Marshal(r)
	
	ret, err := rt.CallFuncOnObj("network", "incomingHttp", reqJson)

	if err != nil {
		has.writeError(w, 500, err.Error())
		return
	}
	
	rStr := ret.(string)
	hr := &HttpResp{}
	errJson := json.Unmarshal([]byte(rStr), hr)
	
	if errJson != nil {
		has.writeError(w, 500, errJson.Error())
		return
	}
	
	has.writeReq(hr,w)
}

func (has *HttpAPIServer) writeReq(resp *HttpResp, w http.ResponseWriter) {
	w.WriteHeader(resp.Status)
	for k, v := range resp.Header {
		w.Header().Set(k,v)
	}
	w.Write([]byte(resp.Body))
}

func (has *HttpAPIServer) writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, msg)
}
