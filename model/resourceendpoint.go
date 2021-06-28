package model

import "time"

type ResourceEndpoint struct {
	Type         string      `json:"@type"`
	Id           string      `json:"id"`
	DateCreated  time.Time   `json:"dateCreated"`
	DateModified time.Time   `json:"dateModified"`
	ResourceType string      `json:"resourceType"`
	HttpEndpoint string      `json:"httpEndpoint"`
	AuthType     string      `json:"authType"`
	AuthContext  interface{} `json:"authContext"`
}
