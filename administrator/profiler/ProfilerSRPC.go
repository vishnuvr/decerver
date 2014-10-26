package profiler

import (
	"runtime"
	"github.com/golang/glog"
	"errors"
	"time"
	"github.com/eris-ltd/deCerver/administrator/server"
)

type ProfilerSRPCFactory struct {
	serviceName string
}

func NewProfilerSRPCFactory() *ProfilerSRPCFactory {
	fact := &ProfilerSRPCFactory{serviceName: "ProfilerSRPC"}
	return fact
}

func (fact *ProfilerSRPCFactory) Init() {

}

func (fact *ProfilerSRPCFactory) ServiceName() string {
	return fact.serviceName
}

func (fact *ProfilerSRPCFactory) CreateService() server.SRPCService {
	service := newProfilerSRPC()
	service.name = fact.serviceName
	return service
}

type ProfilerSRPC struct {
	name     string
	mappings map[string]server.RPCMethod
	conn     *server.WsConn
	forwardLog bool
}

// Create a new handler
func newProfilerSRPC() *ProfilerSRPC {
	srh := &ProfilerSRPC{}

	srh.mappings = make(map[string]server.RPCMethod)
	srh.mappings["MemStats"] = srh.MemStats
	srh.mappings["ForwardLog"] = srh.ForwardLog
	
	return srh
}

func (psrpc *ProfilerSRPC) SetConnection(wsConn *server.WsConn) {
	psrpc.conn = wsConn
}

func (psrpc *ProfilerSRPC) Init() {
}

func (psrpc *ProfilerSRPC) Close() {
}

func (psrpc *ProfilerSRPC) Name() string {
	return psrpc.name
}

func (psrpc *ProfilerSRPC) HandleRPC(rpcReq *server.RPCRequest) (*server.RPCResponse, error) {
	methodName := rpcReq.Method
	resp := &server.RPCResponse{}
	if psrpc.mappings[methodName] == nil {
		glog.Errorf("Method not supported: %s\n", methodName)
		return nil, errors.New("SRPC Method not supported.")
	}

	// Run the method.
	psrpc.mappings[methodName](rpcReq, resp)
	// Add a timestamp.
	resp.Timestamp = int(time.Now().In(time.UTC).Unix())
	// The ID is the method being called, for now.
	resp.Id = methodName

	return resp, nil
}

// Add a new method
func (psrpc *ProfilerSRPC) AddMethod(methodName string, method server.RPCMethod, replaceOld bool) error {
	if psrpc.mappings[methodName] != nil {
		if !replaceOld {
			return errors.New("Tried to overwrite an already existing method.")
		} else {
			glog.Info("Overwriting old method for '" + methodName + "'.")
		}

	}
	psrpc.mappings[methodName] = method
	return nil
}

// Remove a method
func (psrpc *ProfilerSRPC) RemoveMethod(methodName string) {
	if psrpc.mappings[methodName] == nil {
		glog.Info("Removal failed. There is no handler for '" + methodName + "'.")
	} else {
		delete(psrpc.mappings, methodName)
	}
	return
}

func (psrpc *ProfilerSRPC) MemStats(req *server.RPCRequest, resp *server.RPCResponse) {
	
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	// We want to be able to send a number of different errors at some point.
	resp.Result = ms
}

func (psrpc *ProfilerSRPC) ForwardLog(req *server.RPCRequest, resp *server.RPCResponse) {
	
}
