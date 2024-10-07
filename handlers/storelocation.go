package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/barweiss/go-tuple"
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

// VGetStoreLocationsHandler handles the store location list page.
func (env *Env) VGetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storelocationindex(c, w)

	return nil
}

// VCreateStoreLocationHandler handles the store location creation page.
func (env *Env) VCreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storelocationcreate(c, w)

	return nil
}

/*
	REST handlers
*/

func (env *Env) GetStoreLocationsBSTABLEHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoreLocationsHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetStorelocations("http://localhost/?"+r.URL.RawQuery, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetStorelocations",
		}
	}

	// hack to return the response expected by bootstrap table.
	var (
		zmqResponse tuple.T2[[]models.StoreLocation, int]
	)
	if err = json.Unmarshal([]byte(jsonRawMessage), &zmqResponse); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error decoding zmqclient.DBGetStorelocations response",
		}
	}

	if err = json.NewEncoder(w).Encode(models.StoreLocationsResp{Rows: zmqResponse.V1, Total: zmqResponse.V2}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetStoreLocationsHandler godoc
// @Summary Get store locations. Only store locations visible by the authenticated user are returned.
// @tags store_location
// @Produce json
// @Success 200 {object} models.StoreLocationsResp
// @Failure 500
// @Failure 403
// @Router /store_locations/ [get].
func (env *Env) GetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoreLocationsHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetStorelocations("http://localhost/?"+r.URL.RawQuery, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetStorelocations",
		}
	}

	// Convert Rust response to former one.
	var tuple tuple.T2[interface{}, int]
	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return &models.AppError{
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
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling response",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// GetStoreLocationHandler returns a json of the store location with the requested id.
func (env *Env) GetStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	store_location, err := env.DB.GetStoreLocation(id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the store location",
		}
	}

	logger.Log.WithFields(logrus.Fields{"store_location": store_location}).Debug("GetStoreLocationHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(store_location); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreateStoreLocationHandler creates the store location from the request form.
func (env *Env) CreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateStoreLocationHandler")

	var (
		sl  models.StoreLocation
		err error
		id  int64
	)

	if err = json.NewDecoder(r.Body).Decode(&sl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"sl": sl}).Debug("CreateStoreLocationHandler")

	if id, err = env.DB.CreateStoreLocation(sl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "create store location error",
			Code:          http.StatusInternalServerError,
		}
	}
	sl.StoreLocationID = models.NullInt64{NullInt64: sql.NullInt64{Valid: true, Int64: id}}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(sl); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdateStoreLocationHandler updates the store location from the request form.
func (env *Env) UpdateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
		sl  models.StoreLocation
	)

	if err = json.NewDecoder(r.Body).Decode(&sl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"sl": sl}).Debug("UpdateStoreLocationHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	updatedsl, err := env.DB.GetStoreLocation(id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "get store location error",
			Code:          http.StatusInternalServerError,
		}
	}
	updatedsl.StoreLocationName = sl.StoreLocationName
	updatedsl.StoreLocationColor = sl.StoreLocationColor
	updatedsl.StoreLocationCanStore = sl.StoreLocationCanStore
	updatedsl.StoreLocation = sl.StoreLocation
	updatedsl.Entity = sl.Entity

	logger.Log.WithFields(logrus.Fields{"updatedsl": updatedsl}).Debug("UpdateStoreLocationHandler")

	if err := env.DB.UpdateStoreLocation(updatedsl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update store location error",
			Code:          http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(updatedsl); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// DeleteStoreLocationHandler deletes the store location with the requested id.
func (env *Env) DeleteStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	if err = env.DB.DeleteStoreLocation(id); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
