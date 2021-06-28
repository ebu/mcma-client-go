package mcmaclient

type AuthenticatorFactory interface {
	Get(authContext interface{}) Authenticator
}
