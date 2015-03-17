package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

const PARSE_ERROR = -32700
const INVALID_REQUEST = -32600
const METHOD_NOT_FOUND = -32601
const INVALID_PARAMS = -32602
const INTERNAL_ERROR = -32603

type JsonRpcHandler func(*Request, *Response)

// The websocket server handles connections.
// NOTE currently all sessions use the same handlers, since you
// can't run multiple EPMs (or you can, but there's no guarantee
// that they don't try and compete for the same database).
type WsService struct {
	maxConnections uint32
	idPool         *IdPool
	sessions       map[uint32]*Session
	handlers       map[string]JsonRpcHandler
}

func NewWsService(maxConnections uint32) *WsService {
	srv := &WsService{}
	srv.sessions = make(map[uint32]*Session)
	srv.maxConnections = maxConnections
	srv.idPool = NewIdPool(maxConnections)
	srv.handlers = make(map[string]JsonRpcHandler)
	
	srv.handlers["echo"] = srv.echo
	return srv
}

/***************************** Handlers ********************************/

func (this *WsService) echo(req *Request, resp *Response) {
	sVal := &StringValue{}
	err := json.Unmarshal([]byte(*req.Params), &sVal)
	if err != nil {
		resp.Error = Error(INVALID_PARAMS, "Echo requires a string parameter.")
	}
	logger.Printf("Echo: %s", sVal.SVal)
	resp.Result = &StringValue{sVal.SVal}
}

/***********************************************************************/

func (this *WsService) CurrentActiveConnections() uint32 {
	return uint32(len(this.sessions))
}

func (this *WsService) MaxConnections() uint32 {
	return this.maxConnections
}

// This is passed to the Martini server to handle websocket requests.
func (this *WsService) handleWs(w http.ResponseWriter, r *http.Request) {

	// TODO check scheme first.
	logger.Println("New websocket connection.")

	if uint32(len(this.sessions)) == this.maxConnections {
		logger.Println("Connection failed: Already at capacity.")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Printf("Failed to upgrade to websocket (%s)\n", err.Error())
		return
	}

	ss := this.newSession(conn)

	go writer(ss)
	reader(ss)
	ss.writeMsgChannel <- &Message{Data: nil}
	ss.closeSession()
}

// Only called by the 'handleWs' function.
func (this *WsService) newSession(conn *websocket.Conn) *Session {
	ss := &Session{}
	ss.conn = conn
	ss.server = this
	id := this.idPool.GetId()
	ss.sessionId = id
	ss.writeMsgChannel = make(chan *Message, 256)
	ss.writeCloseChannel = make(chan *Message, 256)

	this.sessions[id] = ss

	return ss
}

func (this *WsService) deleteSession(sessionId uint32) {
	if this.sessions[sessionId] == nil {
		logger.Printf("Attempted to remove a session that does not exist (id: %d).", sessionId)
		return
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

func (ss *Session) closeSession() {
	// Deregister ourselves.
	ss.server.deleteSession(ss.sessionId)
	logger.Printf("Closing session: %d", ss.sessionId)
	logger.Printf("Connections remaining: %d", len(ss.server.sessions))
	if ss.conn != nil {
		err := ss.conn.Close()
		if err != nil {
			logger.Printf("Failed to close websocket connection, already removed: %d\n", ss.sessionId)
		}
	}
}

func (ss *Session) handleRequest(msg []byte) {

	req := &Request{}

	err := json.Unmarshal(msg, req)

	var resp *Response

	if err != nil {
		// Can't really say if it's bad json or not a proper request
		// without looking at the error.
		resp = ErrorResp(INVALID_REQUEST, err.Error())
	} else {
		handler, hExists := ss.server.handlers[req.Method]
		if !hExists {
			resp = ErrorResp(-32601, "Method not found: "+req.Method)
		} else {
			resp = &Response{}
			resp.ID = req.ID
			resp.JsonRpc = "2.0"
			handler(req, resp)
		}
	}

	ss.WriteJson(resp)
}

// Get an Error Object
func Error(code int, err string) *ErrorObject {
	return &ErrorObject{code, err}
}

// Get an error response with the fields already filled out.
func ErrorResp(code int, err string) *Response {
	return &Response{
		-1,
		"2.0",
		nil,
		&ErrorObject{code, err},
	}
}

type( 
	Request struct {
		ID      interface{}     `json:"id"`
		JsonRpc string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  *json.RawMessage `json:"params"`
	}
	
	Response struct {
		ID      interface{}  `json:"id"`
		JsonRpc string       `json:"jsonrpc"`
		Result  interface{}  `json:"result"`
		Error   *ErrorObject `json:"error"`
	}
	
	ErrorObject struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

type StringValue struct {
	SVal string
}
