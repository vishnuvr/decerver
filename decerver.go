package decerver

import (
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/modules"
	"github.com/eris-ltd/decerver/ate"
	"github.com/eris-ltd/decerver/dappregistry"
	"github.com/eris-ltd/decerver/events"
	"github.com/eris-ltd/decerver/moduleregistry"
	"github.com/eris-ltd/decerver/server"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
)

var logger *log.Logger = core.NewLogger("Decerver Core")

type Paths struct {
	mutex       *sync.Mutex
	root        string
	modules     string
	log         string
	blockchains string
	filesystems string
	dapps       string
	system      string
	adminpages  string
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

func (p *Paths) Dapps() string {
	return p.dapps
}

func (p *Paths) Blockchains() string {
	return p.blockchains
}

func (p *Paths) Filesystems() string {
	return p.filesystems
}

func (p *Paths) System() string {
	return p.system
}

func (p *Paths) Adminpages() string {
	return p.system
}

// Thread safe read file function. Reads an entire file and returns the bytes.
func (p *Paths) ReadFile(directory, name string) ([]byte, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return ioutil.ReadFile((path.Join(directory, name)))
}

// Thread safe write file function. Writes the provided byte slice into the file 'name'
// in directory 'directory'. Uses filemode 0600.
func (p *Paths) WriteFile(directory, name string, data []byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return ioutil.WriteFile((path.Join(directory, name)), data, 0600)
}

// Creates a new directory for a module, and returns the path.
func (p *Paths) CreateDirectory(moduleName string) string {
	dir := p.modules + "/" + moduleName
	InitDir(dir)
	return dir
}

type DeCerver struct {
	config         *core.DCConfig
	paths          *Paths
	ep             *events.EventProcessor
	ate            *ate.Ate
	webServer      *server.WebServer
	moduleRegistry *moduleregistry.ModuleRegistry
	dappRegistry   *dappregistry.DappRegistry
	isStarted      bool
}

func NewDeCerver() *DeCerver {
	dc := &DeCerver{}
	logger.Println("Starting decerver bootstrapping sequence.")
	dc.ReadConfig("")
	dc.createPaths()
	dc.WriteConfig(dc.config)
	dc.createModuleRegistry()
	dc.createEventProcessor()
	dc.createAte()
	dc.createNetwork()
	dc.createDappRegistry()
	return dc
}

func (dc *DeCerver) Init() {
	err := dc.moduleRegistry.Init()
	if err != nil {
		logger.Printf("Module failed to initialize: %s. Shutting down.\n", err.Error())
		os.Exit(-1)
	}
	dc.initDapps()
}

func (dc *DeCerver) Start() {
	dc.webServer.Start()
	logger.Println("Server started.")

	err := dc.moduleRegistry.Start()
	if err != nil {
		logger.Printf("Module failed to start: %s. Shutting down.\n", err.Error())
		os.Exit(-1)
	}

	// Now everything is registered.
	dc.isStarted = true

	logger.Println("Running...")
	// Just block for now.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	logger.Println("Shutting down.")
	dc.moduleRegistry.Shutdown()
}

func (dc *DeCerver) createPaths() {
	
	dc.paths = &Paths{}
	dc.paths.mutex = &sync.Mutex{}
	
	dc.paths.root = dc.config.RootDir
	InitDir(dc.paths.root)
	dc.paths.log = dc.paths.root + "/logs"
	InitDir(dc.paths.log)
	dc.paths.modules = dc.paths.root + "/modules"
	InitDir(dc.paths.modules)
	dc.paths.dapps = dc.paths.root + "/dapps"
	InitDir(dc.paths.dapps)
	dc.paths.filesystems = dc.paths.root + "/filesystems"
	InitDir(dc.paths.filesystems)
	dc.paths.blockchains = dc.paths.root + "/blockchains"
	InitDir(dc.paths.blockchains)
	dc.paths.system = dc.paths.root + "/system"
	InitDir(dc.paths.system)
	dc.paths.adminpages = dc.paths.system + "/adminpages"
	InitDir(dc.paths.adminpages)
}

func (dc *DeCerver) createNetwork() {
	dc.webServer = server.NewWebServer(uint32(dc.config.MaxClients), dc.paths, dc.config.Port, dc.ate, dc)
}

func (dc *DeCerver) createEventProcessor() {
	dc.ep = events.NewEventProcessor(dc.moduleRegistry)
}

func (dc *DeCerver) createAte() {
	dc.ate = ate.NewAte(dc.ep)
}

func (dc *DeCerver) IsStarted() bool {
	return dc.isStarted
}

func (dc *DeCerver) createModuleRegistry() {
	dc.moduleRegistry = moduleregistry.NewModuleRegistry()
}

func (dc *DeCerver) createDappRegistry() {
	dc.dappRegistry = dappregistry.NewDappRegistry(dc.ate, dc.webServer, dc.moduleRegistry)
	dc.webServer.AddDappRegistry(dc.dappRegistry)
}

func (dc *DeCerver) LoadModule(md modules.Module) {
	// TODO re-add
	md.Register(dc.paths, dc.ate, dc.ep)
	dc.moduleRegistry.Add(md)
	logger.Printf("Registering module '%s'.\n", md.Name())
}

func (dc *DeCerver) initDapps() {
	err := dc.dappRegistry.RegisterDapps(dc.paths.Dapps(), dc.paths.System())

	if err != nil {
		logger.Println("Error loading dapps: " + err.Error())
		os.Exit(0)
	}
}

func (dc *DeCerver) GetFileIO() core.FileIO {
	return dc.paths
}
