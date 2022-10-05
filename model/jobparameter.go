package model

import "encoding/json"

type JobParameter struct {
	ParameterName string
	ParameterType string
}

type jobParameterJson struct {
	ParameterName *string `json:"parameterName"`
	ParameterType *string `json:"parameterType"`
}

func NewJobParameter(parameterName, parameterType string) JobParameter {
	return JobParameter{
		ParameterName: parameterName,
		ParameterType: parameterType,
	}
}

func (jp JobParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jobParameterJson{
		ParameterName: stringPtrOrNull(jp.ParameterName),
		ParameterType: stringPtrOrNull(jp.ParameterType),
	})
}

func (jp JobParameter) UnmarshalJSON(data []byte) error {
	var tmp jobParameterJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	jp.ParameterName = stringOrEmpty(tmp.ParameterName)
	jp.ParameterType = stringOrEmpty(tmp.ParameterType)

	return nil
}
