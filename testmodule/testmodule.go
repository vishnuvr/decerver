package testmodule

import (
	"github.com/eris-ltd/deCerver-interfaces/api"
	"github.com/eris-ltd/deCerver-interfaces/core"
	"github.com/robertkrimen/otto"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct{
	ate	core.ScriptEngine
}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	//ret := h.ate.RunMethod("","test","world")
	//logger.Printf("Ret: %v\n",ret)
	//reply.Message = ret[0]
	reply.Message = "testvalue"
	return nil
}

type TestModule struct {
	helloService *HelloService
}

func (tm *TestModule) Init(se core.ScriptEngine) {
	logger = log.New(os.Stdout, "TestModule", 5)
	
	se.InjectFunction("test",func (otto.FunctionCall) otto.Value {
		val, _ := otto.ToValue("It's working")
		return val
	})
		
	tm.helloService = &HelloService{}
	tm.helloService.ate = se
}

func (tm *TestModule) Logger() *log.Logger {
	return logger
}

func (tm *TestModule) Name() string {
	return "TestModule"
}

func (tm *TestModule) HttpAPIServices() []interface{} {
	ret := make([]interface{},1)
	ret[0] = tm.helloService
	return ret
}

func (tm *TestModule) WsAPIServiceFactories() []api.WsAPIServiceFactory {
	return nil
}

func (tm *TestModule) Shutdown() {
	logger.Println("Goodbye from TestModule")
}


