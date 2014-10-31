package server

import (
	"github.com/eris-ltd/deCerver-interfaces/api"
	"github.com/go-martini/martini"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"log"
)

var logger *log.Logger

type WebServer struct {
	Martini               *martini.ClassicMartini
	maxConnections        uint32
	httpAPIServices       []interface{}
	wsAPIServiceFactories []api.WsAPIServiceFactory
}

func NewWebServer(maxConnections uint32, log *log.Logger) *WebServer {
	logger = log
	ws := &WebServer{}
	ws.maxConnections = maxConnections
	ws.httpAPIServices = make([]interface{}, 0)
	ws.wsAPIServiceFactories = make([]api.WsAPIServiceFactory, 0)
	return ws
}

func (ws *WebServer) RegisterHttpAPIService(service interface{}) {
	ws.httpAPIServices = append(ws.httpAPIServices, service)
}

func (ws *WebServer) RegisterWsAPIServiceFactory(factory api.WsAPIServiceFactory) {
	ws.wsAPIServiceFactories = append(ws.wsAPIServiceFactories, factory)
}

func (ws *WebServer) Start() {

	ws.Martini = martini.Classic()
	// TODO make this settable
	ws.Martini.Use(martini.Static("./web"))

	// Change to production environment.
	// martini.Env = martini.Prod

	// JSON RPC
	if len(ws.httpAPIServices) > 0 {
		rpcs := rpc.NewServer()
		rpcs.RegisterCodec(json2.NewCodec(), "application/json")
		for _, service := range ws.httpAPIServices {
			rpcs.RegisterService(service, "")
			logger.Printf("Say: %b\n",rpcs.HasMethod( "Say" ))
		}
		ws.Martini.Post("/httpapi", rpcs.ServeHTTP)
	}

	// JSON Socket RPC
	if len(ws.wsAPIServiceFactories) > 0 {
		wsapis := NewWsAPIServer(ws.maxConnections)
		for _, factory := range ws.wsAPIServiceFactories {
			wsapis.RegisterServiceFactory(factory, "")
		}
		ws.Martini.Get("/wsapi", wsapis.handleWs)
	}

	go func() {
		ws.Martini.RunOnAddr("localhost:3000")
	}()

}
