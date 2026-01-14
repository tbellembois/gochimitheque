package handlers

import (
	"net/http"

	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	views handlers
*/

// VGetProductsHandler handles the store location list page.
func (env *Env) VGetProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Productindex(c, w)

	return nil
}

// VCreateProductHandler handles the store location creation page.
func (env *Env) VCreateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Productcreate(c, w)

	return nil
}

func (env *Env) VPubchemHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Productpubchem(c, w)

	return nil
}
