package model

import "encoding/json"

type Notification struct {
	Type    string
	Source  string
	Content interface{}
	Custom  map[string]interface{}
}

type notificationJson struct {
	Type    *string                `json:"@type"`
	Source  *string                `json:"source"`
	Content interface{}            `json:"content"`
	Custom  map[string]interface{} `json:"custom"`
}

var NotificationType = "Notification"

func NewNotification(source string, content interface{}) Notification {
	return Notification{
		Type:    NotificationType,
		Source:  source,
		Content: content,
	}
}

func (n Notification) MarshalJSON() ([]byte, error) {
	return json.Marshal(&notificationJson{
		Type:    &NotificationType,
		Source:  stringPtrOrNull(n.Source),
		Content: n.Content,
	})
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	var tmp notificationJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	n.Type = NotificationType
	n.Source = stringOrEmpty(tmp.Source)
	n.Content = tmp.Content

	return nil
}
