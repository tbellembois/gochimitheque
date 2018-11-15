package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	DB        models.Datastore              // application DB connection
	Templates map[string]*template.Template // application templates
	ProxyPath string                        // application proxy path if behind a proxy
}

// FakeHandler returns true
func (env *Env) FakeHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("true")
	return nil
}
