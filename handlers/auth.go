package handlers

import (
	"net/http"

	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	Views handler.
*/

// VSearchHandler return the search div.
func (env *Env) VSearchHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Search(c, w)

	return nil
}

// VMenuHandler return the menu div.
func (env *Env) VMenuHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Menu(c, w)

	return nil
}

// VLoginHandler return the login page.
func (env *Env) VLoginHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Login(c, w)

	return nil
}
