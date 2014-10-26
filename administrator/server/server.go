package server

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type WebServer struct {
	Martini              *martini.ClassicMartini
	maxConnections       uint32
	rpcServices          []interface{}
	sRpcServiceFactories []SRPCServiceFactory
}

func NewWebServer(maxConnections uint32) *WebServer {
	ws := &WebServer{}
	ws.maxConnections = maxConnections
	ws.rpcServices = make([]interface{}, 0)
	ws.sRpcServiceFactories = make([]SRPCServiceFactory, 0)
	return ws
}

func (ws *WebServer) RegisterRPCServices(services ...interface{}) {
	ws.rpcServices = services
}

func (ws *WebServer) RegisterSocketRPCServices(factories ...SRPCServiceFactory) {
	ws.sRpcServiceFactories = factories
}

func (ws *WebServer) Start() {

	ws.Martini = martini.Classic()
	ws.Martini.Use(martini.Static("./web"))
	
	//martini.Env = martini.Prod

	// JSON RPC
	if len(ws.rpcServices) > 0 {
		rpcs := rpc.NewServer()
		rpcs.RegisterCodec(json2.NewCodec(), "application/json")
		for _, service := range ws.rpcServices {
			rpcs.RegisterService(service, "")
		}
		ws.Martini.Post("/status", rpcs.ServeHTTP)
	}

	// JSON Socket RPC
	if len(ws.sRpcServiceFactories) > 0 {
		srpcs := NewSRPCServer(ws.maxConnections)
		for _, factory := range ws.sRpcServiceFactories {
			srpcs.RegisterServiceFactory(factory, "")
		}
		ws.Martini.Get("/srpc", srpcs.handleWs)
	}

	// Graceful shutdown
	gracefulShutdown := &GracefulShutdown{timeout: time.Duration(2) * time.Second}
	ws.Martini.Use(gracefulShutdown.Handler)

	go func() { ws.Martini.RunOnAddr("localhost:3000") }()

	// We just wait for a signal to close down the server.
	err := gracefulShutdown.WaitForSignal(syscall.SIGTERM, syscall.SIGINT)
	if err != nil {
		log.Println(err)
	}
}

type GracefulShutdown struct {
	timeout time.Duration
	wg      sync.WaitGroup
}

func NewGracefulShutdown(t time.Duration) *GracefulShutdown {
	return &GracefulShutdown{timeout: t}
}

func (g *GracefulShutdown) Handler(c martini.Context) {
	g.wg.Add(1)
	c.Next()
	g.wg.Done()
}

func (g *GracefulShutdown) WaitForSignal(signals ...os.Signal) error {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, signals...)
	<-sigchan

	log.Println("Waiting for all requests to finish")

	waitChan := make(chan struct{})
	go func() {
		g.wg.Wait()
		waitChan <- struct{}{}
	}()

	select {
	case <-time.After(g.timeout):
		return fmt.Errorf("timed out waiting %v for shutdown", g.timeout)
	case <-waitChan:
		return nil
	}
}
