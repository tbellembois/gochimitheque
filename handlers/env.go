package handlers

import (
	"github.com/casbin/casbin/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/tbellembois/gochimitheque/datastores"
	"golang.org/x/oauth2"
)

// Env is used to pass variables throughout the application.
type Env struct {
	DB datastores.Datastore

	Enforcer *casbin.Enforcer

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
	OIDCDiscoverURL        string
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
