package deCerver

import (
	"github.com/eris-ltd/deCerver-interfaces/modules"
	"github.com/eris-ltd/deCerver/ate"
	"github.com/eris-ltd/deCerver/server"
	"os"
	"os/signal"
)

type Paths struct {
	root        string
	databases   string
	fileSystems string
	log         string
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

type DeCerver struct {
	config    *DCConfig
	logSys    *LogSystem
	paths     *Paths
	ate       *ate.Ate
	webServer *server.WebServer
}

func NewDeCerver() *DeCerver {
	dc := &DeCerver{}
	return dc
}

func (dc *DeCerver) Init() {
	dc.ReadConfig("")
	dc.initPaths()
	dc.initLogSystem()
	dc.logSys.DCLogger.Print("Logger set.")
	dc.initNetwork()
	dc.initAte()
}

func (dc *DeCerver) Run() {
	dc.webServer.Start()
	dc.logSys.DCLogger.Print("Server started.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	logger.Println("Shutting down")
}

func (dc *DeCerver) initPaths() {
	dc.paths = &Paths{}
	dc.paths.root = dc.config.RootDir
	InitDirs(dc.paths.root)
	dc.paths.log = dc.paths.root + "/logs"
	InitDirs(dc.paths.log)
	dc.paths.databases = dc.paths.root + "/databases"
	InitDirs(dc.paths.databases)
	dc.paths.fileSystems = dc.paths.root + "/filesystems"
	InitDirs(dc.paths.fileSystems)
}

func (dc *DeCerver) initNetwork() {
	dc.webServer = server.NewWebServer(uint32(dc.config.MaxClients), logger)
}

func (dc *DeCerver) initAte() {
	dc.ate = ate.NewAte(logger)
}

func (dc *DeCerver) AddModule(id string, md modules.Module) {
	md.Init(dc.ate)
	dc.logSys.AddLogger(id, md.Logger())
	if md.HttpAPIServices() != nil {
		for _, sv := range md.HttpAPIServices() {
			dc.webServer.RegisterHttpAPIService(sv)	
		}
	}
	if md.WsAPIServiceFactories() != nil {
		for _, sf := range md.WsAPIServiceFactories() {
			dc.webServer.RegisterWsAPIServiceFactory(sf)	
		}
	}
	dc.ate.AddModule(id, md)
}