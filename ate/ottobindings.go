package ate

import (
	"encoding/hex"
	"fmt"
	"github.com/obscuren/sha3"
	"github.com/robertkrimen/otto"
	//"github.com/eris-ltd/decerver-interfaces/events"
	"math/big"
)

// Enables math operations on strings using big.Int.
var BZERO *big.Int = big.NewInt(0)

func isZero(i *big.Int) bool {
	return i.Cmp(BZERO) == 0
}

func BindDefaults(vm *otto.Otto) {
	var err error
	
	// Networking.
	_ , err = vm.Run(`
	
		var jsonErrors = {
			"E_PARSE"       : -32700,
			"E_INVALID_REQ" : -32600,
			"E_NO_METHOD"   : -32601,
			"E_BAD_PARAMS"  : -32602,
			"E_INTERNAL"    : -32603,
			"E_SERVER"      : -32000
		};
	
		// Network is an object that encapsulates all networking activity.
		var network = {};
		
		network.incomingHttpCallback = function(){};
		
		// Used internally.
		network.handleIncomingHttp = function(httpReqAsJson){
			var httpReq = JSON.parse(reqAsJson);
			this.incomingHttpCallback(httpReq);
		};
		
		network.registerIncomingHttpCallback = function(callback){
			if(typeof handler !== "function"){
				throw Error("Attempting to register a non-function as incoming http handler");
			}
			network.httpHandler = handler;
		}
		
		// Websockets
		
		// Each session has a handler.
		network.wsHandlers = {};
		
		network.newWsCallback = function(sessionObj){
			return function (){
				console.log("No callback registered for new websocket connections.");
			};
		};
		
		network.newWsSession = function(sessionObj){
			console.log("Adding new session: " + sessionObj.SessionId());
			network.wsHandlers[sessionObj.SessionId()] = network.newWsCallback(sessionObj);
		}
		
		network.incomingWsMsg = function(sessionId, reqJson) {
			console.log("[Otto] Incoming websocket message.");
			try {
				var request = JSON.parse(reqJson);
				if (typeof(request.Method) === "undefined" || request.Method === ""){
					return JSON.stringify(network.getWsError(jsonErrors.E_NO_METHOD, "No method in request."));
				} else {
					var handler = network.wsHandlers[sessionId];
					if (typeof handler !== "function"){
						return JSON.stringify(network.getWsError(jsonErrors.E_SERVER, "Handler not registered for websocket session: " + sessionId.toString()) );
					}
					var response = handler(request);
					if(response === null){
						return null;
					}
					var respStr;
					try {
						respStr = JSON.stringify(response);
					} catch (err) {
						return JSON.stringify(network.getWsError(jsonErrors.E_INTERNAL, "Failed to marshal response: " + err));
					}
					return respStr;
				}
			} catch (err){
				response = JSON.stringify(network.getWsError(jsonErrors.E_PARSE, err));
			}
		}
		
		network.newWsRequest = function(){
			return {
				"Jsonrpc" : 2.0,
				"Id" : "",
				"Method" : "",
				"Params" : "",
			};
		}
		
		network.getWsResponse = function(){
			return {
				"Jsonrpc" : 2.0,
				"Id" : "",
				"Result" : "",
				"Error" : "",
			};
		}
		
		network.getWsError = function(code, message){
			return {
				"Jsonrpc" : 2.0,
				"Id" : "",
				"Result" : "",
				"Error" : {
					"Code" : code,
					"Message" : message,
					"Data" : null
				  }
			};
		}
		
	`)
	
	if err != nil {
		fmt.Println("[Atë] Error while bootstrapping js networking: " + err.Error())
	} else {
		fmt.Println("[Atë] Networking script loaded.")
	}
	
	// Event manager.
	_ , err = vm.Run(`
	
		var events = {};
		
		events.callbacks = {};
		
		events.registerCallback = function(eventSource, eventType, target, callbackName, callbackFn){
			if(typeof(callbackFn) !== "function"){
				throw Error("Trying to register a non callback function as callback.");
			}
			/*
			val res = RegEvtSub(eventSource,eventType);
			if (res) {
				this.callbacks[callbackFn.toString()] = callbackFn;
			} else {
				console.log("Failed to register event. (Source: " + eventSource + ") (" + eventType + ") (Callback: " + callbackFn.toString() + ")");
			}
			return res;
			*/
		}
		
		events.unregisterCallback = function(callbackName){
			events.callbacks[callbackName] = null;
		}
		
		// Called by the go event processor.
		events.post = function(eventJson){
			var event = JSON.parse(eventJson);
			var fnName = event.Id;
			this.callbacks[fnName](event);
		}
	`)

	if err != nil {
		fmt.Println("[Atë] Error while bootstrapping js event-processing: " + err.Error())
	} else {
		fmt.Println("[Atë] Event processing script loaded.")
	}
	
	bindHelpers(vm)
}

func bindHelpers(vm *otto.Otto) {
	vm.Set("Add", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		result, _ := vm.ToValue("0x" + p0.Add(p0, p1).String())
		return result
	})

	vm.Set("Sub", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		p0.Sub(p0, p1)
		if p0.Sign() < 0 {
			otto.NaNValue() // TODO
		}
		result, _ := vm.ToValue("0x" + p0.String())
		return result
	})

	vm.Set("Mul", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		result, _ := vm.ToValue("0x" + p0.Mul(p0, p1).String())
		return result
	})

	vm.Set("Div", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		if isZero(p1) {
			return otto.NaNValue()
		}
		result, _ := vm.ToValue("0x" + p0.Div(p0, p1).String())
		return result
	})

	vm.Set("Mod", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		if isZero(p1) {
			return otto.NaNValue()
		}
		result, _ := vm.ToValue("0x" + p0.Mod(p0, p1).String())
		return result
	})

	vm.Set("Exp", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		result, _ := vm.ToValue("0x" + p0.Exp(p0, p1, nil).String())
		return result
	})

	vm.Set("IsZero", func(call otto.FunctionCall) otto.Value {
		prm, err0 := call.Argument(0).ToString()
		if err0 != nil {
			return otto.UndefinedValue()
		}
		isZero := prm == "0" || prm == "0x" || prm == "0x0"
		result, _ := vm.ToValue(isZero)

		return result
	})

	// Crypto
	vm.Set("SHA3", func(call otto.FunctionCall) otto.Value {
		prm, err0 := call.Argument(0).ToString()
		if err0 != nil {
			return otto.UndefinedValue()
		}
		h, err := hex.DecodeString(prm)
		if err != nil {
			return otto.UndefinedValue()
		}
		d := sha3.NewKeccak256()
		d.Write(h)
		result, _ := vm.ToValue(hex.EncodeToString(d.Sum(nil)))

		return result
	})
}

func parseUn(call otto.FunctionCall) (*big.Int, error) {
	str, err0 := call.Argument(0).ToString()
	if err0 != nil {
		return nil, err0
	}
	val := atob(str)
	return val, nil
}

func parseBin(call otto.FunctionCall) (*big.Int, *big.Int, error) {
	left, err0 := call.Argument(0).ToString()
	if err0 != nil {
		return nil, nil, err0
	}
	right, err1 := call.Argument(1).ToString()

	if err1 != nil {
		return nil, nil, err1
	}
	p0 := atob(left)
	p1 := atob(right)
	return p0, p1, nil
}

func atob(str string) *big.Int {
	i := new(big.Int)
	fmt.Sscan(str, i)
	return i
}
