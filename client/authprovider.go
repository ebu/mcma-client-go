package mcmaclient

import (
	"fmt"
	"strings"
)

type AuthProvider struct {
	authenticators map[string]Authenticator
}

func (authProvider *AuthProvider) Add(authType string, authenticator Authenticator) {
	authProvider.authenticators[authType] = authenticator
}

func (authProvider *AuthProvider) Get(authType string) (Authenticator, error) {
	var authenticator Authenticator
	for key, a := range authProvider.authenticators {
		if strings.EqualFold(key, authType) {
			authenticator = a
			break
		}
	}
	if authenticator == nil {
		return nil, fmt.Errorf("no authenticators registered for auth type '%s'", authType)
	}

	return authenticator, nil
}

func newAuthProvider() *AuthProvider {
	return &AuthProvider{
		authenticators: make(map[string]Authenticator),
	}
}
