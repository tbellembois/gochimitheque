package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func buildEndpoints() (router *mux.Router) {

	router = mux.NewRouter()

	commonChain := alice.New(env.ContextMiddleware, env.HeadersMiddleware, env.LogingMiddleware)
	securechain := alice.New(env.ContextMiddleware, env.HeadersMiddleware, env.LogingMiddleware, env.AuthenticateMiddleware, env.AuthorizeMiddleware)

	router.Handle("/login", commonChain.Then(env.AppMiddleware(env.VLoginHandler))).Methods("GET")
	router.Handle("/menu", commonChain.Then(env.AppMiddleware(env.VMenuHandler))).Methods("GET")
	router.Handle("/search", commonChain.Then(env.AppMiddleware(env.VSearchHandler))).Methods("GET")
	router.Handle("/get-token", commonChain.Then(env.AppMiddleware(env.GetTokenHandler))).Methods("POST")
	router.Handle("/reset-password", commonChain.Then(env.AppMiddleware(env.ResetPasswordHandler))).Methods("POST")
	router.Handle("/reset", commonChain.Then(env.AppMiddleware(env.ResetHandler))).Methods("GET")
	router.Handle("/captcha", commonChain.Then(env.AppMiddleware(env.CaptchaHandler))).Methods("GET")
	router.Handle("/delete-token", commonChain.Then(env.AppMiddleware(env.DeleteTokenHandler))).Methods("GET")
	router.Handle("/about", commonChain.Then(env.AppMiddleware(env.AboutHandler))).Methods("GET")

	// products public
	if *paramPublicProductsEndpoint {
		router.Handle("/e/{item:products}", commonChain.Then(env.AppMiddleware(env.GetExposedProductsHandler))).Methods("GET")
	}

	// developper tests
	router.Handle("/v/test", securechain.Then(env.AppMiddleware(env.VTestHandler))).Methods("GET")

	// home page
	router.Handle("/", commonChain.Then(env.AppMiddleware(env.HomeHandler))).Methods("GET")

	// welcome announce
	router.Handle("/{view:v}/{item:welcomeannounce}", securechain.Then(env.AppMiddleware(env.VWelcomeAnnounceHandler))).Methods("GET")
	router.Handle("/{item:welcomeannounce}", securechain.Then(env.AppMiddleware(env.UpdateWelcomeAnnounceHandler))).Methods("PUT")
	router.Handle("/{item:welcomeannounce}", commonChain.Then(env.AppMiddleware(env.GetWelcomeAnnounceHandler))).Methods("GET")

	// entities
	router.Handle("/{view:v}/{item:entities}", securechain.Then(env.AppMiddleware(env.VGetEntitiesHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:entities}", securechain.Then(env.AppMiddleware(env.VCreateEntityHandler))).Methods("GET")
	router.Handle("/{item:entities}", securechain.Then(env.AppMiddleware(env.GetEntitiesHandler))).Methods("GET")
	router.Handle("/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.GetEntityHandler))).Methods("GET")
	router.Handle("/{item:entities}/{id}/people", securechain.Then(env.AppMiddleware(env.GetEntityPeopleHandler))).Methods("GET")
	router.Handle("/{item:entities}", securechain.Then(env.AppMiddleware(env.CreateEntityHandler))).Methods("POST")
	router.Handle("/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.UpdateEntityHandler))).Methods("PUT")
	router.Handle("/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.DeleteEntityHandler))).Methods("DELETE")
	router.Handle("/entities/{item:stocks}/{id}", securechain.Then(env.AppMiddleware(env.GetEntityStockHandler))).Methods("GET")

	router.Handle("/f/{view:v}/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}/{id}/people", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	router.Handle("/f/entities/{item:stocks}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	// people
	router.Handle("/{view:v}/{item:people}", securechain.Then(env.AppMiddleware(env.VGetPeopleHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:people}", securechain.Then(env.AppMiddleware(env.VCreatePersonHandler))).Methods("GET")
	router.Handle("/{view:vu}/{item:peoplepass}", securechain.Then(env.AppMiddleware(env.VUpdatePersonPasswordHandler))).Methods("GET")
	router.Handle("/{item:people}", securechain.Then(env.AppMiddleware(env.GetPeopleHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.GetPersonHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}/entities", securechain.Then(env.AppMiddleware(env.GetPersonEntitiesHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}/manageentities", securechain.Then(env.AppMiddleware(env.GetPersonManageEntitiesHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}/permissions", securechain.Then(env.AppMiddleware(env.GetPersonPermissionsHandler))).Methods("GET")
	router.Handle("/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.UpdatePersonHandler))).Methods("PUT")
	router.Handle("/{item:people}", securechain.Then(env.AppMiddleware(env.CreatePersonHandler))).Methods("POST")
	router.Handle("/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.DeletePersonHandler))).Methods("DELETE")
	router.Handle("/{item:peoplep}", securechain.Then(env.AppMiddleware(env.UpdatePersonpHandler))).Methods("POST")

	router.Handle("/f/{view:v}/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}/entities", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}/manageentities", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}/permissions", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	// store locations
	router.Handle("/{view:v}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.VGetStoreLocationsHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.VCreateStoreLocationHandler))).Methods("GET")
	router.Handle("/{item:storelocations}", securechain.Then(env.AppMiddleware(env.GetStoreLocationsHandler))).Methods("GET")
	router.Handle("/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.GetStoreLocationHandler))).Methods("GET")
	router.Handle("/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.UpdateStoreLocationHandler))).Methods("PUT")
	router.Handle("/{item:storelocations}", securechain.Then(env.AppMiddleware(env.CreateStoreLocationHandler))).Methods("POST")
	router.Handle("/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.DeleteStoreLocationHandler))).Methods("DELETE")

	router.Handle("/f/{view:v}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")

	// TODO: merge with wasm code
	router.Handle("/{item:products}/l2eformula/{f}", securechain.Then(env.AppMiddleware(env.ConvertProductEmpiricalToLinearFormulaHandler))).Methods("GET")
	router.Handle("/{view:v}/{item:products}", securechain.Then(env.AppMiddleware(env.VGetProductsHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:products}", securechain.Then(env.AppMiddleware(env.VCreateProductHandler))).Methods("GET")
	router.Handle("/{item:products}", securechain.Then(env.AppMiddleware(env.GetProductsHandler))).Methods("GET")
	router.Handle("/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.GetProductHandler))).Methods("GET")
	router.Handle("/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.UpdateProductHandler))).Methods("PUT")
	router.Handle("/{item:products}", securechain.Then(env.AppMiddleware(env.CreateProductHandler))).Methods("POST")
	router.Handle("/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.DeleteProductHandler))).Methods("DELETE")
	router.Handle("/{item:bookmarks}/{id}", securechain.Then(env.AppMiddleware(env.ToogleProductBookmarkHandler))).Methods("PUT")

	router.Handle("/{item:products}/casnumbers/", securechain.Then(env.AppMiddleware(env.GetProductsCasNumbersHandler))).Methods("GET")
	router.Handle("/{item:products}/casnumbers/{id}", securechain.Then(env.AppMiddleware(env.GetProductsCasNumberHandler))).Methods("GET")

	router.Handle("/{item:products}/cenumbers/", securechain.Then(env.AppMiddleware(env.GetProductsCeNumbersHandler))).Methods("GET")

	router.Handle("/{item:products}/names/", securechain.Then(env.AppMiddleware(env.GetProductsNamesHandler))).Methods("GET")
	router.Handle("/{item:products}/names/{id}", securechain.Then(env.AppMiddleware(env.GetProductsNameHandler))).Methods("GET")

	router.Handle("/{item:products}/linearformulas/", securechain.Then(env.AppMiddleware(env.GetProductsLinearFormulasHandler))).Methods("GET")

	router.Handle("/{item:products}/empiricalformulas/", securechain.Then(env.AppMiddleware(env.GetProductsEmpiricalFormulasHandler))).Methods("GET")
	router.Handle("/{item:products}/empiricalformulas/{id}", securechain.Then(env.AppMiddleware(env.GetProductsEmpiricalFormulaHandler))).Methods("GET")

	router.Handle("/{item:products}/physicalstates/", securechain.Then(env.AppMiddleware(env.GetProductsPhysicalStatesHandler))).Methods("GET")

	router.Handle("/{item:products}/signalwords/", securechain.Then(env.AppMiddleware(env.GetProductsSignalWordsHandler))).Methods("GET")
	router.Handle("/{item:products}/signalwords/{id}", securechain.Then(env.AppMiddleware(env.GetProductsSignalWordHandler))).Methods("GET")

	router.Handle("/{item:products}/synonyms/", securechain.Then(env.AppMiddleware(env.GetProductsSynonymsHandler))).Methods("GET")

	router.Handle("/{item:products}/symbols/", securechain.Then(env.AppMiddleware(env.GetProductsSymbolsHandler))).Methods("GET")
	router.Handle("/{item:products}/symbols/{id}", securechain.Then(env.AppMiddleware(env.GetProductsSymbolHandler))).Methods("GET")

	router.Handle("/{item:products}/classofcompounds/", securechain.Then(env.AppMiddleware(env.GetProductsClassOfCompoundsHandler))).Methods("GET")

	router.Handle("/{item:products}/hazardstatements/", securechain.Then(env.AppMiddleware(env.GetProductsHazardStatementsHandler))).Methods("GET")
	router.Handle("/{item:products}/hazardstatements/{id}", securechain.Then(env.AppMiddleware(env.GetProductsHazardStatementHandler))).Methods("GET")

	router.Handle("/{item:products}/precautionarystatements/", securechain.Then(env.AppMiddleware(env.GetProductsPrecautionaryStatementsHandler))).Methods("GET")
	router.Handle("/{item:products}/precautionarystatements/{id}", securechain.Then(env.AppMiddleware(env.GetProductsPrecautionaryStatementHandler))).Methods("GET")

	router.Handle("/{item:products}/producerrefs/", securechain.Then(env.AppMiddleware(env.GetProductsProducerRefsHandler))).Methods("GET")
	router.Handle("/{item:products}/producers/", securechain.Then(env.AppMiddleware(env.GetProductsProducersHandler))).Methods("GET")
	router.Handle("/{item:products}/supplierrefs/", securechain.Then(env.AppMiddleware(env.GetProductsSupplierRefsHandler))).Methods("GET")
	router.Handle("/{item:products}/suppliers/", securechain.Then(env.AppMiddleware(env.GetProductsSuppliersHandler))).Methods("GET")
	router.Handle("/{item:products}/categories/", securechain.Then(env.AppMiddleware(env.GetProductsCategoriesHandler))).Methods("GET")
	router.Handle("/{item:products}/tags/", securechain.Then(env.AppMiddleware(env.GetProductsTagsHandler))).Methods("GET")

	router.Handle("/{item:products}/producers", securechain.Then(env.AppMiddleware(env.CreateProducerHandler))).Methods("POST")
	router.Handle("/{item:products}/suppliers", securechain.Then(env.AppMiddleware(env.CreateSupplierHandler))).Methods("POST")

	router.Handle("/f/{view:v}/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{view:vc}/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	// storages
	router.Handle("/{view:v}/{item:storages}", securechain.Then(env.AppMiddleware(env.VGetStoragesHandler))).Methods("GET")
	router.Handle("/{view:vc}/{item:storages}", securechain.Then(env.AppMiddleware(env.VCreateStorageHandler))).Methods("GET")
	router.Handle("/{item:storages}", securechain.Then(env.AppMiddleware(env.GetStoragesHandler))).Methods("GET")
	router.Handle("/{item:storages}/others", securechain.Then(env.AppMiddleware(env.GetOtherStoragesHandler))).Methods("GET")
	router.Handle("/{item:storages}/suppliers", securechain.Then(env.AppMiddleware(env.GetStoragesSuppliersHandler))).Methods("GET")
	router.Handle("/{item:storages}/units", securechain.Then(env.AppMiddleware(env.GetStoragesUnitsHandler))).Methods("GET")
	router.Handle("/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.GetStorageHandler))).Methods("GET")
	router.Handle("/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.UpdateStorageHandler))).Methods("PUT")
	router.Handle("/{item:storages}", securechain.Then(env.AppMiddleware(env.CreateStorageHandler))).Methods("POST")
	router.Handle("/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.DeleteStorageHandler))).Methods("DELETE")
	router.Handle("/{item:storages}/{id}/a", securechain.Then(env.AppMiddleware(env.ArchiveStorageHandler))).Methods("DELETE")
	router.Handle("/{item:storages}/{id}/r", securechain.Then(env.AppMiddleware(env.RestoreStorageHandler))).Methods("PUT")
	router.Handle("/{item:borrowings}", securechain.Then(env.AppMiddleware(env.ToogleStorageBorrowingHandler))).Methods("PUT")

	router.Handle("/f/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	router.Handle("/f/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	router.Handle("/f/{item:storages}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	router.Handle("/f/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")

	// validators
	router.Handle("/{item:validate}/entity/{id}/name/", securechain.Then(env.AppMiddleware(env.ValidateEntityNameHandler))).Methods("POST")
	router.Handle("/{item:validate}/person/{id}/email/", securechain.Then(env.AppMiddleware(env.ValidatePersonEmailHandler))).Methods("POST")
	router.Handle("/{item:validate}/product/{id}/casnumber/", securechain.Then(env.AppMiddleware(env.ValidateProductCasNumberHandler))).Methods("POST")
	router.Handle("/{item:validate}/product/{id}/cenumber/", securechain.Then(env.AppMiddleware(env.ValidateProductCeNumberHandler))).Methods("POST")
	router.Handle("/{item:validate}/product/{id}/name/", securechain.Then(env.AppMiddleware(env.ValidateProductNameHandler))).Methods("POST")
	router.Handle("/{item:validate}/product/{id}/empiricalformula/", securechain.Then(env.AppMiddleware(env.ValidateProductEmpiricalFormulaHandler))).Methods("POST")

	// formatters
	router.Handle("/{item:format}/product/{id}/empiricalformula/", securechain.Then(env.AppMiddleware(env.FormatProductEmpiricalFormulaHandler))).Methods("POST")

	// export download
	router.Handle("/{item:download}/{id}", securechain.Then(env.AppMiddleware(env.DownloadExportHandler))).Methods("GET")

	return router

}
