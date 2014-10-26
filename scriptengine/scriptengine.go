package scriptengine

import (
	"encoding/hex"
	"fmt"
	"github.com/eris-ltd/thelonious/ethchain"
	"github.com/eris-ltd/thelonious/monk"
	"github.com/obscuren/sha3"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"strconv"
)

// ScriptEngine
type ScriptEngine struct {
	otto     *otto.Otto
	ethChain *monk.EthChain
	genDoug  string // Gendoug address
	// TODO change when using a name->hash contract.
	models map[string]string
}

func NewScriptEngine(ec *monk.EthChain) *ScriptEngine {
	vm := otto.New()
	se := &ScriptEngine{otto: vm, ethChain: ec}
	se.models = make(map[string]string)
	se.init()
	return se
}

// This code should be used when adding objects and stuff that
// will be used when doing calls.
func (se *ScriptEngine) init() {
	// Address of gendoug.
	gd := hex.EncodeToString(ethchain.GENDOUG)
	se.genDoug = gd
	se.otto.Set("GENDOUG", gd)

	// Define GetStorageAt(account,addr), so that it is accessible from inside the vm.
	se.otto.Set("GetStorageAt", func(call otto.FunctionCall) otto.Value {
		account, err0 := call.Argument(0).ToString()
		if err0 != nil {
			return otto.UndefinedValue()
		}
		
		address, err1 := call.Argument(1).ToString()
		if err1 != nil {
			return otto.UndefinedValue()
		}
		
		ret := se.ethChain.GetStorageAt(account, address)
		if ret != "0x" {
			ret = "0x" + ret
		}
		result, err := se.otto.ToValue(ret)

		if err != nil {
			return otto.UndefinedValue()
		}
		return result
	})
	
	// Inject the math stuff.
	InjectSMath(se.otto)
}

// Later, the hash will be checked against a hash in a contract.
func (se *ScriptEngine) LoadModelFromFile(name, fileName string) {
	bytes, lErr := se.LoadJSFile(fileName)
	if lErr != nil {
		fmt.Printf(lErr.Error())
		return
	}
	_ , err := se.otto.Run(bytes)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	
	d := sha3.NewKeccak256()
	d.Write(bytes)
	sha := hex.EncodeToString(d.Sum(nil))
	se.models[name] = sha
	fmt.Printf("Script hash: %s\n", sha)
	
	se.otto.Run("actionModels['" + sha + "'] = new Model()")
	se.otto.Set("Model",nil)
}

// Get an action model hash from its name.
func (se *ScriptEngine) GetModelHash(modelName string) otto.Value {
	// TODO make sure that each contract stores the action model name in the same place.
	//modelName := se.ethChain.GetStorageAt(account, "0x19")
	
	// TODO this would normally be a call to the model name->hash contract.
	modelHash, ok := se.models[modelName]
	if !ok {
		// TODO Handle. This is always an error.
		return otto.UndefinedValue()
	}
	
	mh, _ := se.otto.ToValue(modelHash)
	
	return mh
}

// Run a command using the c3d tree parser.
// path is the tree path. null, or [] would be genDoug. ["dougId"] doug, and so forth.
// an issue comment would be: ["dougId","orgId","repoId","issueId","commentId"]
// fName is the name of the function, or action, that should be taken.
// params are the function parameters.
// The result is returned as an array of strings.
func (se *ScriptEngine) RunAction(path []string, fName string, param interface{}) ([]string, error) {

	var prm string
	var errPConv error
	if param != nil {
		// Convert params to a string
		prm, errPConv = se.convertParam(param)

		if errPConv != nil {
			return nil, errPConv
		}
	} else {
		prm = "null"
	}

	pt := se.convertPath(path);

	// TODO add the interrupt channel and code.
	result, err := se.otto.Run("treeParser.run(" + pt + ",'" + fName + "'," + prm + ");")
	if err != nil {
		return nil, err
	}

	// If the return value is nil, that indicates a "non success"
	if result == otto.NullValue() {
		return nil, nil
	}

	// Take the result and turn it into a go value.
	obj, expErr := result.Export()
	if expErr != nil {
		return nil, expErr
	}
	return se.convertObj(obj)
}

// Convert. We allow response to be a string, boolean, int, or an array.
// The array may be of mixed types, so long as they are strings, booleans or ints.
func (se *ScriptEngine) convertObj(obj interface{}) ([]string, error) {

	switch val := obj.(type) {
	case string:
		return []string{val}, nil
	case int:
		return []string{strconv.Itoa(val)}, nil
	case bool:
		if val == true {
			return []string{"true"}, nil
		} else {
			return []string{"false"}, nil
		}
	case []interface{}:
		strArr := []string{}
		for idx, v := range val {
			var s string
			switch u := v.(type) {
			case string:
				s = u
				break
			case int:
				s = strconv.Itoa(u)
				break
			case bool:
				if u == true {
					s = "true"
				} else {
					s = "false"
				}
				break
			default:
				return nil, fmt.Errorf("Error in return value: not a string, boolean or int. Idx: %d, Val: %v", idx, u)
			}
			strArr = append(strArr, s)
		}
		return strArr, nil
	default:
		return nil, fmt.Errorf("Error: Return-value not a string, bool, int or array. Val: %v", val)
	}
	// This never happens.
	return nil, nil
}

// Makes a string array into a string str = "['arr[0]','arr[1]',...]"
func (se *ScriptEngine) convertPath(path []string) string {
	if path == nil || len(path) == 0 {
		return "[]";
	}
	pt := "["
	
	for _, s := range path {
		pt = pt + "'" + s + "',"
	}
	// Shave off the last ','
	pt = pt[0 : len(pt)-1]
	pt += "]"
	return pt
}

// Convert. We allow indata to be a string, boolean, int, or an array.
// The array may be of mixed types, so long as they are strings, booleans or ints.
func (se *ScriptEngine) convertParam(param interface{}) (string, error) {

	switch val := param.(type) {
	case string:
		return "'" + val + "'", nil
	case int:
		return strconv.Itoa(val), nil
	case bool:
		if val == true {
			return "true", nil
		} else {
			return "false", nil
		}
	case []interface{}:
		prm := "["
		if len(val) > 0 {
			for idx, v := range val {
				var s string
				switch u := v.(type) {
				case string:
					s = "'" + u + "'"
					break
				case int:
					s = strconv.Itoa(u)
					break
				case bool:
					if u == true {
						s = "true"
					} else {
						s = "false"
					}
					break
				default:
					return "", fmt.Errorf("Error in params: not a string, boolean or int. Idx: %d, Val: %v", idx, u)
				}
				prm = prm + s + ","
			}
			// Shave off the last ','
			prm = prm[0 : len(prm)-1]
		}
		return prm, nil
	default:
		return "", fmt.Errorf("Error: Params not a string, bool, int or array. Val: %v", val)
	}
	// This never happens.
	return "", nil
}

// Loads and compiles javascript.
// TODO Add a folder where scripts are supposed to be.
func (se *ScriptEngine) LoadJSFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return bytes, nil
}

// Use this to set up a new runtime. Should re-do init() and load every
// model from the model->hash contract.
func (se *ScriptEngine) Recover() {

}
