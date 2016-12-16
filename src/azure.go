package main

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest/azure"
)

func getServicePricipalToken(config Config) (*azure.ServicePrincipalToken, error) {
	c := map[string]string{
		"AZURE_CLIENT_ID":       config.AzureClientID,
		"AZURE_CLIENT_SECRET":   config.AzureClientSecret,
		"AZURE_SUBSCRIPTION_ID": config.AzureSubscriptionID,
		"AZURE_TENANT_ID":       config.AzureTenantID,
	}

	spt, err := newServicePrincipalTokenFromCredentials(c, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil, err
	}

	return spt, nil
}

func getServers(config Config) (*[]resources.GenericResource, error) {
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

func getStorageFile(config Config, store string, file storage.File) (*storage.FileStream, error) {
	client, err := storage.NewBasicClient("test", "test")
	if err != nil {
		log.Printf("Error in azure getStorageFile: %s", err)
		return nil, err
	}
	fileService := client.GetFileService()

	fileStream, err2 := fileService.GetFile(store+"/"+file.Name, nil)
	if err2 != nil {
		log.Printf("Error in azure getStorageFile")
		return nil, err
	}

	return fileStream, nil
}

func getStorageFiles(config Config, store string) ([]storage.File, error) {
	client, err := storage.NewBasicClient("test", "test")
	if err != nil {
		log.Printf("Error in azure getStorageFiles: %s", err)
		return nil, err
	}
	fileService := client.GetFileService()

	exists, err2 := fileService.ShareExists(store)
	if err2 != nil || !exists {
		log.Printf("Error in azure getStorageFiles: %s", err)
		return nil, err
	}

	params := storage.ListDirsAndFilesParameters{}
	files, err3 := fileService.ListDirsAndFiles(store, params)
	if err3 != nil {
		log.Printf("Error in azure getStorageFiles: %s", err)
		return nil, err
	}

	return files.Files, nil
}

func newServicePrincipalTokenFromCredentials(c map[string]string, scope string) (*azure.ServicePrincipalToken, error) {
	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(c["AZURE_TENANT_ID"])
	if err != nil {
		panic(err)
	}
	return azure.NewServicePrincipalToken(*oauthConfig, c["AZURE_CLIENT_ID"], c["AZURE_CLIENT_SECRET"], scope)
}
