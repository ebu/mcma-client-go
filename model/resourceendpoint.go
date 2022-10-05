package model

import "encoding/json"

type ResourceEndpoint struct {
	Type         string
	ResourceType string
	HttpEndpoint string
	AuthType     string
}

type resourceEndpointJson struct {
	Type         *string `json:"@type"`
	ResourceType *string `json:"resourceType"`
	HttpEndpoint *string `json:"httpEndpoint"`
	AuthType     *string `json:"authType"`
}

var ResourceEndpointType = "ResourceEndpoint"

func NewResourceEndpoint(resourceType, httpEndpoint string) ResourceEndpoint {
	return ResourceEndpoint{
		Type:         ResourceEndpointType,
		ResourceType: resourceType,
		HttpEndpoint: httpEndpoint,
	}
}

func NewResourceEndpointWithAuth(resourceType, httpEndpoint, authType string) ResourceEndpoint {
	return ResourceEndpoint{
		Type:         ResourceEndpointType,
		ResourceType: resourceType,
		HttpEndpoint: httpEndpoint,
		AuthType:     authType,
	}
}

func (n ResourceEndpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(&resourceEndpointJson{
		Type:         &ResourceEndpointType,
		ResourceType: stringPtrOrNull(n.ResourceType),
		HttpEndpoint: stringPtrOrNull(n.HttpEndpoint),
		AuthType:     stringPtrOrNull(n.AuthType),
	})
}

func (n ResourceEndpoint) UnmarshalJSON(data []byte) error {
	var tmp resourceEndpointJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	n.Type = ResourceEndpointType
	n.ResourceType = stringOrEmpty(tmp.ResourceType)
	n.HttpEndpoint = stringOrEmpty(tmp.HttpEndpoint)
	n.AuthType = stringOrEmpty(tmp.AuthType)

	return nil
}
