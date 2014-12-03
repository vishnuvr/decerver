package dappregistry

import (
	// "path/filepath"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/eris-ltd/decerver-interfaces/api"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/dapps"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"log"
)

var logger *log.Logger = core.NewLogger("Dapp Registry")

type Dapp struct {
	models      map[string]string
	path        string
	packageFile *dapps.PackageFile
}

func NewDapp() *Dapp {
	dapp := &Dapp{models: make(map[string]string)}
	return dapp
}

type DappRegistry struct {
	mutex  *sync.Mutex
	keys   map[string]string
	dapps  map[string]*Dapp
	ate    core.RuntimeManager
	server api.Server
	hashDB *leveldb.DB
}

func NewDappRegistry(ate core.RuntimeManager, server api.Server) *DappRegistry {
	dr := &DappRegistry{}
	dr.keys = make(map[string]string)
	dr.dapps = make(map[string]*Dapp)
	dr.mutex = &sync.Mutex{}
	dr.ate = ate
	dr.server = server
	return dr
}

func (dc *DappRegistry) LoadDapps(directory, dbDir string) error {
	dbDir = path.Join(dbDir,"dapp_stored_hashes")
	dc.hashDB, _ = leveldb.OpenFile(dbDir,nil)
	defer dc.hashDB.Close()
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
			// This is a dapp
			pth := path.Join(directory, fileInfo.Name())
			dc.LoadDapp(pth)
		}
	}
	logger.Println("Done loading dapps.")
	return nil
}

// TODO check dependencies.
func (dc *DappRegistry) LoadDapp(dir string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
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

	// Hash the dapp files and check.
	hash := dc.HashApp(modelDir)

	if hash == nil {
		logger.Println("Failed to get hash of dapp files, skipping. Dapp: " + dir)
		return
	}
	
	oldHash, errH := dc.hashDB.Get([]byte(dapp.path), nil)
	
	if errH != nil {
		// this is a new dapp
		
		/*
		confChan := make(chan bool)
		
		ret <- confChan
		if ret {
			
		} else {
			logger.Println("User denied, skipping. Dapp: " + dir)
			return	
		}
		*/
		logger.Printf("Adding new hash '%s' to folder '%s'.\n",hex.EncodeToString(hash),dir)
		dc.hashDB.Put([]byte(dapp.path),hash,nil)
	}
	
	if errH == nil && !bytes.Equal(hash,oldHash) {
		// TODO this is an old but updated dapp.
		logger.Printf("Hash mismatch: New: '%s', Old: '%s'.\n",hex.EncodeToString(hash),hex.EncodeToString(oldHash))
		dc.hashDB.Put([]byte(dapp.path),hash,nil)
	} else {
		logger.Printf("Hash of '%s' matches the stored value: '%s'.\n", dir, hex.EncodeToString(hash))
	}

	dc.dapps[packageFile.Id] = dapp

	// Create a new javascript runtime
	id := dapp.packageFile.Id

	rt := dc.ate.CreateRuntime(id)
	dc.server.RegisterDapp(id)

	for _, fileInfo := range files {
		fp := path.Join(modelDir, fileInfo.Name())
		if fileInfo.IsDir() {
			logger.Println("Action models are not searched for recursively (yet), skipping directory: " + fp)
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

		addErr := rt.AddScript(jsFile)

		if addErr != nil {
			logger.Printf("Error running javascript file: \n%s\n%s\n", jsFile, addErr.Error())
			continue
		}

		logger.Printf("Loaded javascript file '%s'\n", path.Base(fp))

	}

	return
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
