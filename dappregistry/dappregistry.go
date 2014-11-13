package dappregistry

import (
	// "path/filepath"
	"fmt"
	"io/ioutil"
	"os"
	//"crypto/sha1"
	//"bytes"
	"encoding/json"
	"github.com/eris-ltd/deCerver-interfaces/core"
	"github.com/eris-ltd/deCerver-interfaces/dapps"
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
	mutex *sync.Mutex
	keys  map[string]string
	dapps map[string]*Dapp
	ate   core.Runtime
}

func NewDappRegistry(ate core.Runtime) *DappRegistry {
	dr := &DappRegistry{}
	dr.keys = make(map[string]string)
	dr.dapps = make(map[string]*Dapp)
	dr.mutex = &sync.Mutex{}
	dr.ate = ate
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
	/*
		idxDir := path.Join(dir,dapps.INDEX_FILE_NAME)
		_ , err1 = os.Stat(idxDir)
		if err1 != nil {
			fmt.Printf("Error loading 'index.html' for dapp '%s'. Skipping...\n", dir)
			fmt.Println(err1.Error())
			return
		}
	*/

	/*
		mdFi , err3 := os.Stat(mdDir)
		if err3 != nil || os.IsNotExist(mdFi) {
			fmt.Printf("Dapp '%s' does not have a 'package.json' file. Skipping...\n", dir)
			return err2
		}*/

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

	dapp := NewDapp()
	dapp.path = dir
	dapp.packageFile = packageFile

	// TODO the hashing thing.

	dc.dapps[packageFile.Name] = dapp

	modelDir := path.Join(dir, dapps.MODELS_FOLDER_NAME)

	modelFi, errMfi := os.Stat(modelDir)

	if errMfi != nil {
		fmt.Printf("[Dapp Registry] Error loading 'Models' directory for dapp '%s'. Skipping models\n.", dir)
		fmt.Println(errMfi.Error())
		return
	}

	if !modelFi.IsDir() {
		fmt.Printf("[Dapp Registry] Error loading 'Models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	fmt.Println("[Dapp Registry] Dapp module directory: " + modelDir)
	files, err := ioutil.ReadDir(modelDir)
	if err != nil {
		fmt.Printf("[Dapp Registry] Error loading 'Models' directory for dapp '%s': Not a directory.\n", dir)
		return
	}

	if len(files) == 0 {
		fmt.Printf("[Dapp Registry] No models in model dir for app '%s', skipping.\n", dir)
		return
	}

	for _, fileInfo := range files {
		fp := path.Join(modelDir, fileInfo.Name())
		if fileInfo.IsDir() {
			fmt.Println("[Dapp Registry] Action models are not searched for recursively (yet), skipping directory: " + fp)
			// Skip for now.
			continue
		}
		
		if strings.ToLower(path.Ext(fp)) != ".js" {
			fmt.Println("[Dapp Registry] Skipping non .js file: " + fp)
			continue
		}

		fileBts, errFile := ioutil.ReadFile(fp)
		if errFile != nil {
			fmt.Println("[Dapp Registry] Error reading javascript file: " + fp)
		}
		
		jsFile := string(fileBts)
		
		
		
		parseErr := dc.ate.ParseScript(jsFile)
		
		if parseErr != nil {
			fmt.Printf("[Dapp Registry] Error parsing javascript file: %s\nDUMP: \n%s\n", jsFile, parseErr.Error())
			continue
		}
		
		addErr := dc.ate.AddScript(jsFile)
		
		if addErr != nil {
			fmt.Printf("[Dapp Registry] Error running javascript file: %s\nDUMP: \n%s\n", jsFile, addErr.Error())
		}
		
		result, erk := dc.ate.RunFunction("Shitty.CreateFile","nothing");
		
		if erk != nil {
			fmt.Println(erk.Error())	
		} else {
			fmt.Printf("%v\n",result)
		}
	}

	return
}

func (dc *DappRegistry) HashApp(dir string) []byte {

	return nil
}
