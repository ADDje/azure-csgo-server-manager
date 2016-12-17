package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const TEMPLATE_DIRECTORY = "./templates/"
const TEMPLATE_FILE_STORE = "templates"

type DeploymentTemplate struct {
	Template   interface{}
	Parameters interface{}
}

func GetTemplatesFromFile() (map[string]*DeploymentTemplate, error) {
	files, err := ioutil.ReadDir(TEMPLATE_DIRECTORY)

	if err != nil {
		return nil, err
	}

	templates := make(map[string]*DeploymentTemplate)

	for _, file := range files {
		var name = file.Name()
		// Not a json file
		if len(name) < 4 || name[len(name)-5:] != ".json" {
			log.Printf("File %s is not json and shouldn't be in the template directory", name)
			continue
		}

		// Will be picked up when parsing the matching template if it exists
		if len(name) > 16 && name[len(name)-16:] == ".parameters.json" {
			continue
		}

		var templateName = name[:len(name)-5]
		// Check that parameters file exists
		var paramFile []byte
		paramFile, err = loadTemplateParamsFromFile(templateName)

		if err != nil {
			log.Printf("Matching parameters not found for %s", name)
			continue
		}

		var templateFile []byte
		templateFile, err = loadTemplateFromFile(templateName)

		if err != nil {
			log.Printf("Could not open template file %s: %s", templateName, err)
			continue
		}

		template := DeploymentTemplate{}
		err = json.Unmarshal(templateFile, &template.Template)

		if err != nil {
			log.Printf("Invalid JSON for template %s: %s", templateName, err)
			continue
		}

		err = json.Unmarshal(paramFile, &template.Parameters)
		if err != nil {
			log.Printf("Invalid JSON for parameters %s: %s", templateName, err)
		}

		if err == nil {
			templates[templateName] = &template
		} else {
			log.Printf("Error Reading Config %s: %s", file.Name(), err)
		}
	}

	return templates, nil
}

func GetTemplatesFromAzure() (map[string]*DeploymentTemplate, error) {
	files, err := GetStorageFiles(config, TEMPLATE_FILE_STORE)

	if err != nil {
		return nil, err
	}

	templates := make(map[string]*DeploymentTemplate)

	for _, file := range files {
		var name = file.Name
		// Not a json file
		if len(name) < 4 || name[len(name)-5:] != ".json" {
			log.Printf("File %s is not json and shouldn't be in the template store", name)
			continue
		}

		// Will be picked up when parsing the matching template if it exists
		if len(name) > 16 && name[len(name)-16:] == ".parameters.json" {
			continue
		}

		var templateName = name[:len(name)-5]
		// Check that parameters file exists
		var paramFile string
		paramFile, err = loadTemplateParamsFromStorage(templateName)

		if err != nil {
			log.Printf("Matching parameters not found for %s", name)
			continue
		}

		var templateFile string
		templateFile, err = loadTemplateFromStorage(templateName)

		if err != nil {
			log.Printf("Could not open template file %s: %s", templateName, err)
			continue
		}

		log.Printf("Template: '%s'", templateFile)

		template := DeploymentTemplate{}
		err = json.Unmarshal([]byte(templateFile), &template.Template)
		if err != nil {
			log.Printf("Invalid JSON for template %s: %s", templateName, err)
			continue
		}

		err = json.Unmarshal([]byte(paramFile), &template.Parameters)
		if err != nil {
			log.Printf("Invalid JSON for parameters %s: %s", templateName, err)
		}

		if err == nil {
			templates[templateName] = &template
		} else {
			log.Printf("Error Reading Config %s: %s", file.Name, err)
		}
	}

	return templates, nil
}

func loadTemplateFromFile(name string) ([]byte, error) {
	return ioutil.ReadFile(TEMPLATE_DIRECTORY + name + ".json")
}

func loadTemplateParamsFromFile(name string) ([]byte, error) {
	return ioutil.ReadFile(TEMPLATE_DIRECTORY + name + ".parameters.json")
}

func loadTemplateFromStorage(name string) (string, error) {
	return GetStorageFileText(config, TEMPLATE_FILE_STORE, name+".json")
}

func loadTemplateParamsFromStorage(name string) (string, error) {
	return GetStorageFileText(config, TEMPLATE_FILE_STORE, name+".parameters.json")
}

// CheckTemplateValid returns whether or not a template file is valid
func CheckTemplateValid(template []byte) (bool, error) {
	//TODO

	log.Printf("Content: %s", string(template))
	return true, nil
}
