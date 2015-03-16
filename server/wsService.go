package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"path"
	"sync"
	"encoding/json"
)

// The websocket server handles connections.
type WsService struct {
	maxConnections    uint32
	idPool            *IdPool
	sessions          map[uint32]*Session
	sMutex            *sync.Mutex
}

func NewWsService(maxConnections uint32) *WsService {
	srv := &WsService{}
	srv.sessions = make(map[uint32]*Session)
	srv.maxConnections = maxConnections
	srv.idPool = NewIdPool(maxConnections)
	return srv
}

func (this *WsAPIServer) CurrentActiveConnections() uint32 {
	return len(this.sessions)
}

func (this *WsAPIServer) MaxConnections() uint32 {
	return this.maxConnections
}

// This is passed to the Martini server to handle websocket requests.
func (this *WsService) handleWs(w http.ResponseWriter, r *http.Request) {
	
	// TODO check scheme first.	
	logger.Println("New websocket connection.")
	
	if len(this.sessions) == this.maxConnections {
		logger.Println("Connection failed: Already at capacity.")
	}
	u := r.URL
	p := u.Path
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Printf("Failed to upgrade to websocket (%s)\n", err.Error())
		return
	}

	ss := this.createSession(conn)
	
	go writer(ss)
	reader(ss)
	ss.writeMsgChannel <- &Message{Data: nil}
	ss.Close()
}

// Only called by the 'handleWs' function.
func (this *WsService) newSession(conn *websocket.Conn) (newSession *Session, err error) {
	newSession := &Session{}
	ss.conn = conn
	ss.server = this
	id := srv.idPool.GetId()
	ss.wsConn.sessionId = id
	srv.sessions[id] = ss
	ss.writeMsgChannel =   make(chan *Message, 256),
	ss.writeCloseChannel: make(chan *Message, 256),
	return ss
}

func (this *WsService) deleteSession(sessionId uint32) {
	if this.sessions[sessionId] == nil {
		logger.Printf("Attempted to remove a session that does not exist (id: %d).",sessionId)
		return;
	}
	delete(this.sessions, sessionId)
	this.idPool.ReleaseId(sessionId)
}

type Session struct {
	conn              *websocket.Conn
	server            *WsService
	writeMsgChannel   chan *Message
	writeCloseChannel chan *Message
	sessionId         uint32
}

func (ss *Session) SessionId() uint32 {
	return ss.sessionId
}

func (ss *Session) WriteJson(obj interface{}) {
	msg, err := json.Marshal(obj)
	if err != nil {
		// TODO Protocol stuff.	
		ss.WriteCloseMsg()
	} else {
		ss.writeMsgChannel <- &Message{Data: msg, Type: websocket.TextMessage}
	}
}

// We don't call ss.conn.SendMessage right away as that will not register
// with the writer. Instead we use the close channel.
func (ss *Session) WriteCloseMsg() {
	ss.writeCloseChannel <- &Message{Data: nil, Type: websocket.CloseMessage}
}

func (ss *Session) Close() {
	logger.Printf("CLOSING SESSION: %d\n", ss.wsConn.sessionId)
	// Deregister ourselves.
	ss.server.RemoveSession(ss)
	ss.runtime.CallFuncOnObj("network", "deleteWsSession", int(ss.SessionId()))
	if ss.wsConn.conn != nil {
		err := ss.wsConn.conn.Close()
		if err != nil {
			logger.Printf("Failed to close websocket connection, already removed: %d\n", ss.wsConn.sessionId)
		}
	}
}

func (ss *Session) handleRequest(rpcReq string) {
	logger.Println("RPC Message: " + rpcReq)
	
	// TODO the protocol stuff.
	
	// If there is a return value, pass to the write channel.
	ss.wsConn.writeMsgChannel <- &Message{Data: []byte(retStr), Type: websocket.TextMessage}
}