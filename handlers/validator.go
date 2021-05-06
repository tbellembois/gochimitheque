package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque-utils/sort"
	"github.com/tbellembois/gochimitheque-utils/validator"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

func sendResponse(w http.ResponseWriter, response string) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// TODO: check error
	_ = json.NewEncoder(w).Encode(response)

}

// ValidatePersonEmailHandler checks that the person email does not already exist
// if an id is given is the request the validator ignore the email of the person with this id
func (env *Env) ValidatePersonEmailHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	vars := mux.Vars(r)

	var (
		err       error
		aerr      *models.AppError
		res       bool
		resp      string
		person    models.Person
		person_id int
		dspp      models.DbselectparamPerson
	)

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)

	// init db request parameters
	if dspp, aerr = models.NewdbselectparamPerson(r, nil); aerr != nil {

		logger.Log.Error("NewdbselectparamPerson error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}
	dspp.SetLoggedPersonID(c.PersonID)

	// converting the id
	if person_id, err = strconv.Atoi(vars["id"]); err != nil {

		logger.Log.Error("strconv error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}

	// getting the email
	if err = r.ParseForm(); err != nil {

		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}
	dspp.SetSearch(r.Form.Get("person_email"))

	// getting the people matching the email
	people, count, err := env.DB.GetPeople(dspp)
	if err != nil {

		logger.Log.Error("GetPeople error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}

	if count == 0 {
		res = false
	} else if person_id == -1 {
		res = (count == 1)
	} else {
		// getting the person
		if person, err = env.DB.GetPerson(person_id); err != nil {

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
// if an id != -1 is given is the request the validator ignore the name of the entity with this id
func (env *Env) ValidateEntityNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	vars := mux.Vars(r)

	var (
		err       error
		aerr      *models.AppError
		res       bool
		resp      string
		entity    models.Entity
		entity_id int
		dspe      models.DbselectparamEntity
	)

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)

	// init db request parameters
	if dspe, aerr = models.NewdbselectparamEntity(r, nil); aerr != nil {

		logger.Log.Error("NewdbselectparamEntity error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}
	dspe.SetLoggedPersonID(c.PersonID)

	// converting the id
	if entity_id, err = strconv.Atoi(vars["id"]); err != nil {

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
	dspe.SetSearch(r.Form.Get("entity_name"))

	// getting the entities matching the name
	entities, count, err := env.DB.GetEntities(dspe)
	if err != nil {

		logger.Log.Error("GetEntities error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}

	if count == 0 {
		res = false
	} else if entity_id == -1 {
		res = (count == 1)
	} else {
		// getting the entity
		if entity, err = env.DB.GetEntity(entity_id); err != nil {

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
// FIXME: not used yet
func (env *Env) ValidateProductNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	sendResponse(w, "true")
	return nil

}

// ValidateProductEmpiricalFormulaHandler checks that the product empirical formula is valid
func (env *Env) ValidateProductEmpiricalFormulaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		err  error
		resp string
	)

	// getting the empirical formula
	if err = r.ParseForm(); err != nil {

		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "empiricalformula_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}

	// validating it
	_, err = sort.SortEmpiricalFormula(r.Form.Get("empiricalformula"))
	if err != nil {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "empiricalformula_validate", PluralCount: 1})
	} else {
		resp = "true"
	}

	sendResponse(w, resp)
	return nil

}

// FormatProductEmpiricalFormulaHandler returns the sorted formula
func (env *Env) FormatProductEmpiricalFormulaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		err  error
		resp string
	)

	// getting the empirical formula
	if err = r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	// validating it
	resp, err = sort.SortEmpiricalFormula(r.Form.Get("empiricalformula"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
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

// ValidateProductCasNumberHandler checks that:
// - the cas number is valid
// - a product with the cas number and specificity does not already exist
func (env *Env) ValidateProductCasNumberHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	vars := mux.Vars(r)

	var (
		err        error
		resp       string
		cas        models.CasNumber
		nbProducts int
		aerr       *models.AppError
		dspp       models.DbselectparamProduct
		product_id int
	)

	// getting the cas number
	if err = r.ParseForm(); err != nil {

		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_wrongcas", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}
	logger.Log.WithFields(logrus.Fields{"casnumber": r.Form.Get("casnumber")}).Debug("ValidateProductCasNumberHandler")

	// validating it
	v := validator.IsCasNumber(r.Form.Get("casnumber"))

	if v {
		resp = "true"
	} else {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_wrongcas", PluralCount: 1})
	}

	// converting the id
	if product_id, err = strconv.Atoi(vars["id"]); err != nil {

		logger.Log.Error("strconv error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_wrongcas", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}

	// check pair cas/specificity only on create
	if product_id == -1 {

		// get cas number id
		if cas, err = env.DB.GetProductsCasNumberByLabel(r.Form.Get("casnumber")); err != nil {

			logger.Log.Error("GetProductsCasNumberByLabel error")
			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_wrongcas", PluralCount: 1})
			sendResponse(w, resp)
			return nil

		}

		// init db request parameters
		if dspp, aerr = models.NewdbselectparamProduct(r, nil); aerr != nil {

			logger.Log.Error("NewdbselectparamProduct error")
			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_wrongcas", PluralCount: 1})
			sendResponse(w, resp)
			return nil

		}

		if cas.CasNumberID.Valid {
			dspp.SetCasNumber(int(cas.CasNumberID.Int64))
		}
		dspp.SetProductSpecificity(r.Form.Get("product_specificity"))

		// getting the products matching the cas and specificity
		if _, nbProducts, err = env.DB.GetProducts(dspp); err != nil {

			logger.Log.Error("GetProducts error")
			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_wrongcas", PluralCount: 1})
			sendResponse(w, resp)
			return nil

		}

		if nbProducts != 0 {
			resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_casspecificity", PluralCount: 1})
		}
	}

	sendResponse(w, resp)
	return nil

}

// ValidateProductCeNumberHandler checks that:
// - the ce number is valid
func (env *Env) ValidateProductCeNumberHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		err  error
		resp string
	)

	// getting the ce number
	if err = r.ParseForm(); err != nil {

		logger.Log.Error("ParseForm error")
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cenumber_validate", PluralCount: 1})
		sendResponse(w, resp)
		return nil

	}
	logger.Log.WithFields(logrus.Fields{"cenumber": r.Form.Get("cenumber")}).Debug("ValidateProductCeNumberHandler")

	// validating it
	v := validator.IsCeNumber(r.Form.Get("cenumber"))

	if v {
		resp = "true"
	} else {
		resp = locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cenumber_validate", PluralCount: 1})
	}

	sendResponse(w, resp)
	return nil

}
