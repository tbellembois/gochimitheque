package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/barweiss/go-tuple"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

func sendResponse(w http.ResponseWriter, response string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Log.Errorf("sendResponse error: %s", err.Error())
	}
}

// ValidatePersonEmailHandler checks that the person email does not already exist.
// If an id is given is the request, the validator ignore the email of the person with this id.
func (env *Env) ValidatePersonEmailHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err            error
		res            bool
		resp           string
		count          int
		person         models.Person
		personID       int
		personEmail    string
		jsonRawMessage json.RawMessage
	)

	vars := mux.Vars(r)

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	// converting the id
	if personID, err = strconv.Atoi(vars["id"]); err != nil {
		logger.Log.Error("strconv error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	// getting the email
	personEmail = vars["email"]

	// getting the people matching the email
	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/?search="+personEmail, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetPeople",
		}
	}

	// unmarshalling response
	var tuple tuple.T2[[]models.Person, int]
	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	people := tuple.V1
	count = tuple.V2

	if count == 0 {
		res = false
	} else if personID == -1 {
		res = (count == 1)
	} else {
		// getting the person
		if person, err = env.DB.GetPerson(personID); err != nil {
			logger.Log.Error("GetPerson error")
			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
			sendResponse(w, resp)
			return nil
		}
		res = (person.PersonID != people[0].PersonID)
	}

	logger.Log.WithFields(logrus.Fields{"vars": vars, "res": res}).Debug("ValidatePersonEmailHandler")
	if res {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
	} else {
		resp = "true"
	}

	sendResponse(w, resp)
	return nil
}

// ValidateEntityNameHandler checks that the entity name does not already exist
// if an id != 0 is given is the request the validator ignore the name of the entity with this id.
func (env *Env) ValidateEntityNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		err      error
		res      bool
		resp     string
		entity   models.Entity
		entityID int
		filter   zmqclient.RequestFilter
	)

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	// init db request parameters
	// if filter, aerr = request.NewFilter(r); aerr != nil {
	// 	logger.Log.Error("NewdbselectparamEntity error")
	// 	resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
	// 	sendResponse(w, resp)
	// 	return nil
	// }
	if filter, err = zmqclient.RequestFilterFromRawString("http://localhost/?" + r.URL.RawQuery); err != nil {
		logger.Log.Error("error calling zmqclient.Request_filter")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	// converting the id
	if entityID, err = strconv.Atoi(vars["id"]); err != nil {
		logger.Log.Error("strconv error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	// getting the name
	if err = r.ParseForm(); err != nil {
		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	filter.Search = r.Form.Get("entity_name")

	// getting the entities matching the name
	entities, count, err := env.DB.GetEntities(filter, c.PersonID)
	if err != nil {
		logger.Log.Error("GetEntities error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	if count == 0 {
		res = false
	} else if entityID == -1 {
		res = (count == 1)
	} else {
		// getting the entity
		if entity, err = env.DB.GetEntity(entityID); err != nil {
			logger.Log.Error("GetEntity error")
			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
			sendResponse(w, resp)
			return nil
		}
		res = (entity.EntityID != entities[0].EntityID)
	}

	logger.Log.WithFields(logrus.Fields{"vars": vars, "res": res}).Debug("ValidateEntityNameHandler")
	if res {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
	} else {
		resp = "true"
	}

	sendResponse(w, resp)
	return nil
}

// ValidateProductNameHandler checks that the product name is valid
// FIXME: not used yet.
func (env *Env) ValidateProductNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	sendResponse(w, "true")
	return nil
}

// FormatProductEmpiricalFormulaHandler returns the sorted formula.
func (env *Env) FormatProductEmpiricalFormulaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	type EmpiricalFormulaData struct {
		EmpiricalFormula string `json:"empirical_formula"`
	}

	var (
		resp                 string
		empiricalFormulaData EmpiricalFormulaData
		err                  error
	)

	if err = json.NewDecoder(r.Body).Decode(&empiricalFormulaData); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	// validating it
	resp, err = zmqclient.EmpiricalFormulaFromRawString(empiricalFormulaData.EmpiricalFormula)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// ValidateProductEmpiricalFormulaHandler checks that the product empirical formula is valid.
func (env *Env) ValidateProductEmpiricalFormulaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err  error
		resp string
	)

	// getting the empirical formula
	if err = r.ParseForm(); err != nil {
		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "empirical_formula_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	// validating it
	_, err = zmqclient.EmpiricalFormulaFromRawString(r.Form.Get("empirical_formula"))
	if err != nil {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "empirical_formula_validate", PluralCount: 1})
	} else {
		resp = "true"
	}

	sendResponse(w, resp)
	return nil
}

// ValidateProductCasNumberHandler checks that:
// - the cas number is valid
// - a product with the cas number and specificity does not already exist.
func (env *Env) ValidateProductCasNumberHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	// vars := mux.Vars(r)

	var (
		err  error
		resp string
		// cas        models.CasNumber
		// nbProducts int
		// products   []models.Product
		// filter     zmqclient.RequestFilter
		// productID  int
	)

	// c := request.ContainerFromRequestContext(r)

	// getting the cas number
	if err = r.ParseForm(); err != nil {
		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_wrongcas", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	logger.Log.WithFields(logrus.Fields{"cas_number": r.Form.Get("cas_number")}).Debug("ValidateProductCasNumberHandler")

	// validating it
	v, _ := zmqclient.IsCasNumber(r.Form.Get("cas_number"))

	if v {
		resp = "true"
	} else {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_wrongcas", PluralCount: 1})
	}

	// converting the id
	// if productID, err = strconv.Atoi(vars["id"]); err != nil {
	// 	logger.Log.Error("strconv error")
	// 	resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_wrongcas", PluralCount: 1})
	// 	sendResponse(w, resp)
	// 	return nil
	// }

	// // check pair cas/specificity only on create
	// if productID == -1 {
	// 	// get cas number id
	// 	// if cas, err = env.DB.GetCasNumberByLabel(r.Form.Get("cas_number")); err != nil && err != sql.ErrNoRows {
	// 	if cas, err = datastores.GetByText(models.CasNumber{}, env.DB.GetDB(), r.Form.Get("cas_number")); err != nil && err != sql.ErrNoRows {
	// 		logger.Log.Error("GetProductsCasNumberByLabel error:" + err.Error())
	// 		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_wrongcas", PluralCount: 1})
	// 		sendResponse(w, resp)
	// 		return nil
	// 	}
	//
	// 	if err != sql.ErrNoRows {
	// 		// init db request parameters
	// 		// if filter, aerr = request.NewFilter(r); aerr != nil {
	// 		// 	logger.Log.Error("NewdbselectparamProduct error")
	// 		// 	resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_wrongcas", PluralCount: 1})
	// 		// 	sendResponse(w, resp)
	// 		// 	return nil
	// 		// }
	// 		if filter, err = zmqclient.RequestFilterFromRawString("http://localhost/?" + r.URL.RawQuery); err != nil {
	// 			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_wrongcas", PluralCount: 1})
	// 			sendResponse(w, resp)
	// 		}
	//
	// 		// if cas.CasNumberID.Valid {
	// 		// 	filter.CasNumber = int(cas.CasNumberID.Int64)
	// 		// }
	// 		if cas.CasNumberID != nil {
	// 			// filter.CasNumber = int(cas.CasNumberID.Int64)
	// 			filter.CasNumber = int(*cas.CasNumberID)
	// 		}
	//
	// 		filter.ProductSpecificity = r.Form.Get("product_specificity")
	//
	// 		// getting the products matching the cas and specificity
	// 		if products, nbProducts, err = env.DB.GetProducts(filter, c.PersonID, false); err != nil && err != sql.ErrNoRows {
	// 			logger.Log.Error("GetProducts error:" + err.Error())
	// 			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_wrongcas", PluralCount: 1})
	// 			sendResponse(w, resp)
	// 			return nil
	// 		}
	//
	// 		logger.Log.Debug(nbProducts)
	//
	// 		if nbProducts != 0 {
	// 			if filter.ProductSpecificity == "" && products[0].ProductSpecificity != nil && *products[0].ProductSpecificity != "" {
	//
	// 			} else {
	// 				resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cas_number_validate_casspecificity", PluralCount: 1})
	// 			}
	// 		}
	// 	}
	// }

	sendResponse(w, resp)
	return nil
}

// ValidateProductCeNumberHandler checks that:
// - the ce number is valid.
func (env *Env) ValidateProductCeNumberHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err  error
		resp string
	)

	// getting the ce number
	if err = r.ParseForm(); err != nil {
		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "ce_number_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil
	}

	logger.Log.WithFields(logrus.Fields{"ce_number": r.Form.Get("ce_number")}).Debug("ValidateProductCeNumberHandler")

	// validating it
	v, _ := zmqclient.IsCeNumber(r.Form.Get("ce_number"))

	if v {
		resp = "true"
	} else {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "ce_number_validate", PluralCount: 1})
	}

	sendResponse(w, resp)
	return nil
}
