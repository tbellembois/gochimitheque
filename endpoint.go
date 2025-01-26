package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func buildEndpoints(fakeAuth bool) (router *mux.Router) {
	router = mux.NewRouter()

	var secureChain alice.Chain

	commonChain := alice.New(env.HeadersMiddleware, env.ContextMiddleware, env.LogingMiddleware)

	if fakeAuth {
		secureChain = alice.New(env.HeadersMiddleware, env.ContextMiddleware, env.LogingMiddleware, env.FakeMiddleware, env.AuthorizeMiddleware)
	} else {
		secureChain = alice.New(env.HeadersMiddleware, env.ContextMiddleware, env.LogingMiddleware, env.AuthenticateMiddleware, env.AuthorizeMiddleware)
	}

	router.Handle("/login", commonChain.Then(env.AppMiddleware(env.VLoginHandler))).Methods("GET")
	router.Handle("/menu", commonChain.Then(env.AppMiddleware(env.VMenuHandler))).Methods("GET")
	router.Handle("/search", commonChain.Then(env.AppMiddleware(env.VSearchHandler))).Methods("GET")
	router.Handle("/get-token", commonChain.Then(env.AppMiddleware(env.GetTokenHandler))).Methods("GET")
	router.Handle("/callback", commonChain.Then(env.AppMiddleware(env.CallbackHandler))).Methods("GET")
	router.Handle("/delete-token", commonChain.Then(env.AppMiddleware(env.DeleteTokenHandler))).Methods("GET")
	router.Handle("/about", commonChain.Then(env.AppMiddleware(env.AboutHandler))).Methods("GET")
	router.Handle("/{item:userinfo}", secureChain.Then(env.AppMiddleware(env.UserInfoHandler))).Methods("GET")

	// products public
	// if *paramPublicProductsEndpoint {
	// 	router.Handle("/e/{item:products}", commonChain.Then(env.AppMiddleware(env.GetExposedProductsHandler))).Methods("GET")
	// }

	// ping handler returns pong only if the request is correctly authenticated (ie. with a valid JWT token)
	router.Handle("/{item:ping}", secureChain.Then(env.AppMiddleware(env.VPingHandler))).Methods("GET")

	// home page
	router.Handle("/", commonChain.Then(env.AppMiddleware(env.HomeHandler))).Methods("GET")

	// welcome announce
	router.Handle("/{view:v}/{item:welcomeannounce}", secureChain.Then(env.AppMiddleware(env.VWelcomeAnnounceHandler))).Methods("GET")
	router.Handle("/{item:welcomeannounce}", secureChain.Then(env.AppMiddleware(env.UpdateWelcomeAnnounceHandler))).Methods("PUT")
	router.Handle("/{item:welcomeannounce}", commonChain.Then(env.AppMiddleware(env.GetWelcomeAnnounceHandler))).Methods("GET")

	// entities
	router.Handle("/{view:v}/{item:entities}", secureChain.Then(env.AppMiddleware(env.VGetEntitiesHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:entities}", secureChain.Then(env.AppMiddleware(env.VCreateEntityHandler))).Methods("GET")
	router.Handle("/{item:entities}", secureChain.Then(env.AppMiddleware(env.GetEntitiesHandler))).Methods("GET")
	router.Handle("/{item:entities}/{id}", secureChain.Then(env.AppMiddleware(env.GetEntityHandler))).Methods("GET")
	router.Handle("/{item:entities}/{id}/people", secureChain.Then(env.AppMiddleware(env.GetEntityPeopleHandler))).Methods("GET")
	router.Handle("/{item:entities}", secureChain.Then(env.AppMiddleware(env.CreateEntityHandler))).Methods("POST")
	router.Handle("/{item:entities}/{id}", secureChain.Then(env.AppMiddleware(env.UpdateEntityHandler))).Methods("PUT")
	router.Handle("/{item:entities}/{id}", secureChain.Then(env.AppMiddleware(env.DeleteEntityHandler))).Methods("DELETE")
	router.Handle("/entities/{item:stocks}/{id}", secureChain.Then(env.AppMiddleware(env.GetEntityStockHandler))).Methods("GET")

	router.Handle("/f/{view:v}/{item:entities}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:entities}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}/{id}/people", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:entities}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:entities}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	router.Handle("/f/entities/{item:stocks}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	// people
	router.Handle("/{view:v}/{item:people}", secureChain.Then(env.AppMiddleware(env.VGetPeopleHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:people}", secureChain.Then(env.AppMiddleware(env.VCreatePersonHandler))).Methods("GET")
	router.Handle("/{item:people}", secureChain.Then(env.AppMiddleware(env.GetPeopleHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}", secureChain.Then(env.AppMiddleware(env.GetPersonHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}/entities", secureChain.Then(env.AppMiddleware(env.GetPersonEntitiesHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}/manageentities", secureChain.Then(env.AppMiddleware(env.GetPersonManageEntitiesHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}/permissions", secureChain.Then(env.AppMiddleware(env.GetPersonPermissionsHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}", secureChain.Then(env.AppMiddleware(env.UpdatePersonHandler))).Methods("PUT")
	// router.Handle("/{item:people}", secureChain.Then(env.AppMiddleware(env.CreatePersonHandler))).Methods("POST")
	router.Handle("/{item:people}/{id}", secureChain.Then(env.AppMiddleware(env.DeletePersonHandler))).Methods("DELETE")
	router.Handle("/f/{view:v}/{item:people}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:people}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}/entities", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}/manageentities", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}/permissions", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:people}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:people}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")

	// store locations
	router.Handle("/{view:v}/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.VGetStoreLocationsHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.VCreateStoreLocationHandler))).Methods("GET")

	router.Handle("/{item:storelocations}?store_location={id:[0-9]+}", secureChain.Then(env.AppMiddleware(env.GetStoreLocationsHandler))).Methods("GET")
	router.Handle("/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.GetStoreLocationsHandler))).Methods("GET")
	router.Handle("/{item:storelocations}/{id}", secureChain.Then(env.AppMiddleware(env.UpdateStoreLocationHandler))).Methods("PUT")
	router.Handle("/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.CreateStoreLocationHandler))).Methods("POST")
	router.Handle("/{item:storelocations}/{id}", secureChain.Then(env.AppMiddleware(env.DeleteStoreLocationHandler))).Methods("DELETE")

	router.Handle("/f/{view:v}/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storelocations}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storelocations}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:storelocations}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:storelocations}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")

	// products
	router.Handle("/{item:products}/l2eformula/{f}", secureChain.Then(env.AppMiddleware(env.ConvertProductEmpiricalToLinearFormulaHandler))).Methods("GET")
	router.Handle("/{view:v}/{item:products}", secureChain.Then(env.AppMiddleware(env.VGetProductsHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:products}", secureChain.Then(env.AppMiddleware(env.VCreateProductHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:pubchem}", secureChain.Then(env.AppMiddleware(env.VPubchemHandler))).Methods("GET")
	router.Handle("/{item:products}", secureChain.Then(env.AppMiddleware(env.GetProductsHandler))).Methods("GET")
	router.Handle("/{item:products}/{id}", secureChain.Then(env.AppMiddleware(env.GetProductHandler))).Methods("GET")
	router.Handle("/{item:products}/{id}", secureChain.Then(env.AppMiddleware(env.UpdateProductHandler))).Methods("PUT")
	router.Handle("/{item:products}", secureChain.Then(env.AppMiddleware(env.CreateProductHandler))).Methods("POST")
	router.Handle("/{item:products}/{id}", secureChain.Then(env.AppMiddleware(env.DeleteProductHandler))).Methods("DELETE")
	router.Handle("/{item:bookmarks}/{id}", secureChain.Then(env.AppMiddleware(env.ToogleProductBookmarkHandler))).Methods("PUT")

	router.Handle("/{item:products}/casnumbers/", secureChain.Then(env.AppMiddleware(env.GetProductsCasNumbersHandler))).Methods("GET")
	router.Handle("/{item:products}/casnumbers/{id}", secureChain.Then(env.AppMiddleware(env.GetProductsCasNumberHandler))).Methods("GET")

	router.Handle("/{item:products}/cenumbers/", secureChain.Then(env.AppMiddleware(env.GetProductsCeNumbersHandler))).Methods("GET")

	router.Handle("/{item:products}/names/", secureChain.Then(env.AppMiddleware(env.GetProductsNamesHandler))).Methods("GET")
	router.Handle("/{item:products}/names/{id}", secureChain.Then(env.AppMiddleware(env.GetProductsNameHandler))).Methods("GET")

	router.Handle("/{item:products}/linearformulas/", secureChain.Then(env.AppMiddleware(env.GetProductsLinearFormulasHandler))).Methods("GET")

	router.Handle("/{item:products}/empiricalformulas/", secureChain.Then(env.AppMiddleware(env.GetProductsEmpiricalFormulasHandler))).Methods("GET")
	router.Handle("/{item:products}/empiricalformulas/{id}", secureChain.Then(env.AppMiddleware(env.GetProductsEmpiricalFormulaHandler))).Methods("GET")

	router.Handle("/{item:products}/physicalstates/", secureChain.Then(env.AppMiddleware(env.GetProductsPhysicalStatesHandler))).Methods("GET")

	router.Handle("/{item:products}/signalwords/", secureChain.Then(env.AppMiddleware(env.GetProductsSignalWordsHandler))).Methods("GET")
	router.Handle("/{item:products}/signalwords/{id}", secureChain.Then(env.AppMiddleware(env.GetProductsSignalWordHandler))).Methods("GET")

	router.Handle("/{item:products}/synonyms/", secureChain.Then(env.AppMiddleware(env.GetProductsSynonymsHandler))).Methods("GET")

	router.Handle("/{item:products}/symbols/", secureChain.Then(env.AppMiddleware(env.GetProductsSymbolsHandler))).Methods("GET")
	router.Handle("/{item:products}/symbols/{id}", secureChain.Then(env.AppMiddleware(env.GetProductsSymbolHandler))).Methods("GET")

	router.Handle("/{item:products}/classofcompounds/", secureChain.Then(env.AppMiddleware(env.GetProductsClassOfCompoundsHandler))).Methods("GET")

	router.Handle("/{item:products}/hazardstatements/", secureChain.Then(env.AppMiddleware(env.GetProductsHazardStatementsHandler))).Methods("GET")
	router.Handle("/{item:products}/hazardstatements/{id}", secureChain.Then(env.AppMiddleware(env.GetProductsHazardStatementHandler))).Methods("GET")

	router.Handle("/{item:products}/precautionarystatements/", secureChain.Then(env.AppMiddleware(env.GetProductsPrecautionaryStatementsHandler))).Methods("GET")
	router.Handle("/{item:products}/precautionarystatements/{id}", secureChain.Then(env.AppMiddleware(env.GetProductsPrecautionaryStatementHandler))).Methods("GET")

	router.Handle("/{item:products}/producerrefs/", secureChain.Then(env.AppMiddleware(env.GetProductsProducerRefsHandler))).Methods("GET")
	router.Handle("/{item:products}/producers/", secureChain.Then(env.AppMiddleware(env.GetProductsProducersHandler))).Methods("GET")
	router.Handle("/{item:products}/supplierrefs/", secureChain.Then(env.AppMiddleware(env.GetProductsSupplierRefsHandler))).Methods("GET")
	router.Handle("/{item:products}/suppliers/", secureChain.Then(env.AppMiddleware(env.GetProductsSuppliersHandler))).Methods("GET")
	router.Handle("/{item:products}/categories/", secureChain.Then(env.AppMiddleware(env.GetProductsCategoriesHandler))).Methods("GET")
	router.Handle("/{item:products}/tags/", secureChain.Then(env.AppMiddleware(env.GetProductsTagsHandler))).Methods("GET")

	router.Handle("/{item:products}/producers", secureChain.Then(env.AppMiddleware(env.CreateProducerHandler))).Methods("POST")
	router.Handle("/{item:products}/suppliers", secureChain.Then(env.AppMiddleware(env.CreateSupplierHandler))).Methods("POST")

	router.Handle("/{item:products}/pubchemautocomplete/{name}", secureChain.Then(env.AppMiddleware(env.PubchemAutocompleteHandler))).Methods("GET")
	router.Handle("/{item:products}/pubchemgetcompoundbyname/{name}", secureChain.Then(env.AppMiddleware(env.PubchemGetCompoundByNameHandler))).Methods("GET")
	router.Handle("/{item:products}/pubchemgetproductbyname/{name}", secureChain.Then(env.AppMiddleware(env.PubchemGetProductByNameHandler))).Methods("GET")
	router.Handle("/{item:products}/pubchemcreateproduct", secureChain.Then(env.AppMiddleware(env.CreateUpdateProductFromPubchemHandler))).Methods("POST")
	router.Handle("/{item:products}/pubchemcreateproduct/{id}", secureChain.Then(env.AppMiddleware(env.CreateUpdateProductFromPubchemHandler))).Methods("POST")

	router.Handle("/f/{view:v}/{item:products}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:products}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:products}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:products}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:products}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:products}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:products}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")

	// storages
	router.Handle("/{view:v}/{item:storages}", secureChain.Then(env.AppMiddleware(env.VGetStoragesHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:storages}", secureChain.Then(env.AppMiddleware(env.VCreateStorageHandler))).Methods("GET")
	router.Handle("/{item:storages}", secureChain.Then(env.AppMiddleware(env.GetStoragesHandler))).Methods("GET")
	router.Handle("/{item:storages}/others", secureChain.Then(env.AppMiddleware(env.GetOtherStoragesHandler))).Methods("GET")
	// router.Handle("/{item:storages}/suppliers", secureChain.Then(env.AppMiddleware(env.GetStoragesSuppliersHandler))).Methods("GET")
	router.Handle("/{item:storages}/units", secureChain.Then(env.AppMiddleware(env.GetStoragesUnitsHandler))).Methods("GET")
	router.Handle("/{item:storages}/{id}", secureChain.Then(env.AppMiddleware(env.GetStorageHandler))).Methods("GET")
	router.Handle("/{item:storages}/{id}", secureChain.Then(env.AppMiddleware(env.UpdateStorageHandler))).Methods("PUT")
	router.Handle("/{item:storages}", secureChain.Then(env.AppMiddleware(env.CreateStorageHandler))).Methods("POST")
	router.Handle("/{item:storages}/{id}", secureChain.Then(env.AppMiddleware(env.DeleteStorageHandler))).Methods("DELETE")
	router.Handle("/{item:storages}/{id}/a", secureChain.Then(env.AppMiddleware(env.ArchiveStorageHandler))).Methods("DELETE")
	router.Handle("/{item:storages}/{id}/r", secureChain.Then(env.AppMiddleware(env.RestoreStorageHandler))).Methods("PUT")
	router.Handle("/{item:borrowings}", secureChain.Then(env.AppMiddleware(env.ToogleStorageBorrowingHandler))).Methods("PUT")

	router.Handle("/f/{item:storages}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storages}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:storages}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:storages}/{id}", secureChain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")

	// validators
	router.Handle("/{item:validate}/entity/{id}/name/", secureChain.Then(env.AppMiddleware(env.ValidateEntityNameHandler))).Methods("POST")
	router.Handle("/{item:validate}/person/{id}/email/{email}", secureChain.Then(env.AppMiddleware(env.ValidatePersonEmailHandler))).Methods("GET")
	router.Handle("/{item:validate}/product/{id}/cas_number/", secureChain.Then(env.AppMiddleware(env.ValidateProductCasNumberHandler))).Methods("POST")
	router.Handle("/{item:validate}/product/{id}/cenumber/", secureChain.Then(env.AppMiddleware(env.ValidateProductCeNumberHandler))).Methods("POST")
	router.Handle("/{item:validate}/product/{id}/name/", secureChain.Then(env.AppMiddleware(env.ValidateProductNameHandler))).Methods("POST")
	router.Handle("/{item:validate}/product/{id}/empiricalformula/", secureChain.Then(env.AppMiddleware(env.ValidateProductEmpiricalFormulaHandler))).Methods("POST")

	// formatters
	router.Handle("/{item:format}/product/empiricalformula/", secureChain.Then(env.AppMiddleware(env.FormatProductEmpiricalFormulaHandler))).Methods("POST")

	// export download
	router.Handle("/{item:download}/{id}", secureChain.Then(env.AppMiddleware(env.DownloadExportHandler))).Methods("GET")

	return router
}
