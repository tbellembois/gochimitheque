package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
	"github.com/tbellembois/gochimitheque/zmqclient"
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
		product_id int
	)

	vars := mux.Vars(r)

	if product_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	if _, err = zmqclient.DBToggleProductBookmark(c.PersonID, product_id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateProductFromPubchem",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode("ok"); err != nil {
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

func (env *Env) ExportProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("ExportProductsHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBExportProducts("http://localhost"+r.RequestURI, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetProducts",
		}
	}

	var (
		csv string
	)

	if err = json.Unmarshal(jsonRawMessage, &csv); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=chimitheque_products.csv")

	w.Write([]byte(csv))

	return nil
}

func (env *Env) GetProductsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetProductsHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetProducts("http://localhost"+r.RequestURI, c.PersonID); err != nil {
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

	if jsonRawMessage, err = zmqclient.DBCreateUpdateProduct(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateProduct",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// UpdateProductHandler updates the product from the request form.
func (env *Env) UpdateProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("UpdateProductHandler")

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

	if jsonRawMessage, err = zmqclient.DBCreateUpdateProduct(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateProduct",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// DeleteProductHandler deletes the store location with the requested id.
func (env *Env) DeleteProductHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		jsonRawMessage json.RawMessage
		id64           int64
		id             int
		err            error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	id64 = int64(id)

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeleteProductHandler")

	if jsonRawMessage, err = zmqclient.DBDeleteProduct(id64); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBDeleteProduct",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

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
