package model

import (
	"encoding/json"
	"time"
)

type Service struct {
	Type            string
	Id              string
	DateCreated     time.Time
	DateModified    time.Time
	Name            string
	AuthType        string
	Resources       []ResourceEndpoint
	JobType         string
	JobProfileIds   []string
	InputLocations  []Locator
	OutputLocations []Locator
}

type serviceJson struct {
	Type            *string            `json:"@type"`
	Id              *string            `json:"id"`
	DateCreated     time.Time          `json:"dateCreated"`
	DateModified    time.Time          `json:"dateModified"`
	Name            *string            `json:"name"`
	AuthType        *string            `json:"authType"`
	Resources       []ResourceEndpoint `json:"resources"`
	JobType         *string            `json:"jobType"`
	JobProfileIds   []string           `json:"jobProfileIds"`
	InputLocations  []Locator          `json:"inputLocations"`
	OutputLocations []Locator          `json:"outputLocations"`
}

var ServiceType = "Service"

func NewService(name, authType string, resources []ResourceEndpoint) Service {
	return Service{
		Type:      ServiceType,
		Name:      name,
		AuthType:  authType,
		Resources: resources,
	}
}

func NewServiceNoAuth(name string, resources []ResourceEndpoint) Service {
	return Service{
		Type:      ServiceType,
		Name:      name,
		Resources: resources,
	}
}

func NewServiceForJobType(name, authType string, resources []ResourceEndpoint, jobType string, jobProfileIds []string) Service {
	return Service{
		Type:          ServiceType,
		Name:          name,
		AuthType:      authType,
		Resources:     resources,
		JobType:       jobType,
		JobProfileIds: jobProfileIds,
	}
}

func NewServiceForJobTypeNoAuth(name string, resources []ResourceEndpoint, jobType string, jobProfileIds []string) Service {
	return Service{
		Type:          ServiceType,
		Name:          name,
		Resources:     resources,
		JobType:       jobType,
		JobProfileIds: jobProfileIds,
	}
}

func NewServiceForJobTypeWithLocations(name, authType string, resources []ResourceEndpoint, jobType string, jobProfileIds []string, inputLocations []Locator, outputLocations []Locator) Service {
	return Service{
		Type:            ServiceType,
		Name:            name,
		AuthType:        authType,
		Resources:       resources,
		JobType:         jobType,
		JobProfileIds:   jobProfileIds,
		InputLocations:  inputLocations,
		OutputLocations: outputLocations,
	}
}

func (s Service) MarshalJSON() ([]byte, error) {
	return json.Marshal(&serviceJson{
		Type:            &ServiceType,
		Id:              stringPtrOrNull(s.Id),
		DateCreated:     s.DateCreated,
		DateModified:    s.DateModified,
		Name:            stringPtrOrNull(s.Name),
		AuthType:        stringPtrOrNull(s.AuthType),
		Resources:       s.Resources,
		JobType:         stringPtrOrNull(s.JobType),
		JobProfileIds:   s.JobProfileIds,
		InputLocations:  s.InputLocations,
		OutputLocations: s.OutputLocations,
	})
}

func (s Service) UnmarshalJSON(data []byte) error {
	var tmp serviceJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	s.Type = ServiceType
	s.Id = stringOrEmpty(tmp.Id)
	s.DateCreated = tmp.DateCreated
	s.DateModified = tmp.DateModified
	s.Name = stringOrEmpty(tmp.Name)
	s.AuthType = stringOrEmpty(tmp.AuthType)
	s.Resources = tmp.Resources
	s.JobType = stringOrEmpty(tmp.JobType)
	s.JobProfileIds = tmp.JobProfileIds
	s.InputLocations = tmp.InputLocations
	s.OutputLocations = tmp.OutputLocations

	return nil
}
