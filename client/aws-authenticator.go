package mcmaclient

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"io"
	"net/http"
	"os"
	"time"
)

type AWS4AuthContext struct {
	AccessKey    string
	SecretKey    string
	Region       string
	SessionToken string
}

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
		body.Seek(0, io.SeekStart)
	}
	header, err := aws4Auth.signer.Sign(req, body, aws4Auth.service, aws4Auth.region, time.Now())
	for key, values := range header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}
	if body != nil {
		body.Seek(0, io.SeekStart)
	}
	return err
}

func NewAWS4Authenticator(authContext AWS4AuthContext) AWS4Authenticator {
	var creds *credentials.Credentials
	if len(authContext.AccessKey) > 0 {
		creds = credentials.NewStaticCredentials(authContext.AccessKey, authContext.SecretKey, authContext.SessionToken)
	} else {
		creds = credentials.NewEnvCredentials()
	}
	if len(authContext.Region) == 0 {
		authContext.Region = os.Getenv("AWS_REGION")
		if len(authContext.Region) == 0 {
			authContext.Region = os.Getenv("AWS_DEFAULT_REGION")
		}
	}
	return AWS4Authenticator{
		signer:  v4.NewSigner(creds),
		region:  authContext.Region,
		service: "execute-api",
	}
}

func (resourceManager *ResourceManager) AddAWS4Auth() {
	resourceManager.AddAuth("AWS4", func(authContext interface{}) (Authenticator, error) {
		var aws4AuthContext AWS4AuthContext
		switch authContext.(type) {
		case string:
			s := authContext.(string)
			if len(s) > 0 {
				if err := json.Unmarshal([]byte(s), &aws4AuthContext); err != nil {
					return nil, fmt.Errorf("Failed to unmarshal json to AWS4 auth context: %v", err)
				}
			}
		case AWS4AuthContext:
			aws4AuthContext = authContext.(AWS4AuthContext)
		default:
			return nil, fmt.Errorf("invalid AWS4 auth context")
		}
		return NewAWS4Authenticator(aws4AuthContext), nil
	})
}
