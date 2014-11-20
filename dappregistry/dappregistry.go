package dappregistry

import (
	// "path/filepath"
	"fmt"
	"io/ioutil"
	"os"
	//"crypto/sha1"
	//"bytes"
	"encoding/json"
	"github.com/eris-ltd/decerver-interfaces/api"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/dapps"
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

func (dc *DappRegistry) LoadDapps(directory string) error {
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
		fmt.Printf("[Dapp Registry] Error loading 'Models' directory for dapp '%s'. Skipping models\n.", dir)
		fmt.Println(errMfi.Error())
		return
	}

	if !modelFi.IsDir() {
		fmt.Printf("[Dapp Registry] Error loading 'Models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	files, err := ioutil.ReadDir(modelDir)
	if err != nil {
		fmt.Printf("[Dapp Registry] Error loading 'Models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	if len(files) == 0 {
		fmt.Printf("[Dapp Registry] No models in model dir for app '%s', skipping.\n", dir)
		return
	}

	dapp := NewDapp()
	dapp.path = dir
	dapp.packageFile = packageFile

	// TODO the hashing thing.

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

	return nil
}
