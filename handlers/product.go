package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// VGetproductsHandler handles the store location list page
func (env *Env) VGetProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := containerFromRequestContext(r)

	if e := env.Templates["productindex"].Execute(w, c); e != nil {
		return &models.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

// VCreateProductHandler handles the store location creation page
func (env *Env) VCreateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := containerFromRequestContext(r)

	if e := env.Templates["productcreate"].Execute(w, c); e != nil {
		return &models.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

/*
	REST handlers
*/

// GetProductsCasNumbersHandler returns a json list of the cas numbers matching the search criteria
func (env *Env) GetProductsCasNumbersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("GetProductsCasNumbersHandler")

	var (
		search string
		order  string
		offset uint64
		limit  uint64
		err    error
	)

	if s, ok := r.URL.Query()["search"]; !ok {
		search = ""
	} else {
		search = s[0]
	}
	if o, ok := r.URL.Query()["order"]; !ok {
		order = "asc"
	} else {
		order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; !ok {
		offset = 0
	} else {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "offset atoi conversion",
			}
		}
		offset = uint64(of)
	}
	if l, ok := r.URL.Query()["limit"]; !ok {
		limit = constants.MaxUint64
	} else {
		var lm int
		if lm, err = strconv.Atoi(l[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		limit = uint64(lm)
	}

	// retrieving the logged user id from request context
	c := containerFromRequestContext(r)
	casnumbers, count, err := env.DB.GetProductsCasNumbers(models.GetCommonParameters{LoggedPersonID: c.PersonID, Search: search, Order: order, Offset: offset, Limit: limit})
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the cas numbers",
		}
	}

	type resp struct {
		Rows  []models.CasNumber `json:"rows"`
		Total int                `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: casnumbers, Total: count})
	return nil
}

// GetProductsNamesHandler returns a json list of the names matching the search criteria
func (env *Env) GetProductsNamesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("GetProductsNamesHandler")

	var (
		search string
		order  string
		offset uint64
		limit  uint64
		err    error
	)

	if s, ok := r.URL.Query()["search"]; !ok {
		search = ""
	} else {
		search = s[0]
	}
	if o, ok := r.URL.Query()["order"]; !ok {
		order = "asc"
	} else {
		order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; !ok {
		offset = 0
	} else {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "offset atoi conversion",
			}
		}
		offset = uint64(of)
	}
	if l, ok := r.URL.Query()["limit"]; !ok {
		limit = constants.MaxUint64
	} else {
		var lm int
		if lm, err = strconv.Atoi(l[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		limit = uint64(lm)
	}

	// retrieving the logged user id from request context
	c := containerFromRequestContext(r)
	names, count, err := env.DB.GetProductsNames(models.GetCommonParameters{LoggedPersonID: c.PersonID, Search: search, Order: order, Offset: offset, Limit: limit})
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the cas numbers",
		}
	}

	type resp struct {
		Rows  []models.Name `json:"rows"`
		Total int           `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: names, Total: count})
	return nil
}

// GetProductsSymbolsHandler returns a json list of the symbols matching the search criteria
func (env *Env) GetProductsSymbolsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("GetProductsSymbolsHandler")

	var (
		search string
		order  string
		offset uint64
		limit  uint64
		err    error
	)

	if s, ok := r.URL.Query()["search"]; !ok {
		search = ""
	} else {
		search = s[0]
	}
	if o, ok := r.URL.Query()["order"]; !ok {
		order = "asc"
	} else {
		order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; !ok {
		offset = 0
	} else {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "offset atoi conversion",
			}
		}
		offset = uint64(of)
	}
	if l, ok := r.URL.Query()["limit"]; !ok {
		limit = constants.MaxUint64
	} else {
		var lm int
		if lm, err = strconv.Atoi(l[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		limit = uint64(lm)
	}

	// retrieving the logged user id from request context
	c := containerFromRequestContext(r)
	symbols, count, err := env.DB.GetProductsSymbols(models.GetCommonParameters{LoggedPersonID: c.PersonID, Search: search, Order: order, Offset: offset, Limit: limit})
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the symbols",
		}
	}

	type resp struct {
		Rows  []models.Symbol `json:"rows"`
		Total int             `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: symbols, Total: count})
	return nil
}

// GetProductsHandler returns a json list of the products matching the search criteria
func (env *Env) GetProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("GetProductsHandler")

	var (
		search string
		order  string
		offset uint64
		limit  uint64
		err    error
	)

	if s, ok := r.URL.Query()["search"]; !ok {
		search = ""
	} else {
		search = s[0]
	}
	if o, ok := r.URL.Query()["order"]; !ok {
		order = "asc"
	} else {
		order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; !ok {
		offset = 0
	} else {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "offset atoi conversion",
			}
		}
		offset = uint64(of)
	}
	if l, ok := r.URL.Query()["limit"]; !ok {
		limit = constants.MaxUint64
	} else {
		var lm int
		if lm, err = strconv.Atoi(l[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		limit = uint64(lm)
	}

	// retrieving the logged user id from request context
	c := containerFromRequestContext(r)
	products, count, err := env.DB.GetProducts(models.GetProductsParameters{GetCommonParameters: models.GetCommonParameters{LoggedPersonID: c.PersonID, Search: search, Order: order, Offset: offset, Limit: limit}})
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the products",
		}
	}

	type resp struct {
		Rows  []models.Product `json:"rows"`
		Total int              `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: products, Total: count})
	return nil
}

// GetProductHandler returns a json of the product with the requested id
func (env *Env) GetProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	product, err := env.DB.GetProduct(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the store location",
		}
	}
	log.WithFields(log.Fields{"product": product}).Debug("GetProductHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
	return nil
}

// CreateProductHandler creates the product from the request form
func (env *Env) CreateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	return nil
}

// UpdateProductHandler updates the product from the request form
func (env *Env) UpdateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		p   models.Product
	)

	if err := r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	// if a new name is entered (ie instead of selecting an existing name)
	// r.Form["name.name_id"] == r.Form["name.name_label"]
	// then modifying the name_id to prevent a form decoding error
	if r.PostForm["name.name_id"][0] == r.PostForm["name.name_label"][0] {
		r.PostForm.Set("name.name_id", "-1")
	}
	// idem for casnumber
	if r.PostForm["casnumber.casnumber_id"][0] == r.PostForm["casnumber.casnumber_label"][0] {
		r.PostForm.Set("casnumber.casnumber_id", "-1")
	}

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&p, r.PostForm); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"p": p}).Debug("UpdateProductHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatedp, _ := env.DB.GetProduct(id)
	updatedp.CasNumber = p.CasNumber
	updatedp.Name = p.Name
	updatedp.ProductSpecificity = p.ProductSpecificity
	updatedp.Symbols = p.Symbols
	log.WithFields(log.Fields{"updatedp": updatedp}).Debug("UpdateProductHandler")

	if err := env.DB.UpdateProduct(updatedp); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "update product error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedp)

	return nil
}

// DeleteProductHandler deletes the store location with the requested id
func (env *Env) DeleteProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	env.DB.DeleteProduct(id)
	return nil
}
