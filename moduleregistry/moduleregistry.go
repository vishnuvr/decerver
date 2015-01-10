package moduleregistry

import (
	"errors"
	"fmt"
	"github.com/eris-ltd/decerver/interfaces/modules"
	"time"
)

// The module registry is where the different modules are kept. Currently, modules has
// to be loaded upon startup, and cannot be unloaded.
type ModuleRegistry struct {
	modules     map[string]modules.Module
	moduleNames []string
}

func NewModuleRegistry() *ModuleRegistry {
	mr := &ModuleRegistry{}
	mr.modules = make(map[string]modules.Module, 1)
	mr.moduleNames = make([]string, 1)
	return mr
}

func (mr *ModuleRegistry) GetModules() map[string]modules.Module {
	return mr.modules
}

func (mr *ModuleRegistry) GetModuleNames() []string {
	return mr.moduleNames
}

func (mr *ModuleRegistry) Add(m modules.Module) error {
	// The name cannot already be taken.
	mod := mr.modules[m.Name()]
	if mod != nil {
		str := "Module '" + m.Name() + "' has already been registered."
		return errors.New(str)
	}
	mr.moduleNames = append(mr.moduleNames, m.Name())
	mr.modules[m.Name()] = m
	return nil
}

func (mr *ModuleRegistry) Init() error {
	for _, md := range mr.modules {
		err := md.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO Sync modules better. This stuff really has to be put in the module 
// wrapper layer. This is a high priority issue.
func (mr *ModuleRegistry) Start() error {
		for _, mod := range mr.modules {
			
			go func() {
				fmt.Println("Loading module: " + mod.Name())
				mod.Start()
			}()
			// TODO get channels from module glue instead.
			time.Sleep(300);
		}
	return nil
}

func (mr *ModuleRegistry) Shutdown() error {
	for _, mod := range mr.modules {
		mod.Shutdown()
	}
	return nil
}
