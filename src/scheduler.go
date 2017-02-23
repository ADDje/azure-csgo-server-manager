package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type ScheduleAction struct {
	Name       string              `json:"name"`
	Action     string              `json:"action"`
	Parameters []ScheduleParameter `json:"parameters"`
	Enabled    bool                `json:"enabled"`
}

type ScheduleParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var supportedActions = [...]string{"start", "stop", "deploy", "delete"}

const ACTIONS_STORE = "schedule"
const ACTIONS_FILE = "actions.json"

var actions map[string]*ScheduleAction

func GetScheduleActions() (map[string]*ScheduleAction, error) {
	if actions == nil {
		err := loadActions()
		if err != nil {
			return nil, err
		}
	}
	return actions, nil
}

func GetScheduleAction(name string) (*ScheduleAction, error) {
	if actions == nil {
		err := loadActions()
		if err != nil {
			return nil, err
		}
	}

	action, ok := actions[name]
	if !ok {
		return nil, errors.New("Could not find Action")
	}
	return action, nil
}

func loadActions() error {
	actions = make(map[string]*ScheduleAction)
	var err error

	if config.UseCloudStorage {
		actionsFile, err := GetStorageFile(config, ACTIONS_STORE, ACTIONS_FILE)
		if err != nil {
			// First time, create ACTIONS_FILE
			if strings.Contains(string(err.Error()), "resource does not exist") {
				log.Printf("Creating Schedule Actions file: %s", ACTIONS_FILE)
				err := CreateStorageFile(config, ACTIONS_STORE, ACTIONS_FILE, []byte("{}"))
				if err != nil {
					return err
				}
			}
			return err
		}
		if actionsFile != nil {
			err = json.NewDecoder(*actionsFile).Decode(&actions)
		}
	} else {
		actionsFile, err := ioutil.ReadFile(ACTIONS_FILE)
		if err != nil {
			log.Printf("Could not read actions file: %s", err)
			return err
		}

		err = json.Unmarshal(actionsFile, &actions)
	}
	if err != nil {
		log.Printf("Invalid JSON in schedule actions: %s", err)
	}
	return nil
}

func AddOrUpdateScheduleAction(name string, action *ScheduleAction) error {
	if actions == nil {
		err := loadActions()
		if err != nil {
			return err
		}
	}

	// Renamed
	if action.Name != name {
		log.Printf("Renaming action %s to %s", name, action.Name)
		delete(actions, name)
	}
	actions[action.Name] = action

	err := saveActions()
	if err != nil {
		return err
	}

	return nil
}

func DeleteScheduleAction(name string) error {
	if actions == nil {
		err := loadActions()
		if err != nil {
			return err
		}
	}
	delete(actions, name)

	err := saveActions()
	if err != nil {
		return err
	}

	return nil
}

func ExecuteScheduleAction(name string, params map[string]string) error {
	if actions == nil {
		err := loadActions()
		if err != nil {
			return err
		}
	}

	action, ok := actions[name]
	if !ok {
		return errors.New("Action Not Found")
	}

	if !action.Enabled {
		return errors.New("Action is disabled")
	}

	switch strings.ToLower(action.Action) {
	case "deploy":

		// For all parameters, first check the url, then the action default
		numberOfServers, err := getNumberOfServers(action, params)
		if err != nil {
			return err
		}

		configFile, ok := params["configFile"]
		if !ok {
			configFile, ok = findParameter(action, "configFile")
			if !ok {
				return errors.New("Missing configFile parameters")
			}
		}

		deploymentTemplate, ok := params["deploymentTemplate"]
		if !ok {
			deploymentTemplate, ok = findParameter(action, "deploymentTemplate")
			if !ok {
				return errors.New("Missing deploymentTemplate parameter")
			}
		}

		azureServerName, ok := params["azureServerName"]
		if !ok {
			azureServerName, _ = findParameter(action, "azureServerName")
		}

		vmUserName, ok := params["vmUsername"]
		if !ok {
			vmUserName, _ = findParameter(action, "vmUsername")
		}

		vmPassword, ok := params["vmPassword"]
		if !ok {
			vmPassword, _ = findParameter(action, "vmPassword")
		}

		DeployXTemplates(numberOfServers, config, azureServerName, vmUserName,
			vmPassword, configFile, deploymentTemplate)

		break
	case "delete":

		serverNameTemplate, ok := params["serverNameTemplate"]
		if !ok {
			serverNameTemplate, ok = findParameter(action, "serverNameTemplate")
			if !ok {
				return errors.New("Missing serverNameTemplate parameter")
			}
		}

		numberOfServers, err := getNumberOfServers(action, params)
		if err != nil {
			return err
		}

		startingNumber, err := getStartingNumber(action, params)
		if err != nil {
			return err
		}

		for i := startingNumber; i < startingNumber+numberOfServers; i++ {
			name, err := replaceParameter(serverNameTemplate, i)
			if err != nil {
				log.Printf("Skipping for deletion, %s, %d", name, i)
			} else {
				go FullDeleteVM(config, name)
			}
		}

		break
	case "start":

		serverNameTemplate, ok := params["serverNameTemplate"]
		if !ok {
			serverNameTemplate, ok = findParameter(action, "serverNameTemplate")
			if !ok {
				return errors.New("Missing serverNameTemplate parameter")
			}
		}

		numberOfServers, err := getNumberOfServers(action, params)
		if err != nil {
			return err
		}

		startingNumber, err := getStartingNumber(action, params)
		if err != nil {
			return err
		}

		for i := startingNumber; i < startingNumber+numberOfServers; i++ {
			name, err := replaceParameter(serverNameTemplate, i)
			if err != nil {
				log.Printf("Skipping for deletion, %s, %d", name, i)
			} else {
				go StartVM(config, name)
			}
		}

		break
	case "stop":

		serverNameTemplate, ok := params["serverNameTemplate"]
		if !ok {
			serverNameTemplate, ok = findParameter(action, "serverNameTemplate")
			if !ok {
				return errors.New("Missing serverNameTemplate parameter")
			}
		}

		numberOfServers, err := getNumberOfServers(action, params)
		if err != nil {
			return err
		}

		startingNumber, err := getStartingNumber(action, params)
		if err != nil {
			return err
		}

		for i := startingNumber; i < startingNumber+numberOfServers; i++ {
			name, err := replaceParameter(serverNameTemplate, i)
			if err != nil {
				log.Printf("Skipping for deletion, %s, %d", name, i)
			} else {
				go DeallocateVM(config, name)
			}
		}

		break
	case "save replays":

		serverNameTemplate, ok := params["serverNameTemplate"]
		if !ok {
			serverNameTemplate, ok = findParameter(action, "serverNameTemplate")
			if !ok {
				return errors.New("Missing serverNameTemplate parameter")
			}
		}

		numberOfServers, err := getNumberOfServers(action, params)
		if err != nil {
			return err
		}

		startingNumber, err := getStartingNumber(action, params)
		if err != nil {
			return err
		}

		replayLabel, ok := params["replayLabel"]
		if !ok {
			replayLabel, ok = findParameter(action, "replayLabel")
			if !ok {
				return errors.New("Missing replayLabel parameter")
			}
		}

		vmUserName, ok := params["vmUsername"]
		if !ok {
			vmUserName, ok = findParameter(action, "vmUsername")
			if !ok {
				return errors.New("Missing vmUsername parameter")
			}
		}

		vmPassword, ok := params["vmPassword"]
		if !ok {
			vmPassword, ok = findParameter(action, "vmPassword")
			if !ok {
				return errors.New("Missing vmPassword parameter")
			}
		}

		for i := startingNumber; i < startingNumber+numberOfServers; i++ {
			name, err := replaceParameter(serverNameTemplate, i)
			if err != nil {
				log.Printf("Skipping for deletion, %s, %d", name, i)
			} else {
				go ExportReplays(config, replayLabel, name, vmUserName, vmPassword)
			}
		}

		break
	default:
		return fmt.Errorf("Invalid Action action: %s", action.Action)
	}

	return nil
}

func getNumberOfServers(action *ScheduleAction, params map[string]string) (int, error) {
	numberOfServers, ok := params["numberOfServers"]
	if !ok {
		numberOfServers, ok = findParameter(action, "numberOfServers")
		if !ok {
			err := errors.New("Missing numberOfServers parameter")
			log.Print(err)
			return 0, err
		}
	}

	iNumberOfServers, err := strconv.Atoi(numberOfServers)
	if err != nil {
		log.Printf("Invalid number of servers: %s", err)
		return 0, err
	}

	return iNumberOfServers, nil
}

func getStartingNumber(action *ScheduleAction, params map[string]string) (int, error) {
	startingNumber, ok := params["startingNumber"]
	if !ok {
		startingNumber, ok = findParameter(action, "startingNumber")
	}

	var iStartingNumber int
	if !ok {
		iStartingNumber = 1
	} else {
		var err error
		iStartingNumber, err = strconv.Atoi(startingNumber)
		if err != nil {
			log.Printf("Invalid startingNumber parameter: %s", err)
			return 0, err
		}
	}

	return iStartingNumber, nil
}

func findParameter(action *ScheduleAction, paramName string) (string, bool) {
	for _, v := range action.Parameters {
		if v.Key == paramName {
			return v.Value, true
		}
	}
	return "", false
}

func saveActions() error {
	if actions == nil {
		return errors.New("Actions not yet loaded, therefore cannot have changed")
	}

	fileBytes, err := json.Marshal(actions)
	if err != nil {
		log.Printf("Cannot serialize actions: %s", err)
		return err
	}

	if config.UseCloudStorage {
		err = UpdateStorageFile(config, ACTIONS_STORE, ACTIONS_FILE, fileBytes)
	} else {
		err = ioutil.WriteFile(ACTIONS_FILE, fileBytes, 0770)
		if err != nil {
			log.Printf("Could not write file: %s", err)
		}
	}
	if err != nil {
		return err
	}

	return nil
}
