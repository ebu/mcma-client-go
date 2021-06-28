package model

type QueryResults struct {
	Results            []interface{} `json:"results"`
	NextPageStartToken string        `json:"nextPageStartToken"`
}
