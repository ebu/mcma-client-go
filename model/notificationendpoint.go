package model

import "encoding/json"

type NotificationEndpoint struct {
	Type         string
	Id           string
	HttpEndpoint string
}

type notificationEndpointJson struct {
	Type         *string `json:"@type"`
	Id           *string `json:"id"`
	HttpEndpoint *string `json:"httpEndpoint"`
}

var NotificationEndpointType = "NotificationEndpoint"

func NewNotificationEndpoint(id, httpEndpoint string) NotificationEndpoint {
	return NotificationEndpoint{
		Type:         NotificationEndpointType,
		Id:           id,
		HttpEndpoint: httpEndpoint,
	}
}

func (ne *NotificationEndpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(&notificationEndpointJson{
		Type:         &NotificationEndpointType,
		Id:           stringPtrOrNull(ne.Id),
		HttpEndpoint: stringPtrOrNull(ne.HttpEndpoint),
	})
}

func (ne *NotificationEndpoint) UnmarshalJSON(data []byte) error {
	var tmp notificationEndpointJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	ne.Type = NotificationEndpointType
	ne.Id = stringOrEmpty(tmp.Id)
	ne.HttpEndpoint = stringOrEmpty(tmp.HttpEndpoint)

	return nil
}
