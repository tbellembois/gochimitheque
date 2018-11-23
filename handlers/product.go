package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tbellembois/gochimitheque/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// VGetproductsHandler handles the store location list page
func (env *Env) VGetProductsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["productindex"].ExecuteTemplate(w, "BASE", c); e != nil {
		return &helpers.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

// VCreateProductHandler handles the store location creation page
func (env *Env) VCreateProductHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["productcreate"].ExecuteTemplate(w, "BASE", c); e != nil {
		return &helpers.AppError{
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

func (env *Env) ToogleProductBookmarkHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	var (
		err        error
		isbookmark bool
	)

	product := models.Product{}
	person := models.Person{}
	vars := mux.Vars(r)

	if product.ProductID, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)
	person.PersonID = c.PersonID

	if isbookmark, err = env.DB.IsProductBookmark(product, person); err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting bookmark status",
		}
	}

	// toggling the bookmark
	if isbookmark {
		err = env.DB.DeleteProductBookmark(product, person)
	} else {
		err = env.DB.CreateProductBookmark(product, person)
	}
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error creating the bookmark",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)
	return nil
}

// GetProductsCasNumbersHandler returns a json list of the cas numbers matching the search criteria
func (env *Env) GetProductsCasNumbersHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsCasNumbersHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	casnumbers, count, err := env.DB.GetProductsCasNumbers(dsp)
	if err != nil {
		return &helpers.AppError{
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

// GetProductsCeNumbersHandler returns a json list of the ce numbers matching the search criteria
func (env *Env) GetProductsCeNumbersHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsCeNumbersHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	cenumbers, count, err := env.DB.GetProductsCeNumbers(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the ce numbers",
		}
	}

	type resp struct {
		Rows  []models.CeNumber `json:"rows"`
		Total int               `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: cenumbers, Total: count})
	return nil
}

// GetProductsPhysicalStatesHandler returns a json list of the physical states matching the search criteria
func (env *Env) GetProductsPhysicalStatesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsPhysicalStatesHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	pstates, count, err := env.DB.GetProductsPhysicalStates(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the physical states",
		}
	}

	type resp struct {
		Rows  []models.PhysicalState `json:"rows"`
		Total int                    `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: pstates, Total: count})
	return nil
}

// GetProductsSignalWordsHandler returns a json list of the signal words matching the search criteria
func (env *Env) GetProductsSignalWordsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsSignalWordsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	swords, count, err := env.DB.GetProductsSignalWords(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the signal words",
		}
	}

	type resp struct {
		Rows  []models.SignalWord `json:"rows"`
		Total int                 `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: swords, Total: count})
	return nil
}

// GetProductsClassOfCompoundsHandler returns a json list of the classes of compounds matching the search criteria
func (env *Env) GetProductsClassOfCompoundsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsClassOfCompoundsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	cocs, count, err := env.DB.GetProductsClassOfCompounds(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the classes of compounds",
		}
	}

	type resp struct {
		Rows  []models.ClassOfCompound `json:"rows"`
		Total int                      `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: cocs, Total: count})
	return nil
}

// GetProductsEmpiricalFormulasHandler returns a json list of the empirical formulas matching the search criteria
func (env *Env) GetProductsEmpiricalFormulasHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsEmpiricalFormulasHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, utils.SortEmpiricalFormula); err != nil {
		return aerr
	}

	eformulas, count, err := env.DB.GetProductsEmpiricalFormulas(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the empirical formulas",
		}
	}

	type resp struct {
		Rows  []models.EmpiricalFormula `json:"rows"`
		Total int                       `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: eformulas, Total: count})
	return nil
}

// GetProductsNamesHandler returns a json list of the names matching the search criteria
func (env *Env) GetProductsNamesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsNamesHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	names, count, err := env.DB.GetProductsNames(dsp)
	if err != nil {
		return &helpers.AppError{
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

// GetProductsNameHandler returns a json of the name matching the id
func (env *Env) GetProductsNameHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsNameHandler")

	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	name, err := env.DB.GetProductsName(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the name",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(name)
	return nil
}

// GetProductsEmpiricalFormulaHandler returns a json of the formula matching the id
func (env *Env) GetProductsEmpiricalFormulaHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsEmpiricalFormulaHandler")

	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	ef, err := env.DB.GetProductsEmpiricalFormula(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the empirical formula",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ef)
	return nil
}

// GetProductsCasNumberHandler returns a json of the formula matching the id
func (env *Env) GetProductsCasNumberHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsCasNumberHandler")

	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	cas, err := env.DB.GetProductsCasNumber(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the cas number",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cas)
	return nil
}

// GetProductsSymbolsHandler returns a json list of the symbols matching the search criteria
func (env *Env) GetProductsSymbolsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsSymbolsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	symbols, count, err := env.DB.GetProductsSymbols(dsp)
	if err != nil {
		return &helpers.AppError{
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

// GetProductsHazardStatementsHandler returns a json list of the hazard statements matching the search criteria
func (env *Env) GetProductsHazardStatementsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsHazardStatementsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	hs, count, err := env.DB.GetProductsHazardStatements(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the hazard statements",
		}
	}

	type resp struct {
		Rows  []models.HazardStatement `json:"rows"`
		Total int                      `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: hs, Total: count})
	return nil
}

// GetProductsPrecautionaryStatementsHandler returns a json list of the precautionary statements matching the search criteria
func (env *Env) GetProductsPrecautionaryStatementsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsPrecautionaryStatementsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	ps, count, err := env.DB.GetProductsPrecautionaryStatements(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the precautionary statements",
		}
	}

	type resp struct {
		Rows  []models.PrecautionaryStatement `json:"rows"`
		Total int                             `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: ps, Total: count})
	return nil
}

// GetProductsSynonymsHandler returns a json list of the symbols matching the search criteria
func (env *Env) GetProductsSynonymsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsSynonymsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	synonyms, count, err := env.DB.GetProductsNames(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the synonyms",
		}
	}

	type resp struct {
		Rows  []models.Name `json:"rows"`
		Total int           `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: synonyms, Total: count})
	return nil
}

// GetProductsHandler returns a json list of the products matching the search criteria
func (env *Env) GetProductsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dspp helpers.DbselectparamProduct
	)

	// init db request parameters
	if dspp, aerr = helpers.NewdbselectparamProduct(r, nil); err != nil {
		return aerr
	}

	products, count, err := env.DB.GetProducts(dspp)
	if err != nil {
		return &helpers.AppError{
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
func (env *Env) GetProductHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	product, err := env.DB.GetProduct(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the product",
		}
	}
	log.WithFields(log.Fields{"product": product}).Debug("GetProductHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
	return nil
}

// CreateProductHandler creates the product from the request form
func (env *Env) CreateProductHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("CreateProductHandler")
	var (
		p   models.Product
		err error
	)
	if err = r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)

	if err = Decoder.Decode(&p, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	// p.ProductCreationDate = time.Now()
	p.PersonID = c.PersonID
	log.WithFields(log.Fields{"p": p}).Debug("CreateProductHandler")

	if err, p.ProductID = env.DB.CreateProduct(p); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "create product error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
	return nil
}

// UpdateProductHandler updates the product from the request form
func (env *Env) UpdateProductHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		p   models.Product
	)

	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)

	if err := Decoder.Decode(&p, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	// p.ProductCreationDate = time.Now()
	p.PersonID = c.PersonID
	log.WithFields(log.Fields{"p": p}).Debug("UpdateProductHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatedp, _ := env.DB.GetProduct(id)
	updatedp.CasNumber = p.CasNumber
	updatedp.CeNumber = p.CeNumber
	updatedp.EmpiricalFormula = p.EmpiricalFormula
	updatedp.ProductLinearFormula = p.ProductLinearFormula
	updatedp.Name = p.Name
	updatedp.ProductSpecificity = p.ProductSpecificity
	updatedp.Symbols = p.Symbols
	updatedp.Synonyms = p.Synonyms
	updatedp.ProductMSDS = p.ProductMSDS
	updatedp.ProductRestricted = p.ProductRestricted
	updatedp.ProductRadioactive = p.ProductRadioactive
	updatedp.ProductLinearFormula = p.ProductLinearFormula
	updatedp.ProductThreeDFormula = p.ProductThreeDFormula
	updatedp.ProductDisposalComment = p.ProductDisposalComment
	updatedp.ProductRemark = p.ProductRemark
	updatedp.PhysicalState = p.PhysicalState
	updatedp.SignalWord = p.SignalWord
	updatedp.ClassOfCompound = p.ClassOfCompound
	updatedp.HazardStatements = p.HazardStatements
	updatedp.PrecautionaryStatements = p.PrecautionaryStatements
	log.WithFields(log.Fields{"updatedp": updatedp}).Debug("UpdateProductHandler")

	if err := env.DB.UpdateProduct(updatedp); err != nil {
		return &helpers.AppError{
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
func (env *Env) DeleteProductHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	env.DB.DeleteProduct(id)
	return nil
}

// ConvertProductEmpiricalToLinearFormulaHandler returns the converted formula
func (env *Env) ConvertProductEmpiricalToLinearFormulaHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		l2ef string
	)

	log.Debug(vars["f"])

	l2ef = utils.LinearToEmpiricalFormula(vars["f"])

	log.Debug(l2ef)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(l2ef)

	return nil
}
