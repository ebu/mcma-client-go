package mcmaclient

import (
	"../model"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var serviceType = reflect.TypeOf(model.Service{}).String()

type ResourceManager struct {
	authProvider          *AuthProvider
	httpClient            *http.Client
	mcmaHttpClient        *McmaHttpClient
	servicesUrl           string
	servicesAuthType      string
	servicesAuthContext   string
	tracker               model.McmaTracker
	services              []*ServiceClient
	serviceRegistryClient *ServiceClient
}

func (resourceManager *ResourceManager) getServiceRegistryData() model.Service {
	return model.Service{
		Name:        "Service Registry",
		AuthType:    resourceManager.servicesAuthType,
		AuthContext: resourceManager.servicesAuthContext,
		Resources: []model.ResourceEndpoint{
			{
				ResourceType: serviceType,
				HttpEndpoint: resourceManager.servicesUrl,
			},
		},
	}
}

func (resourceManager *ResourceManager) getMcmaHttpClient() *McmaHttpClient {
	if resourceManager.mcmaHttpClient == nil {
		resourceManager.mcmaHttpClient = &McmaHttpClient{
			httpClient: resourceManager.httpClient,
		}
	}
	return resourceManager.mcmaHttpClient
}

func (resourceManager *ResourceManager) getServiceRegistryClient() *ServiceClient {
	if resourceManager.serviceRegistryClient == nil {
		resourceManager.serviceRegistryClient = &ServiceClient{
			authProvider: resourceManager.authProvider,
			httpClient:   resourceManager.httpClient,
			service:      resourceManager.getServiceRegistryData(),
			tracker:      resourceManager.tracker,
		}
	}
	return resourceManager.serviceRegistryClient
}

func (resourceManager *ResourceManager) Init() error {
	serviceRegistryClient := resourceManager.getServiceRegistryClient()
	resourceManager.services = append(resourceManager.services[:0], serviceRegistryClient)

	servicesEndpoint, found := serviceRegistryClient.GetResourceEndpointClient(serviceType)
	if !found {
		return fmt.Errorf("service resource endpoint not found")
	}

	serviceQueryResults, err := servicesEndpoint.Query("", nil)
	if err != nil {
		return err
	}

	for _, r := range serviceQueryResults.Results {
		service := r.(model.Service)
		serviceClient := &ServiceClient{
			authProvider: resourceManager.authProvider,
			httpClient:   resourceManager.httpClient,
			service:      service,
			tracker:      resourceManager.tracker,
		}
		resourceManager.services = append(resourceManager.services, serviceClient)
	}

	return nil
}

func (resourceManager *ResourceManager) getResourceEndpoint(url string) (*ResourceEndpointClient, error) {
	if url == "" {
		return nil, nil
	}
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}
	for _, s := range resourceManager.services {
		s.loadResources()
		for _, r := range s.resources {
			if strings.HasPrefix(strings.ToLower(url), strings.ToLower(r.httpEndpoint)) {
				return r, nil
			}
		}
	}
	return nil, nil
}

func (resourceManager *ResourceManager) Query(resourceType string, filter []struct {
	key   string
	value string
}) ([]interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}
	anyMatchingClients := false
	usedHttpEndpoints := make(map[string]struct{})
	var results []interface{}
	var errs []string
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClient(resourceType); matched {
			if _, alreadyUsed := usedHttpEndpoints[resourceEndpointClient.httpEndpoint]; alreadyUsed {
				continue
			}
			anyMatchingClients = true
			if queryResults, err := resourceEndpointClient.Query("", filter); err == nil {
				for _, r := range queryResults.Results {
					results = append(results, r)
				}
				usedHttpEndpoints[resourceEndpointClient.httpEndpoint] = struct{}{}
			} else {
				errs = append(errs, err.Error())
			}
		}
	}
	if !anyMatchingClients {
		return nil, fmt.Errorf("no available resource endpoints for resource of type '%s'", resourceType)
	}
	if len(usedHttpEndpoints) == 0 {
		return nil, fmt.Errorf(strings.Join(errs, "\n"))
	}
	return results, nil
}

func (resourceManager *ResourceManager) Get(resourceType string, resourceId string) (interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClient(resourceType); matched {
			return resourceEndpointClient.Get(resourceId)
		}
	}
	return resourceManager.getMcmaHttpClient().Get(resourceId)
}

func (resourceManager *ResourceManager) Create(resource interface{}) (interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}
	resourceType := reflect.TypeOf(resource).String()
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClient(resourceType); matched {
			return resourceEndpointClient.Post("", resource)
		}
	}
	resourceValue := reflect.ValueOf(resource)
	idField := resourceValue.FieldByName("Id")
	if idField.IsZero() || idField.Kind() != reflect.String {
		return nil, fmt.Errorf("no resource endpoint available for type '%s' and no id on resource", resourceType)
	}
	var jsonBody io.Reader
	var err error
	if jsonBody, err = getJsonBody(resource); err != nil {
		return nil, err
	}
	return resourceManager.getMcmaHttpClient().Post(idField.String(), jsonBody)
}

func (resourceManager *ResourceManager) Update(resource interface{}) (interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}
	resourceType := reflect.TypeOf(resource).String()
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClient(resourceType); matched {
			return resourceEndpointClient.Put("", resource)
		}
	}
	resourceValue := reflect.ValueOf(resource)
	idField := resourceValue.FieldByName("Id")
	if idField.IsZero() || idField.Kind() != reflect.String {
		return nil, fmt.Errorf("no resource endpoint available for type '%s' and no id on resource", resourceType)
	}
	var jsonBody io.Reader
	var err error
	if jsonBody, err = getJsonBody(resource); err != nil {
		return nil, err
	}
	return resourceManager.getMcmaHttpClient().Put(idField.String(), jsonBody)
}

func (resourceManager *ResourceManager) Delete(resourceType string, resourceId string) error {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return err
		}
	}
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClient(resourceType); matched {
			_, err := resourceEndpointClient.Delete(resourceId)
			return err
		}
	}
	_, err := resourceManager.getMcmaHttpClient().Delete(resourceId)
	return err
}

func (resourceManager *ResourceManager) SendNotification(resourceId string, resource interface{}, notificationEndpoint model.NotificationEndpoint) error {
	if notificationEndpoint.HttpEndpoint == "" {
		return nil
	}

	notification := model.Notification{
		Source:  resourceId,
		Content: resource,
	}

	var resourceEndpoint *ResourceEndpointClient
	var err error
	if resourceEndpoint, err = resourceManager.getResourceEndpoint(notificationEndpoint.HttpEndpoint); err != nil {
		return err
	}
	if resourceEndpoint != nil {
		_, err = resourceEndpoint.Post(notificationEndpoint.HttpEndpoint, notification)
	} else {
		jsonBody, err := getJsonBody(notification)
		if err == nil {
			_, err = resourceManager.getMcmaHttpClient().Post(notificationEndpoint.HttpEndpoint, jsonBody)
		}
	}
	return err
}

func (resourceManager *ResourceManager) SetHttpClient(httpClient *http.Client) {
	resourceManager.httpClient = httpClient
	resourceManager.mcmaHttpClient = &McmaHttpClient{httpClient: httpClient}
	resourceManager.services = resourceManager.services[:0]
}

func (resourceManager *ResourceManager) AddAuth(authType string, factory AuthenticatorFactory) {
	resourceManager.authProvider.Add(authType, factory)
}

func NewResourceManager(servicesUrl string, servicesAuthType string, servicesAuthContext string) ResourceManager {
	return ResourceManager{
		authProvider:        newAuthProvider(),
		httpClient:          &http.Client{},
		servicesUrl:         servicesUrl,
		servicesAuthType:    servicesAuthType,
		servicesAuthContext: servicesAuthContext,
	}
}

func NewResourceManagerWithTracker(servicesUrl string, servicesAuthType string, servicesAuthContext string, tracker model.McmaTracker) ResourceManager {
	return ResourceManager{
		authProvider:        newAuthProvider(),
		httpClient:          &http.Client{},
		servicesUrl:         servicesUrl,
		servicesAuthType:    servicesAuthType,
		servicesAuthContext: servicesAuthContext,
		tracker:             tracker,
	}
}

func NewResourceManagerNoAuth(servicesUrl string) ResourceManager {
	return ResourceManager{
		authProvider: newAuthProvider(),
		httpClient:   &http.Client{},
		servicesUrl:  servicesUrl,
	}
}

func NewResourceManagerNoAuthWithTracker(servicesUrl string, tracker model.McmaTracker) ResourceManager {
	return ResourceManager{
		authProvider: newAuthProvider(),
		httpClient:   &http.Client{},
		servicesUrl:  servicesUrl,
		tracker:      tracker,
	}
}
