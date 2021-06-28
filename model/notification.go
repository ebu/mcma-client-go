package model

type Notification struct {
	Type    string      `json:"@type"`
	Source  string      `json:"source"`
	Content interface{} `json:"content"`
}
