package decerver

import (
	"encoding/json"
	"github.com/eris-ltd/decerver/interfaces/core"
	"github.com/eris-ltd/decerver/interfaces/modules"
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

const version = "0.1.0"

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

// Thread safe read file function. Reads an entire file and returns the bytes.
func (p *Paths) ReadFile(directory, name string) ([]byte, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return ioutil.ReadFile((path.Join(directory, name)))
}

// Thread safe read file function. It'll read the given file and attempt to
// unmarshal it into the provided object.
func (p *Paths) UnmarshalJsonFromFile(directory, name string, object interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	bts, err := ioutil.ReadFile((path.Join(directory, name)))
	if err != nil {
		return err
	}
	return json.Unmarshal(bts, object)
}

// Thread safe write file function. Writes the provided byte slice into the file 'name'
// in directory 'directory'. Uses filemode 0600.
func (p *Paths) WriteFile(directory, name string, data []byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return ioutil.WriteFile((path.Join(directory, name)), data, 0600)
}

// Thread safe write file function. Writes the provided object into a file after
// marshaling it into json. Uses filemode 0600.
func (p *Paths) MarshalJsonToFile(directory, name string, object interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	bts, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile((path.Join(directory, name)), bts, 0600)
}

// Creates a new directory for a module, and returns the path.
func (p *Paths) CreateModuleDirectory(moduleName string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	dir := p.modules + "/" + moduleName
	return initDir(dir)
}

// Helper function to create directories.
func (p *Paths) CreateDirectory(dir string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return initDir(dir)
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
	sig := <-c
	logger.Println("Shutting down: " + sig.String())
	dc.moduleRegistry.Shutdown()
	logger.Println("Bye.")
}

func (dc *DeCerver) createPaths() {

	dc.paths = &Paths{}
	dc.paths.mutex = &sync.Mutex{}

	dc.paths.root = dc.config.RootDir
	initDir(dc.paths.root)
	dc.paths.log = dc.paths.root + "/logs"
	initDir(dc.paths.log)
	dc.paths.modules = dc.paths.root + "/modules"
	initDir(dc.paths.modules)
	dc.paths.dapps = dc.paths.root + "/dapps"
	initDir(dc.paths.dapps)
	dc.paths.filesystems = dc.paths.root + "/filesystems"
	initDir(dc.paths.filesystems)
	dc.paths.blockchains = dc.paths.root + "/blockchains"
	initDir(dc.paths.blockchains)
	dc.paths.system = dc.paths.root + "/system"
	initDir(dc.paths.system)
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
	dc.paths.CreateModuleDirectory(md.Name())
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
