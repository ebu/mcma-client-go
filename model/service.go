package model

import "time"

type Service struct {
	Type            string             `json:"@type"`
	Id              string             `json:"id"`
	DateCreated     time.Time          `json:"dateCreated"`
	DateModified    time.Time          `json:"dateModified"`
	Name            string             `json:"name"`
	AuthType        string             `json:"authType"`
	Resources       []ResourceEndpoint `json:"resources"`
	JobType         string             `json:"jobType"`
	JobProfileIds   []string           `json:"jobProfileIds"`
	InputLocations  []Locator          `json:"inputLocations"`
	OutputLocations []Locator          `json:"outputLocations"`
}

func NewService(name, authType string, resources []ResourceEndpoint) Service {
	return Service{
		Type:      "Service",
		Name:      name,
		AuthType:  authType,
		Resources: resources,
	}
}

func NewServiceNoAuth(name string, resources []ResourceEndpoint) Service {
	return Service{
		Type:      "Service",
		Name:      name,
		Resources: resources,
	}
}

func NewServiceForJobType(name, authType string, resources []ResourceEndpoint, jobType string, jobProfileIds []string) Service {
	return Service{
		Type:          "Service",
		Name:          name,
		AuthType:      authType,
		Resources:     resources,
		JobType:       jobType,
		JobProfileIds: jobProfileIds,
	}
}

func NewServiceForJobTypeNoAuth(name string, resources []ResourceEndpoint, jobType string, jobProfileIds []string) Service {
	return Service{
		Type:          "Service",
		Name:          name,
		Resources:     resources,
		JobType:       jobType,
		JobProfileIds: jobProfileIds,
	}
}

func NewServiceForJobTypeWithLocations(name, authType string, resources []ResourceEndpoint, jobType string, jobProfileIds []string, inputLocations []Locator, outputLocations []Locator) Service {
	return Service{
		Type:            "Service",
		Name:            name,
		AuthType:        authType,
		Resources:       resources,
		JobType:         jobType,
		JobProfileIds:   jobProfileIds,
		InputLocations:  inputLocations,
		OutputLocations: outputLocations,
	}
}
