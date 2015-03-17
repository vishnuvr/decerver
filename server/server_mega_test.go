package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const TEST_NUM = 100

// Test sending a http request to the echo endpoint
func TestHttpMegaEcho(t *testing.T) {
	for i := 0; i < TEST_NUM; i++ {
		go func() {
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
		}()
	}
}

// Establish websocket connection and rpc to 'echo'
func TestWsMegaEcho(t *testing.T) {
	doneChan := make(chan bool)
	for i := 0; i < TEST_NUM; i++ {
		go func() {
			origin := "http://localhost/"
			url := "ws://localhost:3000/socket"
			ws, err := websocket.Dial(url, "", origin)
			if err != nil {
				panic(err)
			}
			req := &Request{}
			req.ID = 1
			req.JsonRpc = "2.0"
			req.Method = "echo"

			sVal := &StringValue{"testmessage"}
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

			resp := &Response{}

			respErr := json.Unmarshal(msg[:n], resp)

			if respErr != nil {
				panic(respErr)
			}

			respR := resp.Result.(map[string]interface{})
			retStr := respR["SVal"].(string)
			if retStr != "testmessage" {
				t.Error("Expected: testmessage, Got: " + retStr)
			} else {
				//fmt.Println("Websocket echo test: PASSED")
			}
			ws.Close()
			doneChan <- true
		}()
	}
	ctr := 0
	for ctr < TEST_NUM {
		_ = <-doneChan
		ctr++
	}
	time.Sleep(1 * time.Second)
}