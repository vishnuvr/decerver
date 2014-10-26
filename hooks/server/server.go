package server

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"github.com/golang/glog"
)

// Interface that is used for webhook handlers
type WebHookHandler interface {
	Handle(res http.ResponseWriter, req *http.Request)
}

type WebServer struct {
	Martini     *martini.ClassicMartini
}

func NewWebServer() *WebServer {
	return &WebServer{}
}

// Add the provided handler then start handling posts on /postreceive
func (ws *WebServer) Start(handler interface{}) {

	hdlr := handler.(WebHookHandler)
	ws.Martini = martini.Classic()
	ws.Martini.Post("/postreceive", hdlr.Handle)
	
	// Graceful shutdown
	gracefulShutdown := &GracefulShutdown{timeout: time.Duration(2) * time.Second}
	ws.Martini.Use(gracefulShutdown.Handler)

	go func() { ws.Martini.Run() }()

	// Run until someone shuts the server down from the console.
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

	glog.Info("Waiting for all requests to finish")

	waitChan := make(chan struct{})
	go func() {
		g.wg.Wait()
		waitChan <- struct{}{}
	}()

	select {
	case <-time.After(g.timeout):
		glog.Errorf("timed out waiting %v for shutdown", g.timeout)
		return nil 
	case <-waitChan:
		return nil
	}
}