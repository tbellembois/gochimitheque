package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque-utils/convert"
	"github.com/tbellembois/gochimitheque-utils/sort"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/jade"
)

func sanitizeProduct(p *models.Product) {

	for i := range p.Synonyms {
		p.Synonyms[i].NameLabel = strings.Trim(p.Synonyms[i].NameLabel, " ")
	}
	p.NameLabel = strings.Trim(p.NameLabel, " ")
	p.LinearFormulaLabel.String = strings.Trim(p.LinearFormulaLabel.String, " ")
	p.EmpiricalFormulaLabel.String = strings.Trim(p.EmpiricalFormulaLabel.String, " ")
	p.CasNumberLabel.String = strings.Trim(p.CasNumberLabel.String, " ")
	p.CeNumberLabel.String = strings.Trim(p.CeNumberLabel.String, " ")

}

/*
	views handlers
*/

// VGetProductsHandler handles the store location list page
func (env *Env) VGetProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Productindex(c, w)

	return nil
}

// VCreateProductHandler handles the store location creation page
func (env *Env) VCreateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Productcreate(c, w)

	return nil
}

/*
	REST handlers
*/

// GetProductsProducerRefsHandler returns a json list of the producerref
func (env *Env) GetProductsProducerRefsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsProducerRefsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.DbselectparamProducerRef
	)

	// init db request parameters
	if dsp, aerr = models.NewdbselectparamProducerRef(r, nil); err != nil {
		return aerr
	}

	prefs, count, err := env.DB.GetProductsProducerRefs(dsp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the producerrefs",
		}
	}

	type resp struct {
		Rows  []models.ProducerRef `json:"rows"`
		Total int                  `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: prefs, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsSupplierRefsHandler returns a json list of the producerref
func (env *Env) GetProductsSupplierRefsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSupplierRefsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.DbselectparamSupplierRef
	)

	// init db request parameters
	if dsp, aerr = models.NewdbselectparamSupplierRef(r, nil); err != nil {
		return aerr
	}

	srefs, count, err := env.DB.GetProductsSupplierRefs(dsp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the supplierrefs",
		}
	}

	type resp struct {
		Rows  []models.SupplierRef `json:"rows"`
		Total int                  `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: srefs, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsCategoriesHandler returns a json list of the producer
func (env *Env) GetProductsCategoriesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsCategoriesHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	cats, count, err := env.DB.GetProductsCategories(dsp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the categories",
		}
	}

	type resp struct {
		Rows  []models.Category `json:"rows"`
		Total int               `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: cats, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsTagsHandler returns a json list of the tag
func (env *Env) GetProductsTagsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsTagsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	tags, count, err := env.DB.GetProductsTags(dsp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the tags",
		}
	}

	type resp struct {
		Rows  []models.Tag `json:"rows"`
		Total int          `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: tags, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsProducersHandler returns a json list of the producer
func (env *Env) GetProductsProducersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsProducersHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	prs, count, err := env.DB.GetProductsProducers(dsp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the producers",
		}
	}

	type resp struct {
		Rows  []models.Producer `json:"rows"`
		Total int               `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: prs, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsSuppliersHandler returns a json list of the supplier
func (env *Env) GetProductsSuppliersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSuppliersHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	srs, count, err := env.DB.GetProductsSuppliers(dsp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the suppliers",
		}
	}

	type resp struct {
		Rows  []models.Supplier `json:"rows"`
		Total int               `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: srs, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// ToogleProductBookmarkHandler (un)bookmarks the product with id passed in the request vars
// for the logged user.
func (env *Env) ToogleProductBookmarkHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		err        error
		isbookmark bool
	)

	product := models.Product{}
	person := models.Person{}
	vars := mux.Vars(r)

	if product.ProductID, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)
	person.PersonID = c.PersonID

	if isbookmark, err = env.DB.IsProductBookmark(product, person); err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting bookmark status",
		}
	}

	// toggling the bookmark
	if isbookmark {
		err = env.DB.DeleteProductBookmark(product, person)
		product.Bookmark = nil
	} else {
		err = env.DB.CreateProductBookmark(product, person)
		product.Bookmark = &models.Bookmark{
			Person:  person,
			Product: product,
		}
	}
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error creating the bookmark",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(product); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsCasNumbersHandler returns a json list of the cas numbers matching the search criteria
func (env *Env) GetProductsCasNumbersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsCasNumbersHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	// copy/paste CAS can send wrong separators (ie "-")
	// we must then rebuild the correct CAS
	cas := dsp.GetSearch()
	rcas := regexp.MustCompile("(?P<groupone>[0-9]{1,7}).{1}(?P<grouptwo>[0-9]{2}).{1}(?P<groupthree>[0-9]{1})")
	// finding group names
	n := rcas.SubexpNames()
	// finding matches
	ms := rcas.FindAllStringSubmatch(cas, -1)

	if len(ms) > 0 {
		m := ms[0]
		// then building a map of matches
		md := map[string]string{}
		for i, j := range m {
			md[n[i]] = j
		}
		dsp.SetSearch(fmt.Sprintf("%s-%s-%s", md["groupone"], md["grouptwo"], md["groupthree"]))
	}

	casnumbers, count, err := env.DB.GetProductsCasNumbers(dsp)
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
	if err = json.NewEncoder(w).Encode(resp{Rows: casnumbers, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsCeNumbersHandler returns a json list of the ce numbers matching the search criteria
func (env *Env) GetProductsCeNumbersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsCeNumbersHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	cenumbers, count, err := env.DB.GetProductsCeNumbers(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: cenumbers, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsPhysicalStatesHandler returns a json list of the physical states matching the search criteria
func (env *Env) GetProductsPhysicalStatesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsPhysicalStatesHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	pstates, count, err := env.DB.GetProductsPhysicalStates(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: pstates, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsSignalWordsHandler returns a json list of the signal words matching the search criteria
func (env *Env) GetProductsSignalWordsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSignalWordsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	swords, count, err := env.DB.GetProductsSignalWords(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: swords, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsClassOfCompoundsHandler returns a json list of the classes of compounds matching the search criteria
func (env *Env) GetProductsClassOfCompoundsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsClassOfCompoundsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	cocs, count, err := env.DB.GetProductsClassOfCompounds(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: cocs, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsEmpiricalFormulasHandler returns a json list of the empirical formulas matching the search criteria
func (env *Env) GetProductsEmpiricalFormulasHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsEmpiricalFormulasHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, sort.SortEmpiricalFormula); aerr != nil {
		return aerr
	}

	eformulas, count, err := env.DB.GetProductsEmpiricalFormulas(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: eformulas, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsLinearFormulasHandler returns a json list of the linear formulas matching the search criteria
func (env *Env) GetProductsLinearFormulasHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsLinearFormulasHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	lformulas, count, err := env.DB.GetProductsLinearFormulas(dsp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the empirical formulas",
		}
	}

	type resp struct {
		Rows  []models.LinearFormula `json:"rows"`
		Total int                    `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: lformulas, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsNamesHandler returns a json list of the names matching the search criteria
func (env *Env) GetProductsNamesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsNamesHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	names, count, err := env.DB.GetProductsNames(dsp)
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
	if err = json.NewEncoder(w).Encode(resp{Rows: names, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsNameHandler returns a json of the name matching the id
func (env *Env) GetProductsNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsNameHandler")

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

	name, err := env.DB.GetProductsName(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the name",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(name); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsEmpiricalFormulaHandler returns a json of the formula matching the id
func (env *Env) GetProductsEmpiricalFormulaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsEmpiricalFormulaHandler")

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

	ef, err := env.DB.GetProductsEmpiricalFormula(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the empirical formula",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(ef); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsCasNumberHandler returns a json of the formula matching the id
func (env *Env) GetProductsCasNumberHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsCasNumberHandler")

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

	cas, err := env.DB.GetProductsCasNumber(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the cas number",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(cas); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsSignalWordHandler returns a json of the signal word matching the id
func (env *Env) GetProductsSignalWordHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSignalWordHandler")

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

	signalword, err := env.DB.GetProductsSignalWord(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the signal word",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(signalword); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsSymbolsHandler returns a json list of the symbols matching the search criteria
func (env *Env) GetProductsSymbolsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSymbolsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	symbols, count, err := env.DB.GetProductsSymbols(dsp)
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
	if err = json.NewEncoder(w).Encode(resp{Rows: symbols, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsSymbolHandler returns a json of the symbol matching the id
func (env *Env) GetProductsSymbolHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSymbolHandler")

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

	symbol, err := env.DB.GetProductsSymbol(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the symbol",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(symbol); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsHazardStatementsHandler returns a json list of the hazard statements matching the search criteria
func (env *Env) GetProductsHazardStatementsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsHazardStatementsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	hs, count, err := env.DB.GetProductsHazardStatements(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: hs, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsHazardStatementHandler returns a json of the hazardstatement matching the id
func (env *Env) GetProductsHazardStatementHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsHazardStatementHandler")

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

	hs, err := env.DB.GetProductsHazardStatement(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the hazardstatement",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(hs); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsPrecautionaryStatementsHandler returns a json list of the precautionary statements matching the search criteria
func (env *Env) GetProductsPrecautionaryStatementsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsPrecautionaryStatementsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	ps, count, err := env.DB.GetProductsPrecautionaryStatements(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: ps, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsPrecautionaryStatementHandler returns a json of the precautionarystatement matching the id
func (env *Env) GetProductsPrecautionaryStatementHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsPrecautionaryStatementHandler")

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

	ps, err := env.DB.GetProductsPrecautionaryStatement(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the precautionarystatement",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(ps); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsSynonymsHandler returns a json list of the symbols matching the search criteria
func (env *Env) GetProductsSynonymsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSynonymsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	synonyms, count, err := env.DB.GetProductsNames(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: synonyms, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetExposedProductsHandler returns a json of the product with the requested id
func (env *Env) GetExposedProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetExposedProductsHandler")

	var (
		err error
	)

	products, count, err := env.DB.GetExposedProducts()
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
	if err = json.NewEncoder(w).Encode(resp{Rows: products, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsHandler returns a json list of the products matching the search criteria
func (env *Env) GetProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsHandler")

	var (
		err      error
		aerr     *models.AppError
		dspp     models.DbselectparamProduct
		exportfn string
	)

	// init db request parameters
	if dspp, aerr = models.NewdbselectparamProduct(r, nil); err != nil {
		return aerr
	}

	products, count, err := env.DB.GetProducts(dspp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the products",
		}
	}

	// export?
	if _, export := r.URL.Query()["export"]; export {
		exportfn = models.ProductsToCSV(products)
		// emptying results on exports
		products = []models.Product{}
		count = 0
	}

	type resp struct {
		Rows     []models.Product `json:"rows"`
		Total    int              `json:"total"`
		ExportFN string           `json:"exportfn"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: products, Total: count, ExportFN: exportfn}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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
			Message: "error getting the product",
		}
	}
	logger.Log.WithFields(logrus.Fields{"product": product}).Debug("GetProductHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(product); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreateProductHandler creates the product from the request form
func (env *Env) CreateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateProductHandler")
	var (
		p   models.Product
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)

	p.PersonID = c.PersonID
	logger.Log.WithFields(logrus.Fields{"p": fmt.Sprintf("%+v", p)}).Debug("CreateProductHandler")

	sanitizeProduct(&p)
	if p.ProductID, err = env.DB.CreateProduct(p); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "create product error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(p); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)

	// p.ProductCreationDate = time.Now()
	p.PersonID = c.PersonID
	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("UpdateProductHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatedp, _ := env.DB.GetProduct(id)
	updatedp.CasNumber = p.CasNumber
	updatedp.CeNumber = p.CeNumber
	updatedp.EmpiricalFormula = p.EmpiricalFormula
	updatedp.LinearFormula = p.LinearFormula
	updatedp.Name = p.Name
	updatedp.ProductSpecificity = p.ProductSpecificity
	updatedp.Symbols = p.Symbols
	updatedp.Synonyms = p.Synonyms
	updatedp.ProductMSDS = p.ProductMSDS
	updatedp.ProductRestricted = p.ProductRestricted
	updatedp.ProductRadioactive = p.ProductRadioactive
	updatedp.LinearFormula = p.LinearFormula
	updatedp.ProductThreeDFormula = p.ProductThreeDFormula
	updatedp.ProductTwoDFormula = p.ProductTwoDFormula
	updatedp.ProductMolFormula = p.ProductMolFormula
	updatedp.ProductDisposalComment = p.ProductDisposalComment
	updatedp.ProductRemark = p.ProductRemark
	updatedp.ProductNumberPerCarton = p.ProductNumberPerCarton
	updatedp.ProductNumberPerBag = p.ProductNumberPerBag
	updatedp.PhysicalState = p.PhysicalState
	updatedp.SignalWord = p.SignalWord
	updatedp.ClassOfCompound = p.ClassOfCompound
	updatedp.HazardStatements = p.HazardStatements
	updatedp.PrecautionaryStatements = p.PrecautionaryStatements
	updatedp.Tags = p.Tags
	updatedp.Category = p.Category
	updatedp.ProducerRef = p.ProducerRef
	updatedp.SupplierRefs = p.SupplierRefs
	updatedp.ProductSheet = p.ProductSheet
	updatedp.ProductTemperature = p.ProductTemperature
	updatedp.UnitTemperature = p.UnitTemperature
	logger.Log.WithFields(logrus.Fields{"updatedp": updatedp}).Debug("UpdateProductHandler")

	sanitizeProduct(&updatedp)
	if err := env.DB.UpdateProduct(updatedp); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "update product error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(updatedp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

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

	if err := env.DB.DeleteProduct(id); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "delete product error",
			Code:    http.StatusInternalServerError}
	}

	return nil
}

// ConvertProductEmpiricalToLinearFormulaHandler returns the converted formula
func (env *Env) ConvertProductEmpiricalToLinearFormulaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		l2ef string
		err  error
	)

	l2ef = convert.LinearToEmpiricalFormula(vars["f"])

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(l2ef); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// CreateSupplierHandler creates the supplier from the request form
func (env *Env) CreateSupplierHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateSupplierHandler")
	var (
		sup models.Supplier
		err error
		id  int
	)

	if err = json.NewDecoder(r.Body).Decode(&sup); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// if err := r.ParseForm(); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form parsing error",
	// 		Code:    http.StatusBadRequest}
	// }

	// if err := globals.Decoder.Decode(&sup, r.PostForm); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form decoding error",
	// 		Code:    http.StatusBadRequest}
	// }

	if id, err = env.DB.CreateSupplier(sup); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "create supplier error",
			Code:    http.StatusInternalServerError}
	}
	sup.SupplierID = sql.NullInt64{Valid: true, Int64: int64(id)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(sup); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreateProducerHandler creates the producer from the request form
func (env *Env) CreateProducerHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateProducerHandler")
	var (
		pr  models.Producer
		err error
		id  int
	)

	if err = json.NewDecoder(r.Body).Decode(&pr); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// if err := r.ParseForm(); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form parsing error",
	// 		Code:    http.StatusBadRequest}
	// }

	// if err := globals.Decoder.Decode(&pr, r.PostForm); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form decoding error",
	// 		Code:    http.StatusBadRequest}
	// }

	if id, err = env.DB.CreateProducer(pr); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "create producer error",
			Code:    http.StatusInternalServerError}
	}
	pr.ProducerID = sql.NullInt64{Valid: true, Int64: int64(id)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(pr); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}
