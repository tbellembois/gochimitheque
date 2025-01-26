package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/barweiss/go-tuple"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

// Convert a JSON response from the chimitheque_db Rust library to a store location.
func ConvertDBJSONToStorelocation(jsonRawMessage json.RawMessage) (models.StoreLocation, error) {

	logger.Log.Debug("ConvertDBJSONToStorelocation")
	logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToStorelocation")

	var (
		tuple tuple.T2[[]models.StoreLocation, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return models.StoreLocation{}, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	resp := tuple.V1[0]

	logger.Log.WithFields(logrus.Fields{"resp": fmt.Sprintf("%+v", resp)}).Debug("ConvertDBJSONToStorelocation")

	return resp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a store location JSON.
func ConvertDBJSONToStorelocationJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	logger.Log.Debug("ConvertDBJSONToStorelocationJSON")
	logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToStorelocationJSON")

	var (
		storeLocation models.StoreLocation
		err           error
	)

	if storeLocation, err = ConvertDBJSONToStorelocation(jsonRawMessage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling to store location",
		}
	}

	logger.Log.WithFields(logrus.Fields{"storeLocation": fmt.Sprintf("%+v", storeLocation)}).Debug("ConvertDBJSONToStorelocationJSON")

	var jsonresp []byte
	if jsonresp, err = json.Marshal(storeLocation); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling response",
		}
	}

	return jsonresp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a BootstrapTable JSON.
func ConvertDBJSONToBSTableJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	var (
		tuple tuple.T2[interface{}, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	resp := struct {
		Rows     interface{} `json:"rows"`
		Total    int         `json:"total"`
		Exportfn string      `json:"exportfn"`
	}{
		Rows:     tuple.V1,
		Total:    tuple.V2,
		Exportfn: "",
	}

	var jsonresp []byte
	if jsonresp, err = json.Marshal(resp); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling response",
		}
	}

	return jsonresp, nil
}
