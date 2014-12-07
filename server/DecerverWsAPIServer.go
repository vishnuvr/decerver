package server

/*
import (
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/api"
	"github.com/eris-ltd/decerver-interfaces/util"
	"github.com/gorilla/websocket"
	"github.com/eris-ltd/decerver-interfaces/core"
	"net/http"
	"path"
)

// The websocket server handles connections.
type DecerverWsAPIServer struct {
	decerver          core.DeCerver
	activeConnections uint32
	maxConnections    uint32
	idPool            *util.IdPool
	sessions          map[uint32]*Session
}

func NewDecerverWsAPIServer(dc core.DeCerver) *DecerverWsAPIServer {
	srv := &DecerverWsAPIServer{}
	srv.sessions = make(map[uint32]*Session)
	srv.maxConnections = 10 // TODO ?
	srv.idPool = util.NewIdPool(srv.maxConnections)
	srv.decerver = dc
	return srv
}

func (ds *DecerverWsAPIServer) CurrentActiveConnections() uint32 {
	return ds.activeConnections
}

func (ds *DecerverWsAPIServer) MaxConnections() uint32 {
	return ds.maxConnections
}

func (ds *DecerverWsAPIServer) RemoveSession(ss *Session) {
	ds.activeConnections--
	ds.idPool.ReleaseId(ss.wsConn.SessionId())
	delete(ds.sessions, ss.wsConn.SessionId())
}

func (ds *DecerverWsAPIServer) CreateSession(caller string, dc core.DeCerver, wsConn *WsConn) *Session {
	ss := &Session{}
	ss.wsConn = wsConn
	ss.server = ds
	ss.caller = caller
	ds.activeConnections++
	id := srv.idPool.GetId()
	ss.wsConn.sessionId = id
	ds.sessions[id] = ss
	fmt.Printf("ACTIVE DECERVER CONNECTIONS: %v\n", ds.sessions)
	return ss
}

// This is passed to the Martini server.
// Find out what endpoint they called and create a session based on that.
func (srs *WsAPIServer) handleWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New connection.")
	if srs.activeConnections == srs.maxConnections {
		fmt.Println("Connection failed: Already at capacity.")
	}

	u := r.URL
	p := u.Path
	caller := path.Base(p)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade to websocket (%s)\n", err.Error())
		return
	}

	wsConn := &WsConn{
		conn: conn,
		writeMsgChannel:   make(chan *Message, 256),
		writeCloseChannel: make(chan *Message, 256),
	}

	ss := srs.CreateSession(caller, srs.ate.GetRuntime(caller), wsConn)
	go writer(ss)
	reader(ss)
	ss.wsConn.writeMsgChannel <- &Message{Data: nil}
	ss.Close()
}

type Session struct {
	caller   string
	runtime  core.Runtime
	server   *WsAPIServer
	wsConn   *WsConn
}

func (ss *Session) SessionId() uint32 {
	return ss.wsConn.sessionId
}

func (ss *Session) WriteJsonMsg(obj interface{}) {
	ss.wsConn.WriteJsonMsg(obj)
}

func (ss *Session) WriteCloseMsg() {
	ss.wsConn.WriteCloseMsg()
}

func (ss *Session) Close() {
	fmt.Printf("CLOSING SESSION: %d\n", ss.wsConn.SessionId)
	// Deregister ourselves.
	ss.server.RemoveSession(ss)
	if ss.wsConn.conn != nil {
		err := ss.wsConn.conn.Close()
		if err != nil {
			fmt.Printf("Failed to close websocket connection, already removed: %d\n", ss.wsConn.sessionId)
		}
	}
}

func (ss *Session) handleRequest(rpcReq map[string]interface{}) {

	mtd := rpcReq["Method"]
	if mtd == "" {
		err := &api.Error{
			Code:    api.E_NO_METHOD,
			Message: "Method name is empty.",
			Data:    rpcReq,
		}
		ss.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}
	
	ss.runtime.CallFuncOnObj("network","incomingWsMsg",rpcReq)

	
	if serviceName == "" || ss.services[serviceName] == nil {
		err := &api.Error{
			Code:    api.E_NO_METHOD,
			Message: "No service with name: " + serviceName,
			Data:    rpcReq,
		}
		ss.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}

	if handleErr != nil {
		err := &api.Error{
			Code:    api.E_NO_METHOD,
			Message: "No method with name: " + methodName,
			Data:    rpcReq,
		}
		ss.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}
	
	if rpcResp.Result != nil {
		// If there is a return value, pass to the write channel.
		ss.wsConn.writeMsgChannel <- &Message{Data: rpcResp, Type: websocket.TextMessage}
	}
	
}
*/