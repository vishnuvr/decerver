package ate

import (
	"fmt"
	"github.com/eris-ltd/deCerver-interfaces/modules"
	"github.com/eris-ltd/deCerver-interfaces/core"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"strconv"
)

var logger *log.Logger

type Ate struct {
	vm      *otto.Otto
	modules map[string]modules.Module
}

func NewAte(lgr *log.Logger) *Ate {
	logger = lgr
	vm := otto.New()
	ate := &Ate{}
	ate.vm = vm
	ate.modules = make(map[string]modules.Module)
	ate.init()
	return ate
}

func (ate *Ate) ShutDown() {
	logger.Print("Atë shut down.")
}

// Initialize the vm. Add some helper functions and other things.
// TODO set up the interrupt channel.
func (ate *Ate) init() {
	LoadHelpers(ate.vm)
	logger.Print("Atë started")
}

func (ate *Ate) AddModule(id string, md modules.Module) {
	ate.modules[id] = md
}

func (ate *Ate) LoadScript(fileName string) {

}

func (ate *Ate) InjectFunction(fName string, fun core.AteFunc){
	ate.vm.Set(fName,fun)
}

func (ate *Ate) RunAction(path []string, actionName string, params interface{}) []string {
	return nil
}

func (ate *Ate) RunMethod(nameSpace, funcName string, params interface{}) []string {
	
	var prm string
	var errPConv error

	if params != nil {
		// Convert params to a string
		prm, errPConv = ate.convertParam(params)

		if errPConv != nil {
			return nil
		}
	} else {
		prm = "null"
	}
	val, _ := ate.vm.Run(funcName + "(" + prm + ")")

	// Take the result and turn it into a go value.
	obj, expErr := val.Export()

	if expErr != nil {
		return nil
	}

	ret, _ := ate.convertObj(obj)

	return ret
}

// Convert. We allow response to be a string, boolean, int, or an array.
// The array may be of mixed types, so long as they are strings, booleans or ints.
func (ate *Ate) convertObj(obj interface{}) ([]string, error) {

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
func (ate *Ate) convertPath(path []string) string {
	if path == nil || len(path) == 0 {
		return "[]"
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
func (ate *Ate) convertParam(param interface{}) (string, error) {

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
func (ate *Ate) LoadJSFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}
	return bytes, nil
}

// Use this to set up a new runtime. Should re-do init() and load every
// model from the model->hash contract.
func (ate *Ate) Recover() {

}
