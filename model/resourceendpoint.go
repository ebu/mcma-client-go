package model

type ResourceEndpoint struct {
	Type         string       `json:"@type"`
	ResourceType string       `json:"resourceType"`
	HttpEndpoint string       `json:"httpEndpoint"`
	AuthType     *string      `json:"authType"`
	AuthContext  *interface{} `json:"authContext"`
}
