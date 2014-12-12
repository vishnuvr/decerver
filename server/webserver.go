package server

import (
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver/dappregistry"
	"github.com/go-martini/martini"
)

const DEFAULT_PORT = 3000 // For communicating with dapps (the atom browser).
const DECERVER_PORT = 3005 // For communication with the atom client back-end.

const HTTP_BASE = "/http/"
const WS_BASE = "/ws/"

type endpoint struct {
	method  string
	route   string
	handler interface{}
}

type WebServer struct {
	webServer             *martini.ClassicMartini
	maxConnections        uint32
	appsDirectory         string
	port                  int
	ate                   core.RuntimeManager
	decerver              core.DeCerver
	was					  *WsAPIServer
	has                   *HttpAPIServer
	das					  *DecerverAPIServer
	dr					  *dappregistry.DappRegistry
}

func NewWebServer(maxConnections uint32, appDir string, port int, ate core.RuntimeManager, dc core.DeCerver) *WebServer {
	ws := &WebServer{}
	
	ws.maxConnections = maxConnections
	ws.appsDirectory = appDir
	if port <= 0 {
		port = DEFAULT_PORT
	}
	ws.port = port
	ws.ate = ate
	ws.decerver = dc
	
	ws.was = NewWsAPIServer(ws.ate, ws.maxConnections)
	ws.has = NewHttpAPIServer(ws.ate)
	
	ws.webServer = martini.Classic()
	// TODO remember to change to martini.Prod 
	martini.Env = martini.Dev
	
	return ws
}

func (ws *WebServer) RegisterDapp(dappId string){
	fmt.Println("Adding routes for: " + dappId)
	ws.webServer.Any(HTTP_BASE + dappId, ws.has.handleHttp)
	ws.webServer.Get(WS_BASE + dappId, ws.was.handleWs)
}

func (ws *WebServer) AddDappRegistry(dr *dappregistry.DappRegistry){
	ws.dr = dr
}

func (ws *WebServer) Start() error {
	
	ws.webServer.Use(martini.Static(ws.appsDirectory))

	das := NewDecerverAPIServer(ws.decerver, ws.dr)

	// Decerver ready
	ws.webServer.Get("/admin/ready", das.handleReadyGET)

	// Decerver configuration
	ws.webServer.Get("/admin/decerver", das.handleDecerverGET)
	ws.webServer.Post("/admin/decerver", das.handleDecerverPOST)

	// Module configuration
	ws.webServer.Get("/admin/modules/(.*)", das.handleModuleGET)
	ws.webServer.Post("/admin/modules/(.*)", das.handleModulePOST)

	// Decerver configuration
	ws.webServer.Get("/admin/switch/(.*)", das.handleDappSwitch)

	go func() {
		ws.webServer.RunOnAddr("localhost:" + fmt.Sprintf("%d", ws.port))
	}()

	return nil
}
