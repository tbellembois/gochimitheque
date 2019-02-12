package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	DB        models.Datastore              // application DB connection
	Localizer *i18n.Localizer               // application i18n localizer
	Templates map[string]*template.Template // application templates
}

// FakeHandler returns true
func (env *Env) FakeHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("true")
	return nil
}
