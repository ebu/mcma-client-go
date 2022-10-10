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

func TestServiceUnmarshalJSON(t *testing.T) {
	j := "{ \"@type\": \"Service\", \"id\": \"abcd-efgh-0123-4567\", \"name\": \"test\" }"
	s := &Service{}
	err := json.Unmarshal([]byte(j), s)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("Type = %s", s.Type)
	t.Logf("Id = %s", s.Id)
	t.Logf("Name = %s", s.Name)
}
