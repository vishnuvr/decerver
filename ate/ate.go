package ate

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/events"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"strings"
	"sync"
)

type AteEventProcessor struct {
	er events.EventRegistry
}

type JsObj struct {
	Name   string
	Object interface{}
}

type Ate struct {
	runtimes  map[string]*JsRuntime
	apiObjs   []*JsObj
	apiScript []string
	er        events.EventRegistry
}

func NewAte(er events.EventRegistry) *Ate {
	return &Ate{
		make(map[string]*JsRuntime),
		make([]*JsObj,0),
		make([]string,0),
		er,
	}
}

func (ate *Ate) ShutdownRuntimes() {
	for _, rt := range ate.runtimes {
		rt.Shutdown()
	}
}

func (ate *Ate) CreateRuntime(name string) core.Runtime {
	rt := newJsRuntime(name, ate.er)
	ate.runtimes[name] = rt
	rt.jsrEvents = NewJsrEvents(rt)
	// TODO add a "runtime" or "os" object with more stuff in it?

	rt.Init(name)
	for _, jo := range ate.apiObjs {
		err := rt.BindScriptObject(jo.Name, jo.Object)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	for _, s := range ate.apiScript {
		err := rt.AddScript(s)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Printf("Creating new runtime: " + name)
	// DEBUG
	fmt.Printf("Runtimes: %v\n", ate.runtimes)
	return rt
}

func (ate *Ate) GetRuntime(name string) core.Runtime {
	rt, ok := ate.runtimes[name]
	if ok {
		return rt
	} else {
		return nil
	}
}

func (ate *Ate) RemoveRuntime(name string) {
	rt, ok := ate.runtimes[name]
	if ok {
		delete(ate.runtimes,name)
		rt.Shutdown()
	}
}

func (ate *Ate) RegisterApiObject(objectname string, api interface{}) {
	ate.apiObjs = append(ate.apiObjs, &JsObj{objectname, api})
}

func (ate *Ate) RegisterApiScript(script string) {
	ate.apiScript = append(ate.apiScript, script)
}

type JsRuntime struct {
	vm        *otto.Otto
	er        events.EventRegistry
	name      string
	jsrEvents *JsrEvents
	mutex     *sync.Mutex
	lockLvl   int
}

func newJsRuntime(name string, er events.EventRegistry) *JsRuntime {
	vm := otto.New()
	jsr := &JsRuntime{}
	jsr.vm = vm
	jsr.er = er
	jsr.name = name
	jsr.mutex = &sync.Mutex{}
	return jsr
}

func (jsr *JsRuntime) Shutdown() {
	fmt.Println("Runtime shut down: " + jsr.name)
	// TODO implement
}

// TODO set up the interrupt channel.
func (jsr *JsRuntime) Init(name string) {
	jsr.vm.Set("jsr_events", jsr.jsrEvents)
	jsr.BindScriptObject("RuntimeId", name)
	BindDefaults(jsr)
}

func (jsr *JsRuntime) LoadScriptFile(fileName string) error {
	jsr.mutex.Lock()
	defer jsr.mutex.Unlock()
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
	jsr.mutex.Lock()
	defer jsr.mutex.Unlock()
	err := jsr.vm.Set(name, val)
	return err
}

func (jsr *JsRuntime) AddScript(script string) error {
	jsr.mutex.Lock()
	defer jsr.mutex.Unlock()
	_, err := jsr.vm.Run(script)
	return err
}

func (jsr *JsRuntime) RunFunction(funcName string, params []string) (interface{}, error) {
	jsr.mutex.Lock()
	defer jsr.mutex.Unlock()
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
	jsr.mutex.Lock()
	defer jsr.mutex.Unlock()
	ob, err := jsr.vm.Get(objName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	
	val, callErr := ob.Object().Call(funcName, param...)

	if callErr != nil {
		fmt.Println(callErr.Error())
		return nil, err
	}

	// Take the result and turn it into a go value.
	obj, expErr := val.Export()

	if expErr != nil {
		return nil, fmt.Errorf("Error when exporting returned value: %s\n", expErr.Error())
	}
	return obj, nil
}

func (jsr *JsRuntime) CallFunc(funcName string, param ...interface{}) (interface{}, error) {
	jsr.mutex.Lock()
	defer jsr.mutex.Unlock()
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

// Used to call the event processor from inside the javascript vm
type JsrEvents struct {
	jsr *JsRuntime
}

func NewJsrEvents(jsr *JsRuntime) *JsrEvents {
	return &JsrEvents{jsr}
}

func (jsre *JsrEvents) Subscribe(evtSource, evtType, evtTarget, subId string) {
	sub := NewAteSub(evtSource, evtType, evtTarget, subId, jsre.jsr)
	jsre.jsr.er.Subscribe(sub)
	// Launch the sub channel.
	go func(s *AteSub) {
		// DEBUG
		fmt.Println("Starting event loop for atesub: " + s.id)
		for {
			evt, ok := <-s.eventChan
			if !ok {
				fmt.Println("[Atë] Close message received.")
				return
			}
			fmt.Println("[Atë] stuff coming in from event processor: " + evt.Event)

			jsonString, err := json.Marshal(evt)
			// _ , err := json.Marshal(evt)
			if err != nil {
				fmt.Println("Error when posting event to ate: " + err.Error())
			}
			s.rt.CallFuncOnObj("events", "post", string(jsonString))
		}
	}(sub)
}

func (jsre *JsrEvents) Unsubscribe(subId string) {
	jsre.jsr.er.Unsubscribe(subId)
}

type AteSub struct {
	eventChan chan events.Event
	closeChan chan bool
	source    string
	tpe       string
	tgt       string
	id        string
	rt        core.Runtime
}

func NewAteSub(eventSource, eventType, eventTarget, subId string, rt core.Runtime) *AteSub {
	as := &AteSub{}
	// TODO Optimize later.
	as.eventChan = make(chan events.Event, 10)
	as.closeChan = make(chan bool)
	as.source = eventSource
	as.tpe = eventType
	as.tgt = eventTarget
	as.id = subId
	as.rt = rt
	return as
}

func (as *AteSub) Channel() chan events.Event {
	return as.eventChan
}

func (as *AteSub) SetChannel(ec chan events.Event) {
	as.eventChan = ec
}

func (as *AteSub) Source() string {
	return as.source
}

func (as *AteSub) Id() string {
	return as.id
}

func (as *AteSub) Target() string {
	return as.tgt
}

func (as *AteSub) Event() string {
	return as.tpe
}
