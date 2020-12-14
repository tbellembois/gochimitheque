package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/tbellembois/gochimitheque/models"
	. "github.com/tbellembois/gochimitheque/test"
)

func TestMain(m *testing.M) {
	TestInit()
	r := m.Run()
	TestOut()
	os.Exit(r)
}

type GetEntitiesResp struct {
	Rows  []Entity `json:"rows"`
	Total int      `json:"total"`
}

type GetEntitiesReq struct {
	p                  Person
	expectedStatus     int
	expectedTotal      int
	expectedEntityName map[string]string
}

func GetEntities(t *testing.T, e GetEntitiesReq) {

	var (
		response *GetEntitiesResp
		rr       *httptest.ResponseRecorder
		req      *http.Request
		token    string
		err      error
	)

	// middleware chain
	securechain := alice.New(TEnv.ContextMiddleware, TEnv.HeadersMiddleware, TEnv.LogingMiddleware, TEnv.AuthenticateMiddleware, TEnv.AuthorizeMiddleware)

	// getting a JWT token
	if token, err = Authenticate(e.p.PersonEmail, e.p.PersonPassword); err != nil {
		t.Fatal("authenticate: " + err.Error())
	}

	// requests definition
	rr = httptest.NewRecorder()
	r := mux.NewRouter()
	r.Handle("/{item:entities}", securechain.Then(TEnv.AppMiddleware(TEnv.GetEntitiesHandler))).Methods("GET")
	if req, err = http.NewRequest("GET", "/entities", nil); err != nil {
		t.Fatalf("GET /entities: " + err.Error())
	}

	// adding the JWT token
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	// performing the request
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != e.expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, e.expectedStatus)
		return
	}

	// check the response body
	respbody := rr.Body.Bytes()
	response = &GetEntitiesResp{}
	if err = json.Unmarshal(respbody, response); err != nil {
		t.Errorf("json unmarshall error: %v", err)
	}
	if response.Total != e.expectedTotal {
		t.Errorf("wrong total: got %d want %d", response.Total, e.expectedStatus)
	}
	for _, ent := range response.Rows {
		if _, ok := e.expectedEntityName[ent.EntityName]; !ok {
			t.Errorf("wrong entity: not found %s", ent.EntityName)
		}
	}
}

func TestGetEntities(t *testing.T) {

	reqs := []GetEntitiesReq{
		{
			p:                  Admin,
			expectedEntityName: map[string]string{"sample entity": "", "e1": "", "e2": ""},
			expectedStatus:     http.StatusOK,
			expectedTotal:      3,
		},
		{
			p:                  Man1,
			expectedEntityName: map[string]string{"e1": ""},
			expectedStatus:     http.StatusOK,
			expectedTotal:      1,
		},
		{
			p:                  Mickey1,
			expectedEntityName: map[string]string{"e1": ""},
			expectedStatus:     http.StatusOK,
			expectedTotal:      1,
		},
	}

	for _, req := range reqs {
		GetEntities(t, req)
	}
}
