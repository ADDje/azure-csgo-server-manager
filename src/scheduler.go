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
	actions[name] = action

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

	switch strings.ToLower(action.Action) {
	case "deploy":

		// For all parameters, first check the url, then the action default
		numberOfServers, ok := params["numberOfServers"]
		if !ok {
			numberOfServers, ok = findParameter(action, "numberOfServers")
			if !ok {
				return errors.New("Missing numberOfServers parameter")
			}
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

		// Sanity check
		iNumberOfServers, err := strconv.Atoi(numberOfServers)
		if err != nil {
			log.Printf("Invalid number of servers: %s", err)
			return err
		}

		DeployXTemplates(iNumberOfServers, config, azureServerName, vmUserName,
			vmPassword, configFile, deploymentTemplate)

		break
	case "delete":

		nameTemplate, ok := params["nameTemplate"]
		if !ok {
			nameTemplate, ok = findParameter(action, "nameTemplate")
			if !ok {
				return errors.New("Missing nameTemplate parameter")
			}
		}

		numberOfServers, ok := params["numberOfServers"]
		if !ok {
			numberOfServers, ok = findParameter(action, "numberOfServers")
			if !ok {
				return errors.New("Missing numberOfServers parameter")
			}
		}

		startingNumber, ok := params["startingNumber"]
		if !ok {
			startingNumber, ok = findParameter(action, "startingNumber")
		}

		// Sanity check
		iNumberOfServers, err := strconv.Atoi(numberOfServers)
		if err != nil {
			log.Printf("Invalid number of servers: %s", err)
			return err
		}

		var iStartingNumber int
		if !ok {
			iStartingNumber = 1
		} else {
			var err error
			iStartingNumber, err = strconv.Atoi(startingNumber)
			if err != nil {
				log.Printf("Invalid startingNumber parameter: %s", err)
				return err
			}
		}

		for i := iStartingNumber; i < iStartingNumber+iNumberOfServers; i++ {
			name, err := replaceParameter(nameTemplate, i)
			if err != nil {
				log.Printf("Skipping for deletion, %s, %d", name, i)
			} else {
				FullDeleteVM(config, name)
			}
		}

		break
	case "start":

		nameTemplate, ok := params["nameTemplate"]
		if !ok {
			nameTemplate, ok = findParameter(action, "nameTemplate")
			if !ok {
				return errors.New("Missing nameTemplate parameter")
			}
		}

		numberOfServers, ok := params["numberOfServers"]
		if !ok {
			numberOfServers, ok = findParameter(action, "numberOfServers")
			if !ok {
				return errors.New("Missing numberOfServers parameter")
			}
		}

		startingNumber, ok := params["startingNumber"]
		if !ok {
			startingNumber, ok = findParameter(action, "startingNumber")
		}

		// Sanity check
		iNumberOfServers, err := strconv.Atoi(numberOfServers)
		if err != nil {
			log.Printf("Invalid number of servers: %s", err)
			return err
		}

		var iStartingNumber int
		if !ok {
			iStartingNumber = 1
		} else {
			var err error
			iStartingNumber, err = strconv.Atoi(startingNumber)
			if err != nil {
				log.Printf("Invalid startingNumber parameter: %s", err)
				return err
			}
		}

		for i := iStartingNumber; i < iStartingNumber+iNumberOfServers; i++ {
			name, err := replaceParameter(nameTemplate, i)
			if err != nil {
				log.Printf("Skipping for deletion, %s, %d", name, i)
			} else {
				StartVM(config, name)
			}
		}

		break
	case "stop":

		nameTemplate, ok := params["nameTemplate"]
		if !ok {
			nameTemplate, ok = findParameter(action, "nameTemplate")
			if !ok {
				return errors.New("Missing nameTemplate parameter")
			}
		}

		numberOfServers, ok := params["numberOfServers"]
		if !ok {
			numberOfServers, ok = findParameter(action, "numberOfServers")
			if !ok {
				return errors.New("Missing numberOfServers parameter")
			}
		}

		startingNumber, ok := params["startingNumber"]
		if !ok {
			startingNumber, ok = findParameter(action, "startingNumber")
		}

		// Sanity check
		iNumberOfServers, err := strconv.Atoi(numberOfServers)
		if err != nil {
			log.Printf("Invalid number of servers: %s", err)
			return err
		}

		var iStartingNumber int
		if !ok {
			iStartingNumber = 1
		} else {
			var err error
			iStartingNumber, err = strconv.Atoi(startingNumber)
			if err != nil {
				log.Printf("Invalid startingNumber parameter: %s", err)
				return err
			}
		}

		for i := iStartingNumber; i < iStartingNumber+iNumberOfServers; i++ {
			name, err := replaceParameter(nameTemplate, i)
			if err != nil {
				log.Printf("Skipping for deletion, %s, %d", name, i)
			} else {
				DeallocateVM(config, name)
			}
		}

		break
	default:
		return fmt.Errorf("Invalid Action action: %s", action.Action)
	}

	return nil
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
