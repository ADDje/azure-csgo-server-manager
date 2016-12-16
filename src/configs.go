package main

import (
	"io/ioutil"
	"log"

	"github.com/MetalMichael/go-csgo-cfg"
)

const (
	CONFIG_DIRECTORY = "./configs/"
)

// Loads cfg file from the config directory
func loadConfig(filename interface{}) (*CsgoServerSettings, error) {
	outModel := new(CsgoServerSettings)
	config, err := csgo_cfg.Load(filename)

	if err != nil {
		return nil, err
	}

	err = config.MapTo(outModel)

	return outModel, err
}

// Loads cfg file from the config directory as raw text
func loadConfigText(filename string) (string, error) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(file), nil
}

func GetServerConfigsFromAzure() (map[string]*CsgoServerSettings, error) {
	configFiles, err := getStorageFiles(config, "Configs")

	if err != nil {
		return nil, err
	}

	configs := make(map[string]*CsgoServerSettings)

	for _, file := range configFiles {
		azureFile, err := getStorageFile(config, "Configs", file)
		if err != nil {
			return nil, err
		}

		buffer := make([]byte, azureFile.Properties.ContentLength)
		_, err2 := azureFile.Body.Read(buffer)
		if err2 != nil {
			return nil, err2
		}

		config, err3 := loadConfig(buffer)

		if err3 == nil {
			configs[file.Name] = config
		} else {
			log.Printf("Error Reading Config %s: %s", file.Name, err3)
		}
	}

	return configs, nil
}

// GetServerConfigsFromFile loads all configs in the config directory
// and returns them in a map indexed by filename
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

// GetServerConfigTextFromFile loads a config file by name and returns
func GetServerConfigTextFromFile(name string) (string, error) {
	config, err := loadConfigText(CONFIG_DIRECTORY + name)

	if err != nil {
		log.Printf("Error Reading Config File Text %s: %s", name, err)

		return "", err
	}

	return config, nil
}

// CheckConfigValid returns whether or not a config file is valid
func CheckConfigValid(config []byte) (bool, error) {
	_, err := csgo_cfg.Load(config)

	if err != nil {
		return false, err
	}

	return true, nil
}
