package mcmaclient

import (
	"encoding/json"
	"fmt"
	"strings"
)

type cacheKey struct {
	authType    string
	authContext string
}

type AuthProvider struct {
	authenticatorFactories map[string]AuthenticatorFactory
	cache                  map[cacheKey]Authenticator
}

func (authProvider *AuthProvider) Add(authType string, factory AuthenticatorFactory) {
	authProvider.authenticatorFactories[authType] = factory
}

func (authProvider *AuthProvider) Get(authType string, authContext interface{}) (Authenticator, error) {
	var authContextStr string
	switch v := authContext.(type) {
	case string:
		authContextStr = v
	case nil:
		authContextStr = ""
	default:
		authContextJson, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		authContextStr = string(authContextJson)
	}

	cacheKey := cacheKey{authType, authContextStr}
	cacheItem, found := authProvider.cache[cacheKey]
	if found {
		return cacheItem, nil
	}

	var authenticatorFactory AuthenticatorFactory
	for key, af := range authProvider.authenticatorFactories {
		if strings.EqualFold(key, authType) {
			authenticatorFactory = af
			break
		}
	}
	if authenticatorFactory == nil {
		return nil, fmt.Errorf("no authenticators registered for auth type '%s'", authType)
	}

	authenticator := authenticatorFactory.Get(authContext)
	authProvider.cache[cacheKey] = authenticator

	return authProvider.cache[cacheKey], nil
}

func newAuthProvider() *AuthProvider {
	return &AuthProvider{
		authenticatorFactories: make(map[string]AuthenticatorFactory),
		cache:                  make(map[cacheKey]Authenticator),
	}
}
