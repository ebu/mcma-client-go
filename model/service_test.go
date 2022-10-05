package model

import (
	"encoding/json"
	"testing"
)

func TestServiceMarshalJSON(t *testing.T) {
	s := NewService("test", "", make([]ResourceEndpoint, 0))
	data, err := json.Marshal(s)
	if err != nil {
		t.Errorf("%v", err)
	}
	j := string(data)
	t.Log(j)
}
