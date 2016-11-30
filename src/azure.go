package main

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest/azure"
)

func getServers(config Config) (*[]resources.GenericResource, error) {

	c := map[string]string{
		"AZURE_CLIENT_ID":       config.AzureClientID,
		"AZURE_CLIENT_SECRET":   config.AzureClientSecret,
		"AZURE_SUBSCRIPTION_ID": config.AzureSubscriptionID,
		"AZURE_TENANT_ID":       config.AzureTenantID,
	}

	log.Print(c)

	spt, err := newServicePrincipalTokenFromCredentials(c, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil, err
	}

	client := resources.NewClient(c["AZURE_SUBSCRIPTION_ID"])

	client.Authorizer = spt

	results, err := client.List("resourceGroup eq '"+config.ResourceGroup+"' and resourceType eq 'Microsoft.Compute/virtualMachines'", "", nil)

	if err != nil {
		log.Printf("Error in azure getServers: %s", err)
		return nil, err
	}

	return results.Value, nil
}

func newServicePrincipalTokenFromCredentials(c map[string]string, scope string) (*azure.ServicePrincipalToken, error) {
	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(c["AZURE_TENANT_ID"])
	if err != nil {
		panic(err)
	}
	return azure.NewServicePrincipalToken(*oauthConfig, c["AZURE_CLIENT_ID"], c["AZURE_CLIENT_SECRET"], scope)
}
