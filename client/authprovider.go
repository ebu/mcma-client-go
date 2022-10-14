package mcmaclient

import (
	"fmt"
	"strings"
)

type AuthProvider struct {
	authenticators       map[string]Authenticator
	defaultAuthenticator *Authenticator
}

func (authProvider *AuthProvider) Add(authType string, authenticator Authenticator) {
	authProvider.authenticators[authType] = authenticator
	if len(authProvider.authenticators) == 1 {
		authProvider.defaultAuthenticator = &authenticator
	} else {
		authProvider.defaultAuthenticator = nil
	}
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

func (authProvider *AuthProvider) GetDefault() *Authenticator {
	return authProvider.defaultAuthenticator
}

func newAuthProvider() *AuthProvider {
	return &AuthProvider{
		authenticators: make(map[string]Authenticator),
	}
}
