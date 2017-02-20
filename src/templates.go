package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

const TEMPLATE_DIRECTORY = "./templates/"
const TEMPLATE_FILE_STORE = "templates"

type DeploymentTemplate struct {
	Template   map[string]interface{}
	Parameters TemplateParameterFile
}

type TemplateParameterFile struct {
	Schema         string                       `json:"$schema"`
	ContentVersion string                       `json:"contentVersion"`
	Parameters     map[string]TemplateParameter `json:"parameters"`
}

type TemplateParameter struct {
	Value interface{} `json:"value"`
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

		var templateName = getBlobName(name[:len(name)-5])
		// Check that parameters file exists
		paramFile, err := loadTemplateParamsFromFile(templateName)

		if err != nil {
			log.Printf("Matching parameters not found for %s", name)
			continue
		}

		templateFile, err := loadTemplateFromFile(templateName)

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

		var templateName = getBlobName(name[:len(name)-5])
		// Check that parameters file exists
		paramFile, err := loadTemplateParamsFromStorage(templateName)

		if err != nil {
			log.Printf("Matching parameters not found for %s", name)
			continue
		}

		templateFile, err := loadTemplateFromStorage(templateName)

		if err != nil {
			log.Printf("Could not open template file %s: %s", templateName, err)
			continue
		}

		template := DeploymentTemplate{}
		err = json.NewDecoder(*templateFile).Decode(&template.Template)
		if err != nil {
			log.Printf("Invalid JSON for template %s: %s", templateName, err)
			continue
		}

		err = json.NewDecoder(*paramFile).Decode(&template.Parameters)
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

func loadTemplateFromStorage(name string) (*io.ReadCloser, error) {
	f, err := GetStorageFile(config, TEMPLATE_FILE_STORE, name+".json")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func loadTemplateParamsFromStorage(name string) (*io.ReadCloser, error) {
	f, err := GetStorageFile(config, TEMPLATE_FILE_STORE, name+".parameters.json")
	if err != nil {
		return nil, err
	}
	return f, nil
}

// CheckTemplateValid returns map if valid json
func CheckTemplateValid(template []byte) (*TemplateParameterFile, error) {

	myMap := TemplateParameterFile{}
	err := json.Unmarshal(template, &myMap)

	if err != nil {
		log.Printf("Template not valid: %s", err)
		return nil, err
	}
	return &myMap, nil
}

// CheckParametersValid returns struct if valid parameters json
func CheckParametersValid(parameters []byte) (*TemplateParameterFile, error) {

	myMap := TemplateParameterFile{}
	err := json.Unmarshal(parameters, &myMap)

	if err != nil {
		log.Printf("Template not valid: %s", err)
		return nil, err
	}
	return &myMap, nil
}

// GetDefaultParametersFile returns default empty parameters file with schema and stuff
func GetDefaultParametersFile() []byte {
	return []byte(`{
    "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentParameters.json#",
    "contentVersion": "1.0.0.0",
    "parameters": {}
}`)
}
