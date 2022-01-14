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
		return nil, fmt.Errorf("failed to get json body: %v", err)
	}
	return bytes.NewReader(bodyJson), nil
}

func readJsonRespBody(resp *http.Response, t reflect.Type) (interface{}, error) {
	if resp.StatusCode == 404 {
		return nil, nil
	}
	var getVal func() interface{}
	var getPtr func() interface{}
	if t.Kind() != reflect.Map {
		tmp := reflect.New(t)
		getVal = func() interface{} { return tmp.Elem().Interface() }
		getPtr = func() interface{} { return tmp.Interface() }
	} else {
		m := make(map[string]interface{})
		mPtr := &m
		getVal = func() interface{} { return *mPtr }
		getPtr = func() interface{} { return mPtr }
	}
	if resp.ContentLength == 0 {
		return getVal(), nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp body: %v", err)
	}
	err = json.Unmarshal(body, getPtr())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json resp body: %v", err)
	}

	return getVal(), nil
}
