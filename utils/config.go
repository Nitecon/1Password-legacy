package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"runtime"
	"sync"
)

var (
	config      *Configuration
	configLock  = new(sync.RWMutex)
	baseCfgPath = fmt.Sprintf("%q%s%q%s", os.PathSeparator, "Nitecon", os.PathSeparator, "1Password.json")
	baseConfig  = `{"initial_setup": false, "main_location":""}`
)

// Configuration is the main configuration struct that is used by the application and can be pulled with utils.GetConfig().
type Configuration struct {
	InitialConfig bool   `json:"initial_setup"`
	MainLocation  string `json:"main_location"`
}

func getUserConfigLoc() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("LOCALAPPDATA")
		if _, err := os.Stat(appData + baseCfgPath); os.IsNotExist(err) {
			err := os.MkdirAll(appData+path.Dir(baseCfgPath), 0775)
			if err != nil {
				panic("You do not have access to create configuration data as your user something is really wrong exiting...\n" + err.Error())
			}
			err = ioutil.WriteFile(appData+baseCfgPath, []byte(baseConfig), 0664)
			if err != nil {
				panic("Could not write default config as your user something is really wrong exiting...\n" + err.Error())
			}

		}
		return appData + baseCfgPath
	}
	usr, err := user.Current()
	if err != nil {
		panic("You are not an actual user so I can't help you...\n" + err.Error())
	}
	cfgLoc := usr.HomeDir + "/.config"
	if _, err := os.Stat(cfgLoc + baseCfgPath); os.IsNotExist(err) {
		err := os.MkdirAll(cfgLoc+path.Dir(baseCfgPath), 0775)
		if err != nil {
			panic("You do not have access to create configuration data as your user something is really wrong exiting...\n" + err.Error())
		}
		err = ioutil.WriteFile(cfgLoc+baseCfgPath, []byte(baseConfig), 0664)
		if err != nil {
			panic("Could not write default config as your user something is really wrong exiting...\n" + err.Error())
		}
	}
	return cfgLoc + baseCfgPath
}

// SetConfig is used to load the user's configuration from disk which specifies where the vault lives etc.
func SetConfig() (err error) {
	configFile, err := ioutil.ReadFile(getUserConfigLoc())
	if err != nil {
		return
	}
	tempConf := new(Configuration)
	err = json.Unmarshal(configFile, tempConf)
	if err != nil {
		return
	}
	configLock.Lock()
	config = tempConf
	configLock.Unlock()
	return
}

// GetConfig function returns the currently set active configuration to be used.
func GetConfig() *Configuration {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}
