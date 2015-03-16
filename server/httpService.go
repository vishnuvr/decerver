package server

import (
	"fmt"
	"net/http"
)

const EPM_HELP = "Wrong command, asshole."

type HttpService struct {
	// Maybe keep track of some statistics if this is used to create chains 
	// via some Eris web service later, like it works with the compilers.
}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (this *HttpService) handleNotFound(w http.ResponseWriter, r *http.Request) {
	this.writeError(w, 404, EPM_HELP)
}

func (this *HttpService) writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, msg)
}
