package mcmaclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ebu/mcma-libraries-go/model"
)

type McmaHttpClient struct {
	httpClient    *http.Client
	authenticator *Authenticator
	tracker       model.McmaTracker
}

type nopCloser struct {
	io.ReadSeeker
}

func (nopCloser) Close() error {
	return nil
}

func (client *McmaHttpClient) Get(url string, throwOn404 bool) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Send(req, throwOn404)
}

func (client *McmaHttpClient) Post(url string, body *bytes.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nopCloser{body})
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Send(req, true)
}

func (client *McmaHttpClient) Put(url string, body *bytes.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, nopCloser{body})
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Send(req, true)
}

func (client *McmaHttpClient) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Send(req, true)
}

func (client *McmaHttpClient) Send(req *http.Request, throwOn404 bool) (*http.Response, error) {
	backOff := 100 * time.Millisecond
	var err error
	for {
		if &client.tracker != nil {
			tracker := req.Header.Get("mcma-tracker")
			if tracker != "" {
				req.Header.Del("mcma-tracker")
			}
			trackerJson, err := json.Marshal(&client.tracker)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal MCMA tracker to json: %v", err)
			}
			trackerBase64 := base64.StdEncoding.EncodeToString(trackerJson)
			req.Header.Set("mcma-tracker", trackerBase64)
		}

		if client.authenticator != nil {
			authenticator := *client.authenticator
			err := authenticator.Authenticate(req)
			if err != nil {
				return nil, err
			}
		}

		var resp *http.Response
		resp, err = client.httpClient.Do(req)

		if err == nil && resp.StatusCode < 500 {
			if resp.StatusCode < 200 || (resp.StatusCode >= 300 && (resp.StatusCode != 404 || throwOn404)) {
				var errorBody bytes.Buffer
				if resp.Body != nil {
					errorBody.ReadFrom(resp.Body)
				}
				return nil, fmt.Errorf("received resp %v doing %v to %v: %v", resp.Status, req.Method, req.URL, errorBody.String())
			} else {
				return resp, nil
			}
		}

		if err == nil {
			err = fmt.Errorf("received resp %v", resp.Status)
		}

		if backOff*2 >= time.Minute {
			break
		}

		backOff *= 2
		time.Sleep(backOff)
	}

	return nil, fmt.Errorf("failed to do %v to %v after %v ms: %v", req.Method, req.URL, backOff, err)
}
