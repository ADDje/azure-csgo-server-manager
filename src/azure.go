package main

import (
	"bytes"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest/azure"
)

func getServicePricipalToken(config Config) (*azure.ServicePrincipalToken, error) {
	spt, err := newServicePrincipalTokenFromCredentials(config, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil, err
	}

	return spt, nil
}

// GetServers Get csgo servers/vms
func GetServers(config Config) (*[]resources.GenericResource, error) {
	client := resources.NewClient(config.AzureSubscriptionID)

	spt, err := getServicePricipalToken(config)
	client.Authorizer = spt

	results, err := client.List("resourceGroup eq '"+config.ResourceGroup+"' and resourceType eq 'Microsoft.Compute/virtualMachines'", "", nil)

	if err != nil {
		log.Printf("Error in azure getServers: %s", err)
		return nil, err
	}

	return results.Value, nil
}

// GetStorageFile Returns file from cloud storage by name and store
func GetStorageFile(config Config, store string, file string) (*storage.FileStream, error) {
	client, err := getStorageClient(config)
	if err != nil {
		return nil, err
	}

	fileStream, err2 := client.GetFile(store+"/"+file, nil)
	if err2 != nil {
		log.Printf("Error in azure getStorageFile: %s", err2)
		return nil, err2
	}

	return fileStream, nil
}

// GetStorageFileText Returns file contents from cloud storage by name and store
func GetStorageFileText(config Config, store string, file string) (string, error) {
	fileStream, err := GetStorageFile(config, store, file)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, fileStream.Properties.ContentLength)
	if fileStream.Properties.ContentLength > 0 {
		_, err3 := fileStream.Body.Read(buffer)
		if err3 != nil && err3 != io.EOF {
			log.Printf("Error in azure GetStorageFileText: %s", err3)
			return "", err3
		}
	}

	return string(buffer), nil
}

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

func getStorageClient(c Config) (*storage.FileServiceClient, error) {
	client, err := storage.NewBasicClient(config.AzureStorageServer, config.AzureStorageKey)
	if err != nil {
		log.Printf("Error in azure getStorageClient: %s", err)
		return nil, err
	}
	fs := client.GetFileService()
	return &fs, nil
}

func newServicePrincipalTokenFromCredentials(c Config, scope string) (*azure.ServicePrincipalToken, error) {
	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(config.AzureTenantID)
	if err != nil {
		panic(err)
	}
	return azure.NewServicePrincipalToken(*oauthConfig, config.AzureClientID, config.AzureClientSecret, scope)
}
