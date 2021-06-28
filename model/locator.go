package model

type Locator struct {
	Type string `json:"@type"`
	Url  string `json:"url"`
}
