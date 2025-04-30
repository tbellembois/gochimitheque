package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

func sanitizeProduct(p *models.Product) {
	for i := range p.Synonyms {
		p.Synonyms[i].NameLabel = strings.Trim(p.Synonyms[i].NameLabel, " ")
	}

	p.NameLabel = strings.Trim(p.NameLabel, " ")

	if p.LinearFormulaLabel != nil {
		var LinearFormulaLabelPointer *string
		LinearFormulaLabelPointer = new(string)
		*LinearFormulaLabelPointer = strings.Trim(*p.LinearFormulaLabel, " ")
		p.LinearFormulaLabel = LinearFormulaLabelPointer
	}

	if p.EmpiricalFormulaLabel != nil {
		var EmpiricalFormulaLabelPointer *string
		EmpiricalFormulaLabelPointer = new(string)
		*EmpiricalFormulaLabelPointer = strings.Trim(*p.EmpiricalFormulaLabel, " ")
		p.EmpiricalFormulaLabel = EmpiricalFormulaLabelPointer
	}
	// p.CasNumberLabel.String = strings.Trim(p.CasNumberLabel.String, " ")

	if p.CasNumberLabel != nil {
		var CasNumberLabelPointer *string
		CasNumberLabelPointer = new(string)
		*CasNumberLabelPointer = strings.Trim(*p.CasNumberLabel, " ")
		p.CasNumberLabel = CasNumberLabelPointer
	}

	if p.CeNumberLabel != nil {
		var CeNumberLabelPointer *string
		CeNumberLabelPointer = new(string)
		*CeNumberLabelPointer = strings.Trim(*p.CeNumberLabel, " ")
		p.CeNumberLabel = CeNumberLabelPointer
	}
}

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

/*
	REST handlers
*/

// GetProductsProducerRefsHandler returns a json list of the producer_ref.
func (env *Env) GetProductsProducerRefsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	logger.Log.Debug("GetProductsProducerRefsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetProducerrefs("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetProducerrefs",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsSupplierRefsHandler returns a json list of the producer_ref.
func (env *Env) GetProductsSupplierRefsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSupplierRefsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetSupplierrefs("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetSupplierrefs",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

func (env *Env) PubchemGetProductByNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetCompoundByNameHandler")

	vars := mux.Vars(r)

	var (
		err     error
		product zmqclient.PubchemProduct
	)

	if product, err = zmqclient.PubchemGetProductByName(vars["name"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.GetCompoundByName",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(product); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (env *Env) PubchemGetCompoundByNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetCompoundByNameHandler")

	vars := mux.Vars(r)

	var (
		err       error
		compounds zmqclient.Compounds
	)

	if compounds, err = zmqclient.PubchemGetCompoundByName(vars["name"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.GetCompoundByName",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(compounds); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// PubchemAutocompleteHandler calls the autocomplete Pubchem API.
func (env *Env) PubchemAutocompleteHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("PubchemAutocompleteHandler")

	vars := mux.Vars(r)

	var (
		err          error
		autocomplete zmqclient.PubchemAutocomplete
	)

	if autocomplete, err = zmqclient.PubchemAutocompleteProductName(vars["name"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.AutocompleteProductName",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(autocomplete); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetProductsCategoriesHandler returns a json list of the producer.
func (env *Env) GetProductsCategoriesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsCategoriesHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetCategories("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetCategories",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsTagsHandler returns a json list of the tag.
func (env *Env) GetProductsTagsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsTagsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetTags("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.GetProductsTagsHandler",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsProducersHandler returns a json list of the producer.
func (env *Env) GetProductsProducersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsProducersHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetProducers("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetProducers",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsSuppliersHandler returns a json list of the supplier.
func (env *Env) GetProductsSuppliersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSuppliersHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetSuppliers("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetSuppliers",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

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
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)
	person.PersonID = c.PersonID

	if isbookmark, err = env.DB.IsProductBookmark(product, person); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting bookmark status",
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
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error creating the bookmark",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(product); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetProductsCasNumbersHandler returns a json list of the cas numbers matching the search criteria.
func (env *Env) GetProductsCasNumbersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	logger.Log.Debug("GetProductsCasNumbersHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetCasnumbers("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetCasnumbers",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil

}

// GetProductsCeNumbersHandler returns a json list of the ce numbers matching the search criteria.
func (env *Env) GetProductsCeNumbersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsCeNumbersHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetCenumbers("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetCenumbers",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsPhysicalStatesHandler returns a json list of the physical states matching the search criteria.
func (env *Env) GetProductsPhysicalStatesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsPhysicalStatesHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetPhysicalstates("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetPhysicalstates",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsSignalWordsHandler returns a json list of the signal words matching the search criteria.
func (env *Env) GetProductsSignalWordsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSignalWordsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetSignalwords("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetSignalwords",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsClassOfCompoundsHandler returns a json list of the classes of compounds matching the search criteria.
func (env *Env) GetProductsClassOfCompoundsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsClassOfCompoundsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetClassesofcompound("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetClassesofcompound",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsEmpiricalFormulasHandler returns a json list of the empirical formulas matching the search criteria.
func (env *Env) GetProductsEmpiricalFormulasHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsEmpiricalFormulasHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetEmpiricalformulas("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetEmpiricalformulas",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil

}

// GetProductsLinearFormulasHandler returns a json list of the linear formulas matching the search criteria.
func (env *Env) GetProductsLinearFormulasHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsLinearFormulasHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetLinearformulas("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetLinearformulas",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsSymbolsHandler returns a json list of the symbols matching the search criteria.
func (env *Env) GetProductsSymbolsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSymbolsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetSymbols("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetSymbols",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsHazardStatementsHandler returns a json list of the hazard statements matching the search criteria.
func (env *Env) GetProductsHazardStatementsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsHazardStatementsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetHazardstatements("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetHazardstatements",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsPrecautionaryStatementsHandler returns a json list of the precautionary statements matching the search criteria.
func (env *Env) GetProductsPrecautionaryStatementsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsPrecautionaryStatementsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetPrecautionarystatements("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetPrecautionarystatements",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetProductsNamesHandler returns a json list of the names matching the search criteria.
func (env *Env) GetProductsNamesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsNamesHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetNames("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetNames",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil

}

// GetProductsSynonymsHandler returns a json list of the symbols matching the search criteria.
func (env *Env) GetProductsSynonymsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsSynonymsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetNames("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetNames",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

func (env *Env) GetProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetProducts("http://localhost/?"+r.URL.RawQuery, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetProducts",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// CreateProductHandler creates the product from the request form.
func (env *Env) CreateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateProductHandler")

	var (
		p   models.Product
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	p.PersonID = c.PersonID

	logger.Log.WithFields(logrus.Fields{"p": fmt.Sprintf("%+v", p)}).Debug("CreateProductHandler")

	sanitizeProduct(&p)

	var pid int64

	if pid, err = env.DB.CreateUpdateProduct(p, false); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "create product error",
			Code:          http.StatusInternalServerError,
		}
	}

	p.ProductID = int(pid)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(p); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdateProductHandler updates the product from the request form.
func (env *Env) UpdateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
		p   models.Product
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	// p.ProductCreationDate = time.Now()
	p.PersonID = c.PersonID

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	var updatedp *models.Product

	// getting the product
	var (
		jsonRawMessage json.RawMessage
	)
	if jsonRawMessage, err = zmqclient.DBGetProducts("http://localhost/?product="+strconv.Itoa(id), p.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "zmqclient.DBGetProducts",
			Code:          http.StatusInternalServerError,
		}

	}

	if updatedp, err = zmqclient.ConvertDBJSONToProduct(jsonRawMessage); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "ConvertDBJSONToProduct",
			Code:          http.StatusInternalServerError,
		}
	}

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
	// updatedp.ProductMolFormula = p.ProductMolFormula
	updatedp.ProductInchi = p.ProductInchi
	updatedp.ProductInchikey = p.ProductInchikey
	updatedp.ProductCanonicalSmiles = p.ProductCanonicalSmiles
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
	updatedp.UnitMolecularWeight = p.UnitMolecularWeight

	logger.Log.WithFields(logrus.Fields{"updatedp": fmt.Sprintf("%+v", updatedp)}).Debug("UpdateProductHandler")

	sanitizeProduct(updatedp)
	if _, err := env.DB.CreateUpdateProduct(*updatedp, true); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update product error",
			Code:          http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(updatedp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// DeleteProductHandler deletes the store location with the requested id.
func (env *Env) DeleteProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if err := env.DB.DeleteProduct(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "delete product error",
			Code:          http.StatusInternalServerError,
		}
	}

	return nil
}

// CreateSupplierHandler creates the supplier from the request form.
func (env *Env) CreateSupplierHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateSupplierHandler")

	var (
		jsonRawMessage json.RawMessage
		body           []byte
		err            error
	)

	if body, err = io.ReadAll(r.Body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading request body",
		}
	}
	logger.Log.Debug("body " + string(body))

	if jsonRawMessage, err = zmqclient.DBCreateUpdateSupplier(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateSupplier",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// CreateProducerHandler creates the producer from the request form.
func (env *Env) CreateProducerHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateProducerHandler")

	var (
		jsonRawMessage json.RawMessage
		body           []byte
		err            error
	)

	if body, err = io.ReadAll(r.Body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading request body",
		}
	}
	logger.Log.Debug("body " + string(body))

	if jsonRawMessage, err = zmqclient.DBCreateUpdateProducer(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateProducer",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}
