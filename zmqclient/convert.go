package zmqclient

import (
	"encoding/json"
	"net/http"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque/models"
)

// Convert a JSON response from the chimitheque_db Rust library to a person.
func ConvertDBJSONToPerson(jsonRawMessage json.RawMessage) (*models.Person, error) {

	// logger.Log.Debug("ConvertDBJSONToPerson")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToPerson")

	var (
		tuple tuple.T2[[]models.Person, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	var resp *models.Person
	if len(tuple.V1) == 0 {
		resp = nil
	} else {
		resp = &tuple.V1[0]
	}

	// logger.Log.WithFields(logrus.Fields{"resp": fmt.Sprintf("%+v", resp)}).Debug("ConvertDBJSONToPerson")

	return resp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a person JSON.
func ConvertDBJSONToPersonJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	// logger.Log.Debug("ConvertDBJSONToPersonJSON")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToPersonJSON")

	var (
		person *models.Person
		err    error
	)

	if person, err = ConvertDBJSONToPerson(jsonRawMessage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	// logger.Log.WithFields(logrus.Fields{"person": fmt.Sprint(litter.Sdump(person))}).Debug("ConvertDBJSONToPersonJSON")

	var jsonresp []byte
	if jsonresp, err = json.Marshal(person); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling person",
		}
	}

	return jsonresp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a slice of entity.
func ConvertDBJSONToEntities(jsonRawMessage json.RawMessage) ([]models.Entity, error) {

	// logger.Log.Debug("ConvertDBJSONToEntities")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToEntities")

	var (
		tuple tuple.T2[[]models.Entity, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return []models.Entity{}, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	resp := tuple.V1

	// logger.Log.WithFields(logrus.Fields{"resp": fmt.Sprintf("%+v", resp)}).Debug("ConvertDBJSONToEntities")

	return resp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a slice of entity JSON.
func ConvertDBJSONToEntitiesJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	// logger.Log.Debug("ConvertDBJSONToEntitiesJSON")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToEntitiesJSON")

	var (
		entities []models.Entity
		err      error
	)

	if entities, err = ConvertDBJSONToEntities(jsonRawMessage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	// logger.Log.WithFields(logrus.Fields{"entity": fmt.Sprintf("%+v", entities)}).Debug("ConvertDBJSONToEntitiesJSON")

	var jsonresp []byte
	if jsonresp, err = json.Marshal(entities); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling entities",
		}
	}

	return jsonresp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a entity.
func ConvertDBJSONToEntity(jsonRawMessage json.RawMessage) (*models.Entity, error) {

	// logger.Log.Debug("ConvertDBJSONToEntity")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToEntity")

	var (
		tuple tuple.T2[[]models.Entity, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	var resp *models.Entity
	if len(tuple.V1) == 0 {
		resp = nil
	} else {
		resp = &tuple.V1[0]
	}

	// logger.Log.WithFields(logrus.Fields{"resp": fmt.Sprintf("%+v", resp)}).Debug("ConvertDBJSONToEntity")

	return resp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a entity JSON.
func ConvertDBJSONToEntityJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	// logger.Log.Debug("ConvertDBJSONToEntityJSON")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToEntityJSON")

	var (
		entity *models.Entity
		err    error
	)

	if entity, err = ConvertDBJSONToEntity(jsonRawMessage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	// logger.Log.WithFields(logrus.Fields{"entity": fmt.Sprintf("%+v", entity)}).Debug("ConvertDBJSONToStorelocationJSON")

	var jsonresp []byte
	if jsonresp, err = json.Marshal(entity); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling entity",
		}
	}

	return jsonresp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a store location.
func ConvertDBJSONToStorelocation(jsonRawMessage json.RawMessage) (*models.StoreLocation, error) {

	// logger.Log.Debug("ConvertDBJSONToStorelocation")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToStorelocation")

	var (
		tuple tuple.T2[[]models.StoreLocation, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	var resp *models.StoreLocation
	if len(tuple.V1) == 0 {
		resp = nil
	} else {
		resp = &tuple.V1[0]
	}

	// logger.Log.WithFields(logrus.Fields{"resp": fmt.Sprintf("%+v", resp)}).Debug("ConvertDBJSONToStorelocation")

	return resp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a store location JSON.
func ConvertDBJSONToStorelocationJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	// logger.Log.Debug("ConvertDBJSONToStorelocationJSON")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToStorelocationJSON")

	var (
		storeLocation *models.StoreLocation
		err           error
	)

	if storeLocation, err = ConvertDBJSONToStorelocation(jsonRawMessage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	// logger.Log.WithFields(logrus.Fields{"storeLocation": fmt.Sprintf("%+v", storeLocation)}).Debug("ConvertDBJSONToStorelocationJSON")

	var jsonresp []byte
	if jsonresp, err = json.Marshal(storeLocation); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling storelocation",
		}
	}

	return jsonresp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a storage.
func ConvertDBJSONToStorage(jsonRawMessage json.RawMessage) (*models.Storage, error) {

	// logger.Log.Debug("ConvertDBJSONToStorage")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToStorage")

	var (
		tuple tuple.T2[[]models.Storage, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	var resp *models.Storage
	if len(tuple.V1) == 0 {
		resp = nil
	} else {
		resp = &tuple.V1[0]
	}

	// logger.Log.WithFields(logrus.Fields{"resp": fmt.Sprintf("%+v", resp)}).Debug("ConvertDBJSONToStorage")

	return resp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a Storage JSON.
func ConvertDBJSONToStorageJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	// logger.Log.Debug("ConvertDBJSONToStorageJSON")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToStorageJSON")

	var (
		storage *models.Storage
		err     error
	)

	if storage, err = ConvertDBJSONToStorage(jsonRawMessage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	// logger.Log.WithFields(logrus.Fields{"storage": fmt.Sprintf("%+v", storage)}).Debug("ConvertDBJSONToStorageJSON")

	var jsonresp []byte
	if jsonresp, err = json.Marshal(storage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling storage",
		}
	}

	return jsonresp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a product.
func ConvertDBJSONToProduct(jsonRawMessage json.RawMessage) (*models.Product, error) {

	// logger.Log.Debug("ConvertDBJSONToProduct")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToProduct")

	var (
		tuple tuple.T2[[]models.Product, int]
		err   error
	)

	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	var resp *models.Product
	if len(tuple.V1) == 0 {
		resp = nil
	} else {
		resp = &tuple.V1[0]
	}

	// logger.Log.WithFields(logrus.Fields{"resp": fmt.Sprintf("%+v", resp)}).Debug("ConvertDBJSONToProduct")

	return resp, nil
}

// Convert a JSON response from the chimitheque_db Rust library to a Product JSON.
func ConvertDBJSONToProductJSON(jsonRawMessage json.RawMessage) ([]byte, *models.AppError) {

	// logger.Log.Debug("ConvertDBJSONToProductJSON")
	// logger.Log.WithFields(logrus.Fields{"jsonRawMessage": fmt.Sprintf("%+v", jsonRawMessage)}).Debug("ConvertDBJSONToProductJSON")

	var (
		product *models.Product
		err     error
	)

	if product, err = ConvertDBJSONToProduct(jsonRawMessage); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	// logger.Log.WithFields(logrus.Fields{"product": fmt.Sprintf("%+v", product)}).Debug("ConvertDBJSONToProductJSON")

	var jsonresp []byte
	if jsonresp, err = json.Marshal(product); err != nil {
		return nil, &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error marshalling product",
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
