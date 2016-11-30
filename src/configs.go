package main

import (
	"io/ioutil"
	"log"

	"../../go-csgo-cfg"
)

const (
	CONFIG_DIRECTORY = "./configs/"
)

// Loads cfg file from the config directory
func loadConfig(filename string) (*CsgoServerSettings, error) {
	outModel := new(CsgoServerSettings)
	config, err := csgo_cfg.Load(filename)

	if err != nil {
		return nil, err
	}

	err = config.MapTo(outModel)

	return outModel, err
}

// GetServerConfigsFromFile loads all configs in the config directory
// and returns them in a map indexed by filename
// TODO: Cache
func GetServerConfigsFromFile() (map[string]*CsgoServerSettings, error) {
	files, err := ioutil.ReadDir(CONFIG_DIRECTORY)

	if err != nil {
		return nil, err
	}

	configs := make(map[string]*CsgoServerSettings)

	for _, file := range files {
		config, err := loadConfig(CONFIG_DIRECTORY + file.Name())

		if err == nil {
			configs[file.Name()] = config
		} else {
			log.Printf("Error Reading Config %s: %s", file.Name(), err)
		}
	}

	return configs, nil
}

// GetServerConfigFromFile loads a config file by name and returns as struct
func GetServerConfigFromFile(name string) (*CsgoServerSettings, error) {
	config, err := loadConfig(CONFIG_DIRECTORY + name)

	if err != nil {
		log.Printf("Error Reading Config File %s: %s", name, err)

		return nil, err
	}

	return config, nil
}
