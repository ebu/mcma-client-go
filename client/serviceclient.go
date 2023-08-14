package mcmaclient

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/ebu/mcma-libraries-go/model"
)

type ServiceClient struct {
	authProvider    *AuthProvider
	httpClient      *http.Client
	service         model.Service
	tracker         *model.McmaTracker
	resources       []*ResourceEndpointClient
	resourcesByType map[string]*ResourceEndpointClient
}

func (serviceClient *ServiceClient) loadResources() {
	if serviceClient.resourcesByType != nil {
		return
	}
	serviceClient.resourcesByType = make(map[string]*ResourceEndpointClient)
	for _, r := range serviceClient.service.Resources {
		resourceEndpointClient := &ResourceEndpointClient{
			authProvider:     serviceClient.authProvider,
			httpClient:       serviceClient.httpClient,
			resourceEndpoint: r,
			serviceAuthType:  serviceClient.service.AuthType,
			tracker:          serviceClient.tracker,
		}
		serviceClient.resources = append(serviceClient.resources, resourceEndpointClient)
		serviceClient.resourcesByType[r.ResourceType] = resourceEndpointClient
	}
}

func (serviceClient *ServiceClient) GetResourceEndpointClientByType(t reflect.Type) (*ResourceEndpointClient, bool) {
	resourceTypeParts := strings.Split(t.String(), ".")
	resourceType := resourceTypeParts[len(resourceTypeParts)-1]
	return serviceClient.GetResourceEndpointClientByTypeName(resourceType)
}

func (serviceClient *ServiceClient) GetResourceEndpointClientByTypeName(resourceType string) (*ResourceEndpointClient, bool) {
	serviceClient.loadResources()
	resourceEndpoint, found := serviceClient.resourcesByType[resourceType]
	return resourceEndpoint, found
}

func (serviceClient *ServiceClient) GetResourceEndpointClientByTypeAndUrl(t reflect.Type, url string) (*ResourceEndpointClient, bool) {
	resourceEndpoint, found := serviceClient.GetResourceEndpointClientByType(t)
	if !found {
		return nil, false
	}
	if !resourceEndpoint.hasMatchingHttpEndpoint(url) {
		return nil, false
	}
	return resourceEndpoint, true
}

func (serviceClient *ServiceClient) GetResourceEndpointClientByTypeNameAndUrl(resourceType string, url string) (*ResourceEndpointClient, bool) {
	resourceEndpoint, found := serviceClient.GetResourceEndpointClientByTypeName(resourceType)
	if !found {
		return nil, false
	}
	if !resourceEndpoint.hasMatchingHttpEndpoint(url) {
		return nil, false
	}
	return resourceEndpoint, true
}
