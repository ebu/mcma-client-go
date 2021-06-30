package model

import (
	"encoding/json"
	"reflect"
)

type QueryResults struct {
	Results            []interface{} `json:"results"`
	NextPageStartToken string        `json:"nextPageStartToken"`
}

func (qr QueryResults) GetResults(t reflect.Type) ([]interface{}, error) {
	var results []interface{}
	for _, r := range qr.Results {
		resultVal := reflect.New(t)
		resultPtr := resultVal.Interface()
		rJson, err := json.Marshal(r)
		if err != nil {
			return results, err
		}
		err = json.Unmarshal(rJson, resultPtr)
		if err != nil {
			return results, err
		}
		result := resultVal.Elem().Interface()
		results = append(results, result)
	}
	return results, nil
}
