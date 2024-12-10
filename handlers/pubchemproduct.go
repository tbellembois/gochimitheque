package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

// CreateSupplierHandler creates the supplier from the request form.
func (env *Env) CreateUpdateProductFromPubchemHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateUpdateProductFromPubchemHandler")

	vars := mux.Vars(r)

	var (
		jsonRawMessage json.RawMessage
		body           []byte
		err            error
		product_id     *int
	)

	if len(vars["id"]) == 0 {
		product_id = nil
	} else {
		var id int
		if id, err = strconv.Atoi(vars["id"]); err != nil {
			return &models.AppError{
				OriginalError: err,
				Message:       "id atoi conversion",
				Code:          http.StatusInternalServerError,
			}
		}

		product_id = &id
	}

	c := request.ContainerFromRequestContext(r)

	if body, err = io.ReadAll(r.Body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading request body",
		}
	}
	logger.Log.Debug("body " + string(body))

	if jsonRawMessage, err = zmqclient.DBCreateUpdateProductFromPubchem(body, c.PersonID, product_id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateProductFromPubchem",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}
