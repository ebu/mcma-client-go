package mcmaclient

import (
	"github.com/ebu/mcma-libraries-go/model"
	"net/http"
	"reflect"
	"strings"
)

type ServiceClient struct {
	authProvider    *AuthProvider
	httpClient      *http.Client
	service         model.Service
	tracker         model.McmaTracker
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
			authProvider:       serviceClient.authProvider,
			httpClient:         serviceClient.httpClient,
			resourceEndpoint:   r,
			serviceAuthType:    serviceClient.service.AuthType,
			serviceAuthContext: serviceClient.service.AuthContext,
			tracker:            serviceClient.tracker,
		}
		serviceClient.resources = append(serviceClient.resources, resourceEndpointClient)
		serviceClient.resourcesByType[r.ResourceType] = resourceEndpointClient
	}
}

func (serviceClient *ServiceClient) GetResourceEndpointClient(t reflect.Type) (*ResourceEndpointClient, bool) {
	resourceTypeParts := strings.Split(t.String(), ".")
	resourceType := resourceTypeParts[len(resourceTypeParts)-1]
	serviceClient.loadResources()
	resourceEndpoint, found := serviceClient.resourcesByType[resourceType]
	return resourceEndpoint, found
}
