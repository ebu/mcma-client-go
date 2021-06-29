package mcmaclient

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"io"
	"net/http"
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
	var seekableBody io.ReadSeeker
	if req.Body != nil {
		var canSeek bool
		seekableBody, canSeek = req.Body.(io.ReadSeeker)
		if !canSeek {
			return fmt.Errorf("must support Seek() for AWS signing")
		}
	}
	_, err := aws4Auth.signer.Sign(req, seekableBody, aws4Auth.region, aws4Auth.service, time.Now())
	return err
}

func NewAWS4Authenticator(authContext AWS4AuthContext) AWS4Authenticator {
	return AWS4Authenticator{
		signer:  v4.NewSigner(credentials.NewStaticCredentials(authContext.AccessKey, authContext.SecretKey, authContext.SessionToken)),
		region:  authContext.Region,
		service: "execute-api",
	}
}

func (resourceManager *ResourceManager) AddAWS4Auth() {
	resourceManager.AddAuth("AWS4", func(authContext interface{}) (Authenticator, error) {
		var ac interface{}
		switch authContext.(type) {
		case string:
			if err := json.Unmarshal([]byte(authContext.(string)), ac); err != nil {
				return nil, err
			}
		default:
			ac = authContext
		}
		aws4AuthContext, isAws4AuthContext := ac.(AWS4AuthContext)
		if !isAws4AuthContext {
			return nil, fmt.Errorf("invalid AWS4 auth context")
		}
		return NewAWS4Authenticator(aws4AuthContext), nil
	})
}
