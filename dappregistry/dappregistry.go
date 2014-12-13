package dappregistry

import (
	// "path/filepath"
	//"bytes"
	"crypto/sha1"
	//"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/eris-ltd/decerver-interfaces/api"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/dapps"
	"github.com/eris-ltd/decerver-interfaces/modules"
	// "github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"log"
	//"time"
)

var logger *log.Logger = core.NewLogger("Dapp Registry")

// const REG_URL = "http://localhost:9999" 

type Dapp struct {
	models      []string
	path        string
	packageFile *dapps.PackageFile
}

func (dapp *Dapp) GetModels() []string {
	return dapp.models
}

func (dapp *Dapp) GetPath() string {
	return dapp.path
}

func (dapp *Dapp) GetPackageFile() *dapps.PackageFile {
	return dapp.packageFile
}

func NewDapp() *Dapp {
	dapp := &Dapp{models: make([]string,0)}
	return dapp
}

type DappRegistry struct {
	mutex  *sync.Mutex
	keys   map[string]string
	dapps  map[string]*Dapp
	ate    core.RuntimeManager
	server api.Server
//	hashDB *leveldb.DB
	runningDapp *Dapp
	moduleReg modules.ModuleRegistry
}

func NewDappRegistry(ate core.RuntimeManager, server api.Server, mr modules.ModuleRegistry) *DappRegistry {
	dr := &DappRegistry{}
	dr.keys = make(map[string]string)
	dr.dapps = make(map[string]*Dapp)
	dr.mutex = &sync.Mutex{}
	dr.ate = ate
	dr.server = server
	dr.moduleReg = mr
	return dr
}

func (dc *DappRegistry) RegisterDapps(directory, dbDir string) error {
//	dbDir = path.Join(dbDir,"dapp_stored_hashes")
//	dc.hashDB, _ = leveldb.OpenFile(dbDir,nil)
//	defer dc.hashDB.Close()
	logger.Println("Loading dapps")
	files, err := ioutil.ReadDir(directory)
	
	if err != nil {
		return err
	}
	
	if len(files) == 0 {
		logger.Println("No dapps has been downloaded.")
		return nil
	}
	
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			pth := path.Join(directory, fileInfo.Name())
			dc.RegisterDapp(pth)
		}
	}
	logger.Println("Done loading dapps.")
	return nil
}

func (dc *DappRegistry) RegisterDapp(dir string) {
	
	pkDir := path.Join(dir, dapps.PACKAGE_FILE_NAME)
	_, errPfn := os.Stat(pkDir)
	if errPfn != nil {
		logger.Printf("Error loading 'package.json' for dapp '%s'. Skipping...\n", dir)
		logger.Println(errPfn.Error())
		return
	}
	
	pkBts, errP := ioutil.ReadFile(pkDir)

	if errP != nil {
		logger.Printf("Error loading 'package.json' for dapp '%s'. Skipping...\n", dir)
		logger.Println(errP.Error())
		return
	}

	packageFile := &dapps.PackageFile{}
	pkUnmErr := json.Unmarshal(pkBts, packageFile)

	if pkUnmErr != nil {
		logger.Printf("The 'package.json' file for dapp '%s' is corrupted. Skipping...\n", dir)
		logger.Println(pkUnmErr.Error())
		return
	}

	idxDir := path.Join(dir, dapps.INDEX_FILE_NAME)
	_, errIf := os.Stat(idxDir)
	
	if errIf != nil {
		logger.Printf("Cannot find an 'index.html' file for dapp '%s'. Skipping...\n", dir)
		logger.Println(errIf.Error())
		return
	}

	modelDir := path.Join(dir, dapps.MODELS_FOLDER_NAME)

	modelFi, errMfi := os.Stat(modelDir)
	logger.Print("## Loading dapp: " + packageFile.Name + " ##")
	if errMfi != nil {
		logger.Printf("Error loading 'models' directory for dapp '%s'. Skipping.\n", dir)
		logger.Println(errMfi.Error())
		return
	}

	if !modelFi.IsDir() {
		logger.Printf("Error loading 'models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	files, err := ioutil.ReadDir(modelDir)
	if err != nil {
		logger.Printf("Error loading 'models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	if len(files) == 0 {
		logger.Printf("No models in model dir for app '%s', skipping.\n", dir)
		return
	}
	
	dapp := NewDapp()
	dapp.path = dir
	dapp.packageFile = packageFile
		
	/*
	
	TODO: Dapp verification has been postponed. Leaving this in, but
	      
	
	// Hash the dapp files and check.
	hash := dc.HashApp(modelDir)
	
	if hash == nil {
		logger.Println("Failed to get hash of dapp files, skipping. Dapp: " + dir)
		return
	}
	
	oldHash, errH := dc.hashDB.Get([]byte(dapp.path), nil)
	
	if errH != nil {
		verify(dapp.path);
		logger.Printf("Adding new hash '%s' to folder '%s'.\n",hex.EncodeToString(hash),dir)
		dc.hashDB.Put([]byte(dapp.path),hash,nil)
	}
	
	if errH == nil && !bytes.Equal(hash,oldHash) {
		// TODO this is an old but updated dapp.
		logger.Printf("Hash mismatch: New: '%s', Old: '%s'.\n",hex.EncodeToString(hash),hex.EncodeToString(oldHash))
		verify(dapp.path);
		dc.hashDB.Put([]byte(dapp.path),hash,nil)
		dc.hashDB.Delete(oldHash,nil)
	} else {
		logger.Printf("Hash of '%s' matches the stored value: '%s'.\n", dir, hex.EncodeToString(hash))
	}
	*/

	dc.dapps[packageFile.Id] = dapp

	// TODO recursively and perhaps also a require.js type load file
	// to ensure the proper loading order.
	for _, fileInfo := range files {
		fp := path.Join(modelDir, fileInfo.Name())
		if fileInfo.IsDir() {
			logger.Println("Action models are not gotten recursively (yet). Skipping directory: " + fp)
			// Skip for now.
			continue
		}

		if strings.ToLower(path.Ext(fp)) != ".js" {
			//fmt.Println("[Dapp Registry] Skipping non .js file: " + fp)
			continue
		}

		fileBts, errFile := ioutil.ReadFile(fp)
		if errFile != nil {
			logger.Println("Error reading javascript file: " + fp)
		}

		jsFile := string(fileBts)

		logger.Printf("Loaded javascript file '%s'\n", path.Base(fp))
		
		dapp.models = append(dapp.models,jsFile);

	}
	
	// Register the handlers right away.
	dc.server.RegisterDapp(dapp.packageFile.Id)

	return
}

// TODO check dependencies.
func (dc *DappRegistry) LoadDapp(dappId string) error {
	
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	dapp, ok := dc.dapps[dappId]
	if (!ok){
		return errors.New("Error loading dapp: " + dappId + ". No dapp with that name has been registered.")
	}
	
	if dc.runningDapp != nil {
		if dc.runningDapp.packageFile.Id == dappId {
			return errors.New("Error loading dapp - already running: " + dappId)
		}
		dc.UnloadDapp(dc.runningDapp)
	}
	
	logger.Println("Loading dapp: " + dappId)
	
	rt := dc.ate.CreateRuntime(dappId)
	
	for _, js := range dapp.models {
		rt.AddScript(js)
	}
	
	// Monk hack until we script	
	deps := dapp.packageFile.ModuleDependencies
	
	if deps != nil {
		for _, d := range deps {
			if d.Name == "monk" {
				mData := d.Data
				if mData != nil {
					monkData := &dapps.MonkData{}
					err := json.Unmarshal(mData, monkData)
					if err != nil {
						logger.Fatal("Blockchain will not work. Chain data for monk not available in dapp package file: " + dapp.packageFile.Name);
					}
					monkMod, ok := dc.moduleReg.GetModules()["monk"]
					if !ok {
						logger.Fatal("Blockchain will not work. There is no Monk module.");
					}
					psAddr := monkData.PeerServerAddress
					addAndPort := strings.Split(psAddr,":")
					if len(addAndPort) != 2 {
						logger.Fatal("Blockchain will not work. Malformed peerserver url: " + psAddr);
					}
					
					port, pErr := strconv.Atoi(addAndPort[1])
					if pErr != nil {
						logger.Fatal("Blockchain will not work. Malformed peerserver url (port not an integer)");
					}
					
					monkMod.SetProperty("RemoteHost",addAndPort[0])
					monkMod.SetProperty("RemotePort",port)
					monkMod.SetProperty("ChainId",monkData.ChainId)
					logger.Println("Calling restart on monk")
					cr := make(chan bool)
					go func(){
						monkMod.Restart()
						cr <- true
					}()
					<- cr
					logger.Print("Started");
					rt.BindScriptObject("RootContract",monkData.RootContract)
				} else {
					logger.Fatal("Blockchain will not work. Chain data for monk not available in dapp package file: " + dapp.packageFile.Name);
				}
			}
		}
	}
	
	dc.runningDapp = dapp
	return nil
}

func (dc *DappRegistry) UnloadDapp(dapp *Dapp){
	// TODO cleanup
	dappId := dapp.packageFile.Id
	logger.Println("Unregistering dapp: " + dappId)
	dc.ate.RemoveRuntime(dappId);
}

func (dc *DappRegistry) HashApp(dir string) []byte {
	logger.Println("Hashing models folder: " + dir)
	hashes := dc.HashDir(dir)
	if hashes == nil {
		return nil
	}
	hash := sha1.Sum(hashes)
	return hash[:]
}

func (dc *DappRegistry) HashDir(directory string) []byte {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		logger.Println(err.Error())
		return nil
	}
	if len(files) == 0 {
		logger.Println("No files in directory: " + directory)
		return nil
	}
	hashes := make([]byte, 0)
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			// This is a dapp
			pth := path.Join(directory, fileInfo.Name())
			hs := dc.HashDir(pth)
			if hs != nil {
				hashes = append(hashes, hs...)
			} else {
				return nil
			}
		} else {
			fBts, errF := ioutil.ReadFile(path.Join(directory, fileInfo.Name()))
			if errF != nil {
				logger.Printf(" Error loading '%s', skipping...\n", fileInfo.Name())
				logger.Println(errF.Error())
				return nil
			}
			hash := sha1.Sum(fBts)
			hashes = append(hashes, hash[:]...)
		}
	}
	return hashes
}

func (dc *DappRegistry) GetDappList() []*dapps.DappInfo{
	arr := make([]*dapps.DappInfo,len(dc.dapps))
	ctr := 0;
	for _ , dapp := range dc.dapps {
		arr[ctr] = dapps.DappInfoFromPackageFile(dapp.packageFile)
		ctr++
	}
	return arr
}

/*
func getVerification(string dappName) bool {
	
}
*/