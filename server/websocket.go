// This file contains code for handling websocket specific stuff, such as
// message types, connection objects and channels, and settings. It is the
// bridge between the server and the SRPC handling code.
package server

import (
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/deCerver-interfaces/api"
	"github.com/gorilla/websocket"
	"time"
	//"bytes"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8192
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
}

// Base message type we pass to writer. Text, ping and close.
type Message struct {
	Data interface{}
	Type int
}

// A connection's writer can only be used by one process at a time.
// To avoid any problems, no external process is allowed access
// to the websocket connection object itself. All they can do is pass
// messages via the WriteMsgChannel that the WsConn 'middle-man' provides.
// Note that ping and close is not public, because no external process
// should ever use it, but those messages still conflict with text
// messages and must therefore be passed to the write-routine in the
// same manner.
//
// All text messages must be json formatted strings.

func GetPingMessage() *Message {
	return &Message{Type: websocket.PingMessage}
}

func GetCloseMessage() *Message {
	return &Message{Type: websocket.CloseMessage}
}

func GetJsonMessage(data interface{}) *Message {
	return &Message{Data: data, Type: websocket.TextMessage}
}

type WsConn struct {
	sessionId         uint32
	conn              *websocket.Conn
	writeMsgChannel   chan *Message
	writeCloseChannel chan *Message
}

func (wc *WsConn) SessionId() uint32 {
	return wc.sessionId
}

func (wc *WsConn) Connection() *websocket.Conn {
	return wc.conn
}

func (wc *WsConn) WriteTextMsg(msg interface{}) {
	wc.writeMsgChannel <- &Message{Data: msg, Type: websocket.TextMessage}
}

func (wc *WsConn) WriteCloseMsg() {
	wc.writeCloseChannel <- &Message{Data: "", Type: websocket.CloseMessage}
}

// Handle the reader
func reader(sh *SessionHandler) {
	conn := sh.wsConn.conn
	defer func() {
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	//wsc.conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetReadDeadline(time.Time{})
	//wsc.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// TODO add back the reader timeout.
	for {
		fmt.Println("Waiting to read socket.")

		mType, message, err := conn.NextReader()

		if err != nil {
			sh.wsConn.writeCloseChannel <- GetCloseMessage()
			return
		}

		if mType == websocket.TextMessage {
			rpcReq := &api.Request{}
			umErr := json.NewDecoder(message).Decode(rpcReq)
			if umErr == nil {
				sh.handleRequest(rpcReq)
			} else {
				fmt.Println("Failed to unmarshal message from client.")
				sh.wsConn.writeCloseChannel <- GetCloseMessage()
				return
			}
		} else if mType == websocket.CloseMessage {
			return
		}

	}
}

// Handle the writer
func writer(sh *SessionHandler) {
	conn := sh.wsConn.conn
	defer func() {
		conn.Close()
	}()
	fmt.Println("Waiting to write to socket.")
	for {
		message, ok := <-sh.wsConn.writeMsgChannel

		if !ok {
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if message.Data == nil {
			return
		}

		if message.Type == websocket.CloseMessage {
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := conn.WriteJSON(message.Data); err != nil {
			return
		}

	}
}
