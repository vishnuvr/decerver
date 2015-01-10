package server

import (
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/decerver/interfaces/core"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpRespProxy struct {
	URL *url.URL
	Method string
	Host string
	Header http.Header
	Body string
}

func ProxyFromHttpReq(r *http.Request) (*HttpRespProxy, error) {
	p := &HttpRespProxy{}
	p.Method = r.Method
	r.Host = r.Host
	p.URL = r.URL
	p.Header = r.Header
	bts, err := ioutil.ReadAll(r.Body);
	if err != nil {
		return nil, err;
	} else {
		p.Body = string(bts)
		return p, nil
	} 
}

type HttpResp struct {
	Status int               `json:"status"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
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
	caller := strings.Split(strings.TrimLeft(p,"/"),"/")[1];
	
	rt := has.ate.GetRuntime(caller)
	// TODO Update this. It's basically how we check if dapp is ready now.
	if rt == nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, "Dapp not in focus")
		return
	}
	
	// logger.Println("Incoming: %v\n", r)
	
	prx, errpr := ProxyFromHttpReq(r)
	// Shouldn't happen.
	if errpr != nil {
		has.writeError(w, 400, errpr.Error())
		return
	}
	reqJson, errM := json.Marshal(prx)

	if errM != nil {
		logger.Println("Error when marshalling http request (this really should not happen) : " + errM.Error())
	}
	// logger.Println("Http request json: " + string(reqJson))
	ret, err := rt.CallFuncOnObj("network", "handleIncomingHttp", string(reqJson))

	if err != nil {
		has.writeError(w, 500, err.Error())
		return
	}

	rStr, sOk := ret.(string)
	if !sOk {
		has.writeError(w, 500, "Passing non string as return value from otto.")
		return
	}
	hr := &HttpResp{}
	errJson := json.Unmarshal([]byte(rStr), hr)

	if errJson != nil {
		has.writeError(w, 500, errJson.Error())
		return
	}

	has.writeReq(hr, w)
}

func (has *HttpAPIServer) writeReq(resp *HttpResp, w http.ResponseWriter) {
	logger.Printf("Response status message: %d\n", resp.Status)
	logger.Printf("Response header stuff: %v\n", resp.Header)
	w.WriteHeader(resp.Status)
	for k, v := range resp.Header {
		w.Header().Set(k, v)
	}
	w.Write([]byte(resp.Body))
}

func (has *HttpAPIServer) writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, msg)
}
