package handlers

import (
	"crypto/rand"
	"errors"

	"github.com/casbin/casbin/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/tbellembois/gochimitheque/datastores"
	"golang.org/x/oauth2"
)

// https://github.com/northbright/Notes/blob/master/jwt/generate_hmac_secret_key_for_jwt.md
func genSymmetricKey(bits int) (k []byte, err error) {
	if bits <= 0 || bits%8 != 0 {
		return nil, errors.New("key size error")
	}

	size := bits / 8
	k = make([]byte, size)
	if _, err = rand.Read(k); err != nil {
		return nil, err
	}

	return k, nil
}

// Env is used to pass variables throughout the application.
type Env struct {
	DB datastores.Datastore

	Enforcer *casbin.Enforcer

	// AutoCreateUser is used with a proxy authentication
	AutoCreateUser bool
	// AppPath is the application proxy path if behind a proxy
	// "/"" by default
	AppPath string
	// AppURL is application base url
	// "http://localhost:8081" by default
	AppURL string
	// AppFullURL is application full url
	// "ProxyURL + ProxyPath"
	AppFullURL string
	// DockerPort if used
	DockerPort int
	// BuildID is a compile time variable
	BuildID string
	// DisableCache disables the views cache
	DisableCache bool

	// OIDC parameters
	OIDCServer             string
	OIDCClientID           string
	OIDCClientSecret       string
	OIDCProvider           *oidc.Provider
	OIDCVerifier           *oidc.IDTokenVerifier
	OIDCConfig             *oidc.Config
	OAuth2Config           oauth2.Config
	OIDCEndSessionEndpoint string
}

func NewEnv() Env {
	return Env{}
}
