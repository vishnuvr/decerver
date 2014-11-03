package deCerver

import (
	"github.com/eris-ltd/deCerver-interfaces/modules"
	"github.com/eris-ltd/deCerver/ate"
	"github.com/eris-ltd/deCerver/server"
	"github.com/golang/glog"
	"os"
	"os/signal"
)

type Paths struct {
	root        string
	databases   string
	fileSystems string
	log         string
	apps        string
}

func (p *Paths) Root() string {
	return p.root
}

func (p *Paths) Databases() string {
	return p.databases
}

func (p *Paths) FileSystems() string {
	return p.fileSystems
}

func (p *Paths) Log() string {
	return p.log
}

func (p *Paths) Apps() string {
	return p.apps
}

type DeCerver struct {
	config    *DCConfig
	paths     *Paths
	ate       *ate.Ate
	webServer *server.WebServer
	modules map[string]modules.Module
}

func NewDeCerver() *DeCerver {
	dc := &DeCerver{}
	dc.modules = make(map[string]modules.Module)
	return dc
}

func (dc *DeCerver) Init() {
	
	glog.Infoln("Initializing decerver")
	dc.ReadConfig("")
	dc.initPaths()
	dc.initNetwork()
	dc.initAte()
	
	for _ , mod := range dc.modules {
		glog.Infof("Initializing module: %s\n",mod.Name())
		mod.Init()
	}
	
}

func (dc *DeCerver) Start() {
	dc.webServer.Start()
	glog.Infof("Server started.")
	
	for _ , mod := range dc.modules {
		glog.Infof("Starting module: %s\n",mod.Name())
		mod.Start()
	}
	
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	glog.Infof("Shutting down")
}

func (dc *DeCerver) initPaths() {
	dc.paths = &Paths{}
	dc.paths.root = dc.config.RootDir
	InitDir(dc.paths.root)
	dc.paths.log = dc.paths.root + "/logs"
	InitDir(dc.paths.log)
	dc.paths.databases = dc.paths.root + "/databases"
	InitDir(dc.paths.databases)
	dc.paths.fileSystems = dc.paths.root + "/filesystems"
	InitDir(dc.paths.fileSystems)
	dc.paths.apps = dc.paths.root + "/apps"
	InitDir(dc.paths.apps)
}

func (dc *DeCerver) initNetwork() {
	dc.webServer = server.NewWebServer(uint32(dc.config.MaxClients),dc.paths.Apps())
}

func (dc *DeCerver) initAte() {
	dc.ate = ate.NewAte()
}

func (dc *DeCerver) AddModule(id string, md modules.Module) {
	md.Register(nil,dc.webServer,dc.ate)
	dc.modules[md.Name()] = md
	glog.Infof("Registering module '%s'.",md.Name())
}