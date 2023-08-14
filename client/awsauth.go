package mcmaclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

type AWS4Authenticator struct {
	signer  *v4.Signer
	region  string
	service string
}

func (aws4Auth AWS4Authenticator) Authenticate(req *http.Request) error {
	var body io.ReadSeeker
	if req.Body != nil {
		var canSeek bool
		if body, canSeek = req.Body.(io.ReadSeeker); !canSeek {
			return fmt.Errorf("body must be seekable for AWS4 auth")
		}
	}
	if body != nil {
		_, err := body.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek to start of request body for AWS auth: %v", err)
		}
	}
	header, err := aws4Auth.signer.Sign(req, body, aws4Auth.service, aws4Auth.region, time.Now())
	for key, values := range header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}
	if body != nil {
		_, err := body.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek to start of request body for AWS auth: %v", err)
		}
	}
	return err
}

func newAWS4Authenticator(creds *credentials.Credentials, region string) AWS4Authenticator {
	if len(region) == 0 {
		region = os.Getenv("AWS_REGION")
		if len(region) == 0 {
			region = os.Getenv("AWS_DEFAULT_REGION")
		}
	}
	return AWS4Authenticator{
		signer:  v4.NewSigner(creds),
		region:  region,
		service: "execute-api",
	}
}

func NewAWS4AuthenticatorFromKeys(accessKey, secretKey, sessionToken, region string) AWS4Authenticator {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, sessionToken)
	return newAWS4Authenticator(creds, region)
}

func NewAWS4AuthenticatorFromProfile(profile, region string) AWS4Authenticator {
	creds := credentials.NewSharedCredentials("", profile)
	return newAWS4Authenticator(creds, region)
}

func NewAWS4AuthenticatorFromEnvVars() AWS4Authenticator {
	creds := credentials.NewEnvCredentials()
	return newAWS4Authenticator(creds, "")
}

func (resourceManager *ResourceManager) AddAWS4AuthFromKeys(accessKey, secretKey, sessionToken, region string) {
	resourceManager.AddAWS4Auth(NewAWS4AuthenticatorFromKeys(accessKey, secretKey, sessionToken, region))
}

func (resourceManager *ResourceManager) AddAWS4AuthFromProfile(profile, region string) {
	resourceManager.AddAWS4Auth(NewAWS4AuthenticatorFromProfile(profile, region))
}

func (resourceManager *ResourceManager) AddAWS4AuthFromEnvVars() {
	resourceManager.AddAWS4Auth(NewAWS4AuthenticatorFromEnvVars())
}

func (resourceManager *ResourceManager) AddAWS4Auth(authenticator AWS4Authenticator) {
	resourceManager.AddAuth("AWS4", authenticator)
}
