package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/models"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	// DB is the database connection
	DB datastores.Datastore
	//Templates map[string]*template.Template // application templates
}

// FakeHandler returns true
func (env *Env) FakeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode("true"); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}
