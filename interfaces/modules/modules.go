package modules

import (
	"github.com/eris-ltd/decerver/interfaces/events"
	"github.com/eris-ltd/decerver/interfaces/files"
	"github.com/eris-ltd/decerver/interfaces/types"
)

// typedef for javascript objects.
type JsObject map[string]interface{}

type (
	ModuleInfo struct {
		Name       string      `json:"name"`
		Version    string      `json:"version"`
		Author     *AuthorInfo `json:"author"`
		Licence    string      `json:"licence"`
		Repository string      `json:"repository"`
	}

	AuthorInfo struct {
		Name  string `json:"name"`
		EMail string `json:"e-mail"`
	}
)

type (
	Module interface {
		// For registering with decerver.
		Register(dc DecerverModuleApi) error
		Init() error
		Start() error
		Restart() error
		Shutdown() error
		Name() string
		Subscribe(name, event, target string) error
		UnSubscribe(name string)

		SetProperty(name string, data interface{})
		Property(name string) interface{}
	}

	// Interface for the module manager.
	ModuleManager interface {
		Modules() map[string]Module
		ModuleNames() []string
		Add(m Module) error
		Init() error
		Start() error
		Shutdown() error
	}
	
	// This is the functionality that decerver exports to modules
	// when they register.
	DecerverModuleApi interface {
		// register an object with the script runtime manager (Atë).
		RegisterRuntimeObject(string,interface{})
		// Register script in the form of a string
		RegisterRuntimeScript(string)
		// Post an event.
		PostEvent(events.Event)
		// File and folder management tool.
		FileIO() files.FileIO
	}
)

// Converts a data and an error values into a javascript ready object.

// If there is no error
func JsReturnValNoErr(data interface{}) JsObject {
	ret := make(JsObject)
	ret["Error"] = ""
	ret["Data"] = types.ToJsValue(data)
	ret["Status"] = 0
	return ret
}

// If there is an error
func JsReturnValErr(err error, statusCode int) JsObject {
	ret := make(JsObject)
	ret["Data"] = nil
	ret["Error"] = err.Error()
	ret["Status"] = statusCode
	return ret
}