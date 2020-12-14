package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/globals"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/utils"
)

// IsCeNumber returns true if c is a valid ce number
func IsCeNumber(c string) bool {

	var (
		err                error
		checkdigit, checkd int
	)

	if c == "000-000-0" {
		return true
	}

	// compiling regex
	r := regexp.MustCompile("^(?P<groupone>[0-9]{3})-(?P<grouptwo>[0-9]{3})-(?P<groupthree>[0-9]{1})$")
	// finding group names
	n := r.SubexpNames()
	// finding matches
	ms := r.FindAllStringSubmatch(c, -1)
	if len(ms) == 0 {
		return false
	}
	m := ms[0]
	// then building a map of matches
	md := map[string]string{}
	for i, j := range m {
		md[n[i]] = j
	}

	if len(m) > 0 {
		numberpart := md["groupone"] + md["grouptwo"]

		// converting the check digit into int
		if checkdigit, err = strconv.Atoi(string(md["groupthree"])); err != nil {
			return false
		}

		// calculating the check digit
		counter := 1  // loop counter
		currentd := 0 // current processed digit in c

		for i := 0; i < len(numberpart); i++ {
			// converting digit into int
			if currentd, err = strconv.Atoi(string(numberpart[i])); err != nil {
				return false
			}
			checkd += counter * currentd
			counter++
			//fmt.Printf("counter: %d currentd: %d checkd: %d\n", counter, currentd, checkd)
		}
	}

	return checkd%11 == checkdigit
}

// IsCasNumber returns true if c is a valid cas number
func IsCasNumber(c string) bool {

	var (
		err                error
		checkdigit, checkd int
	)

	if c == "0000-00-0" {
		return true
	}

	// compiling regex
	r := regexp.MustCompile("^(?P<groupone>[0-9]{1,7})-(?P<grouptwo>[0-9]{2})-(?P<groupthree>[0-9]{1})$")
	// finding group names
	n := r.SubexpNames()
	// finding matches
	ms := r.FindAllStringSubmatch(c, -1)
	if len(ms) == 0 {
		return false
	}
	m := ms[0]
	// then building a map of matches
	md := map[string]string{}
	for i, j := range m {
		md[n[i]] = j
	}

	if len(m) > 0 {
		numberpart := md["groupone"] + md["grouptwo"]

		// converting the check digit into int
		if checkdigit, err = strconv.Atoi(string(md["groupthree"])); err != nil {
			return false
		}
		//fmt.Printf("checkdigit: %d\n", checkdigit)

		// calculating the check digit
		counter := 1  // loop counter
		currentd := 0 // current processed digit in c

		for i := len(numberpart) - 1; i >= 0; i-- {
			// converting digit into int
			if currentd, err = strconv.Atoi(string(numberpart[i])); err != nil {
				return false
			}
			checkd += counter * currentd
			counter++
			//fmt.Printf("counter: %d currentd: %d checkd: %d\n", counter, currentd, checkd)
		}
	}
	return checkd%10 == checkdigit
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
	if dspp, aerr = models.NewdbselectparamPerson(r, nil); err != nil {
		return aerr
	}
	dspp.SetLoggedPersonID(c.PersonID)

	// converting the id
	if person_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// getting the email
	if err = r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	dspp.SetSearch(r.Form.Get("person_email"))

	// getting the people matching the email
	people, count, err := env.DB.GetPeople(dspp)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the people",
		}
	}

	if count == 0 {
		res = false
	} else if person_id == -1 {
		res = (count == 1)
	} else {
		// getting the person
		if person, err = env.DB.GetPerson(person_id); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for person by email",
			}
		}
		res = (person.PersonID != people[0].PersonID)
	}

	globals.Log.WithFields(logrus.Fields{"vars": vars, "res": res}).Debug("ValidatePersonEmailHandler")
	if res {
		resp = globals.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "person_emailexist_validate", PluralCount: 1})
	} else {
		resp = "true"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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
	if dspe, aerr = models.NewdbselectparamEntity(r, nil); err != nil {
		return aerr
	}
	dspe.SetLoggedPersonID(c.PersonID)

	// converting the id
	if entity_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// getting the name
	if err = r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	dspe.SetSearch(r.Form.Get("entity_name"))

	// getting the entities matching the name
	entities, count, err := env.DB.GetEntities(dspe)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entities",
		}
	}

	if count == 0 {
		res = false
	} else if entity_id == -1 {
		res = (count == 1)
	} else {
		// getting the entity
		if entity, err = env.DB.GetEntity(entity_id); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for entity by id",
			}
		}
		res = (entity.EntityID != entities[0].EntityID)
	}

	globals.Log.WithFields(logrus.Fields{"vars": vars, "res": res}).Debug("ValidateEntityNameHandler")
	if res {
		resp = globals.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "entity_nameexist_validate", PluralCount: 1})
	} else {
		resp = "true"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// ValidateProductNameHandler checks that the product name is valid
// FIXME: not used yet
func (env *Env) ValidateProductNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	resp := "true"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	// validating it
	_, err = utils.SortEmpiricalFormula(r.Form.Get("empiricalformula"))

	if err != nil {
		resp = globals.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "empiricalformula_validate", PluralCount: 1})
	} else {
		resp = "true"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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
	resp, err = utils.SortEmpiricalFormula(r.Form.Get("empiricalformula"))

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
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	// validating it
	v := IsCasNumber(r.Form.Get("casnumber"))

	if v {
		resp = "true"
	} else {
		resp = globals.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_wrongcas", PluralCount: 1})
	}

	// returning true on 0000-00-0 cas (ie. no cas number)
	// if r.Form.Get("casnumber") == "0000-00-0" {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(http.StatusOK)
	// 	if err = json.NewEncoder(w).Encode(resp); err != nil {
	// 		return &models.AppError{
	// 			Code:    http.StatusInternalServerError,
	// 			Message: err.Error(),
	// 		}
	// 	}
	// 	return nil
	// }

	// converting the id
	if product_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// check pair cas/specificity only on create
	if product_id == -1 {
		// get cas number id
		if cas, err = env.DB.GetProductsCasNumberByLabel(r.Form.Get("casnumber")); err != nil {
			return &models.AppError{
				Error:   err,
				Message: "get cas number",
				Code:    http.StatusInternalServerError}
		}
		// init db request parameters
		if dspp, aerr = models.NewdbselectparamProduct(r, nil); err != nil {
			return aerr
		}
		if cas.CasNumberID.Valid {
			dspp.SetCasNumber(int(cas.CasNumberID.Int64))
		}
		dspp.SetProductSpecificity(r.Form.Get("product_specificity"))
		// getting the products matching the cas and specificity
		if _, nbProducts, err = env.DB.GetProducts(dspp); err != nil {
			return &models.AppError{
				Error:   err,
				Message: "get products",
				Code:    http.StatusInternalServerError}
		}

		if nbProducts != 0 {
			resp = globals.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "casnumber_validate_casspecificity", PluralCount: 1})
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	// validating it
	v := IsCeNumber(r.Form.Get("cenumber"))

	if v {
		resp = "true"
	} else {
		resp = globals.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cenumber_validate", PluralCount: 1})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}
