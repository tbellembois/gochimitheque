package handlers

import (
	"crypto/rand"
	"errors"

	"github.com/casbin/casbin/v2"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/ldap"
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
	// TokenSignKey is the JWT token signing key
	TokenSignKey []byte
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
	// LDAP connection
	LDAPConnection *ldap.LDAPConnection
}

func NewEnv() Env {
	var (
		env Env
		err error
	)

	// Generate JWT signing key.
	if env.TokenSignKey, err = genSymmetricKey(64); err != nil {
		panic(err)
	}

	return env
}
