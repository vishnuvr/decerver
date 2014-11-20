package ate

import (
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/events"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"strings"
)

type AteEventProcessor struct {
	er events.EventRegistry
}

type Ate struct {
	runtimes map[string]*JsRuntime
	apis     map[string]interface{}
	er       events.EventRegistry
}

func NewAte(er events.EventRegistry) *Ate {
	return &Ate{make(map[string]*JsRuntime), make(map[string]interface{}), er}
}

func (ate *Ate) ShutdownRuntimes() {
	for _, rt := range ate.runtimes {
		rt.Shutdown()
	}
}

func (ate *Ate) CreateRuntime(name string) core.Runtime {
	rt := newJsRuntime(name, ate.er)
	ate.runtimes[name] = rt
	for k, v := range ate.apis {
		// TODO error checking
		rt.BindScriptObject(k, v)
	}
	fmt.Printf("Regging new runtime: " + name)
	fmt.Printf("Runtimes: %v\n", ate.runtimes)
	return rt
}

func (ate *Ate) GetRuntime(name string) core.Runtime {
	fmt.Println(name)
	fmt.Printf("Ate: %v\n", ate)
	return ate.runtimes[name]
}

func (ate *Ate) RemoveRuntime(name string) {
	ate.runtimes[name] = nil
}

func (ate *Ate) RegisterApi(name string, api interface{}) {
	ate.apis[name] = api
}

type JsRuntime struct {
	vm        *otto.Otto
	subChan   chan events.Event
	closeChan chan bool
	er        events.EventRegistry
	name      string
}

func newJsRuntime(name string, er events.EventRegistry) *JsRuntime {
	vm := otto.New()
	jsr := &JsRuntime{}
	jsr.vm = vm
	jsr.subChan = make(chan events.Event)
	jsr.er = er
	jsr.name = name
	jsr.Init()
	return jsr
}

func (jsr *JsRuntime) Shutdown() {
	fmt.Println("Runtime shut down: " + jsr.name)
	jsr.closeChan <- true
}

// TODO set up the interrupt channel.
func (jsr *JsRuntime) Init() {
	BindDefaults(jsr.vm)
}

func (jsr *JsRuntime) LoadScriptFile(fileName string) error {
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	_, err = jsr.vm.Run(bytes)

	return err
}

func (jsr *JsRuntime) LoadScriptFiles(fileName ...string) error {
	for _, sf := range fileName {
		err := jsr.LoadScriptFile(sf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (jsr *JsRuntime) BindScriptObject(name string, val interface{}) error {
	return jsr.vm.Set(name, val)
}

func (jsr *JsRuntime) AddScript(script string) error {
	_, err := jsr.vm.Run(script)
	return err
}

func (jsr *JsRuntime) RunFunction(funcName string, params []string) (interface{}, error) {

	cmd := funcName + "("

	paramStr := ""
	for _, p := range params {
		paramStr += p + ","
	}
	paramStr = strings.Trim(paramStr, ",")
	cmd += paramStr + ");"

	fmt.Println("Running function: " + cmd)
	val, runErr := jsr.vm.Run(cmd)

	if runErr != nil {
		return nil, fmt.Errorf("Error when running function '%s': %s\n", funcName, runErr.Error())
	}

	// Take the result and turn it into a go value.
	obj, expErr := val.Export()

	if expErr != nil {
		return nil, fmt.Errorf("Error when exporting returned value: %s\n", expErr.Error())
	}

	return obj, nil
}

func (jsr *JsRuntime) CallFuncOnObj(objName, funcName string, param ...interface{}) (interface{}, error) {
	ob, err := jsr.vm.Get(objName)
	if err != nil {
		fmt.Println(err.Error())
	}
	val, callErr := ob.Object().Call(funcName, param...)

	if callErr != nil {
		fmt.Println(callErr.Error())
	}
	// Take the result and turn it into a go value.
	obj, expErr := val.Export()

	if expErr != nil {
		return nil, fmt.Errorf("Error when exporting returned value: %s\n", expErr.Error())
	}

	return obj, nil
}

func (jsr *JsRuntime) CallFunc(funcName string, param ...interface{}) (interface{}, error) {
	val, callErr := jsr.vm.Call(funcName, nil, param)

	if callErr != nil {
		fmt.Println(callErr.Error())
		return nil, callErr
	}

	fmt.Printf("%v\n", val)

	// Take the result and turn it into a go value.
	obj, expErr := val.Export()

	if expErr != nil {
		return nil, fmt.Errorf("Error when exporting returned value: %s\n", expErr.Error())
	}

	return obj, nil
}

// Use this to set up a new runtime. Should re-do init().
// TODO implement
func (jsr *JsRuntime) Recover() {
}

func (jsr *JsRuntime) Channel() chan events.Event {
	return jsr.subChan
}

func (jsr *JsRuntime) Id() string {
	return "Ate"
}

func (jsr *JsRuntime) Source() string {
	return "*"
}

func (jsr *JsRuntime) Close() {

}
