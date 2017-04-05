package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

var (
	config      *Configuration
	configLock  = new(sync.RWMutex)
	baseCfgPath = fmt.Sprintf("/%s/%s", "Nitecon", "1Password.json")
	baseConfig  = `{"initial_setup": false, "main_location":""}`
)

// Configuration is the main configuration struct that is used by the application and can be pulled with utils.GetConfig().
type Configuration struct {
	InitialConfig bool   `json:"initial_setup"`
	MainLocation  string `json:"main_location"`
}

func getConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("LOCALAPPDATA")
		configPath := filepath.ToSlash(appData + path.Dir(baseCfgPath))
		return configPath, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	cfgLoc := usr.HomeDir + "/.config"
	configPath := filepath.ToSlash(cfgLoc + path.Dir(baseCfgPath))
	return configPath, nil
}

func getUserConfigLoc() string {
	cfgDir, err := getConfigDir()
	if err != nil {
		// If this isn't a real user we panic out
		panic(err)
	}
	cfgLoc := filepath.ToSlash(cfgDir + "/1Password.json")
	if _, err := os.Stat(cfgLoc); os.IsNotExist(err) {
		err := os.MkdirAll(cfgDir, 0775)
		if err != nil {
			panic("You do not have access to create configuration data as your user something is really wrong exiting...\n" + err.Error())
		}
		err = ioutil.WriteFile(cfgLoc, []byte(baseConfig), 0664)
		if err != nil {
			panic("Could not write default config as your user something is really wrong exiting...\n" + err.Error())
		}

	}
	return cfgLoc
}

// UpdateConfig writes the updated configuration data to the storage location.
func UpdateConfig(conf Configuration) (err error) {
	cfgDir, _ := getConfigDir()
	d, err := json.Marshal(&conf)
	if err != nil {
		return errors.New("Could not save new json config")
	}
	ioutil.WriteFile(cfgDir+"/1Password.json", d, 0664)
	return nil
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
