package deCerver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

var (
	GoPath = os.Getenv("GOPATH")
	usr, _ = user.Current() // error?!
)

type DCConfig struct {
	ConfigFile string `json:"config_file"`
	RootDir    string `json:"root_dir"`
	LogFile    string `json:"log_file"`
	LogLevel   int    `json:"log_level"`
	MaxClients int    `json:max_clients`
	Port       int    `json:port`
}

// set default config object
var DefaultConfig = &DCConfig{
	ConfigFile: "config",
	RootDir:    path.Join(usr.HomeDir, ".deCerver"),
	LogFile:    "",
	LogLevel:   5,
	MaxClients : 10,
}

// can these methods be functions in decerver that take the modules as argument?
func (dc *DeCerver) WriteConfig(config_file string) {
	b, err := json.Marshal(dc.config)
	if err != nil {
		fmt.Println("error marshalling config:", err)
		return
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	ioutil.WriteFile(config_file, out.Bytes(), 0600)
}

func (dc *DeCerver) ReadConfig(config_file string) {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		fmt.Println("could not read config", err)
		fmt.Println("resorting to defaults")
		dc.config = DefaultConfig
		dc.WriteConfig(config_file)
		return
	}
	var config DCConfig
	err = json.Unmarshal(b, &config)
	if err != nil {
		fmt.Println("error unmarshalling config from file:", err)
		fmt.Println("resorting to defaults")
		dc.config = DefaultConfig
		return
	}
	dc.config = &config
}

func (dc *DeCerver) SetConfig(config interface{}) error {
	if s, ok := config.(string); ok {
		dc.ReadConfig(s)
	} else if s, ok := config.(DCConfig); ok {
		dc.config = &s
	} else {
		return errors.New("could not set config")
	}
	return nil
}
