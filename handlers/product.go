package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/tbellembois/gochimitheque/jade"
	"github.com/tbellembois/gochimitheque/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// VGetProductsHandler handles the store location list page
func (env *Env) VGetProductsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Productindex(c, w)

	return nil
}

// VCreateProductHandler handles the store location creation page
func (env *Env) VCreateProductHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Productcreate(c, w)

	return nil
}

/*
	REST handlers
*/

// MagicHandler handles the magical selector.
func (env *Env) MagicHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("MagicHandler")

	rhs := regexp.MustCompile("((?:EU){0,1}H[0-9]{3}[FfDdAi]{0,2})")
	rps := regexp.MustCompile("(P[0-9]{3})")

	// form receiver
	type magic struct {
		MSDS string
	}
	// response
	type Resp struct {
		HS []models.HazardStatement        `json:"hs"`
		PS []models.PrecautionaryStatement `json:"ps"`
	}

	var (
		err  error
		m    magic
		hs   models.HazardStatement
		ps   models.PrecautionaryStatement
		resp Resp
	)

	if err = r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	if err = global.Decoder.Decode(&m, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}

	shs := rhs.FindAllStringSubmatch(m.MSDS, -1)
	sps := rps.FindAllStringSubmatch(m.MSDS, -1)

	for _, h := range shs {
		// silent db errors
		hs, err = env.DB.GetProductsHazardStatementByReference(h[1])
		if err != sql.ErrNoRows {
			resp.HS = append(resp.HS, hs)
		}
	}
	for _, p := range sps {
		// silent db errors
		ps, err = env.DB.GetProductsPrecautionaryStatementByReference(p[1])
		if err != sql.ErrNoRows {
			resp.PS = append(resp.PS, ps)
		}
	}

	log.WithFields(log.Fields{"m.msds": m.MSDS, "shs": shs, "sps": sps}).Debug("MagicHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(resp)

	return nil
}

// ToogleProductBookmarkHandler (un)bookmarks the product with id passed in the request vars
// for the logged user.
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

	// copy/paste CAS can send wrong separators (ie "-")
	// we must then rebuild the correct CAS
	cas := dsp.GetSearch()
	rcas := regexp.MustCompile("(?P<groupone>[0-9]{1,7}).{1}(?P<grouptwo>[0-9]{2}).{1}(?P<groupthree>[0-9]{1})")
	// finding group names
	n := rcas.SubexpNames()
	// finding matches
	ms := rcas.FindAllStringSubmatch(cas, -1)
	log.Debug(cas)
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

// GetProductsLinearFormulasHandler returns a json list of the linear formulas matching the search criteria
func (env *Env) GetProductsLinearFormulasHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsLinearFormulasHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	lformulas, count, err := env.DB.GetProductsLinearFormulas(dsp)
	if err != nil {
		return &helpers.AppError{
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
	json.NewEncoder(w).Encode(resp{Rows: lformulas, Total: count})
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

// GetProductsSignalWordHandler returns a json of the signal word matching the id
func (env *Env) GetProductsSignalWordHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsSignalWordHandler")

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

	signalword, err := env.DB.GetProductsSignalWord(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the signal word",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(signalword)
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

// GetProductsSymbolHandler returns a json of the symbol matching the id
func (env *Env) GetProductsSymbolHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsSymbolHandler")

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

	symbol, err := env.DB.GetProductsSymbol(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the symbol",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(symbol)
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

// GetProductsHazardStatementHandler returns a json of the hazardstatement matching the id
func (env *Env) GetProductsHazardStatementHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsHazardStatementHandler")

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

	hs, err := env.DB.GetProductsHazardStatement(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the hazardstatement",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hs)
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

// GetProductsPrecautionaryStatementHandler returns a json of the precautionarystatement matching the id
func (env *Env) GetProductsPrecautionaryStatementHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetProductsPrecautionaryStatementHandler")

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

	ps, err := env.DB.GetProductsPrecautionaryStatement(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the precautionarystatement",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ps)
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
		err      error
		aerr     *helpers.AppError
		dspp     helpers.DbselectparamProduct
		exportfn string
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
	json.NewEncoder(w).Encode(resp{Rows: products, Total: count, ExportFN: exportfn})
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

	if err = global.Decoder.Decode(&p, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	// p.ProductCreationDate = time.Now()
	p.PersonID = c.PersonID
	log.WithFields(log.Fields{"p": p}).Debug("CreateProductHandler")

	if p.ProductID, err = env.DB.CreateProduct(p); err != nil {
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

	if err := global.Decoder.Decode(&p, r.PostForm); err != nil {
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
	updatedp.ProductMolFormula = p.ProductMolFormula
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

	if err := env.DB.DeleteProduct(id); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "delete product error",
			Code:    http.StatusInternalServerError}
	}

	return nil
}

// ConvertProductEmpiricalToLinearFormulaHandler returns the converted formula
func (env *Env) ConvertProductEmpiricalToLinearFormulaHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		l2ef string
	)

	l2ef = utils.LinearToEmpiricalFormula(vars["f"])

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(l2ef)

	return nil
}
