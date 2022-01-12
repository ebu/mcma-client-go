package mcmaclient

import (
	"bytes"
	"fmt"
	"github.com/ebu/mcma-libraries-go/model"
	"net/http"
	"reflect"
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
	if authContext == nil {
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

func (resourceEndpointClient ResourceEndpointClient) getFullUrl(url string) string {
	if url == "" {
		return resourceEndpointClient.resourceEndpoint.HttpEndpoint
	}
	if strings.HasPrefix(strings.ToLower(url), strings.ToLower(resourceEndpointClient.resourceEndpoint.HttpEndpoint)) {
		return url
	}
	return strings.TrimSuffix(resourceEndpointClient.resourceEndpoint.HttpEndpoint, "/") + "/" + url
}

func (resourceEndpointClient ResourceEndpointClient) hasMatchingHttpEndpoint(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), strings.ToLower(resourceEndpointClient.resourceEndpoint.HttpEndpoint))
}

func (resourceEndpointClient ResourceEndpointClient) getHttpEndpoint() string {
	return resourceEndpointClient.resourceEndpoint.HttpEndpoint
}

func (resourceEndpointClient *ResourceEndpointClient) execute(t reflect.Type, url string, body interface{}, execute func(mcmaHttpClient *McmaHttpClient, url string, body *bytes.Reader) (*http.Response, error)) (interface{}, error) {
	mcmaHttpClient, err := resourceEndpointClient.getMcmaHttpClient()
	if err != nil {
		return nil, err
	}

	url = resourceEndpointClient.getFullUrl(url)
	reqBody, err := getJsonReqBody(body)
	if err != nil {
		return nil, err
	}
	resp, err := execute(mcmaHttpClient, url, reqBody)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	return readJsonRespBody(resp, t)
}

func (resourceEndpointClient *ResourceEndpointClient) Query(t reflect.Type, url string, queryParameters []struct {
	key   string
	value string
}) (model.QueryResults, error) {
	var queryResults model.QueryResults
	mcmaHttpClient, err := resourceEndpointClient.getMcmaHttpClient()
	if err != nil {
		return queryResults, err
	}

	url = resourceEndpointClient.getFullUrl(url)
	if len(queryParameters) > 0 {
		url += "?"
		for _, p := range queryParameters {
			url += fmt.Sprintf("%s=%s", p.key, p.value)
		}
	}

	getResp, err := mcmaHttpClient.Get(url, true)
	if err != nil {
		return queryResults, fmt.Errorf("failed to query %v: %v", url, err)
	}

	body, err := readJsonRespBody(getResp, reflect.TypeOf(queryResults))
	if err != nil {
		return queryResults, fmt.Errorf("failed to get query results for %v: %v", url, err)
	}

	queryResults = body.(model.QueryResults)

	results, err := queryResults.GetResults(t)
	if err != nil {
		return queryResults, fmt.Errorf("failed to get typed query results for %v: %v", url, err)
	}

	queryResults.Results = results

	return queryResults, err
}

func (resourceEndpointClient *ResourceEndpointClient) QueryMaps(url string, queryParameters []struct {
	key   string
	value string
}) (model.QueryResults, error) {
	var queryResults model.QueryResults
	mcmaHttpClient, err := resourceEndpointClient.getMcmaHttpClient()
	if err != nil {
		return queryResults, err
	}

	url = resourceEndpointClient.getFullUrl(url)
	if len(queryParameters) > 0 {
		url += "?"
		for _, p := range queryParameters {
			url += fmt.Sprintf("%s=%s", p.key, p.value)
		}
	}

	getResp, err := mcmaHttpClient.Get(url, true)
	if err != nil {
		return queryResults, fmt.Errorf("failed to query %v: %v", url, err)
	}

	body, err := readJsonRespBody(getResp, reflect.TypeOf(queryResults))
	if err != nil {
		return queryResults, fmt.Errorf("failed to get query results for %v: %v", url, err)
	}

	return body.(model.QueryResults), err
}

func (resourceEndpointClient *ResourceEndpointClient) Get(t reflect.Type, url string) (interface{}, error) {
	return resourceEndpointClient.execute(t, url, nil, func(client *McmaHttpClient, url string, body *bytes.Reader) (*http.Response, error) {
		return client.Get(url, false)
	})
}

func (resourceEndpointClient *ResourceEndpointClient) GetResource(url string) (map[string]interface{}, error) {
	var m map[string]interface{}
	mi, err := resourceEndpointClient.Get(reflect.TypeOf(m), url)
	if err != nil {
		return nil, err
	}
	return mi.(map[string]interface{}), nil
}

func (resourceEndpointClient *ResourceEndpointClient) Post(t reflect.Type, url string, body interface{}) (interface{}, error) {
	return resourceEndpointClient.execute(t, url, body, func(client *McmaHttpClient, url string, body *bytes.Reader) (*http.Response, error) {
		return client.Post(url, body)
	})
}

func (resourceEndpointClient *ResourceEndpointClient) PostResource(url string, body map[string]interface{}) (interface{}, error) {
	return resourceEndpointClient.Post(reflect.TypeOf(body), url, body)
}

func (resourceEndpointClient *ResourceEndpointClient) Put(t reflect.Type, url string, body interface{}) (interface{}, error) {
	return resourceEndpointClient.execute(t, url, body, func(client *McmaHttpClient, url string, body *bytes.Reader) (*http.Response, error) {
		return client.Put(url, body)
	})
}

func (resourceEndpointClient *ResourceEndpointClient) PutResource(url string, body map[string]interface{}) (interface{}, error) {
	return resourceEndpointClient.Put(reflect.TypeOf(body), url, body)
}

func (resourceEndpointClient *ResourceEndpointClient) Delete(url string) error {
	_, err := resourceEndpointClient.execute(nil, url, nil, func(client *McmaHttpClient, url string, body *bytes.Reader) (*http.Response, error) {
		return client.Delete(url)
	})
	return err
}
