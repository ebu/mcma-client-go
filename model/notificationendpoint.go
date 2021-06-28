package model

type NotificationEndpoint struct {
	Type         string `json:"@type"`
	Id           string `json:"id"`
	HttpEndpoint string `json:"httpEndpoint"`
}
