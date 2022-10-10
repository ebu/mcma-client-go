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

func (re ResourceEndpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(&resourceEndpointJson{
		Type:         &ResourceEndpointType,
		ResourceType: stringPtrOrNull(re.ResourceType),
		HttpEndpoint: stringPtrOrNull(re.HttpEndpoint),
		AuthType:     stringPtrOrNull(re.AuthType),
	})
}

func (re *ResourceEndpoint) UnmarshalJSON(data []byte) error {
	var tmp resourceEndpointJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	re.Type = ResourceEndpointType
	re.ResourceType = stringOrEmpty(tmp.ResourceType)
	re.HttpEndpoint = stringOrEmpty(tmp.HttpEndpoint)
	re.AuthType = stringOrEmpty(tmp.AuthType)

	return nil
}
