package model

import "testing"

func TestServiceMarshalJSON(t *testing.T) {
	s := NewService("test", "", make([]ResourceEndpoint, 0))
	data, err := s.MarshalJSON()
	if err != nil {
		t.Errorf("%v", err)
	}
	json := string(data)
	t.Log(json)
}
