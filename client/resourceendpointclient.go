package mcmaclient

import (
	"../model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ResourceEndpointClient struct {
	authProvider       *AuthProvider
	httpClient         *http.Client
	resourceEndpoint   model.ResourceEndpoint
	serviceAuthType    string
	serviceAuthContext interface{}
	tracker            model.McmaTracker
	mcmaHttpClient     *McmaHttpClient
	httpEndpoint       string
}

func (resourceEndpointClient *ResourceEndpointClient) getMcmaHttpClient() (*McmaHttpClient, error) {
	if resourceEndpointClient.mcmaHttpClient != nil {
		return resourceEndpointClient.mcmaHttpClient, nil
	}

	authType := resourceEndpointClient.resourceEndpoint.AuthType
	if authType == "" {
		authType = resourceEndpointClient.serviceAuthType
	}

	authContext := resourceEndpointClient.resourceEndpoint.AuthContext
	if authContext == "" {
		authContext = resourceEndpointClient.serviceAuthContext
	}

	var authenticator Authenticator
	var err error
	if resourceEndpointClient.authProvider != nil && authType != "" {
		authenticator, err = resourceEndpointClient.authProvider.Get(authType, authContext)
		if err != nil {
			return nil, err
		}
	}

	resourceEndpointClient.mcmaHttpClient = &McmaHttpClient{
		httpClient:    resourceEndpointClient.httpClient,
		authenticator: &authenticator,
		tracker:       resourceEndpointClient.tracker,
	}

	return resourceEndpointClient.mcmaHttpClient, nil
}

func (resourceEndpointClient *ResourceEndpointClient) getFullUrl(url string) string {
	if url == "" {
		return resourceEndpointClient.httpEndpoint
	}
	if strings.HasPrefix(strings.ToLower(url), strings.ToLower(resourceEndpointClient.httpEndpoint)) {
		return url
	}
	return strings.TrimSuffix(resourceEndpointClient.httpEndpoint, "/") + "/" + url
}

func (resourceEndpointClient *ResourceEndpointClient) Query(url string, queryParameters []struct {
	key   string
	value string
}) (model.QueryResults, error) {
	mcmaHttpClient, err := resourceEndpointClient.getMcmaHttpClient()
	if err != nil {
		return model.QueryResults{}, err
	}

	url = resourceEndpointClient.getFullUrl(url)
	if len(queryParameters) > 0 {
		url += "?"
		for _, p := range queryParameters {
			url += fmt.Sprintf("%s=%s", p.key, p.value)
		}
	}

	getResp, err := mcmaHttpClient.Get(url)
	if err != nil {
		return model.QueryResults{}, err
	}

	var body []byte
	_, err = getResp.Body.Read(body)
	if err != nil {
		return model.QueryResults{}, err
	}

	queryResults := model.QueryResults{}
	err = json.Unmarshal(body, &queryResults)

	return queryResults, err
}

func (resourceEndpointClient *ResourceEndpointClient) execute(url string, body interface{}, execute func(mcmaHttpClient *McmaHttpClient, url string, body io.Reader) (interface{}, error)) (interface{}, error) {
	mcmaHttpClient, err := resourceEndpointClient.getMcmaHttpClient()
	if err != nil {
		return nil, err
	}

	url = resourceEndpointClient.getFullUrl(url)
	reqBody, err := getJsonBody(body)
	if err != nil {
		return nil, err
	}

	return execute(mcmaHttpClient, url, reqBody)
}

func (resourceEndpointClient *ResourceEndpointClient) Get(url string) (interface{}, error) {
	return resourceEndpointClient.execute(url, nil, func(client *McmaHttpClient, url string, body io.Reader) (interface{}, error) {
		return client.Get(url)
	})
}

func (resourceEndpointClient *ResourceEndpointClient) Post(url string, body interface{}) (interface{}, error) {
	return resourceEndpointClient.execute(url, body, func(client *McmaHttpClient, url string, body io.Reader) (interface{}, error) {
		return client.Post(url, body)
	})
}

func (resourceEndpointClient *ResourceEndpointClient) Put(url string, body interface{}) (interface{}, error) {
	return resourceEndpointClient.execute(url, body, func(client *McmaHttpClient, url string, body io.Reader) (interface{}, error) {
		return client.Put(url, body)
	})
}

func (resourceEndpointClient *ResourceEndpointClient) Delete(url string) (interface{}, error) {
	return resourceEndpointClient.execute(url, nil, func(client *McmaHttpClient, url string, body io.Reader) (interface{}, error) {
		return client.Delete(url)
	})
}
