package mcmaclient

import (
	"net/http"
)

type McmaApiKeyAuthenticator struct {
	apiKey string
}

func (mcmaApiKeyAuth McmaApiKeyAuthenticator) Authenticate(req *http.Request) error {
	req.Header.Set("x-mcma-api-key", mcmaApiKeyAuth.apiKey)
	return nil
}

func NewMcmaApiKeyAuthenticator(apiKey string) McmaApiKeyAuthenticator {
	return McmaApiKeyAuthenticator{
		apiKey: apiKey,
	}
}

func (resourceManager *ResourceManager) AddMcmaApiKeyAuth(apiKey string) {
	resourceManager.AddAuth("McmaApiKey", NewMcmaApiKeyAuthenticator(apiKey))
}
