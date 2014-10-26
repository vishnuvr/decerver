package server

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"strings"
)

var null = json.RawMessage([]byte("null"))

// A JSON-RPC message received by the server.
type RPCRequest struct {
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// An Array of objects to pass as arguments to the method.
	Params *json.RawMessage `json:"params"`
	// Timestamp
	Timestamp int `json:"timestamp"`
}

// A JSON-RPC message sent by the server.
type RPCResponse struct {
	// The name of the Object that was returned by the invoke method.
	Id string `json:"id"`
	// The Object that was returned by the invoked method. This must be null
	// in case there was an error invoking the method.
	Result interface{} `json:"result"`
	// An Error object if there was an error invoking the method. It must be
	// null if there was no error.
	Error interface{} `json:"error"`
	// Timestamp
	Timestamp int `json:"timestamp"`
}

type ErrorCode int

const (
	E_PARSE       ErrorCode = -32700
	E_INVALID_REQ ErrorCode = -32600
	E_NO_METHOD   ErrorCode = -32601
	E_BAD_PARAMS  ErrorCode = -32602
	E_INTERNAL    ErrorCode = -32603
	E_SERVER      ErrorCode = -32000
)

type Error struct {
	// A Number that indicates the error type that occurred.
	Code ErrorCode `json:"code"` /* required */

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"` /* required */

	// A Primitive or Structured value that contains additional information about the error.
	Data interface{} `json:"data"` /* optional */
}

func (e *Error) Error() string {
	return e.Message
}

func getErrorResponse(err *Error) *RPCResponse {
	return &RPCResponse{Error: err}
}

// The socket RPC server

type SRPCService interface {
	Name() string
	Init()
	HandleRPC(*RPCRequest) (*RPCResponse, error)
	SetConnection(wsConn *WsConn)
	Close()
}

type SRPCServiceFactory interface {
	CreateService() SRPCService
	Init()
	ServiceName() string
}

// Function that handles a specific type of posting
type RPCMethod func(*RPCRequest, *RPCResponse)

type SRPCServer struct {
	activeConnections uint32
	maxConnections    uint32
	idPool            *IdPool
	serviceFactories  map[string]SRPCServiceFactory
	activeHandlers    map[uint32]*SessionHandler
}

func NewSRPCServer(maxConnections uint32) *SRPCServer {
	srv := &SRPCServer{}
	srv.serviceFactories = make(map[string]SRPCServiceFactory)
	srv.activeHandlers = make(map[uint32]*SessionHandler)
	srv.maxConnections = maxConnections
	srv.idPool = NewIdPool(maxConnections)
	return srv
}

func (srv *SRPCServer) CurrentActiveConnections() uint32 {
	return srv.activeConnections
}

func (srv *SRPCServer) MaxConnections() uint32 {
	return srv.maxConnections
}

func (srv *SRPCServer) RemoveSessionHandler(sh *SessionHandler) {
	srv.activeConnections--
	srv.idPool.ReturnId(sh.wsConn.SessionId)
	delete(srv.activeHandlers, sh.wsConn.SessionId)
}

func (srv *SRPCServer) CreateSessionHandler(wsConn *WsConn) *SessionHandler {
	sh := &SessionHandler{}
	sh.wsConn = wsConn
	sh.services = make(map[string]SRPCService)
	for _, v := range srv.serviceFactories {
		sc := v.CreateService()
		sc.SetConnection(wsConn)
		sc.Init()
		sh.services[v.ServiceName()] = sc
	}
	sh.server = srv;
	srv.activeConnections++
	id := srv.idPool.GetId()
	sh.wsConn.SessionId = id;
	srv.activeHandlers[id] = sh
	fmt.Printf("ACTIVE CONNECTIONS: %v\n",srv.activeHandlers);
	return sh
}

func (srv *SRPCServer) RegisterServiceFactory(factory SRPCServiceFactory, factoryName string) {
	sname := factoryName
	// If name is "", get the name from the object type.
	if factoryName == "" {
		sname = reflect.Indirect(reflect.ValueOf(factory)).Type().Name()
	}

	fmt.Println("SERVICE FACTORY ADDED: " + sname)
	var ok bool
	srv.serviceFactories[sname] = factory
	factory.Init()

	if !ok {
		return
	}
}

// TODO do this properly.
func (srv *SRPCServer) DeregisterServiceFactory(serviceFactoryName string) {
	delete(srv.serviceFactories, serviceFactoryName)
}

// This is passed to the Martini server. We only allow 1 connection at this time.
func (srs *SRPCServer) handleWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New connection.")
	if srs.activeConnections == srs.maxConnections {
		fmt.Println("Connection failed: Already at capacity.")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Error(err)
		return
	}

	wsConn := &WsConn{conn: conn,
		WriteMsgChannel:   make(chan *Message, 256),
		writeCloseChannel: make(chan *Message, 256),
	}
	sh := srs.CreateSessionHandler(wsConn)
	go writer(sh)
	reader(sh)
	sh.wsConn.WriteMsgChannel <- &Message{Data : nil}
	sh.Close();
}

type SessionHandler struct {
	server    *SRPCServer
	services  map[string]SRPCService
	wsConn    *WsConn
}

func (sh *SessionHandler) Close() {
	fmt.Printf("CLOSING HANDLER: %d\n", sh.wsConn.SessionId)
	// TODO run Close() on all services (and add that command).
	for _ , srvc := range sh.services {
		srvc.Close()
	}
	sh.services = nil
	// Deregister ourselves.
	sh.server.RemoveSessionHandler(sh)
}

func (sh *SessionHandler) handleRequest(rpcReq *RPCRequest) {

	mtd := rpcReq.Method
	if mtd == "" {
		err := &Error{
			Code:    E_NO_METHOD,
			Message: "Method name is empty.",
			Data:    rpcReq,
		}
		sh.wsConn.WriteMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	mtdfrm := strings.Split(mtd, ".")
	if len(mtdfrm) != 2 {
		err := &Error{
			Code:    E_NO_METHOD,
			Message: "Method name malformed. Need to be on the form 'Service.Method': " + mtd,
			Data:    rpcReq,
		}
		sh.wsConn.WriteMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	serviceName := mtdfrm[0]

	if serviceName == "" || sh.services[serviceName] == nil {
		err := &Error{
			Code:    E_NO_METHOD,
			Message: "No service with name: " + serviceName,
			Data:    rpcReq,
		}
		sh.wsConn.WriteMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	methodName := mtdfrm[1]
	rpcReq.Method = methodName // Replace now that we know which service to use.
	service := sh.services[serviceName]

	// TODO add a way to divide responses up into multiple messages without
	// passing the connection object.
	rpcResp, handleErr := service.HandleRPC(rpcReq)

	if handleErr != nil {
		err := &Error{
			Code:    E_NO_METHOD,
			Message: "No method with name: " + methodName,
			Data:    rpcReq,
		}
		sh.wsConn.WriteMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	if rpcResp.Result != nil {
		// If there is a return value, pass to the write channel.
		sh.wsConn.WriteMsgChannel <- &Message{Data: rpcResp, Type: websocket.TextMessage}
	}
}
