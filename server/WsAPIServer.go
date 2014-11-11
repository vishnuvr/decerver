package server

import (
	"fmt"
	"github.com/eris-ltd/deCerver-interfaces/api"
	"github.com/eris-ltd/deCerver-interfaces/util"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

func getErrorResponse(err *api.Error) *api.Response {
	rsp := &api.Response{}
	rsp.Error = err
	return rsp
}

// The websocket RPC server
type WsAPIServer struct {
	activeConnections uint32
	maxConnections    uint32
	idPool            *util.IdPool
	serviceFactories  map[string]api.WsAPIServiceFactory
	activeHandlers    map[uint32]*SessionHandler
}

func NewWsAPIServer(maxConnections uint32) *WsAPIServer {
	srv := &WsAPIServer{}
	srv.serviceFactories = make(map[string]api.WsAPIServiceFactory)
	srv.activeHandlers = make(map[uint32]*SessionHandler)
	srv.maxConnections = maxConnections
	srv.idPool = util.NewIdPool(maxConnections)
	return srv
}

func (srv *WsAPIServer) CurrentActiveConnections() uint32 {
	return srv.activeConnections
}

func (srv *WsAPIServer) MaxConnections() uint32 {
	return srv.maxConnections
}

func (srv *WsAPIServer) RemoveSessionHandler(sh *SessionHandler) {
	srv.activeConnections--
	srv.idPool.ReleaseId(sh.wsConn.SessionId())
	delete(srv.activeHandlers, sh.wsConn.SessionId())
}

func (srv *WsAPIServer) CreateSessionHandler(wsConn *WsConn) *SessionHandler {
	sh := &SessionHandler{}
	sh.wsConn = wsConn
	
	sh.server = srv
	srv.activeConnections++
	id := srv.idPool.GetId()
	sh.wsConn.sessionId = id
	srv.activeHandlers[id] = sh
	fmt.Printf("ACTIVE CONNECTIONS: %v\n", srv.activeHandlers)
	
	sh.services = make(map[string]api.WsAPIService)
	
	for _, v := range srv.serviceFactories {
		
		sc := v.CreateService()
		sc.SetConnection(wsConn)
		sc.Init()
		sh.services[v.ServiceName()] = sc
		fmt.Printf("Adding service factory '%s' to session handler.\n")
	}
	
	return sh
}

func (srv *WsAPIServer) RegisterServiceFactory(factory api.WsAPIServiceFactory) {
	
	srv.serviceFactories[factory.ServiceName()] = factory
	factory.Init()
}

// TODO do this properly.
func (srv *WsAPIServer) DeregisterServiceFactory(serviceFactoryName string) {
	delete(srv.serviceFactories, serviceFactoryName)
}

// This is passed to the Martini server.
func (srs *WsAPIServer) handleWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New connection.")
	if srs.activeConnections == srs.maxConnections {
		fmt.Println("Connection failed: Already at capacity.")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade to websockets (%s)\n", err.Error())
		return
	}

	wsConn := &WsConn{conn: conn,
		writeMsgChannel:   make(chan *Message, 256),
		writeCloseChannel: make(chan *Message, 256),
	}
	sh := srs.CreateSessionHandler(wsConn)
	go writer(sh)
	reader(sh)
	sh.wsConn.writeMsgChannel <- &Message{Data: nil}
	sh.Close()
}

type SessionHandler struct {
	server   *WsAPIServer
	services map[string]api.WsAPIService
	wsConn   *WsConn
}

func (sh *SessionHandler) Close() {
	fmt.Printf("CLOSING HANDLER: %d\n", sh.wsConn.SessionId)
	
	for _, srvc := range sh.services {
		srvc.Shutdown()
	}	
	sh.services = nil
	// Deregister ourselves.
	sh.server.RemoveSessionHandler(sh)
	if sh.wsConn.conn != nil {
		err := sh.wsConn.conn.Close()
		if err != nil {
			fmt.Printf("Failed to close websocket connection, already removed: %d\n", sh.wsConn.sessionId)
		}
	}
}

func (sh *SessionHandler) handleRequest(rpcReq *api.Request) {

	mtd := rpcReq.Method
	if mtd == "" {
		err := &api.Error{
			Code:    api.E_NO_METHOD,
			Message: "Method name is empty.",
			Data:    rpcReq,
		}
		sh.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	mtdfrm := strings.Split(mtd, ".")
	if len(mtdfrm) != 2 {
		err := &api.Error{
			Code:    api.E_NO_METHOD,
			Message: "Method name malformed. Need to be on the form 'Service.Method': " + mtd,
			Data:    rpcReq,
		}
		sh.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	serviceName := mtdfrm[0]

	if serviceName == "" || sh.services[serviceName] == nil {
		err := &api.Error{
			Code:    api.E_NO_METHOD,
			Message: "No service with name: " + serviceName,
			Data:    rpcReq,
		}
		sh.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	methodName := mtdfrm[1]
	rpcReq.Method = methodName // Replace now that we know which service to use.
	service := sh.services[serviceName]

	// TODO add a way to divide responses up into multiple messages without
	// passing the connection object.
	rpcResp, handleErr := service.HandleRPC(rpcReq)

	if handleErr != nil {
		err := &api.Error{
			Code:    api.E_NO_METHOD,
			Message: "No method with name: " + methodName,
			Data:    rpcReq,
		}
		sh.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}
	if rpcResp.Result != nil {
		// If there is a return value, pass to the write channel.
		sh.wsConn.writeMsgChannel <- &Message{Data: rpcResp, Type: websocket.TextMessage}
	}
}
