package dappregistry

import (
	// "path/filepath"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/api"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/dapps"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

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
	fmt.Println("[Dapp Registry] Loading dapps")
	files, err := ioutil.ReadDir(directory)
	
	if err != nil {
		return err
	}
	
	if len(files) == 0 {
		fmt.Println("[Dapp Registry] No dapps has been downloaded")
		return nil
	}
	
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			// This is a dapp
			pth := path.Join(directory, fileInfo.Name())
			dc.LoadDapp(pth)
		}
	}
	fmt.Println("[Dapp Registry] Done loading.")
	return nil
}

// TODO check dependencies.
func (dc *DappRegistry) LoadDapp(dir string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	pkDir := path.Join(dir, dapps.PACKAGE_FILE_NAME)
	_, err1 := os.Stat(pkDir)
	if err1 != nil {
		fmt.Printf("[Dapp Registry] Error loading 'package.json' for dapp '%s'. Skipping...\n", dir)
		fmt.Println(err1.Error())
		return
	}

	pkBts, errP := ioutil.ReadFile(pkDir)

	if errP != nil {
		fmt.Printf("[Dapp Registry] Error loading 'package.json' for dapp '%s'. Skipping...\n", dir)
		fmt.Println(errP.Error())
		return
	}

	packageFile := &dapps.PackageFile{}
	pkUnmErr := json.Unmarshal(pkBts, packageFile)

	if pkUnmErr != nil {
		fmt.Printf("[Dapp Registry] The 'package.json' file for dapp '%s' is corrupted. Skipping...\n", dir)
		fmt.Println(pkUnmErr.Error())
	}

	idxDir := path.Join(dir, dapps.INDEX_FILE_NAME)
	_, err1 = os.Stat(idxDir)
	if err1 != nil {
		fmt.Printf("Cannot find an 'index.html' file for dapp '%s'. Skipping...\n", dir)
		fmt.Println(err1.Error())
		return
	}

	modelDir := path.Join(dir, dapps.MODELS_FOLDER_NAME)

	modelFi, errMfi := os.Stat(modelDir)
	fmt.Println("[Dapp Registry] ## Loading dapp: " + packageFile.Name + " ##")
	if errMfi != nil {
		fmt.Printf("[Dapp Registry] Error loading 'models' directory for dapp '%s'. Skipping.\n", dir)
		fmt.Println(errMfi.Error())
		return
	}

	if !modelFi.IsDir() {
		fmt.Printf("[Dapp Registry] Error loading 'models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	files, err := ioutil.ReadDir(modelDir)
	if err != nil {
		fmt.Printf("[Dapp Registry] Error loading 'models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	if len(files) == 0 {
		fmt.Printf("[Dapp Registry] No models in model dir for app '%s', skipping.\n", dir)
		return
	}

	dapp := NewDapp()
	dapp.path = dir
	dapp.packageFile = packageFile

	// Hash the dapp files and check.
	hash := dc.HashApp(dir)

	if hash == nil {
		fmt.Println("Failed to get hash of dapp files, skipping. Dapp: " + dir)
		return
	}
	
	oldHash, errH := dc.hashDB.Get([]byte(dapp.path), nil)
	
	if errH != nil {
		// TODO this is an old dapp
		fmt.Printf("Adding new hash '%s' to folder '%s'.\n",hex.EncodeToString(hash),dir)
		dc.hashDB.Put([]byte(dapp.path),hash,nil)
	}
	
	if errH == nil && !bytes.Equal(hash,oldHash) {
		// TODO this is an old but updated dapp.
		fmt.Printf("Hash mismatch: New: '%s', Old: '%s'.\n",hex.EncodeToString(hash),hex.EncodeToString(oldHash))
	} else {
		fmt.Printf("Hash of '%s' matches the stored value: '%s'.\n", dir, hex.EncodeToString(hash))
	}

	dc.dapps[packageFile.Id] = dapp

	// Create a new javascript runtime
	id := dapp.packageFile.Id

	rt := dc.ate.CreateRuntime(id)
	dc.server.RegisterDapp(id)

	for _, fileInfo := range files {
		fp := path.Join(modelDir, fileInfo.Name())
		if fileInfo.IsDir() {
			fmt.Println("[Dapp Registry] Action models are not searched for recursively (yet), skipping directory: " + fp)
			// Skip for now.
			continue
		}

		if strings.ToLower(path.Ext(fp)) != ".js" {
			//fmt.Println("[Dapp Registry] Skipping non .js file: " + fp)
			continue
		}

		fileBts, errFile := ioutil.ReadFile(fp)
		if errFile != nil {
			fmt.Println("[Dapp Registry] Error reading javascript file: " + fp)
		}

		jsFile := string(fileBts)

		addErr := rt.AddScript(jsFile)

		if addErr != nil {
			fmt.Printf("[Dapp Registry] Error running javascript file: \n%s\n%s\n", jsFile, addErr.Error())
			continue
		}

		fmt.Printf("[Dapp Registry] Loaded javascript file '%s'\n", path.Base(fp))

	}

	return
}

func (dc *DappRegistry) HashApp(dir string) []byte {
	fmt.Println("[Dapp Registry] Hashing app folder: " + dir)
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
		fmt.Println(err.Error())
		return nil
	}
	if len(files) == 0 {
		fmt.Println("No files in directory: " + directory)
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
				fmt.Printf("[Dapp Registry] Error loading '%s', skipping...\n", fileInfo.Name())
				fmt.Println(errF.Error())
				return nil
			}
			hash := sha1.Sum(fBts)
			hashes = append(hashes, hash[:]...)
		}
	}
	return hashes
}
