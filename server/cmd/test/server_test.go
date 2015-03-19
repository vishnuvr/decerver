package server

import (
	"fmt"
	"bytes"
	"encoding/json"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"
	"github.com/eris-ltd/decerver/server"
)

const TEST_NUM = 10

var srvr *server.Server

func init() {
	rootPath, _ := filepath.Abs("/public")
	srvr = server.NewServer("localhost", 3000, TEST_NUM, rootPath)
	go func(){
		srvr.Start()
	}()
	time.Sleep(1 * time.Second)
}

// Test sending a http request to the echo endpoint
func TestHttpEcho(t *testing.T) {
	fmt.Println("Start http echo test.")
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:3000/echo/testmessage", bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}
	resp, err2 := client.Do(req)

	if err2 != nil {
		panic(err2)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	retStr := string(body)
	if retStr != "testmessage" {
		t.Error("Expected: testmessage, Got: " + retStr)
	} else {
		fmt.Println("Http echo test: PASSED")
	}
}

// Establish websocket connection and rpc to 'echo'
func TestWsEcho(t *testing.T) {
	fmt.Println("Start websocket echo test.")
	origin := "http://localhost/"
	url := "ws://localhost:3000/websocket"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	req := &server.Request{}
	req.ID = 1
	req.JsonRpc = "2.0"
	req.Method = "echo"

	sVal := &server.StringValue{"testmessage"}
	bts, _ := json.Marshal(sVal)
	raw := json.RawMessage(bts)
	req.Params = &raw

	bts, errJson := json.Marshal(req)
	if errJson != nil {
		panic(errJson)
	}
	if _, err := ws.Write(bts); err != nil {
		panic(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		panic(err)
	}

	resp := &server.Response{}

	respErr := json.Unmarshal(msg[:n], resp)

	if respErr != nil {
		panic(respErr)
	}

	respR, ok := resp.Result.(map[string]interface{})
	if(!ok){
		t.Error("Response result cannot be cast to map")	
	}
	retStr := respR["value"].(string)
	if retStr != "testmessage" {
		t.Error("Expected: testmessage, Got: " + retStr)
	} else {
		fmt.Println("Websocket echo test: PASSED")
	}
	ws.Close()
		
	time.Sleep(1*time.Second)
}
