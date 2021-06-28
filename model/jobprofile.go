package model

import "time"

type JobProfile struct {
	Type                    string                 `json:"@type"`
	Id                      string                 `json:"id"`
	DateCreated             time.Time              `json:"dateCreated"`
	DateModified            time.Time              `json:"dateModified"`
	Name                    string                 `json:"name"`
	InputParameters         []JobParameter         `json:"inputParameters"`
	OutputParameters        []JobParameter         `json:"outputParameters"`
	OptionalInputParameters []JobParameter         `json:"optionalInputParameters"`
	CustomProperties        map[string]interface{} `json:"customProperties"`
}

func NewJobProfile(name string) JobProfile {
	return JobProfile{
		Type: "JobProfile",
		Name: name,
	}
}
