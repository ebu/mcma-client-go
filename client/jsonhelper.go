package mcmaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	var r reflect.Value
	if t.Kind() != reflect.Map {
		r = reflect.New(t)
	} else {
		var m map[string]interface{}
		r = reflect.MakeMap(reflect.TypeOf(m))
	}
	if resp.ContentLength == 0 {
		return r.Elem().Interface(), nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp body: %v", err)
	}
	err = json.Unmarshal(body, r.Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json resp body: %v", err)
	}
	return r.Elem().Interface(), nil
}
