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
	tracker       *model.McmaTracker
}

type nopCloser struct {
	io.ReadSeeker
}

func (nopCloser) Close() error {
	return nil
}

func getHttpErrorResponse(req *http.Request, resp *http.Response) error {
	// return an error with details from the body if possible
	var errorBody bytes.Buffer
	if resp.Body != nil {
		_, _ = errorBody.ReadFrom(resp.Body)
	}
	return fmt.Errorf("%v %v returned %v: %v", req.Method, req.URL, resp.Status, errorBody.String())
}

func newHttpRequest(method string, url string, body *bytes.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nopCloser{body})
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.ContentLength = int64(body.Len())
		snapshot := *body
		req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return nopCloser{&r}, nil
		}
	}
	return req, nil
}

func (client *McmaHttpClient) Get(url string, throwOn404 bool) (*http.Response, error) {
	return client.GetWithRetries(url, throwOn404, DefaultRetryOptions)
}
func (client *McmaHttpClient) GetWithRetries(url string, throwOn404 bool, retryOpts RetryOptions) (*http.Response, error) {
	req, err := newHttpRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.SendWithRetries(req, throwOn404, retryOpts)
}

func (client *McmaHttpClient) Post(url string, body *bytes.Reader) (*http.Response, error) {
	return client.PostWithRetries(url, body, DefaultRetryOptions)
}
func (client *McmaHttpClient) PostWithRetries(url string, body *bytes.Reader, retryOpts RetryOptions) (*http.Response, error) {
	req, err := newHttpRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.SendWithRetries(req, true, retryOpts)
}

func (client *McmaHttpClient) Put(url string, body *bytes.Reader) (*http.Response, error) {
	return client.PutWithRetries(url, body, DefaultRetryOptions)
}
func (client *McmaHttpClient) PutWithRetries(url string, body *bytes.Reader, retryOpts RetryOptions) (*http.Response, error) {
	req, err := newHttpRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.SendWithRetries(req, true, retryOpts)
}

func (client *McmaHttpClient) Delete(url string) (*http.Response, error) {
	return client.DeleteWithRetries(url, DefaultRetryOptions)
}
func (client *McmaHttpClient) DeleteWithRetries(url string, retryOpts RetryOptions) (*http.Response, error) {
	req, err := newHttpRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return client.SendWithRetries(req, true, retryOpts)
}

func (client *McmaHttpClient) Send(req *http.Request, throwOn404 bool) (*http.Response, error) {
	return client.SendWithRetries(req, throwOn404, DefaultRetryOptions)
}

func (client *McmaHttpClient) SendWithRetries(req *http.Request, throwOn404 bool, retryOpts RetryOptions) (*http.Response, error) {
	start := time.Now()

	if client.tracker != nil {
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

	done, resp, err := ExecuteWithRetries(client.httpClient, req, retryOpts)

	// connectivity/network or code error
	if err != nil {
		return resp, err
	}

	// we retried until we hit the limit
	if !done {
		lastRespErr := getHttpErrorResponse(req, resp)
		return resp, fmt.Errorf("failed to do %v to %v after %v ms - last err: %v", req.Method, req.URL, start.UnixMilli()-time.Now().UnixMilli(), lastRespErr)
	}

	// non-error response (or possible explicit exception for 404)
	if resp.StatusCode < 400 || (resp.StatusCode == 404 && !throwOn404) {
		return resp, nil
	}

	return resp, getHttpErrorResponse(req, resp)
}
