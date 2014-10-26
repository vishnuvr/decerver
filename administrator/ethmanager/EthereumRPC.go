package ethmanager

/*
import (
	"encoding/hex"
	"fmt"
	"github.com/obscuren/mutan"
	"github.com/ethereum/eth-go"
	"github.com/ethereum/eth-go/ethutil"
	"net/http"
	"strings"
)

type EthereumRPC struct {
	EthCli   *eth.Ethereum
	Compiler *mutan.Compiler
	CurBlock int
}

func (ethRPC *EthereumRPC) CompileMutan(r *http.Request, args *VString, reply *CompilerResult) error {
	reader := strings.NewReader(args.SVal)
	bytecode, errs := ethRPC.Compiler.Compile(reader)
	if len(errs) > 0 {
		reply.Errors = make([]string, len(errs))
		for idx := range errs {
			reply.Errors[idx] = errs[idx].Error()
		}
		reply.Success = false
	} else {
		reply.Bytes = hex.EncodeToString(bytecode)
		reply.Success = true
	}
	return nil
}

func (ethRPC *EthereumRPC) IsContract(r *http.Request, args *VString, reply *VBool) error {
	reply.BVal = isContract(ethRPC.EthCli,args.SVal)
	return nil
}

func (ethRPC *EthereumRPC) BalanceAt(r *http.Request, args *VString, reply *VString) error {
	sHex, err := hex.DecodeString(args.SVal)
	if(err != nil){
		fmt.Println(err.Error())
		reply.SVal = ERR_MALFORMED_ADDRESS
		return nil
	}
	balance := ethRPC.EthCli.StateManager().CurrentState().GetBalance(sHex)
	reply.SVal = balance.String()
	return nil
}

func (ethRPC *EthereumRPC) MyBalance(r *http.Request, args *NoArgs, reply *VString) error {
	myAddr := ethRPC.EthCli.KeyManager().Address();
	balance := ethRPC.EthCli.StateManager().CurrentState().GetBalance(myAddr)
	reply.SVal = balance.String()
	return nil
}

func (ethRPC *EthereumRPC) StorageAt(r *http.Request, args *StateAtArgs, reply *VString) error {
	stateobj := getStateObject(ethRPC.EthCli,args.Address);
	if(stateobj == nil){
		reply.SVal = ERR_NO_SUCH_ADDRESS
		return nil
	}
	storage := stateobj.GetStorage(ethutil.Big(args.Storage))
	if(storage == nil){
		reply.SVal = ERR_STATE_NO_STORAGE
		return nil
	}
	reply.SVal = storage.String();
	return nil
}

// TODO min gascost is hardcoded in block_chain NewBlock. Will this vary at some point?
func (ethRPC *EthereumRPC) MinGascost(r *http.Request, args *NoArgs, reply *VString) error {
	reply.SVal = "10000000000000"
	return nil
}

func (ethRPC *EthereumRPC) StartMining(r *http.Request, args *NoArgs, reply *VBool) error {
	reply.BVal = startMining(ethRPC.EthCli)
	return nil
}

func (ethRPC *EthereumRPC) StopMining(r *http.Request, args *NoArgs, reply *VBool) error {
	reply.BVal = stopMining(ethRPC.EthCli)
	return nil
}

func (ethRPC *EthereumRPC) IsMining(r *http.Request, args *NoArgs, reply *VBool) error {
	reply.BVal = ethRPC.EthCli.IsMining()
	return nil
}

func (ethRPC *EthereumRPC) MyAddress(r *http.Request, args *NoArgs, reply *VString) error {
	reply.SVal = hex.EncodeToString(ethRPC.EthCli.KeyManager().Address())
	return nil
}

func (ethRPC *EthereumRPC) BlockGenesis(r *http.Request, args *NoArgs, reply *BlockData) error {

	block := ethRPC.EthCli.BlockChain().Genesis()

	// Block Number
	reply.Number = block.Number.String()

	// Block Time
	reply.Time = int(block.Time)

	// Block Nonce
	reply.Nonce = hex.EncodeToString(block.Nonce)

	// Block Transactions (hashes)
	trsct := block.Transactions()

	reply.Transactions = make([]*Transaction, len(trsct))
	if trsct != nil {
		for idx := range trsct {
			reply.Transactions[idx] = &Transaction{}
			transactionFromTx(reply.Transactions[idx],trsct[idx])
		}
	}

	// Block Hash (from args)
	reply.Hash = hex.EncodeToString(block.Hash())

	// Block PrevHash
	reply.PrevHash = ""

	// Block Difficulty
	reply.Difficulty = block.Difficulty.String()

	// Block Coinbase
	reply.Coinbase = ""

	// Uncles
	reply.Uncles = make([]string, 0)

	return nil
}

func (ethRPC *EthereumRPC) BlockLatest(r *http.Request, args *NoArgs, reply *BlockData) error {
	addr, _ := hex.DecodeString("29c8e2e2a699ed64296025795b5dca20647c66de")
	acc := ethRPC.EthCli.StateManager().CurrentState().GetAccount(addr)
	if acc == nil {
		fmt.Println("No such account.");
	}
	fmt.Printf("%x\n",acc.CodeHash);
	lbh := ethRPC.EthCli.BlockChain().LastBlockHash
	argz := &VString{SVal: hex.EncodeToString(lbh)}
	ethRPC.BlockByHash(r, argz, reply)
	return nil
}

func (ethRPC *EthereumRPC) BlockMiniByHash(r *http.Request, args *VString, reply *BlockMiniData) error {
	// Get the block.
	bts, err := hex.DecodeString(args.SVal)

	if err != nil {
		fmt.Println(err.Error())
		reply.Hash = ERR_MALFORMED_BLOCK_HASH
		return nil
	}

	block := ethRPC.EthCli.BlockChain().GetBlock(bts)

	if block == nil {
		reply.Hash = ERR_NO_SUCH_BLOCK
		return nil
	}
	getBlockMiniDataFromBlock(reply,block)

	return nil
}

func (ethRPC *EthereumRPC) Sync(r *http.Request, args *VInteger, reply *Sync) error {
	lastnum := ethRPC.EthCli.BlockChain().LastBlockNumber

	// Difference between available blocks and client latest blocks.
	ctr := int(lastnum) - args.IVal
	if ctr > 0 {
		reply.NewBlocks = make([]*BlockMiniData, ctr)
		lastHash := ethRPC.EthCli.BlockChain().LastBlockHash
		block := ethRPC.EthCli.BlockChain().CurrentBlock;
		ctr = ctr - 1
		reply.NewBlocks[ctr] = &BlockMiniData{}
		getBlockMiniDataFromBlock(reply.NewBlocks[ctr],block)

		for ctr > 0 {
			ctr = ctr - 1
			reply.NewBlocks[ctr] = &BlockMiniData{}
			lastHash = block.PrevHash
			block = ethRPC.EthCli.BlockChain().GetBlock(lastHash)
			getBlockMiniDataFromBlock(reply.NewBlocks[ctr],block)
		}
	}

	reply.IsMining = ethRPC.EthCli.IsMining()
	reply.PreTxs = popTxPre()
	reply.PostTxs = popTxPost()
	myAddr := ethRPC.EthCli.KeyManager().Address();
	balance := ethRPC.EthCli.StateManager().CurrentState().GetBalance(myAddr)
	reply.Ether = balance.String()

	return nil
}

func (ethRPC *EthereumRPC) BlockByHash(r *http.Request, args *VString, reply *BlockData) error {
	// Get the block.
	bts, err := hex.DecodeString(args.SVal)
	if err != nil {
		fmt.Println(err.Error())
		reply.Hash = ERR_MALFORMED_TX_HASH
		return nil
	}

	block := ethRPC.EthCli.BlockChain().GetBlock(bts)

	if block == nil {
		reply.Hash = ERR_NO_SUCH_BLOCK
		return nil
	}

	getBlockDataFromBlock(reply,block)

	return nil
}

func (ethRPC *EthereumRPC) GetTransaction(r *http.Request, args *TransactionArgs, reply *Transaction) error {
	// Get the block.
	bts, err := hex.DecodeString(args.BlockHash)
	if err != nil {
		reply.Recipient = ERR_MALFORMED_ADDRESS
		return nil
	}

	block := ethRPC.EthCli.BlockChain().GetBlock(bts)

	if block == nil {
		reply.Recipient = ERR_NO_SUCH_BLOCK
		return nil
	}

	txStr, txErr := hex.DecodeString(args.TxHash)
	if txErr != nil {
		reply.Recipient = ERR_MALFORMED_TX_HASH
	}
	tx := block.GetTransaction(txStr)
	if tx == nil {
		reply.Recipient = ERR_NO_SUCH_TX
	}
	transactionFromTx(reply,tx)
	reply.BlockHash = args.BlockHash

	return nil
}

func (ethRPC *EthereumRPC) Transact(r *http.Request, args *TxIndata, reply *TxReceipt) error {
	// TODO check sender.
	err := createTx(ethRPC.EthCli, args.Recipient, args.Value, args.Gas, args.GasCost, args.Data, reply)
	if err != nil {
		reply.Error = err.Error()
	}
	return nil
}

func (ethRPC *EthereumRPC) Account(r *http.Request, args *VString, reply *Account) error {
	getAccount(ethRPC.EthCli,args.SVal,reply)
	return nil
}
*/
