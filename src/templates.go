package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	TEMPLATE_DIRECTORY = "./templates/"
)

type DeploymentTemplate struct {
	Template   interface{}
	Parameters interface{}
}

func GetTemplates() (map[string]*DeploymentTemplate, error) {
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
		paramFile, err = LoadTemplateParamsFromFile(templateName)

		if err != nil {
			log.Printf("Matching parameters not found for %s", name)
			continue
		}

		var templateFile []byte
		templateFile, err = LoadTemplateFromFile(templateName)

		if err != nil {
			log.Printf("Could not open template file %s: %s", templateName, err)
			continue
		}

		template := DeploymentTemplate{}
		err = json.Unmarshal(templateFile, &template.Template)

		if err != nil {
			log.Printf("Invalid JSON for template %s", templateName)
			continue
		}

		err = json.Unmarshal(paramFile, &template.Parameters)

		if err != nil {
			log.Printf("Invalid JSON for parameters %s", templateName)
		}

		if err == nil {
			templates[templateName] = &template
		} else {
			log.Printf("Error Reading Config %s: %s", file.Name(), err)
		}
	}

	return templates, nil
}

func LoadTemplateFromFile(name string) ([]byte, error) {
	return ioutil.ReadFile(TEMPLATE_DIRECTORY + name + ".json")
}

func LoadTemplateParamsFromFile(name string) ([]byte, error) {
	return ioutil.ReadFile(TEMPLATE_DIRECTORY + name + ".parameters.json")
}

// CheckTemplateValid returns whether or not a template file is valid
func CheckTemplateValid(template []byte) (bool, error) {
	//TODO

	log.Printf("Content: %s", string(template))
	return true, nil
}
