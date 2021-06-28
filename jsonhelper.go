package mcmaclient

import (
	"bytes"
	"encoding/json"
	"io"
)

func getJsonBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bodyJson), nil
}
