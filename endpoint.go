package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func buildEndpoints() (router *mux.Router) {
	router = mux.NewRouter()

	commonChain := alice.New(env.HeadersMiddleware, env.ContextMiddleware, env.LogingMiddleware)

	router.Handle("/", commonChain.Then(env.AppMiddleware(env.HomeHandler))).Methods("GET")
	router.Handle("/login", commonChain.Then(env.AppMiddleware(env.VLoginHandler))).Methods("GET")
	router.Handle("/menu", commonChain.Then(env.AppMiddleware(env.VMenuHandler))).Methods("GET")
	router.Handle("/search", commonChain.Then(env.AppMiddleware(env.VSearchHandler))).Methods("GET")
	router.Handle("/about", commonChain.Then(env.AppMiddleware(env.AboutHandler))).Methods("GET")

	//
	// entities
	//
	// views
	router.Handle("/v/{item:entities}", commonChain.Then(env.AppMiddleware(env.VGetEntitiesHandler))).Methods("GET")
	router.Handle("/vc/{item:entities}", commonChain.Then(env.AppMiddleware(env.VCreateEntityHandler))).Methods("GET")

	//
	// people
	//
	// views
	router.Handle("/v/{item:people}", commonChain.Then(env.AppMiddleware(env.VGetPeopleHandler))).Methods("GET")
	router.Handle("/vc/{item:people}", commonChain.Then(env.AppMiddleware(env.VCreatePersonHandler))).Methods("GET")

	//
	// store locations
	//
	// views
	router.Handle("/v/{item:store_locations}", commonChain.Then(env.AppMiddleware(env.VGetStoreLocationsHandler))).Methods("GET")
	router.Handle("/vc/{item:store_locations}", commonChain.Then(env.AppMiddleware(env.VCreateStoreLocationHandler))).Methods("GET")

	//
	// products
	//

	// views
	router.Handle("/v/{item:products}", commonChain.Then(env.AppMiddleware(env.VGetProductsHandler))).Methods("GET")
	router.Handle("/vc/{item:products}", commonChain.Then(env.AppMiddleware(env.VCreateProductHandler))).Methods("GET")
	router.Handle("/vc/{item:pubchem}", commonChain.Then(env.AppMiddleware(env.VPubchemHandler))).Methods("GET")
	//
	// storages
	//
	// views
	router.Handle("/v/{item:storages}", commonChain.Then(env.AppMiddleware(env.VGetStoragesHandler))).Methods("GET")
	router.Handle("/vc/{item:storages}", commonChain.Then(env.AppMiddleware(env.VCreateStorageHandler))).Methods("GET")

	return router
}
