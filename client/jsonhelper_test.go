package mcmaclient

import (
	"bytes"
	"encoding/json"
	"github.com/ebu/mcma-libraries-go/model"
	"net/http"
	"reflect"
	strconv "strconv"
	"testing"
)

func TestGetJsonReqBody(t *testing.T) {
	s := model.NewServiceForJobTypeWithLocations(
		"test",
		"aws4",
		[]model.ResourceEndpoint{
			model.NewResourceEndpoint("JobAssignment", "https://service/api/job-assignments"),
		},
		"AmeJob",
		[]string{
			"https://service-registry/api/profiles/1",
		},
		[]model.Locator{
			{
				Type: "Locator",
				Url:  "https://s3/bucket",
			},
		},
		[]model.Locator{
			{
				Type: "Locator",
				Url:  "https://s3/bucket",
			},
		})
	r, err := getJsonReqBody(s)
	if err != nil {
		t.Errorf("%v", err)
	}
	b := make([]byte, r.Size())
	_, err = r.Read(b)
	if err != nil {
		t.Errorf("%v", err)
	}
	j := string(b)
	println(j)
}

func TestReadyJsonRespBody(t *testing.T) {
	s := model.NewServiceForJobTypeWithLocations(
		"test",
		"aws4",
		[]model.ResourceEndpoint{
			model.NewResourceEndpoint("JobAssignment", "https://service/api/job-assignments"),
		},
		"AmeJob",
		[]string{
			"https://service-registry/api/profiles/1",
		},
		[]model.Locator{
			{
				Type: "Locator",
				Url:  "https://s3/bucket",
			},
		},
		[]model.Locator{
			{
				Type: "Locator",
				Url:  "https://s3/bucket",
			},
		},
	)
	j, err := json.Marshal(s)
	if err != nil {
		t.Errorf("%v", err)
	}
	resp := &http.Response{
		StatusCode:    200,
		ContentLength: int64(len(j)),
		Body: nopCloser{
			bytes.NewReader(j),
		},
	}
	res, err := readJsonRespBody(resp, reflect.TypeOf(model.Service{}))
	if err != nil {
		t.Errorf("%v", err)
	}
	svc := res.(model.Service)
	println("Type: " + svc.Type)
	println("Name: " + svc.Name)
	println("AuthType: " + svc.AuthType)
	println("JobType: " + svc.JobType)
	println("JobProfileIds: " + strconv.Itoa(len(svc.JobProfileIds)))
}
