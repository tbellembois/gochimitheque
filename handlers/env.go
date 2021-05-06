package handlers

import (
	"crypto/rand"
	"errors"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	jsonadapter "github.com/casbin/json-adapter/v2"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/logger"
)

// genSymmetricKey generates a key for the JWT encryption
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

// Env is a structure used to pass variables throughout the application.
type Env struct {
	DB datastores.Datastore

	Enforcer    *casbin.Enforcer
	CasbinModel string

	// TokenSignKey is the JWT token signing key
	TokenSignKey []byte
	// ProxyPath is the application proxy path if behind a proxy
	// "/"" by default
	ProxyPath string
	// ProxyURL is application base url
	// "http://localhost:8081" by default
	ProxyURL string
	// ApplicationFullURL is application full url
	// "http://localhost:8081" by default
	// "ProxyURL + ProxyPath" if behind a proxy
	ApplicationFullURL string

	// BuildID is a compile time variable
	BuildID string
	// DisableCache disables the views cache
	DisableCache bool
}

func NewEnv() Env {

	var (
		err error
		env Env
	)

	// generate JWT signing key
	if env.TokenSignKey, err = genSymmetricKey(64); err != nil {
		panic(err)
	}

	return env

}

func (env *Env) InitCasbinPolicy() {

	var (
		err             error
		jsonAdapterData []byte
	)

	if jsonAdapterData, err = env.DB.ToCasbinJSONAdapter(); err != nil {
		logger.Log.Error("error getting json adapter data: " + err.Error())
		os.Exit(1)
	}

	m, e := model.NewModelFromString(env.CasbinModel)
	if e != nil {
		logger.Log.Error("model creation error: " + e.Error())
		os.Exit(1)
	}

	a := jsonadapter.NewAdapter(&jsonAdapterData)
	if env.Enforcer, err = casbin.NewEnforcer(m, a); err != nil {
		logger.Log.Error("enforcer creation error: " + err.Error())
		os.Exit(1)
	}

	env.Enforcer.AddFunction("matchStorage", env.MatchStorageFunc)
	env.Enforcer.AddFunction("matchStorelocation", env.MatchStorelocationFunc)
	env.Enforcer.AddFunction("matchPeople", env.MatchPeopleFunc)
	env.Enforcer.AddFunction("matchEntity", env.MatchEntityFunc)

	if err = env.Enforcer.LoadPolicy(); err != nil {
		logger.Log.Error("enforcer policy load error: " + err.Error())
		os.Exit(1)
	}

}
