package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

const DEPLOYMENT_NAME = "csgo-server-manager"
const VHD_CONTAINER_NAME = "vhds"
const FILE_CONTAINER_NAME = "server-manager"

func DeployXTemplates(x int, config Config, vmName string, adminUserName string,
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
	err = json.NewDecoder(*parameters).Decode(&parametersJSON)
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

	configLink := GetStorageFileLink(config, CONFIG_FILE_STORE, configName)

	// Add config file, needed for automated config deployment
	parametersJSON.Parameters["configFileUrl"] = TemplateParameter{Value: configLink}
	parametersJSON.Parameters["configFileName"] = TemplateParameter{Value: configName}

	for t := 1; t <= x; t++ {
		log.Printf("Deploying server: %d of %d", t, x)
		go deployTemplate(client, config, t, vmName, adminUserName, adminPassword, configName, templateUri, parametersJSON)
	}

	return nil
}

func deployTemplate(client *resources.DeploymentsClient, config Config, number int, vmName string, adminUserName string,
	adminPassword string, configName string, templateUri string, parametersJSON TemplateParameterFile) error {

	// Replace any variables in parameters with their values
	replacedParametersJson, err := replaceParameterVariables(parametersJSON, number)
	if err != nil {
		return err
	}

	deploymentName := DEPLOYMENT_NAME + "-" + strconv.Itoa(number)

	exportParameters := convertParameters(*replacedParametersJson)

	template := resources.TemplateLink{
		URI: &templateUri,
	}
	properties := resources.DeploymentProperties{
		TemplateLink: &template,
		Parameters:   &exportParameters,
		Mode:         "Incremental",
	}
	deployment := resources.Deployment{Properties: &properties}

	resChan, errChan := client.CreateOrUpdate(config.ResourceGroup, deploymentName, deployment, nil)
	select {
	case err = <-errChan:
		log.Printf("Error creating deployment: %s", err)
		return err
	case _ = <-resChan:
		return nil
	}
}

func StartVM(config Config, vmName string) error {

	log.Printf("Starting VM: %s", vmName)

	client, err := getVMClient(config)
	if err != nil {
		return err
	}

	resChan, errChan := client.Start(config.ResourceGroup, vmName, nil)
	select {
	case err = <-errChan:
		log.Printf("Error Starting VM: %s", err)
		return err
	case _ = <-resChan:
		return nil
	}
}

func DeallocateVM(config Config, vmName string) error {

	log.Printf("Deallocating VM: %s", vmName)

	client, err := getVMClient(config)
	if err != nil {
		return err
	}

	resChan, errChan := client.Deallocate(config.ResourceGroup, vmName, nil)
	select {
	case err = <-errChan:
		log.Printf("Error deallocating VM: %s", err)
		return err
	case _ = <-resChan:
		return nil
	}
}

func DeleteVM(config Config, vmName string) error {
	client, err := getVMClient(config)
	if err != nil {
		return err
	}

	resChan, errChan := client.Delete(config.ResourceGroup, vmName, nil)
	select {
	case err = <-errChan:
		log.Printf("Error deleting VM: %s", err)
		return err
	case _ = <-resChan:
		return nil
	}
}

func DeleteVhd(config Config, vhdUri string) error {
	//client, err := getResourcesClient(config)
	//if err != nil {
	//log.Printf("Delete Vhd Error: %s", err)
	//return err
	//}

	log.Printf("vhd uri: %s", vhdUri)
	//parts := strings.Split(vhdUri, "/")
	//name := parts[len(parts)-1]
	//store := parts[len(parts)-2 : len(parts)-1][0]

	return nil
}

func GetNicDetails(config Config, nicName string) (*network.Interface, error) {
	networkClient, err := getInterfacesClient(config)
	if err != nil {
		return nil, err
	}

	nicDetails, err := networkClient.Get(config.ResourceGroup, nicName, "")
	if err != nil {
		log.Printf("Could not get NIC Details for %s: %s", nicName, err)
		return nil, err
	}

	return &nicDetails, nil
}

func GetIpDetails(config Config, ipId string) (*network.PublicIPAddress, error) {
	ipClient, err := getIpClient(config)
	if err != nil {
		return nil, err
	}

	log.Printf("Getting IP... %s", ipId)
	ipDetails, err := ipClient.Get(config.ResourceGroup, ipId, "")
	if err != nil {
		return nil, err
	}

	if ipDetails.IPAddress == nil {
		err := errors.New("IP Not ready")
		log.Print(err)
		return nil, err
	}

	log.Printf("%s", *ipDetails.IPAddress)

	return &ipDetails, nil
}

func DeleteVMNetworkThings(config Config, vmProps *compute.VirtualMachineProperties) error {

	resourcesClient, err := getResourcesClient(config)
	if err != nil {
		return err
	}

	networkClient, err := getInterfacesClient(config)
	if err != nil {
		return err
	}

	nic := *(*vmProps.NetworkProfile.NetworkInterfaces)[0].ID
	nicParts := strings.Split(nic, "/")
	nicName := nicParts[len(nicParts)-1]

	nicDetails, err := networkClient.Get(config.ResourceGroup, nicName, "")
	if err != nil {
		log.Printf("Could not get NIC Details for %s: %s", nicName, err)
		return err
	}

	ip := (*nicDetails.IPConfigurations)[0]
	ipID := *ip.PublicIPAddress.ID

	netParts := strings.Split(*ip.Subnet.ID, "/")
	net := strings.Join(netParts[0:len(netParts)-2], "/")

	log.Printf("Deleting: \nNIC: %s, \nIP:  %s, \nNet: %s", nic, ipID, net)

	resChan, errChan := resourcesClient.DeleteByID(nic, nil)
	select {
	case err = <-errChan:
		log.Printf("Could not delete NIC %s: %s", nic, err)
		return err
	case _ = <-resChan:
		log.Printf("NIC Deleted: %s", nic)
	}

	resChan, errChan = resourcesClient.DeleteByID(ipID, nil)
	select {
	case err = <-errChan:
		log.Printf("Could not delete IP %s: %s", ipID, err)
		return err
	case _ = <-resChan:
		log.Printf("Deleted IP:  %s", ipID)
	}

	resChan, errChan = resourcesClient.DeleteByID(net, nil)
	select {
	case err = <-errChan:
		log.Printf("Could not delete Network %s: %s", net, err)
		return err
	case _ = <-resChan:
		log.Printf("Deleted Net: %s", net)
	}

	return nil
}

func GetVmIp(config Config, vmName string) (*string, error) {
	vmInfo, err := GetVmProperties(config, vmName)
	if err != nil {
		return nil, err
	}

	return GetVmIpWithProperties(config, vmInfo)
}

func GetVmIpWithProperties(config Config, props *compute.VirtualMachineProperties) (*string, error) {
	nic := *(*props.NetworkProfile.NetworkInterfaces)[0].ID
	nicParts := strings.Split(nic, "/")
	nicName := nicParts[len(nicParts)-1]

	nicDetails, err := GetNicDetails(config, nicName)
	if err != nil {
		return nil, err
	}
	nicDetails2 := *nicDetails

	ip := (*nicDetails2.IPConfigurations)[0]
	pubIP := *ip.PublicIPAddress

	ipParts := strings.Split(*pubIP.ID, "/")
	ipID := ipParts[len(ipParts)-1]

	ipDetails, err := GetIpDetails(config, ipID)
	if err != nil {
		return nil, err
	}

	return ipDetails.IPAddress, nil
}

// FullDeleteVM deletes a VM, its network resources and ip address
// Unfortunately this isn't very dynamic, and may fail using a different template
// or with different resources
func FullDeleteVM(config Config, vmName string) error {

	log.Printf("Full deleting VM: %s", vmName)

	vmDetails, err := GetVmProperties(config, vmName)
	if err != nil {
		return err
	}

	err = DeleteVM(config, vmName)
	if err != nil {
		return err
	}

	err = DeleteVMNetworkThings(config, vmDetails)
	if err != nil {
		return err
	}

	vhdURI := vmDetails.StorageProfile.OsDisk.Vhd.URI
	err = DeleteVhd(config, *vhdURI)
	if err != nil {
		return err
	}

	log.Printf("Finished deleting VM: %s", vmName)

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

	var wg sync.WaitGroup
	wg.Add(len(*results.Value))
	for k, vm := range *results.Value {
		go func(k int, vm compute.VirtualMachine) {
			vmInfo, err := GetVmProperties(config, *vm.Name)
			if err == nil {
				(*results.Value)[k].VirtualMachineProperties.InstanceView = vmInfo.InstanceView
			}

			wg.Done()
		}(k, vm)
	}
	wg.Wait()

	return results.Value, nil
}

func GetVmProperties(config Config, vmName string) (*compute.VirtualMachineProperties, error) {
	client, err := getVMClient(config)
	if err != nil {
		return nil, err
	}

	results, err := client.Get(config.ResourceGroup, vmName, "instanceView")
	if err != nil {
		log.Printf("Error in azure GetVmStatus: %s", err)
		return nil, err
	}

	return results.VirtualMachineProperties, nil
}

// GetStorageFileLink returns an externally accessible link to a storage file
// Uses a secure token so should not be used frivolously
func GetStorageFileLink(config Config, store, file string) string {
	return "https://" + config.AzureStorageServer + ".blob.core.windows.net/" + FILE_CONTAINER_NAME + "/" + store +
		"/" + file + config.AzureSASToken
}

// GetStorageFile Returns file from cloud storage by name and store
func GetStorageFile(config Config, store, file string) (*io.ReadCloser, error) {
	return GetRawStorageFile(config, store+"/"+file)
}

// GetRawStorageFile Returns file from cloud storage using exact name
func GetRawStorageFile(config Config, file string) (*io.ReadCloser, error) {
	client, err := getStorageClient(config)
	if err != nil {
		return nil, err
	}

	container, err := getFileContainer(client)
	if err != nil {
		return nil, err
	}

	options := storage.GetBlobOptions{}

	blob := container.GetBlobReference(file)
	fileStream, err := blob.Get(&options)
	if err != nil {
		log.Printf("Error in azure GetStorageFile: %s", err)
		return nil, err
	}

	return &fileStream, nil
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
// 		r, err := fileStream.Body.Read(buffer)
// 		log.Printf("%d bytes read", r)
// 		if err != nil && err != io.EOF {
// 			log.Printf("Error in azure GetStorageFileText: %s", err)
// 			return "", err
// 		}
// 	}

// 	return string(buffer), nil
// }

// GetStorageFiles Returns files from cloud storage by name and store
func GetStorageFiles(config Config, store string) ([]storage.Blob, error) {
	client, err := getStorageClient(config)
	if err != nil {
		return nil, err
	}

	container, err := getFileContainer(client)
	if err != nil {
		return nil, err
	}

	params := storage.ListBlobsParameters{
		Prefix: store + "/",
	}

	blobs, err := container.ListBlobs(params)
	if err != nil {
		log.Printf("Error in azure GetStorageFiles: %s", err)
		return nil, err
	}

	return blobs.Blobs, nil
}

// DeleteStorageFile Deletes a file in cloud storage using name and store
func DeleteStorageFile(config Config, store string, file string) error {
	client, err := getStorageClient(config)
	if err != nil {
		return err
	}

	fileName := store + "/" + file
	log.Printf("Deleting Azure File: %s", fileName)

	container, err := getFileContainer(client)
	if err != nil {
		return err
	}

	params := storage.DeleteBlobOptions{}

	blob := container.GetBlobReference(fileName)
	err = blob.Delete(&params)
	if err != nil {
		log.Printf("Error in azure DeleteStorageFile delete: %s", err)
		return err
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

	fileName := store + "/" + file
	log.Printf("Updating Azure File: %s", fileName)

	container, err := getFileContainer(client)
	if err != nil {
		return err
	}

	blob := container.GetBlobReference(fileName)
	params := storage.DeleteBlobOptions{}

	// Updating a file in storage is falls back to reading writing individual bytes
	// Probably just easier to delete then add
	err = blob.Delete(&params)
	if err != nil {
		log.Printf("Error in azure UpdateStorageFile delete: %s", err)
		return err
	}

	err = CreateStorageFile(config, store, file, contents)
	if err != nil {
		log.Printf("Error in azure UpdateStorageFile create: %s", err)
		return err
	}

	return nil
}

// CreateStorageFile Creates a file in cloud storage in a specific store
func CreateStorageFile(config Config, store string, file string, contents []byte) error {
	client, err := getStorageClient(config)
	if err != nil {
		return err
	}

	container, err := getFileContainer(client)
	if err != nil {
		return err
	}

	blob := container.GetBlobReference(store + "/" + file)
	params := storage.PutPageOptions{}
	writeRange := storage.BlobRange{Start: 0, End: uint64(len(contents))}

	r := bytes.NewReader(contents)
	err = blob.WriteRange(writeRange, r, &params)
	if err != nil {
		log.Printf("Error in azure CreateStorageFile create: %s", err)
		return err
	}

	return nil
}

func getFileContainer(c *storage.BlobStorageClient) (*storage.Container, error) {
	container := c.GetContainerReference(FILE_CONTAINER_NAME)

	exists, err := container.Exists()
	if err != nil {
		log.Printf("Could not get container. %s", err)
		return nil, err
	} else if !exists {
		err = errors.New(fmt.Sprintf("File container does not exist: %s", FILE_CONTAINER_NAME))
		log.Printf("Could not get container. %s", err)
		return nil, err
	}

	return container, nil
}

func convertParameters(parameters TemplateParameterFile) map[string]interface{} {

	myMap := make(map[string]interface{})
	for k, v := range parameters.Parameters {
		myMap[k] = v
	}

	return myMap

}

func replaceParameter(template string, number int) (string, error) {
	reg, err := regexp.Compile(`(\${n})`)
	if err != nil {
		log.Printf("Invalid regex: %s", err)
		return "", err
	}

	return reg.ReplaceAllString(template, strconv.Itoa(number)), nil
}

func replaceParameterVariables(parametersJSON TemplateParameterFile, number int) (*TemplateParameterFile, error) {
	outParams := TemplateParameterFile{
		Schema:         parametersJSON.Schema,
		ContentVersion: parametersJSON.ContentVersion,
		Parameters:     make(map[string]TemplateParameter),
	}
	reg, err := regexp.Compile(`(\${n})`)
	if err != nil {
		log.Printf("Invalid regex: %s", err)
		return nil, err
	}

	for key, param := range parametersJSON.Parameters {
		value := param.Value.(string)

		outParams.Parameters[key] = TemplateParameter{Value: reg.ReplaceAllString(value, strconv.Itoa(number))}
	}

	return &outParams, nil
}

func getDeploymentClient(c Config) (*resources.DeploymentsClient, error) {
	client := resources.NewDeploymentsClient(c.AzureSubscriptionID)

	spt, err := getServicePrincipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt

	return &client, nil
}

func getResourcesClient(c Config) (*resources.GroupClient, error) {
	client := resources.NewGroupClient(c.AzureSubscriptionID)

	spt, err := getServicePrincipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt

	return &client, nil
}

func getStorageClient(c Config) (*storage.BlobStorageClient, error) {
	client, err := storage.NewBasicClient(c.AzureStorageServer, c.AzureStorageKey)
	if err != nil {
		log.Printf("Error in azure getStorageClient: %s", err)
		return nil, err
	}
	bs := client.GetBlobService()
	return &bs, nil
}

func getInterfacesClient(c Config) (*network.InterfacesClient, error) {
	client := network.NewInterfacesClient(c.AzureSubscriptionID)
	spt, err := getServicePrincipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt
	return &client, nil
}

func getIpClient(c Config) (*network.PublicIPAddressesClient, error) {
	client := network.NewPublicIPAddressesClient(c.AzureSubscriptionID)
	spt, err := getServicePrincipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt
	return &client, nil
}

func getVMClient(c Config) (*compute.VirtualMachinesClient, error) {
	client := compute.NewVirtualMachinesClient(c.AzureSubscriptionID)
	spt, err := getServicePrincipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt
	return &client, nil
}

func getDiskClient(c Config) (*disk.DisksClient, error) {
	client := disk.NewDisksClient(c.AzureSubscriptionID)
	spt, err := getServicePrincipalToken(c)
	if err != nil {
		return nil, err
	}
	client.Authorizer = spt
	return &client, nil
}

func getServicePrincipalToken(config Config) (*autorest.BearerAuthorizer, error) {
	spt, err := newServicePrincipalTokenFromCredentials(config, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Printf("Error getting service principal: %s", err)
		return nil, err
	}

	return autorest.NewBearerAuthorizer(spt), nil
}

func newServicePrincipalTokenFromCredentials(config Config, scope string) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, config.AzureTenantID)
	if err != nil {
		panic(err)
	}
	return adal.NewServicePrincipalToken(*oauthConfig, config.AzureClientID, config.AzureClientSecret, scope)
}

func getBlobName(blob string) string {
	parts := strings.Split(blob, "/")
	return parts[len(parts)-1]
}
