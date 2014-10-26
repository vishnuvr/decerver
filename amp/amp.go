package amp

import (
	"encoding/json"
	"fmt"
	//"github.com/eris-ltd/eth-go-mods/ethchain"
	"github.com/eris-ltd/thelonious/ethreact"
	"github.com/eris-ltd/thelonious/monk"
	"time"
)

const TIMEOUT = 10 * time.Second
const MAX_BLOCKS = 1

// Placeholder that can be used to get an action model (as a byte array) from its id.
var ActionModels map[string][]byte

type ActionModel struct {
	Actions map[string]*Action `json:"actions"`
	Data map[string]string	   `json:"data"`
}

type Action struct {
	PreCall  []string   `json:"precall"`
	Call     [][]string `json:"call"`
	PostCall []string   `json:"postcall"`
	Success  string     `json:"success"`
	Result   []string   `json:"result"`
}

type Result struct {
	Success bool
	Result  []string
	Error   string
}

type ActionModelParser struct {
	ethChain *monk.EthChain
	contract string
	model    *ActionModel
}

func NewActionModelParser(ethChain *monk.EthChain) *ActionModelParser {
	amp := &ActionModelParser{}
	amp.ethChain = ethChain
	return amp
}

func (amp *ActionModelParser) Initialize(contract string, modelName string) error {
	amp.contract = contract
	model, err := amp.getTheModel(modelName)
	if err != nil {
		return err
	}
	amp.model = model
	return nil
}

// Just fetch from that placeholder map for now.
func (amp *ActionModelParser) getTheModel(modelName string) (*ActionModel, error) {
	modelId := modelName //amp.ethChain.GetStorageAt(amp.contract, "0x19")
	if modelId == "0x" || modelId == "0" || modelId == "" {
		return nil, fmt.Errorf("Storage at contract address 0x19 is 0: ", amp.contract)
	}

	if ampBytes, ok := ActionModels[modelId]; ok {
		actionModel := &ActionModel{}
		json.Unmarshal(ampBytes, actionModel)
		return actionModel, nil
	} else {
		return nil, fmt.Errorf("No actionmodel with id: %s", modelId)
	}
}

func (amp *ActionModelParser) PerformAction(action string, params []string) *Result {

	res := &Result{}
	fmt.Printf("%v\n",amp.model.Actions["deadbeef"])
	if actn, ok := amp.model.Actions[action]; ok {
		
		// Set globals and params
		globals := fakeGlobals(amp.contract)
		params := fakeParams(nil,nil)
	
		// Creates a new parser
		parser := NewParser(amp.ethChain,globals,params)
		
		// Precall
		errPre := parser.ParsePP(actn.PreCall)
		
		if errPre != nil {
			res.Error = errPre.Error()
			return res
		}
		
		// No txs
		_ , errCall := parser.ParseCall(actn.Call)
		
		if errCall != nil {
			res.Error = errCall.Error()
			return res
		}
		
		if len(actn.Call) > 0 {
			// Now we wait for the transactions to be mined into a block.
			// Notice this is not 100% safe, as some (or all) of the transactions
			// may fail. 
			//
			// Dropping the txs list since monk.Msg does not return one (TODO remind).
			blockChannel := make(chan ethreact.Event)
			amp.ethChain.Ethereum.Reactor().Subscribe("newBlock", blockChannel)
			blocks := 0
			
			loop:
			for {
				select {
				case <-blockChannel:
					/*
					block, _ := evt.(*ethchain.Block)
					for _, tx := range block.Transactions() {
						hash := tx.Hash()
						if val, ok := amp.txs[hash]; ok {
							delete(amp.txs, hash)
						}
					}
					*/
					blocks++
					if blocks == MAX_BLOCKS {
						break
					}
					/*
					if len(amp.txs) == 0 || blocks == MAX_BLOCKS {
						break
					}
					*/
				case <-time.After(TIMEOUT):
					fmt.Println("Timed out waiting for block to mine.")
					break loop
				}
			}
		}
		
		errPost := parser.ParsePP(actn.PostCall)
		
		if errPost != nil {
			res.Error = errPost.Error()
			return res
		}
		
		succ, errSucc := parser.ParseSuccess(actn.Success)
		
		if errSucc != nil {
			res.Error = errSucc.Error()
			return res
		}
		
		// Running out of words
		rrrrr, errR := parser.ParseResult(actn.Result)
		
		if errR != nil {
			res.Error = errR.Error()
			return res
		}
		
		res.Success = succ
		res.Result = rrrrr
		
		fmt.Printf("Precall: %v\n",parser.precall)
		fmt.Printf("Postcall: %v\n",parser.postcall)
		
	} else {
		res.Success = false
		res.Error = "No such action: " + action
	}
	return res
}

func (amp *ActionModelParser) Error(str string) string {
	return fmt.Sprintf("Error: %s\nACTION MODEL DUMP:\n%v\n", str, amp.model)
}

func fakeGlobals(contract string) map[string]string {
	globals := make(map[string]string)
	globals["doug"] = "0x0000000000000000000000000000000000000000"
	globals["gendoug"] = "0x1111111111111111111111111111111111111111"
	globals["this"] = contract
	return globals
}

func fakeParams(keys, values []string) map[string]string {
	ps := make(map[string]string)
	if keys == nil || len(keys) == 0 {
		return ps
	}
	for idx , key := range keys {
		ps[key] = values[idx]
	}
	return ps
}
