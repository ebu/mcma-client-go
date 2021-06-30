package mcmaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func getJsonReqBody(body interface{}) (*bytes.Reader, error) {
	if body == nil {
		return nil, nil
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Failed to get json body: %v", err)
	}
	return bytes.NewReader(bodyJson), nil
}

func readJsonRespBody(resp *http.Response, t reflect.Type) (interface{}, error) {
	if resp.StatusCode == 404 {
		return nil, nil
	}
	r := reflect.New(t)
	if resp.ContentLength == 0 {
		return r.Elem().Interface(), nil
	}
	var body bytes.Buffer
	var totalBytesRead int64
	for totalBytesRead < resp.ContentLength {
		bytesRead, err := body.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}
		totalBytesRead += bytesRead
	}
	err := json.Unmarshal(body.Bytes(), r.Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json resp body: %v", err)
	}
	return r.Elem().Interface(), nil
}
