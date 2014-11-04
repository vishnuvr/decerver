package server

import (
	"github.com/go-martini/martini"
	"net/http"
)

// Interface that is used for webhook handlers
type WebHookHandler interface {
	Handle(res http.ResponseWriter, req *http.Request)
}

type WebServer struct {
	Martini *martini.ClassicMartini
}

func NewWebServer() *WebServer {
	return &WebServer{}
}

// Add the provided handler then start handling posts on /postreceive
func (ws *WebServer) Start(handler interface{}) {

	hdlr := handler.(WebHookHandler)
	ws.Martini = martini.Classic()
	ws.Martini.Post("/postreceive", hdlr.Handle)

	ws.Martini.Run()

}
