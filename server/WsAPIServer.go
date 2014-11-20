package server

import (
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/api"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/util"
	"github.com/gorilla/websocket"
	"net/http"
	"path"
)

func getErrorResponse(err *api.Error) []byte {
	rsp := &api.Response{}
	rsp.Error = err
	bts, _ := json.Marshal(rsp)
	return bts
}

// The websocket server handles connections.
type WsAPIServer struct {
	ate               core.RuntimeManager
	activeConnections uint32
	maxConnections    uint32
	idPool            *util.IdPool
	sessions          map[uint32]*Session
}

func NewWsAPIServer(ate core.RuntimeManager, maxConnections uint32) *WsAPIServer {
	srv := &WsAPIServer{}
	srv.sessions = make(map[uint32]*Session)
	srv.maxConnections = maxConnections
	srv.idPool = util.NewIdPool(maxConnections)
	srv.ate = ate
	return srv
}

func (srv *WsAPIServer) CurrentActiveConnections() uint32 {
	return srv.activeConnections
}

func (srv *WsAPIServer) MaxConnections() uint32 {
	return srv.maxConnections
}

func (srv *WsAPIServer) RemoveSession(ss *Session) {
	srv.activeConnections--
	srv.idPool.ReleaseId(ss.wsConn.SessionId())
	delete(srv.sessions, ss.wsConn.SessionId())
}

func (srv *WsAPIServer) CreateSession(caller string, rt core.Runtime, wsConn *WsConn) *Session {
	fmt.Printf("Runtime: %v\n", rt)
	ss := &Session{}
	ss.wsConn = wsConn
	ss.server = srv
	ss.caller = caller
	ss.runtime = rt
	srv.activeConnections++
	id := srv.idPool.GetId()
	ss.wsConn.sessionId = id
	srv.sessions[id] = ss
	fmt.Printf("ACTIVE CONNECTIONS: %v\n", srv.sessions)
	return ss
}

// This is passed to the Martini server.
// Find out what endpoint they called and create a session based on that.
func (srv *WsAPIServer) handleWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New connection.")
	if srv.activeConnections == srv.maxConnections {
		fmt.Println("Connection failed: Already at capacity.")
	}
	u := r.URL
	p := u.Path
	caller := path.Base(p)
	fmt.Println("Caller: " + caller)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade to websocket (%s)\n", err.Error())
		return
	}
	wsConn := &WsConn{
		conn:              conn,
		writeMsgChannel:   make(chan *Message, 256),
		writeCloseChannel: make(chan *Message, 256),
	}
	rt := srv.ate.GetRuntime(caller)
	ss := srv.CreateSession(caller, rt, wsConn)
	// We add this session to the callers (dapps) runtime.
	err = rt.BindScriptObject("tempObj", NewSessionJs(ss))
	
	if err != nil {
		panic(err.Error())
	}
	
	// TODO expose toValue in runtime?
	rt.AddScript("network.newWsSession(tempObj); tempObj = null;");
	//rt.CallFuncOnObj("network", "newWsSession", val)
	go writer(ss)
	reader(ss)
	ss.wsConn.writeMsgChannel <- &Message{Data: nil}
	ss.Close()
}

type Session struct {
	caller    string
	runtime   core.Runtime
	server    *WsAPIServer
	wsConn    *WsConn
	sessionJs *SessionJs
}

func (ss *Session) SessionId() uint32 {
	return ss.wsConn.sessionId
}

func (ss *Session) WriteCloseMsg() {
	ss.wsConn.WriteCloseMsg()
}

func (ss *Session) Close() {
	fmt.Printf("CLOSING SESSION: %d\n", ss.wsConn.sessionId)
	// Deregister ourselves.
	ss.server.RemoveSession(ss)
	if ss.wsConn.conn != nil {
		err := ss.wsConn.conn.Close()
		if err != nil {
			fmt.Printf("Failed to close websocket connection, already removed: %d\n", ss.wsConn.sessionId)
		}
	}
}

func (ss *Session) handleRequest(rpcReq string) {
	fmt.Println("RPC Message: " + rpcReq)
	ret, err := ss.runtime.CallFuncOnObj("network", "incomingWsMsg", int(ss.wsConn.sessionId), rpcReq)
	
	if err != nil {
		err := &api.Error{
			Code:    api.E_SERVER,
			Message: "Js runtime error: " + ss.caller,
			Data:    rpcReq,
		}
		ss.wsConn.writeMsgChannel <- &Message{Data: getErrorResponse(err), Type: websocket.TextMessage}
		return
	}
	if ret == nil {
		return
	}
	retStr := ret.(string)
	// If there is a return value, pass to the write channel.
	ss.wsConn.writeMsgChannel <- &Message{Data: []byte(retStr), Type: websocket.TextMessage}
}

type SessionJs struct {
	session *Session
}

func NewSessionJs(ss *Session) *SessionJs {
	return &SessionJs{ss}
}

func (sjs *SessionJs) WriteJson(msg string) {
	sjs.session.wsConn.WriteJsonMsg([]byte(msg))
}

func (sjs *SessionJs) SessionId() int {
	return int(sjs.session.SessionId())
}