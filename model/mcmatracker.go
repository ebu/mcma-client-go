package model

type McmaTracker struct {
	Type   string            `json:"@type"`
	Id     string            `json:"id"`
	Label  string            `json:"label"`
	Custom map[string]string `json:"custom"`
}
