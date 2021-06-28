package mcmaclient

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ebu/mcma-libraries-go/model"
)

type McmaHttpClient struct {
	httpClient    *http.Client
	authenticator *Authenticator
	tracker       model.McmaTracker
}

func (client *McmaHttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Send(req)
}

func (client *McmaHttpClient) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	return client.Send(req)
}

func (client *McmaHttpClient) Put(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	return client.Send(req)
}

func (client *McmaHttpClient) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Send(req)
}

func (client *McmaHttpClient) Send(request *http.Request) (*http.Response, error) {
	if &client.tracker != nil {
		tracker := request.Header.Get("mcma-tracker")
		if tracker != "" {
			request.Header.Del("mcma-tracker")
		}
		trackerJson, err := json.Marshal(&client.tracker)
		if err != nil {
			return nil, err
		}
		trackerBase64 := base64.StdEncoding.EncodeToString(trackerJson)
		request.Header.Set("mcma-tracker", trackerBase64)
	}

	authenticator := *client.authenticator
	if authenticator != nil {
		err := authenticator.Authenticate(request)
		if err != nil {
			return nil, err
		}
	}

	return client.httpClient.Do(request)
}
