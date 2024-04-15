package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

// CreateSupplierHandler creates the supplier from the request form.
func (env *Env) CreateProductFromPubchemHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateProductFromPubchemHandler")

	var (
		jsonRawMessage json.RawMessage
		body           []byte
		err            error
	)

	c := request.ContainerFromRequestContext(r)

	if body, err = io.ReadAll(r.Body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading request body",
		}
	}
	logger.Log.Debug("body " + string(body))

	if jsonRawMessage, err = zmqclient.DBCreateProductFromPubchem(body, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateProductFromPubchem",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}
