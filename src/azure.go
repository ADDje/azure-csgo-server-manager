package main

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest/azure"
)

const DEPLOYMENT_NAME string = "csgo-server-manager"

func DeployTemplate(config Config, number int, vmName string, adminUserName string,
	adminPassword string, configName string, templateName string) error {

	client, err := getDeploymentClient(config)
	if err != nil {
		return err
	}

	templateUri := GetStorageFileLink(config, TEMPLATE_FILE_STORE, templateName+".json")

	parameters, err := GetStorageFile(config, TEMPLATE_FILE_STORE, templateName+".parameters.json")
	if err != nil {
		return err
	}

	// Read the file into a json map
	parametersJSON := TemplateParameterFile{}
	err = json.NewDecoder(parameters.Body).Decode(&parametersJSON)
	if err != nil {
		log.Printf("Invalid parameters json: %s", err)
		return err
	}

	// Replace values from file with overrides
	if vmName != "" {
		parametersJSON.Parameters["vmName"] = TemplateParameter{Value: vmName}
	}
	if adminUserName != "" {
		parametersJSON.Parameters["adminUserName"] = TemplateParameter{Value: adminUserName}
	}
	if adminPassword != "" {
		parametersJSON.Parameters["adminPassword"] = TemplateParameter{Value: adminPassword}
	}

	// Replace any variables in parameters with their values
	err = replaceParameterVariables(&parametersJSON, number)
	if err != nil {
		return err
	}

	configLink := GetStorageFileLink(config, CONFIG_FILE_STORE, configName)

	// Add config file, needed for automated config deployment
	parametersJSON.Parameters["configFileUrl"] = TemplateParameter{Value: configLink}
	parametersJSON.Parameters["configFileName"] = TemplateParameter{Value: configName}

	deploymentName := DEPLOYMENT_NAME + "-" + strconv.Itoa(number)

	exportParameters := convertParameters(parametersJSON)

	log.Printf("%s", exportParameters)

	template := resources.TemplateLink{
		URI: &templateUri,
	}
	properties := resources.DeploymentProperties{
		TemplateLink: &template,
		Parameters:   &exportParameters,
		Mode:         "Incremental",
	}
	deployment := resources.Deployment{Properties: &properties}

	_, err = client.CreateOrUpdate(config.ResourceGroup, deploymentName, deployment, nil)
	if err != nil {
		log.Printf("Error creating deployment: %s", err)
	}

	return nil
}

func StartVM(config Config, vmName string) error {
	client, err := getVMClient(config)
	if err != nil {
		return err
	}

	_, err = client.Start(config.ResourceGroup, vmName, nil)
	if err != nil {
		log.Printf("Error Starting VM: %s", err)
		return err
	}
	return nil
}

func DeallocateVM(config Config, vmName string) error {
	client, err := getVMClient(config)
	if err != nil {
		return err
	}

	_, err = client.Deallocate(config.ResourceGroup, vmName, nil)
	if err != nil {
		log.Printf("Error Starting VM: %s", err)
		return err
	}
	return nil
}

// GetVms Gets VMs from the ResourceGroup defined in the config
func GetVms(config Config) (*[]compute.VirtualMachine, error) {
	client, err := getVMClient(config)
	if err != nil {
		return nil, err
	}

	results, err := client.List(config.ResourceGroup)
	if err != nil {
		log.Printf("Error in azure getServers: %s", err)
		return nil, err
	}

	// VM properties are missing InstanceViewStatus which is where the status comes from
	// TODO: Make parallel

	for k, vm := range *results.Value {
		vmInfo, err := GetVmStatus(config, *vm.Name)
		if err == nil {
			(*results.Value)[k].VirtualMachineProperties.InstanceView = vmInfo
		}
	}

	return results.Value, nil
}

func GetVmStatus(config Config, vmName string) (*compute.VirtualMachineInstanceView, error) {
	client, err := getVMClient(config)
	if err != nil {
		return nil, err
	}

	results, err := client.Get(config.ResourceGroup, vmName, "instanceView")
	if err != nil {
		log.Printf("Error in azure GetVmStatus: %s", err)
		return nil, err
	}

	return results.VirtualMachineProperties.InstanceView, nil
}

func GetStorageFileLink(config Config, store string, file string) string {
	return "https://" + config.AzureStorageServer + ".file.core.windows.net/" + store +
		"/" + file + config.AzureSASToken
}

// GetStorageFile Returns file from cloud storage by name and store
func GetStorageFile(config Config, store string, file string) (*storage.FileStream, error) {
	client, err := getStorageClient(config)
	if err != nil {
		return nil, err
	}

	fileStream, err2 := client.GetFile(store+"/"+file, nil)
	if err2 != nil {
		log.Printf("Error in azure GetStorageFile: %s", err2)
		return nil, err2
	}

	return fileStream, nil
}

// // GetStorageFileText Returns file contents from cloud storage by name and store
// func GetStorageFileText(config Config, store string, file string) (string, error) {
// 	fileStream, err := GetStorageFile(config, store, file)
// 	if err != nil {
// 		return "", err
// 	}

// 	log.Printf("%d", fileStream.Properties.ContentLength)
// 	buffer := make([]byte, fileStream.Properties.ContentLength)
// 	if fileStream.Properties.ContentLength > 0 {
// 		r, err3 := fileStream.Body.Read(buffer)
// 		log.Printf("%d bytes read", r)
// 		if err3 != nil && err3 != io.EOF {
// 			log.Printf("Error in azure GetStorageFileText: %s", err3)
// 			return "", err3
// 		}
// 	}

// 	return string(buffer), nil
// }

// GetStorageFiles Returns files from cloud storage by name and store
func GetStorageFiles(config Config, store string) ([]storage.File, error) {
	client, err := getStorageClient(config)
	if err != nil {
		return nil, err
	}

	exists, err2 := client.ShareExists(store)
	if err2 != nil || !exists {
		log.Printf("Error in azure GetStorageFiles: %s", err)
		return nil, err
	}

	params := storage.ListDirsAndFilesParameters{}
	files, err3 := client.ListDirsAndFiles(store, params)
	if err3 != nil {
		log.Printf("Error in azure GetStorageFiles: %s", err)
		return nil, err
	}

	return files.Files, nil
}

// DeleteStorageFile Deletes a file in cloud storage using name and store
func DeleteStorageFile(config Config, store string, file string) error {
	client, err := getStorageClient(config)
	if err != nil {
		return err
	}

	_, err2 := client.GetFile(store+"/"+file, nil)
	if err2 != nil {
		log.Printf("Error in azure DeleteStorageFile: %s", err2)
		return err2
	}

	// Updating a file in storage is falls back to reading writing individual bytes
	// Probably just easier to delete then add
	err3 := client.DeleteFile(store + "/" + file)
	if err3 != nil {
		log.Printf("Error in azure DeleteStorageFile delete: %s", err3)
		return err3
	}

	return nil
}

// UpdateStorageFile Updates a file in cloud storage using store and name
// Actually just deletes/creates it again
func UpdateStorageFile(config Config, store string, file string, contents []byte) error {
	client, err := getStorageClient(config)
	if err != nil {
		return err
	}

	_, err2 := client.GetFile(store+"/"+file, nil)
	if err2 != nil {
		log.Printf("Error in azure UpdateStorageFile: %s", err2)
		return err2
	}

	// Updating a file in storage is falls back to reading writing individual bytes
	// Probably just easier to delete then add
	err3 := client.DeleteFile(store + "/" + file)
	if err3 != nil {
		log.Printf("Error in azure UpdateStorageFile delete: %s", err3)
		return err3
	}

	err4 := CreateStorageFile(config, store, file, contents)
	if err4 != nil {
		log.Printf("Error in azure UpdateStorageFile create: %s", err4)
		return err4
	}

	return nil
}

// CreateStorageFile Creates a file in cloud storage in a specific store
func CreateStorageFile(config Config, store string, file string, contents []byte) error {
	client, err := getStorageClient(config)
	if err != nil {
		return err
	}

	err2 := client.CreateFile(store+"/"+file, uint64(len(contents)), nil)
	if err2 != nil {
		log.Printf("Error in azure CreateStorageFile create: %s", err2)
		return err2
	}

	if len(contents) <= 0 {
		return nil
	}

	r := bytes.NewReader(contents)
	fileRange := storage.FileRange{
		Start: 0,
		End:   uint64(len(contents)) - 1,
	}
	err3 := client.PutRange(store+"/"+file, r, fileRange)
	if err3 != nil {
		log.Printf("Error in azure CreateStorageFile write: %s", err3)
		return err3
	}

	return nil
}

// This isn't very nice
func convertParameters(parameters TemplateParameterFile) map[string]interface{} {

	myMap := make(map[string]interface{})
	for k, v := range parameters.Parameters {
		myMap[k] = v
	}

	return myMap

	// bytes, err := json.Marshal(parameters)
	// if err != nil {
	// 	log.Printf("Error converting parameters 1: %s", err)
	// 	return nil, err
	// }

	// // And back
	// myMap := make(map[string]interface{})
	// err = json.Unmarshal(bytes, &myMap)
	// if err != nil {
	// 	log.Printf("Error converting parameters 2: %s", err)
	// 	return nil, err
	// }

	// return myMap, nil
}

func replaceParameterVariables(parametersJSON *TemplateParameterFile, number int) error {
	reg, err := regexp.Compile(`(\${n})`)
	if err != nil {
		log.Printf("Invalid regex: %s", err)
		return err
	}

	for key, param := range parametersJSON.Parameters {
		value := param.Value.(string)

		parametersJSON.Parameters[key] = TemplateParameter{Value: reg.ReplaceAllString(value, strconv.Itoa(number))}
	}

	log.Printf("%s", parametersJSON)

	return nil
}

func getDeploymentClient(c Config) (*resources.DeploymentsClient, error) {
	client := resources.NewDeploymentsClient(c.AzureSubscriptionID)

	spt, err := getServicePricipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt

	return &client, nil
}

func getResourcesClient(c Config) (*resources.Client, error) {
	client := resources.NewClient(c.AzureSubscriptionID)

	spt, err := getServicePricipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt

	return &client, nil
}

func getStorageClient(c Config) (*storage.FileServiceClient, error) {
	client, err := storage.NewBasicClient(c.AzureStorageServer, c.AzureStorageKey)
	if err != nil {
		log.Printf("Error in azure getStorageClient: %s", err)
		return nil, err
	}
	fs := client.GetFileService()
	return &fs, nil
}

func getVMClient(c Config) (*compute.VirtualMachinesClient, error) {
	client := compute.NewVirtualMachinesClient(c.AzureSubscriptionID)
	spt, err := getServicePricipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt
	return &client, nil
}

func getServicePricipalToken(config Config) (*azure.ServicePrincipalToken, error) {
	spt, err := newServicePrincipalTokenFromCredentials(config, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Printf("Error getting service principal: %s", err)
		return nil, err
	}

	return spt, nil
}

func newServicePrincipalTokenFromCredentials(c Config, scope string) (*azure.ServicePrincipalToken, error) {
	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(config.AzureTenantID)
	if err != nil {
		panic(err)
	}
	return azure.NewServicePrincipalToken(*oauthConfig, config.AzureClientID, config.AzureClientSecret, scope)
}
