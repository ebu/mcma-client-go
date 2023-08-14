package mcmaclient

import (
	"io"
	"net/http"
	"testing"

	"github.com/ebu/mcma-libraries-go/model"
)

func TestSeekableReqBody(t *testing.T) {
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
	req1, err := http.NewRequest("POST", "test", r)
	if err != nil {
		t.Errorf("%v", err)
	}
	req2, err := http.NewRequest("POST", "test", nopCloser{r})
	if err != nil {
		t.Errorf("%v", err)
	}
	req3, err := newHttpRequest("POST", "test", r)
	if err != nil {
		t.Errorf("%v", err)
	}
	println(req1.ContentLength)
	println(req2.ContentLength)
	println(req3.ContentLength)

	_, isRs1 := req1.Body.(io.ReadSeeker)
	println(isRs1)
	_, isRs2 := req2.Body.(io.ReadSeeker)
	println(isRs2)
	_, isRs3 := req3.Body.(io.ReadSeeker)
	println(isRs3)
}
