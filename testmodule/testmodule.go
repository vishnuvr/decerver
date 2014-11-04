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

type HelloService struct {
	ate core.Runtime
}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	ret, err := h.ate.RunMethod("", "test", "world")
	if err != nil {
		return err
	}
	logger.Printf("Ret: %v\n", ret)
	reply.Message = ret[0]
	//reply.Message = "testvalue"
	return nil
}

type TestModule struct {
	helloService *HelloService
}

func (tm *TestModule) Register(fileIO core.FileIO, registry api.ApiRegistry, runtime core.Runtime) error {
	logger = log.New(os.Stdout, "TestModule", 5)

	runtime.BindScriptObject("test", func(otto.FunctionCall) otto.Value {
		val, _ := otto.ToValue("It's working")
		return val
	})

	tm.helloService = &HelloService{}
	tm.helloService.ate = runtime
	registry.RegisterHttpServices(tm.helloService)
	return nil
}

func (tm *TestModule) Init() error {
	return nil
}

func (tm *TestModule) Start() error {
	return nil
}

func (tm *TestModule) ReadConfig(config_file string) {

}

func (tm *TestModule) WriteConfig(config_file string) {

}

func (tm *TestModule) Name() string {
	return "TestModule"
}

func (tm *TestModule) Shutdown() error {
	logger.Println("Goodbye from TestModule")
	return nil
}
