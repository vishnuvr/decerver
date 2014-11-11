package server

import (
	"fmt"
	"github.com/eris-ltd/deCerver-interfaces/api"
	"github.com/go-martini/martini"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

const defaultPort = 3000

type WebServer struct {
	Martini               *martini.ClassicMartini
	maxConnections        uint32
	httpAPIServices       []interface{}
	wsAPIServiceFactories []api.WsAPIServiceFactory
	appsDirectory         string
	port                  int
}

func NewWebServer(maxConnections uint32, appDir string, port int) *WebServer {
	ws := &WebServer{}
	ws.maxConnections = maxConnections
	ws.httpAPIServices = make([]interface{}, 0)
	ws.wsAPIServiceFactories = make([]api.WsAPIServiceFactory, 0)
	ws.appsDirectory = appDir
	if port <= 0 {
		port = defaultPort
	}
	ws.port = port
	
	return ws
}

func (ws *WebServer) RegisterHttpServices(service ...interface{}) {
	for _, s := range service {
		ws.httpAPIServices = append(ws.httpAPIServices, s)
	}
}

func (ws *WebServer) RegisterWsServiceFactories(factory ...api.WsAPIServiceFactory) {
	for _, f := range factory {
		ws.wsAPIServiceFactories = append(ws.wsAPIServiceFactories, f)
	}
}

func (ws *WebServer) Start() {

	ws.Martini = martini.Classic()
	// TODO make this settable
	//so := martini.StaticOptions{}
	ws.Martini.Use(martini.Static(ws.appsDirectory))

	// Change to production environment.
	// martini.Env = martini.Prod

	// JSON RPC
	if len(ws.httpAPIServices) > 0 {
		rpcs := rpc.NewServer()
		rpcs.RegisterCodec(json2.NewCodec(), "application/json")
		for _, service := range ws.httpAPIServices {
			rpcs.RegisterService(service, "")
		}
		ws.Martini.Post("/httpapi", rpcs.ServeHTTP)
	}

	// JSON Socket RPC
	if len(ws.wsAPIServiceFactories) > 0 {
		wsapis := NewWsAPIServer(ws.maxConnections)
		for _, factory := range ws.wsAPIServiceFactories {
			wsapis.RegisterServiceFactory(factory)
			fmt.Println("Registering new websocket API: " + factory.ServiceName())
		}
		ws.Martini.Get("/wsapi", wsapis.handleWs)
	}
	
	// Decerver communication
	ws.Martini.Get("/admin/decerver",handleDecerverGET)
	ws.Martini.Post("/admin/decerver",handleDecerverPOST)
	
	//"decerver/config/update/"

	go func() {
		ws.Martini.RunOnAddr("localhost:" + fmt.Sprintf("%d",ws.port))
	}()

}
