package dappregistry

import (
	// "path/filepath"
	"io/ioutil"
	//"os"
	"fmt"
)

// Structs that are mapped to the package file.
type (
	PackageFile struct {
		Name               string              `json:"name"`
		Icon               string              `json:"app_icon"`
		Version            string              `json:"version"`
		Homepage           string              `json:"homepage"`
		Author             *Author             `json:"author"`
		Repository         *Repository         `json:"repository"`
		Bugs               *Bugs               `json:"bugs"`
		Licence            *Licence            `json:"licence"`
		ModuleDependencies *ModuleDependencies `json:"moduleDependencies"`
	}

	Author struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	Repository struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	}

	Bugs struct {
		Url string `json:"url"`
	}

	Licence struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	}

	ModuleDependencies struct {
		deps map[string]*Module
	}

	Module struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
)


type Dapp struct {
	path        string
	packageFile *PackageFile
}

func NewDapp() *Dapp {
	dapp := &Dapp{}
	return dapp
}

type DappRegistry struct {
	keys map[string]string
	dapps map[string]*Dapp
}

func NewDappRegistry() *DappRegistry {
	dr := &DappRegistry{}
	dr.dapps = make(map[string]*Dapp)
	return dr
}

func (dc *DappRegistry) LoadDapps(directory string) error {
	fmt.Println("Dapp directory: " + directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	
	if len(files) == 0 {
		fmt.Println("No dapps has been downloaded")
		return nil	
	}
	fmt.Println("Processing dapp directories:")
	for _ , fileInfo := range files {
		//f, fErr := os.Open(file)
		//if fErr != nil {
		//	continue
		//}
    	fmt.Println(fileInfo.Name())	
    	
	}
    
    return nil
}

func (dc *DappRegistry) LoadDapp(dir string) {
}

func (dc *DappRegistry) HashApp(dir string) []byte {
	return nil
}
