package deCerver

import (
	"github.com/eris-ltd/deCerver-interfaces/modules"
	"github.com/eris-ltd/deCerver/ate"
	"github.com/eris-ltd/deCerver/events"
	"github.com/eris-ltd/deCerver/moduleregistry"
	"github.com/eris-ltd/deCerver/server"
	"github.com/golang/glog"
	"os"
	"os/signal"
)

type Paths struct {
	root    string
	modules string
	log     string
	apps    string
}

func (p *Paths) Root() string {
	return p.root
}

func (p *Paths) Modules() string {
	return p.modules
}

func (p *Paths) Log() string {
	return p.log
}

func (p *Paths) Apps() string {
	return p.apps
}

// Creates a new directory for a module, and returns the path.
func (p *Paths) CreateDirectory(moduleName string) string {
	dir := p.modules + "/" + moduleName
	InitDir(dir)
	return dir
}

type DeCerver struct {
	config         *DCConfig
	paths          *Paths
	ep             *events.EventProcessor
	ate            *ate.Ate
	webServer      *server.WebServer
	moduleRegistry *moduleregistry.ModuleRegistry
}

func NewDeCerver() *DeCerver {
	dc := &DeCerver{}
	glog.Infoln("Creating new decerver")
	dc.ReadConfig("")
	dc.createPaths()
	dc.createNetwork()
	dc.createEventProcessor()
	dc.createAte()
	dc.createModuleRegistry()
	return dc
}

func (dc *DeCerver) Init() {
	err := dc.moduleRegistry.Init()
	if err != nil {
		glog.Fatalf("Module failed to load: %s. Shutting down.", err.Error())
		os.Exit(-1)
	}
}

func (dc *DeCerver) Start() {
	dc.webServer.Start()
	glog.Infof("Server started.")

	err := dc.moduleRegistry.Start()
	if err != nil {
		glog.Fatalf("Module failed to start: %s. Shutting down.", err.Error())
		os.Exit(-1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	glog.Infof("Shutting down")
}

func (dc *DeCerver) createPaths() {
	dc.paths = &Paths{}
	dc.paths.root = dc.config.RootDir
	InitDir(dc.paths.root)
	dc.paths.log = dc.paths.root + "/logs"
	InitDir(dc.paths.log)
	dc.paths.modules = dc.paths.root + "/modules"
	InitDir(dc.paths.modules)
	dc.paths.apps = dc.paths.root + "/apps"
	InitDir(dc.paths.apps)
}

func (dc *DeCerver) createNetwork() {
	dc.webServer = server.NewWebServer(uint32(dc.config.MaxClients), dc.paths.Apps())
}

func (dc *DeCerver) createEventProcessor() {
	dc.ep = events.NewEventProcessor()
}

func (dc *DeCerver) createAte() {
	dc.ate = ate.NewAte(dc.ep)
}

func (dc *DeCerver) createModuleRegistry() {
	dc.moduleRegistry = moduleregistry.NewModuleRegistry()
}

func (dc *DeCerver) AddModule(md modules.Module) {
	md.Register(nil, dc.webServer, dc.ate)
	dc.moduleRegistry.Add(md)
	glog.Infof("Registering module '%s'.", md.Name())
}
