package model

import "encoding/json"

type McmaTracker struct {
	Type   string
	Id     string
	Label  string
	Custom map[string]string
}

type mcmaTrackerJson struct {
	Type   *string           `json:"@type"`
	Id     *string           `json:"id"`
	Label  *string           `json:"label"`
	Custom map[string]string `json:"custom"`
}

var McmaTrackerType = "McmaTracker"

func NewTracker(id, label string, custom map[string]string) McmaTracker {
	return McmaTracker{
		Type:   McmaTrackerType,
		Id:     id,
		Label:  label,
		Custom: custom,
	}
}

func (t McmaTracker) MarshalJSON() ([]byte, error) {
	return json.Marshal(&mcmaTrackerJson{
		Type:   &McmaTrackerType,
		Id:     stringPtrOrNull(t.Id),
		Label:  stringPtrOrNull(t.Label),
		Custom: t.Custom,
	})
}

func (t *McmaTracker) UnmarshalJSON(data []byte) error {
	var tmp mcmaTrackerJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	t.Type = McmaTrackerType
	t.Id = stringOrEmpty(tmp.Id)
	t.Label = stringOrEmpty(tmp.Label)
	t.Custom = tmp.Custom

	return nil
}
