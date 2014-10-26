package ponos

import (
	"encoding/json"
	"errors"
	"github.com/eris-ltd/deCerver/administrator/server"
	"github.com/eris-ltd/deCerver/scriptengine"
	"github.com/eris-ltd/thelonious/monk"
	"github.com/eris-ltd/thelonious"
	"github.com/golang/glog"
	"time"
)

type PathRequest struct {
	Path []string
	RootAddress string
}

type PonosSRPCFactory struct {
	serviceName string
	ethereum    *eth.Ethereum
}

func NewPonosSRPCFactory(ethereum *eth.Ethereum) *PonosSRPCFactory {
	fact := &PonosSRPCFactory{serviceName: "PonosSRPC", ethereum: ethereum}
	return fact
}

func (fact *PonosSRPCFactory) Init() {

}

func (fact *PonosSRPCFactory) ServiceName() string {
	return fact.serviceName
}

func (fact *PonosSRPCFactory) CreateService() server.SRPCService {
	ec := monk.NewEth(fact.ethereum)
	ec.Init()
	service := newPonosSRPC()
	service.name = fact.serviceName
	service.tp = scriptengine.NewTreeParser(ec)
	return service
}

type PonosSRPC struct {
	name       string
	mappings   map[string]server.RPCMethod
	conn       *server.WsConn
	forwardLog bool
	tp         *scriptengine.TreeParser
}

// Create a new handler
func newPonosSRPC() *PonosSRPC {
	srh := &PonosSRPC{}

	srh.mappings = make(map[string]server.RPCMethod)
	srh.mappings["GetTree"] = srh.GetTree
	srh.mappings["GetTreeFromPath"] = srh.GetTreeFromPath
	
	return srh
}

func (psrpc *PonosSRPC) SetConnection(wsConn *server.WsConn) {
	psrpc.conn = wsConn
}

func (psrpc *PonosSRPC) Init() {
}

func (psrpc *PonosSRPC) Close() {
}

func (psrpc *PonosSRPC) Name() string {
	return psrpc.name
}

func (psrpc *PonosSRPC) HandleRPC(rpcReq *server.RPCRequest) (*server.RPCResponse, error) {
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
func (psrpc *PonosSRPC) AddMethod(methodName string, method server.RPCMethod, replaceOld bool) error {
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
func (psrpc *PonosSRPC) RemoveMethod(methodName string) {
	if psrpc.mappings[methodName] == nil {
		glog.Info("Removal failed. There is no handler for '" + methodName + "'.")
	} else {
		delete(psrpc.mappings, methodName)
	}
	return
}

func (psrpc *PonosSRPC) GetTree(req *server.RPCRequest, resp *server.RPCResponse) {
	// TODO fix when we have doug of all dougs.
	tree, err := psrpc.tp.ParseTree("ponos", "0x99883cb40909f59f8082128afd4b67a2ca2b206c")
	if err != nil {
		resp.Error = err.Error()
	} else {
		ft := psrpc.tp.FlattenTree(tree)
		resp.Result = ft.Tree
	}
}

func (psrpc *PonosSRPC) GetTreeFromPath(req *server.RPCRequest, resp *server.RPCResponse) {
	params := &PathRequest{}
	err := json.Unmarshal(*req.Params, params)

	if err != nil {
		resp.Error = err.Error()
		return
	}
	
	tree, err := psrpc.tp.ParseTreeByPath(params.Path,params.RootAddress)
	if err != nil {
		resp.Error = err.Error()
	} else {
		ft := psrpc.tp.FlattenTree(tree)
		resp.Result = ft.Tree
	}
}
