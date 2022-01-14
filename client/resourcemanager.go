package mcmaclient

import (
	"bytes"
	"fmt"
	"github.com/ebu/mcma-libraries-go/model"
	"net/http"
	"reflect"
	"strings"
)

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
				ResourceType: "Service",
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
			if r.hasMatchingHttpEndpoint(url) {
				return r, nil
			}
		}
	}
	return nil, nil
}

func (resourceManager *ResourceManager) Init() error {
	serviceRegistryClient := resourceManager.getServiceRegistryClient()
	resourceManager.services = append(resourceManager.services[:0], serviceRegistryClient)

	servicesEndpoint, found := serviceRegistryClient.GetResourceEndpointClientByType(reflect.TypeOf(model.Service{}))
	if !found {
		return fmt.Errorf("service resource endpoint not found")
	}

	serviceQueryResults, err := servicesEndpoint.Query(reflect.TypeOf(model.Service{}), "", nil)
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

func (resourceManager *ResourceManager) Query(t reflect.Type, filter []struct {
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
		if resourceEndpointClient, matched := s.GetResourceEndpointClientByType(t); matched {
			if _, alreadyUsed := usedHttpEndpoints[resourceEndpointClient.getHttpEndpoint()]; alreadyUsed {
				continue
			}
			anyMatchingClients = true
			if queryResults, err := resourceEndpointClient.Query(t, "", filter); err == nil {
				for _, r := range queryResults.Results {
					results = append(results, r)
				}
				usedHttpEndpoints[resourceEndpointClient.getHttpEndpoint()] = struct{}{}
			} else {
				errs = append(errs, err.Error())
			}
		}
	}
	if !anyMatchingClients {
		return nil, fmt.Errorf("no available resource endpoints for resource of type '%s'", t.String())
	}
	if len(usedHttpEndpoints) == 0 {
		return nil, fmt.Errorf(strings.Join(errs, "\n"))
	}
	return results, nil
}

func (resourceManager *ResourceManager) QueryResources(resourceType string, filter []struct {
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
		if resourceEndpointClient, matched := s.GetResourceEndpointClientByTypeName(resourceType); matched {
			if _, alreadyUsed := usedHttpEndpoints[resourceEndpointClient.getHttpEndpoint()]; alreadyUsed {
				continue
			}
			anyMatchingClients = true
			if queryResults, err := resourceEndpointClient.QueryMaps("", filter); err == nil {
				for _, r := range queryResults.Results {
					results = append(results, r)
				}
				usedHttpEndpoints[resourceEndpointClient.getHttpEndpoint()] = struct{}{}
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

func (resourceManager *ResourceManager) GetResource(resourceType string, resourceId string) (map[string]interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClientByTypeName(resourceType); matched {
			return resourceEndpointClient.GetResource(resourceId)
		}
	}
	resp, err := resourceManager.getMcmaHttpClient().Get(resourceId, false)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	resource, err := readJsonRespBody(resp, reflect.TypeOf(m))
	if err != nil {
		return nil, err
	}

	return resource.(map[string]interface{}), nil
}

func (resourceManager *ResourceManager) Get(t reflect.Type, resourceId string) (interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}
	if t.Kind() != reflect.Map {
		for _, s := range resourceManager.services {
			if resourceEndpointClient, matched := s.GetResourceEndpointClientByType(t); matched {
				return resourceEndpointClient.Get(t, resourceId)
			}
		}
	}
	resp, err := resourceManager.getMcmaHttpClient().Get(resourceId, false)
	if err != nil {
		return nil, err
	}

	resource, err := readJsonRespBody(resp, t)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (resourceManager *ResourceManager) Create(resource interface{}) (interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}

	t := reflect.TypeOf(resource)
	var id string
	if t.Kind() != reflect.Map {
		for _, s := range resourceManager.services {
			if resourceEndpointClient, matched := s.GetResourceEndpointClientByType(t); matched {
				return resourceEndpointClient.Post(t, "", resource)
			}
		}

		resourceValue := reflect.ValueOf(resource)
		idField := resourceValue.FieldByName("Id")
		if idField.IsZero() || idField.Kind() != reflect.String {
			return nil, fmt.Errorf("no resource endpoint available for type '%s' and no id on resource", t.String())
		}
		id = idField.String()
	} else {
		resourceMap := resource.(map[string]interface{})
		resourceType, foundType := resourceMap["@type"]
		if !foundType {
			return nil, fmt.Errorf("@type property not found in map")
		}
		for _, s := range resourceManager.services {
			if resourceEndpointClient, matched := s.GetResourceEndpointClientByTypeName(resourceType.(string)); matched {
				return resourceEndpointClient.PostResource("", resourceMap)
			}
		}

		idVal, foundId := resourceMap["id"]
		if foundId {
			return nil, fmt.Errorf("no resource endpoint available for type '%s' and no id on resource", resourceType)
		}
		id = idVal.(string)
	}

	var jsonBody *bytes.Reader
	var err error
	if jsonBody, err = getJsonReqBody(resource); err != nil {
		return nil, err
	}

	resp, err := resourceManager.getMcmaHttpClient().Post(id, jsonBody)
	if err != nil {
		return nil, err
	}

	createdResource, err := readJsonRespBody(resp, t)
	if err != nil {
		return nil, err
	}

	return createdResource, nil

}

func (resourceManager *ResourceManager) Update(resource interface{}) (interface{}, error) {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return nil, err
		}
	}

	t := reflect.TypeOf(resource)
	var id string
	if t.Kind() != reflect.Map {
		for _, s := range resourceManager.services {
			if resourceEndpointClient, matched := s.GetResourceEndpointClientByType(t); matched {
				return resourceEndpointClient.Put(t, "", resource)
			}
		}
		resourceValue := reflect.ValueOf(resource)
		idField := resourceValue.FieldByName("Id")
		if idField.IsZero() || idField.Kind() != reflect.String {
			return nil, fmt.Errorf("no resource endpoint available for type '%s' and no id on resource", t.String())
		}
		id = idField.String()
	} else {
		resourceMap := resource.(map[string]interface{})
		resourceType, foundType := resourceMap["@type"]
		if !foundType {
			return nil, fmt.Errorf("@type property not found in map")
		}
		for _, s := range resourceManager.services {
			if resourceEndpointClient, matched := s.GetResourceEndpointClientByTypeName(resourceType.(string)); matched {
				return resourceEndpointClient.PutResource("", resourceMap)
			}
		}

		idVal, foundId := resourceMap["id"]
		if foundId {
			return nil, fmt.Errorf("no resource endpoint available for type '%s' and no id on resource", resourceType)
		}
		id = idVal.(string)
	}

	var jsonBody *bytes.Reader
	var err error
	if jsonBody, err = getJsonReqBody(resource); err != nil {
		return nil, err
	}

	resp, err := resourceManager.getMcmaHttpClient().Put(id, jsonBody)
	if err != nil {
		return nil, err
	}

	createdResource, err := readJsonRespBody(resp, t)
	if err != nil {
		return nil, err
	}

	return createdResource, nil
}

func (resourceManager *ResourceManager) DeleteResource(resourceType string, resourceId string) error {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return err
		}
	}
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClientByTypeName(resourceType); matched {
			err := resourceEndpointClient.Delete(resourceId)
			return err
		}
	}
	_, err := resourceManager.getMcmaHttpClient().Delete(resourceId)
	return err
}

func (resourceManager *ResourceManager) Delete(t reflect.Type, resourceId string) error {
	if len(resourceManager.services) == 0 {
		err := resourceManager.Init()
		if err != nil {
			return err
		}
	}
	for _, s := range resourceManager.services {
		if resourceEndpointClient, matched := s.GetResourceEndpointClientByType(t); matched {
			err := resourceEndpointClient.Delete(resourceId)
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
		_, err = resourceEndpoint.Post(nil, notificationEndpoint.HttpEndpoint, notification)
	} else {
		jsonBody, err := getJsonReqBody(notification)
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
