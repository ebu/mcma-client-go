package model

import "time"

type Service struct {
	Type            string             `json:"@type"`
	Id              string             `json:"id"`
	DateCreated     time.Time          `json:"dateCreated"`
	DateModified    time.Time          `json:"dateModified"`
	Name            string             `json:"name"`
	AuthType        string             `json:"authType"`
	AuthContext     string             `json:"authContext"`
	Resources       []ResourceEndpoint `json:"resources"`
	JobType         string             `json:"jobType"`
	JobProfileIds   []string           `json:"jobProfileIds"`
	InputLocations  []Locator          `json:"inputLocations"`
	OutputLocations []Locator          `json:"outputLocations"`
}
