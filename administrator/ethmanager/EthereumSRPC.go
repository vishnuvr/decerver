package ethmanager

// This handles socket-based rpc. Part of it is reacting to requests sent from the
// client, and part of it is reacting to changes in the ethereum world state,
// and propagating these.
import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eris-ltd/deCerver/administrator/server"
	"github.com/eris-ltd/thelonious"
	"github.com/eris-ltd/thelonious/ethchain"
	"github.com/eris-ltd/thelonious/ethreact"
	"github.com/eris-ltd/thelonious/monk"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"time"
)

type EthereumSRPCFactory struct {
	ethereum    *eth.Ethereum
	ethLogger   *EthLogger
	serviceName string
}

func NewEthereumSRPCFactory(ethereum *eth.Ethereum) *EthereumSRPCFactory {
	fact := &EthereumSRPCFactory{
		ethereum:    ethereum,
		ethLogger:   NewEthLogger(),
		serviceName: "EthereumSRPC",
	}
	return fact
}

func (fact *EthereumSRPCFactory) Init() {

}

func (fact *EthereumSRPCFactory) ServiceName() string {
	return fact.serviceName
}

func (fact *EthereumSRPCFactory) CreateService() server.SRPCService {
	ec := monk.NewEth(fact.ethereum)
	ec.Init()
	service := newEthereumSRPC(ec)
	service.name = fact.serviceName
	service.ethLogger = fact.ethLogger
	return service
}

type EthereumSRPC struct {
	name     string
	mappings map[string]server.RPCMethod
	// TODO replace with pipe.
	ethChain    *monk.EthChain
	conn        *server.WsConn
	ethLogger   *EthLogger
	ethListener *EthListener
}

// Create a new handler
func newEthereumSRPC(eth *monk.EthChain) *EthereumSRPC {
	esrpc := &EthereumSRPC{}
	esrpc.ethChain = eth

	esrpc.mappings = make(map[string]server.RPCMethod)
	esrpc.mappings["CompileMutan"] = esrpc.CompileMutan
	esrpc.mappings["MyBalance"] = esrpc.MyBalance
	esrpc.mappings["MyAddress"] = esrpc.MyAddress
	esrpc.mappings["StartMining"] = esrpc.StartMining
	esrpc.mappings["StopMining"] = esrpc.StopMining
	esrpc.mappings["LastBlockNumber"] = esrpc.LastBlockNumber
	esrpc.mappings["BlockByHash"] = esrpc.BlockByHash
	esrpc.mappings["Account"] = esrpc.Account
	esrpc.mappings["Transact"] = esrpc.Transact
	esrpc.mappings["WorldState"] = esrpc.WorldState

	return esrpc
}

func (esrpc *EthereumSRPC) SetConnection(wsConn *server.WsConn) {
	esrpc.conn = wsConn
}

func (esrpc *EthereumSRPC) Init() {
	esrpc.ethListener = newEthListener(esrpc)
}

func (esrpc *EthereumSRPC) Close() {
	esrpc.ethListener.Close()
}

func (esrpc *EthereumSRPC) Name() string {
	return esrpc.name
}

func (esrpc *EthereumSRPC) HandleRPC(rpcReq *server.RPCRequest) (*server.RPCResponse, error) {
	methodName := rpcReq.Method
	resp := &server.RPCResponse{}
	if esrpc.mappings[methodName] == nil {
		glog.Errorf("Method not supported: %s\n", methodName)
		return nil, errors.New("SRPC Method not supported.")
	}

	// Run the method.
	esrpc.mappings[methodName](rpcReq, resp)
	// Add a timestamp.
	resp.Timestamp = getTimestamp()
	// The ID is the method being called, for now.
	resp.Id = methodName

	return resp, nil
}

// Add a new method
func (esrpc *EthereumSRPC) AddMethod(methodName string, method server.RPCMethod, replaceOld bool) error {
	if esrpc.mappings[methodName] != nil {
		if !replaceOld {
			return errors.New("Tried to overwrite an already existing method.")
		} else {
			glog.Info("Overwriting old method for '" + methodName + "'.")
		}

	}
	esrpc.mappings[methodName] = method
	return nil
}

// Remove a method
func (esrpc *EthereumSRPC) RemoveMethod(methodName string) {
	if esrpc.mappings[methodName] == nil {
		glog.Info("Removal failed. There is no handler for '" + methodName + "'.")
	} else {
		delete(esrpc.mappings, methodName)
	}
	return
}

func (esrpc *EthereumSRPC) CompileMutan(req *server.RPCRequest, resp *server.RPCResponse) {
	params := &VString{}
	err := json.Unmarshal(*req.Params, params)

	if err != nil {
		resp.Error = err.Error()
		return
	}

	retVal := &CompilerResult{}
	// We want to be able to send a number of different errors at some point.
	bytecode, compErr := esrpc.ethChain.Pipe.CompileMutan(params.SVal)
	if compErr != nil {
		retVal.Errors = make([]string, 1)
		retVal.Errors[0] = compErr.Error()
		retVal.Success = false
	} else {
		retVal.Bytes = hex.EncodeToString(bytecode)
		retVal.Success = true
	}
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) MyBalance(req *server.RPCRequest, resp *server.RPCResponse) {
	retVal := &VString{}
	// TODO Replace with pipe
	myAddr := esrpc.ethChain.Ethereum.KeyManager().Address()
	balance := esrpc.ethChain.Pipe.Balance(myAddr)
	// -----------------
	retVal.SVal = balance.String()
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) MyAddress(req *server.RPCRequest, resp *server.RPCResponse) {
	retVal := &VString{}
	retVal.SVal = hex.EncodeToString(esrpc.ethChain.Ethereum.KeyManager().Address())
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) StartMining(req *server.RPCRequest, resp *server.RPCResponse) {
	retVal := &VBool{}
	retVal.BVal = esrpc.ethChain.StartMining()
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) StopMining(req *server.RPCRequest, resp *server.RPCResponse) {
	retVal := &VBool{}
	retVal.BVal = esrpc.ethChain.StopMining()
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) LastBlockNumber(req *server.RPCRequest, resp *server.RPCResponse) {
	retVal := &VInteger{}
	retVal.IVal = getLastBlockNumber(esrpc.ethChain.Ethereum)
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) BlockByHash(req *server.RPCRequest, resp *server.RPCResponse) {
	params := &VString{}
	err := json.Unmarshal(*req.Params, params)

	if err != nil {
		resp.Error = err.Error()
		return
	}

	retVal := &BlockData{}
	hash, decErr := hex.DecodeString(params.SVal)

	if decErr != nil {
		resp.Error = decErr.Error()
		return
	}

	block := esrpc.ethChain.Pipe.Block(hash)
	if block == nil {
		resp.Error = "No block with hash: " + params.SVal
		return
	}

	getBlockDataFromBlock(retVal, block)
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) Account(req *server.RPCRequest, resp *server.RPCResponse) {
	params := &VString{}
	err := json.Unmarshal(*req.Params, params)

	if err != nil {
		resp.Error = err.Error()
		return
	}

	retVal := &Account{}
	addr, decErr := hex.DecodeString(params.SVal)

	if decErr != nil {
		resp.Error = decErr.Error()
		return
	}

	curBlock := esrpc.ethChain.Ethereum.BlockChain().CurrentBlock
	account := curBlock.State().GetStateObject(addr)
	if account == nil {
		resp.Error = "No account with address: " + params.SVal
		return
	}

	getAccountFromStateObject(retVal, account)
	resp.Result = retVal
}

/*
type TxIndata struct {
	Recipient string
	Value     string
	Gas       string
	GasCost   string
	Data      string
}
*/

func (esrpc *EthereumSRPC) Transact(req *server.RPCRequest, resp *server.RPCResponse) {
	params := &TxIndata{}
	err := json.Unmarshal(*req.Params, params)

	if err != nil {
		resp.Error = err.Error()
		return
	}

	retVal := &TxReceipt{}
	// TODO check sender.
	err = createTx(esrpc.ethChain.Ethereum, params.Recipient, params.Value, params.Gas, params.GasCost, params.Data, retVal)
	if err != nil {
		retVal.Error = err.Error()
	}
	resp.Result = retVal
}

func (esrpc *EthereumSRPC) WorldState(req *server.RPCRequest, resp *server.RPCResponse) {

	blocks := getWorldState(esrpc.ethChain.Ethereum)
	// Let the client know how many blocks there are.
	resp = &server.RPCResponse{}
	resp.Id = "NumBlocks"
	resp.Result = &VInteger{IVal: len(blocks) - 1}
	resp.Timestamp = getTimestamp()
	esrpc.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}

	// Send blocks one at a time.
	for i := 0; i < len(blocks); i++ {
		resp = &server.RPCResponse{}
		resp.Id = "Blocks"
		resp.Result = blocks[i]
		resp.Timestamp = getTimestamp()
		esrpc.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
	}

	accounts := getAccounts(esrpc.ethChain.Ethereum)

	// Let the client know how many accounts there are.
	resp = &server.RPCResponse{}
	resp.Id = "NumAccounts"
	resp.Result = &VInteger{IVal: len(accounts)}
	resp.Timestamp = getTimestamp()
	esrpc.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}

	// Dispatch these one at a time, and also register listeners to all these addresses.
	for i := 0; i < len(accounts); i++ {
		resp = &server.RPCResponse{}
		resp.Id = "Accounts"
		resp.Result = accounts[i]
		resp.Timestamp = getTimestamp()
		esrpc.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
	}

	// Finalize.
	resp = &server.RPCResponse{}
	resp.Id = "WorldStateDone"
	resp.Result = &NoArgs{}
	resp.Timestamp = getTimestamp()
	esrpc.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}

}

type EthListener struct {
	ethRPC            *EthereumSRPC
	txPreChannel      chan ethreact.Event
	txPreFailChannel  chan ethreact.Event
	txPostChannel     chan ethreact.Event
	txPostFailChannel chan ethreact.Event
	blockChannel      chan ethreact.Event
	stopChannel       chan bool
	logSub            *LogSub
}

func newEthListener(ethRPC *EthereumSRPC) *EthListener {
	el := &EthListener{}
	el.ethRPC = ethRPC

	el.logSub = NewStdLogSub()
	el.logSub.SubId = ethRPC.conn.SessionId
	ethRPC.ethLogger.AddSub(el.logSub)
	el.blockChannel = make(chan ethreact.Event, 10)
	el.txPreChannel = make(chan ethreact.Event, 10)
	el.txPreFailChannel = make(chan ethreact.Event, 10)
	el.txPostChannel = make(chan ethreact.Event, 10)
	el.txPostFailChannel = make(chan ethreact.Event, 10)
	el.stopChannel = make(chan bool)
	el.ethRPC.ethChain.Ethereum.Reactor().Subscribe("newBlock", el.blockChannel)
	el.ethRPC.ethChain.Ethereum.Reactor().Subscribe("newTx:pre", el.txPreChannel)
	el.ethRPC.ethChain.Ethereum.Reactor().Subscribe("newTx:pre:fail", el.txPreFailChannel)
	el.ethRPC.ethChain.Ethereum.Reactor().Subscribe("newTx:post", el.txPostChannel)
	el.ethRPC.ethChain.Ethereum.Reactor().Subscribe("newTx:post:fail", el.txPostFailChannel)

	go func(el *EthListener) {
		for {
			select {
			case evt := <-el.blockChannel:
				block, _ := evt.Resource.(*ethchain.Block)
				fmt.Println("Block added")
				resp := &server.RPCResponse{}
				resp.Id = "BlockAdded"
				bd := &BlockMiniData{}
				getBlockMiniDataFromBlock(el.ethRPC.ethChain.Ethereum, bd, block)
				resp.Result = bd
				resp.Timestamp = getTimestamp()
				el.ethRPC.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
			case evt := <-el.txPreChannel:
				tx, _ := evt.Resource.(*ethchain.Transaction)
				resp := &server.RPCResponse{}
				resp.Id = "TxPre"
				trans := &Transaction{}
				getTransactionFromTx(trans, tx)
				resp.Result = trans
				resp.Timestamp = getTimestamp()
				el.ethRPC.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
			case evt := <-el.txPreFailChannel:
				txFail, _ := evt.Resource.(*ethchain.TxFail)
				resp := &server.RPCResponse{}
				resp.Id = "TxPreFail"
				trans := &Transaction{}
				getTransactionFromTx(trans, txFail.Tx)
				trans.Error = txFail.Err.Error()
				resp.Result = trans
				resp.Timestamp = getTimestamp()
				el.ethRPC.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
			case evt := <-el.txPostChannel:
				tx, _ := evt.Resource.(*ethchain.Transaction)
				resp := &server.RPCResponse{}
				resp.Id = "TxPost"
				trans := &Transaction{}
				getTransactionFromTx(trans, tx)
				resp.Result = trans
				resp.Timestamp = getTimestamp()
				el.ethRPC.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
			case evt := <-el.txPostFailChannel:
				txFail, _ := evt.Resource.(*ethchain.TxFail)
				resp := &server.RPCResponse{}
				resp.Id = "TxPostFail"
				trans := &Transaction{}
				getTransactionFromTx(trans, txFail.Tx)
				trans.Error = txFail.Err.Error()
				resp.Result = trans
				resp.Timestamp = getTimestamp()
				el.ethRPC.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
			case txt := <-el.logSub.Channel:
				resp := &server.RPCResponse{}
				resp.Id = "Log"
				resp.Result = &VString{SVal: txt}
				resp.Timestamp = getTimestamp()
				el.ethRPC.conn.WriteMsgChannel <- &server.Message{Data: resp, Type: websocket.TextMessage}
			case <-el.stopChannel:
				// Quit this
				return
			}
		}
	}(el)

	return el
}

func (el *EthListener) Close() {
	rctr := el.ethRPC.ethChain.Ethereum.Reactor()
	rctr.Unsubscribe("newBlock", el.blockChannel)
	rctr.Unsubscribe("newTx:pre", el.txPreChannel)
	rctr.Unsubscribe("newTx:pre:fail", el.txPreFailChannel)
	rctr.Unsubscribe("newTx:post", el.txPostChannel)
	rctr.Unsubscribe("newTx:post:fail", el.txPostFailChannel)
	el.ethRPC.ethLogger.RemoveSub(el.logSub)
}

func getTimestamp() int {
	return int(time.Now().In(time.UTC).UnixNano() >> 6)
}
