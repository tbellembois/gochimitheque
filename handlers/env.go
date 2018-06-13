package handlers

import (
	"html/template"

	"github.com/tbellembois/gochimitheque/models"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	DB        models.Datastore              // application DB connection
	Templates map[string]*template.Template // application templates
	// PersonEmail string                        // connected user email
	// Permissions []models.Permission           // connected user permissions - used by javascript to dynamically show/hide html elements
}
