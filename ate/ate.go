package ate

import (
	"fmt"
	"github.com/eris-ltd/deCerver-interfaces/events"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"strconv"
)

type AteEventProcessor struct {
	er events.EventRegistry
}

type Ate struct {
	vm        *otto.Otto
	subChan   chan events.Event
	closeChan chan bool
}

func NewAte(er events.EventRegistry) *Ate {
	vm := otto.New()
	ate := &Ate{}
	ate.vm = vm
	ate.subChan = make(chan events.Event)

	return ate
}

func (ate *Ate) ShutDown() {
	fmt.Println("Atë shut down.")
	ate.closeChan <- true
}

// Initialize the vm. Add some helper functions and other things.
// TODO set up the interrupt channel.
func (ate *Ate) Init() {
	BindDefaults(ate.vm)
	fmt.Println("Atë started")
}

func (ate *Ate) LoadScriptFile(fileName string) error {
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	_, err = ate.vm.Run(bytes)

	return err
}

func (ate *Ate) LoadScriptFiles(fileName ...string) error {
	for _ , sf := range fileName {
		err := ate.LoadScriptFile(sf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ate *Ate) BindScriptObject(name string, val interface{}) error {
	return ate.vm.Set(name, val)
}

func (ate *Ate) AddScript(script string) error {
	_ , err := ate.vm.Run(script)
	return err
}

func (ate *Ate) RunAction(path []string, actionName string, params interface{}) ([]string, error) {
	return nil, nil
}

func (ate *Ate) RunFunction(funcName string, params interface{}) ([]string, error) {

	var prm string
	var errPConv error

	if params != nil {
		// Convert params to a string
		prm, errPConv = ate.convertParam(params)

		if errPConv != nil {
			return nil, fmt.Errorf("Error when converting parameters: %s\n", errPConv.Error())
		}
	} else {
		prm = "null"
	}
	val, runErr := ate.vm.Run(funcName + "(" + prm + ")")

	if runErr != nil {
		return nil, fmt.Errorf("Error when running function '%s': %s\n", funcName, runErr.Error())
	}

	// Take the result and turn it into a go value.
	obj, expErr := val.Export()

	if expErr != nil {
		return nil, fmt.Errorf("Error when exporting returned value: %s\n", expErr.Error())
	}

	ret, convErr := ate.convertObj(obj)

	if expErr != nil {
		return nil, fmt.Errorf("Error when converting returned value: %s\n", convErr.Error())
	}

	return ret, nil
}

func (ate *Ate) CallFuncOnObj(objName, funcName string, params ... interface{}) {
	val, err := ate.vm.Get(objName)
	if err != nil {
		fmt.Println(err.Error())
	}
	_ , err = val.Object().Call(funcName,params)
	
	if err != nil {
		fmt.Println(err.Error())
	}
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

// Use this to set up a new runtime. Should re-do init().
// TODO implement
func (ate *Ate) Recover() {
	//ate.vm = otto.New()
	//ate.init()
}

func (ate *Ate) Channel() chan events.Event {
	return ate.subChan
}

func (ate *Ate) Id() string {
	return "Ate"
}

func (ate *Ate) Source() string {
	return "*"
}

func (ate *Ate) Close() {
	
}
