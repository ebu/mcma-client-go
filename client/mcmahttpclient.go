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
	start := time.Now()
	backOffTimesInMs := []time.Duration{
		250 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
		45 * time.Second,
		1 * time.Minute,
	}

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
		if authenticator != nil {
			err := authenticator.Authenticate(req)
			if err != nil {
				return nil, err
			}
		}
	}

	trySendReq := func() (bool, *http.Response, error) {
		resp, err := client.httpClient.Do(req)

		if err != nil {
			return true, resp, err
		}

		if resp.StatusCode < 400 {
			return true, resp, nil
		}

		if resp.StatusCode < 500 && resp.StatusCode != 404 && resp.StatusCode != 429 {
			var errorBody bytes.Buffer
			if resp.Body != nil {
				_, _ = errorBody.ReadFrom(resp.Body)
			}
			err = fmt.Errorf("%v %v returned %v: %v", req.Method, req.URL, resp.Status, errorBody.String())
			return true, resp, err
		}

		err = fmt.Errorf("received resp %v", resp.Status)
		return false, resp, err
	}

	done, resp, err := trySendReq()
	if done {
		return resp, err
	}

	for i := 0; i < len(backOffTimesInMs); i++ {
		time.Sleep(backOffTimesInMs[i])

		done, resp, err = trySendReq()
		if done {
			return resp, err
		}
	}

	if resp.StatusCode == 404 && !throwOn404 {
		return resp, nil
	}

	return resp, fmt.Errorf("failed to do %v to %v after %v ms: %v", req.Method, req.URL, start.UnixMilli()-time.Now().UnixMilli(), err)
}
