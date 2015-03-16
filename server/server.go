package server

import (
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout,"[Server] ", log.LstdFlags)

type Server struct {
	maxConnections uint32
	host		   string
	port           uint16
	rootDir        string
	cMartini       *martini.ClassicMartini
	httpService    *HttpServer
	wsService      *WsServer
}

func NewServer(host string, port uint16, maxConnections uint32, rootDir string) *Server {
	
	cMartini := martini.Classic()
	
	// TODO remember to change to martini.Prod
	cMartini.Env = martini.Dev
	
	httpService := NewHttpService()
	wsService := NewWsService(maxConnections)
	
	return &Server{
			maxConnections,
			host,
			port,
			rootDir,
			cMartini,
			httpService,
			wsService,
		}

	return ws
}

func (this *WebServer) Start() error {

	this.cMartini.Use(martini.Static(rootDir))

	m.NotFound(this.httpService.handleNotFound)
	
	// TODO Close down properly. Removed the third party stuff since 
	// it was a mess.
	go func() {
		ws.webServer.RunOnAddr(ws.host + ":" + fmt.Sprintf("%d", ws.port))
	}()
	
	
	
	return nil
}
