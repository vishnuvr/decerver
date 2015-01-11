package scripting

import ()

type(
	// This is the interface for the javascript runtime manager, or 'Atë'.
	RuntimeManager interface {
		GetRuntime(string) Runtime
		CreateRuntime(string) Runtime
		RemoveRuntime(string)
		RegisterApiObject(string, interface{})
		RegisterApiScript(string)
		ShutdownRuntimes()
	}

	// This is the interface for a javascript runtime.
	Runtime interface {
		Init(string)
		Shutdown()
		BindScriptObject(name string, val interface{}) error
		LoadScriptFile(fileName string) error
		LoadScriptFiles(fileName ...string) error
		AddScript(script string) error
		CallFunc(funcName string, param ...interface{}) (interface{}, error)
		CallFuncOnObj(objName, funcName string, param ...interface{}) (interface{}, error)
	}
)