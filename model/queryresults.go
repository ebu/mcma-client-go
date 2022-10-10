package model

import (
	"encoding/json"
	"reflect"
)

type QueryResults struct {
	Results            []interface{}
	NextPageStartToken string
}

type queryResultsJson struct {
	Results            []interface{} `json:"results"`
	NextPageStartToken *string       `json:"nextPageStartToken"`
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

func (qr QueryResults) MarshalJSON() ([]byte, error) {
	return json.Marshal(&queryResultsJson{
		Results:            qr.Results,
		NextPageStartToken: stringPtrOrNull(qr.NextPageStartToken),
	})
}

func (qr *QueryResults) UnmarshalJSON(data []byte) error {
	var tmp queryResultsJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	qr.Results = tmp.Results
	qr.NextPageStartToken = stringOrEmpty(tmp.NextPageStartToken)

	return nil
}
