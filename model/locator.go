package model

import "encoding/json"

var LocatorType = "Locator"

type Locator struct {
	Type string
	Url  string
}

type locatorJson struct {
	Type *string `json:"@type"`
	Url  *string `json:"url"`
}

func NewLocator(url string) Locator {
	return Locator{
		Type: LocatorType,
		Url:  url,
	}
}

func (l Locator) MarshalJSON() ([]byte, error) {
	return json.Marshal(&locatorJson{
		Type: &LocatorType,
		Url:  stringPtrOrNull(l.Url),
	})
}

func (l *Locator) UnmarshalJSON(data []byte) error {
	var tmp locatorJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	l.Type = LocatorType
	l.Url = stringOrEmpty(tmp.Url)

	return nil
}
