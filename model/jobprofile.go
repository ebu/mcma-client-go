package model

import (
	"encoding/json"
	"time"
)

type JobProfile struct {
	Type                    string
	Id                      string
	DateCreated             time.Time
	DateModified            time.Time
	Name                    string
	InputParameters         []JobParameter
	OutputParameters        []JobParameter
	OptionalInputParameters []JobParameter
	CustomProperties        map[string]interface{}
}

type jobProfileJson struct {
	Type                    string                 `json:"@type"`
	Id                      *string                `json:"id"`
	DateCreated             time.Time              `json:"dateCreated"`
	DateModified            time.Time              `json:"dateModified"`
	Name                    *string                `json:"name"`
	InputParameters         []JobParameter         `json:"inputParameters"`
	OutputParameters        []JobParameter         `json:"outputParameters"`
	OptionalInputParameters []JobParameter         `json:"optionalInputParameters"`
	CustomProperties        map[string]interface{} `json:"customProperties"`
}

var JobProfileType = "JobProfile"

func NewJobProfile(name string) JobProfile {
	return JobProfile{
		Type: JobProfileType,
		Name: name,
	}
}

func (jp *JobProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jobProfileJson{
		Type:                    JobProfileType,
		Id:                      stringPtrOrNull(jp.Id),
		DateCreated:             jp.DateCreated,
		DateModified:            jp.DateModified,
		Name:                    stringPtrOrNull(jp.Name),
		InputParameters:         jp.InputParameters,
		OutputParameters:        jp.OutputParameters,
		OptionalInputParameters: jp.OptionalInputParameters,
		CustomProperties:        jp.CustomProperties,
	})
}

func (jp *JobProfile) UnmarshalJSON(data []byte) error {
	var tmp jobProfileJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	jp.Type = JobProfileType
	jp.Id = stringOrEmpty(tmp.Id)
	jp.DateCreated = tmp.DateCreated
	jp.DateModified = tmp.DateModified
	jp.Name = stringOrEmpty(tmp.Name)
	jp.InputParameters = tmp.InputParameters
	jp.OutputParameters = tmp.OutputParameters
	jp.OptionalInputParameters = tmp.OptionalInputParameters
	jp.CustomProperties = tmp.CustomProperties

	return nil
}
