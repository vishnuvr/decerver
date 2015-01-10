package decerver

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/eris-ltd/decerver/interfaces/core"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

var (
	GoPath = os.Getenv("GOPATH")
	usr, _ = user.Current() // error?!
)

// set default config object
var DefaultConfig = &core.DCConfig{
	RootDir:    path.Join(usr.HomeDir, ".decerver"),
	LogFile:    "",
	MaxClients: 10,
	Port:       3000,
}

func (dc *DeCerver) WriteConfig(dcConfig *core.DCConfig) {
	b, err := json.Marshal(dcConfig)

	if err != nil {
		logger.Println("error marshalling config:", err)
		return
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	ioutil.WriteFile((path.Join(dc.paths.Root(), "config.json")), out.Bytes(), 0600)
}

func (dc *DeCerver) GetConfig() *core.DCConfig {
	return dc.config
}

func (dc *DeCerver) ReadConfig(config_file string) {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		logger.Println("could not read config", err)
		logger.Println("resorting to defaults")
		dc.config = DefaultConfig
		return
	}
	config := &core.DCConfig{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		logger.Println("error unmarshalling config from file:", err)
		logger.Println("resorting to defaults")
		dc.config = DefaultConfig
		return
	}
	dc.config = config
}

func (dc *DeCerver) SetConfig(config interface{}) error {
	if s, ok := config.(string); ok {
		dc.ReadConfig(s)
	} else if s, ok := config.(*core.DCConfig); ok {
		dc.config = s
	} else {
		return errors.New("could not set config")
	}
	return nil
}
