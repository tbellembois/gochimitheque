package handlers

import (
	"github.com/tbellembois/gochimitheque/models"
	"html/template"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	DB        models.Datastore              // application DB connection
	Templates map[string]*template.Template // application templates
}
