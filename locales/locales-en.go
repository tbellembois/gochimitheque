package locales

var LOCALES_EN = []byte(`
[test]
	one = "One test"
	other = "Several tests"

[nil]
	one = " "

[project_leader]
	one = "project leader"
[project_site]
	one = "web site"
[project_support]
	one = "support/bug report"
[project_license]
	one = "license"
[project_version]
	one = "version"

[wasm_loading]
	one = "Loading Web Assembly module."
[wasm_loaded]
	one = "Web Assembly module loaded."

[created]
	one = "created"
[modified]
	one = "modified"
[select_all]
	one = "select all"
[none]
	one = "none"
[nocamera]
	one = "no camera detected"

[confirm]
	one = "confirm"
[edit]
	one = "edit"
[delete]
	one = "delete"
[archive]
	one = "archive"
[save]
	one = "save"
[close]
	one = "close"
[list]
	one = "list"
[create]
	one = "create"
[check_all]
	one = "check all"

[required_input]
	one = "required input"
[error_occured]
	one = "an error occured"
[no_result]
	one = "no result"
[no_item]
	one = "no item"
[active_filter]
	one = "active filter(s)"
[no_filter]
	one = "no filter"
[remove_filter]
	one = "remove filter"

[empirical_formula_convert]
	one = "convert to empirical formula"
[no_empirical_formula]
	one = "no empirical formula"
[no_cas_number]
	one = "no CAS number"
[howto_magicalselector]
	one = "how to use the magical selector"

[password]
	one = "password"
[confirm_password]
	one = "confirm password"
[invalid_password]
	one = "invalid password"
[invalid_email]
	one = "invalid email"
[not_same_password]
	one = "you have not entered the same password"

[members]
	one = "members"
[storelocations]
	one = "storelocations"

[magical_selector]
	one = "magical selector"

[nb_duplicate]
	one = "number of items (bottles, boxes...)"
[nb_duplicate_comment]
	one = "will create a storage card per item with a different barecode (except if \"identical barecode\" is checked)"
[identical_barecode]
	one = "identical barecode"
[identical_barecode_comment]
	one = "generate the same barecode for every storage card - scanning a storage card qrcode will also return the storages with the same barecode"

[bt_loadingMessage]
	one = "loading..."
[bt_recordsPerPage]
	one = "records per page"
[bt_showingRowsTotal]
	one = "total records"
[bt_search]
	one = "search"
[bt_noMatches]
	one = "no matches"

[email_placeholder]
	one = "enter your email"
[submitlogin_text]
	one = "enter"
[password_placeholder]
	one = "enter your password"
[resetpassword_text]
	one = "reset password"
[resetpassword2_text]
	one = "reset my password, I am not a robot"
[resetpassword_message_mailsentto]
	one = "a reinitialization link has been sent to %s"
[resetpassword_areyourobot]
	one = "are you a robot?"
[resetpassword_mailbody1]
	one = '''
	This is your new temporary Chimithèque password: %s

	You can change it in the application.
	'''
[resetpassword_mailsubject1]
	one = "Chimithèque new temporary password\r\n"
[resetpassword_mailbody2]
	one = '''
	Click on this link to reinitialize your password: %sreset?token=%s

	You will then receive a new mail with a temporary password.
	'''
[resetpassword_mailsubject2]
	one = "Chimithèque password reset link\r\n"
[resetpassword_done]
	one = "A new temporary password has been sent to %s"

[createperson_mailsubject]
	one = "Chimithèque new account\r\n"
[createperson_mailbody]
	one = '''
	A Chimithèque account has been created for you.

	You can now initialize your password.

	Go to the login page %s, enter you email address %s and click on the "reset password" link.

	You will then receive a temporary password.
	'''

[logo_information1]
	one = "Chimithèque logo designed by "
[logo_information2]
	one = "Do not use or copy without her permission."

[welcomeannounce_text_title]
	one = "Main page additional text"
[welcomeannounce_text_modificationsuccess]
	one = "announce modified"

[s_tags]
	one = "tag(s)"
[s_category]
	one = "category"
[s_entity]
	one = "entity"
[s_storelocation]
	one = "store location"
[s_producerref]
	one = "producer reference number"
[s_custom_name_part_of]
	one = "part of name"
[s_name]
	one = "exact name"
[s_casnumber]
	one = "CAS number"
[s_empiricalformula]
	one = "emp. formula"
[s_storage_barecode]
	one = "barecode"
[s_signalword]
	one = "signal word"
[s_symbols]
	one = "symbol(s)"
[s_hazardstatements]
	one = "hazard statement(s)"
[s_precautionarystatements]
	one = "precautionary statement(s)"
[s_casnumber_cmr]
	one = "CMR"
[s_borrowing]
	one = "borrowed storages"
[s_storage_to_destroy]
	one = "storages to destroy"

[menu_home]
	one = "home"
[menu_bookmark]
	one = "bookmarks"
[menu_scanqr]
	one = "scan"
[menu_borrow]
	one = "my borrowed products"
[menu_create_productcard]
	one = "create a card"
[menu_entity]
	one = "entities"
[menu_storelocation]
	one = "store locations"
[menu_people]
	one = "people"
[menu_welcomeannounce]
	one = "change the login message"
[menu_password]
	one = "change my password"
[menu_logout]
	one = "logout"
[menu_about]
	one = "about"
[menu_account]
	one = "my account"

[clearsearch_text]
	one = "reset filters"
[search_text]
	one = "search"
[advancedsearch_text]
	one = "advanced search"

[chemical_product]
	one = "chemical product"
[biological_product]
	one = "biological reagent"
[consumable_product]
	one = "lab consumable"

[switchproductview_text]
	one = "switch to product view"
[switchstorageview_text]
	one = "switch to storage view"
[export_text]
	one = "export"
[download_export]
	one = "download export"
[export_progress]
	one = "export in progress -  this operation can be long"
[export_done]
	one = "export done"
[showdeleted_text]
	one = "show archives"
[hidedeleted_text]
	one = "hide archives"
[storeagain_text]
	one = "store this product"
[totalstock_text]
	one = "compute total stock"

[unit_label_title]
	one = "unit"
[supplier_label_title]
	one = "supplier"
[supplierref_label_title]
	one = "supplier reference number"

[add_producer_title]
	one = "add a producer to the list"
[producer_added]
	one = "producer added"
[add_supplier_title]
	one = "add a supplier to the list"
[supplier_added]
	one = "supplier added"

[store]
	one = "store"
[storages]
	one = "storages"
[storage]
	one = "storage"
[archives]
	one = "archives"
[ostorages]
	one = "availability"
[storage_create_title]
	one = "store a product"
[storage_update_title]
	one = "update a storage"
[storage_clone]
	one = "clone"
[storage_borrow]
	one = "borrow"
[storage_unborrow]
	one = "unborrow"
[storage_restore]
	one = "restore"
[storage_showhistory]
	one = "show history"
[storage_history]
	one = "history"
[storage_restored_message]
	one = "storage restored"
[storage_trashed_message]
	one = "storage trashed"
[storage_deleted_message]
	one = "storage deleted"
[storage_borrow_updated]
	one = "borrow updated"
[storage_created_message]
	one = "storage created"
[storage_updated_message]
	one = "storage updated"

[storage_storelocation_title]
	one = "store location"
[storage_concentration_title]
	one = "concentration"
[storage_quantity_title]
	one = "quantity"
[storage_barecode_title]
	one = "barecode"
[storage_create_barecode_comment]
	one = "if you leave this field empty a barecode will be auto-generated"
[storage_batchnumber_title]
	one = "batch number"
[storage_entrydate_title]
	one = "entry date"
[storage_exitdate_title]
	one = "exit date"
[storage_openingdate_title]
	one = "opening date"
[storage_expirationdate_title]
	one = "expiration date"
[storage_borrower_title]
	one = "borrower"
[storage_comment_title]
	one = "comment"
[storage_reference_title]
	one = "reference"
[storage_todestroy_title]
	one = "to destroy"
[storage_product_table_header]
	one = "product"
[storage_storelocation_table_header]
	one = "store location"
[storage_quantity_table_header]
	one = "quantity"
[storage_barecode_table_header]
	one = "barecode"
[storage_storelocation_placeholder]
	one = "select a store location"
[storage_borrower_placeholder]
	one = "select a borrower"
[storage_supplier_placeholder]
	one = "select or enter a supplier"
[storage_print_qrcode]
	one = "print qrcode"
[storage_number_of_unit]
	one = "number of unit(s)"
[storage_number_of_bag]
	one = "number of bag(s)"
[storage_number_of_bag_comment]
	one = "only if the number of units per bag for the corresponding product is set"
[storage_number_of_carton]
	one = "number of carton(s)"
[storage_number_of_carton_comment]
	one = "only if the number of units per carton for the corresponding product is set"
[storage_one_number_required]
	one = "at least one of the numbers required"

[stock_storelocation_title]
	one = "in this store location"
[stock_storelocation_sub_title]
	one = "with children store locations"

[empiricalformula_label_title]
	one = "empirical formula"
[cenumber_label_title]
	one = "EC"
[casnumber_label_title]
	one = "CAS"
[casnumber_cmr_title]
	one = "CMR"
[signalword_label_title]
	one = "signal word"
[symbol_label_title]
	one = "symbol(s)"
[linearformula_label_title]
	one = "linear formula"
[hazardstatement_label_title]
	one = "hazard statement(s)"
[precautionarystatement_label_title]
	one = "precautionary statement(s)"
[classofcompound_label_title]
	one = "class(es) of compounds"
[physicalstate_label_title]
	one = "physical state"
[name_label_title]
	one = "name"
[synonym_label_title]
	one = "synonym(s)"

[restricted]
	one = "restricted access"
[bookmark]
	one = "bookmark"
[unbookmark]
	one = "remove bookmark"
[product_create_title]
	one = "create a product card"
[product_update_title]
	one = "update product"
[product_threedformula_title]
	one = "3D formula"
[product_twodformula_title]
	one = "molecule picture"
[product_threedformula_mol_title]
	one = "3D formula MOL file"
[product_msds_title]
	one = "MSDS link"
[product_sheet_title]
	one = "producer product cheet"
[product_temperature_title]
	one = "preconised storage temperature"
[product_number_per_carton_title]
	one = "number of units per carton"
[product_number_per_bag_title]
	one = "number of units per bag"
[producer_label_title]
	one = "producer"
[producerref_label_title]
	one = "producer reference number"
[producerref_create_needproducer]
	one = "to create a new reference select a producer first"
[supplierref_create_needsupplier]
	one = "to create a new reference select a supplier first"
[category_label_title]
	one = "category"
[tag_label_title]
	one = "tag(s)"
[product_disposalcomment_title]
	one = "disposal comment"
[product_remark_title]
	one = "remark"
[product_specificity_title]
	one = "specificity"
[product_radioactive_title]
	one = "radioactive"
[product_restricted_title]
	one = "restricted access"
[product_name_table_header]
	one = "name"
[product_empiricalformula_table_header]
	one = "emp. formula"
[product_cas_table_header]
	one = "CAS"
[product_specificity_table_header]
	one = "spec."
[product_cas_placeholder]
	one = "select or enter a CAS number"
[product_ce_placeholder]
	one = "select or enter an EC number"
[product_physicalstate_placeholder]
	one = "select or enter a physical state"
[product_signalword_placeholder]
	one = "select a signal word"
[product_classofcompound_placeholder]
	one = "select or enter class(es) of compound"
[product_name_placeholder]
	one = "select or enter a name"
[product_synonyms_placeholder]
	one = "select or enter name(s)"
[producerref_placeholder]
	one = "select or enter a reference"
[product_empiricalformula_placeholder]
	one = "select or enter a formula"
[product_linearformula_placeholder]
	one = "select or enter a formula"
[product_symbols_placeholder]
	one = "select symbol(s)"
[product_hazardstatements_placeholder]
	one = "select statement(s)"
[product_precautionarystatements_placeholder]
	one = "select statement(s)"
[product_producer_placeholder]
	one = "select a producer"
[product_producerref_placeholder]
	one = "select or enter a producer reference"
[product_supplier_placeholder]
	one = "select a supplier"
[product_supplierref_placeholder]
	one = "select or enter supplier reference(s)"
[product_tag_placeholder]
	one = "select or enter tag(s)"
[product_category_placeholder]
	one = "select or enter a category"
[product_unit_placeholder]
	one = "select a unit"
[product_deleted_message]
	one = "product deleted"
[product_updated_message]
	one = "product updated"
[product_created_message]
	one = "product created"
[product_flammable]
	one = "flammable"

[person_create_title]
	one = "create a person"
[person_update_title]
	one = "update person"
[person_deleted_message]
	one = "person deleted"
[person_email_title]
	one = "email"
[person_password_title]
	one = "password"
[person_entity_title]
	one = "entity(ies)"
[person_permission_title]
	one = "permissions"
[person_email_table_header]
	one = "email"
[person_can_not_remove_entity_manager]
	one = "this entity can not be removed, the user is one of its manager"
[person_created_message]
	one = "person created"
[person_updated_message]
	one = "person updated"
[person_password_updated_message]
	one = "password updated"
[person_entity_placeholder]
	one = "select entity(ies)"
[person_select_all_none_storage]
	one = "select all 'no permission'"
[person_select_all_r_storage]
	one = "select all 'view only'"
[person_select_all_rw_storage]
	one = "select all 'view, modify, create and delete'"
[person_show_password]
  one = "show password field"
  
[permission_product]
	one = "products"
[permission_rproduct]
	one = "restricted products"
[permission_storages]
	one = "storages"
[permission_none]
	one = "no permission"
[permission_read]
	one = "view only"
[permission_crud]
	one = "view, modify, create and delete"

[storelocation_create_title]
	one = "create a store location"
[storelocation_update_title]
	one = "update store location"
[storelocation_deleted_message]
	one = "store location deleted"
[storelocation_created_message]
	one = "store location created"
[storelocation_updated_message]
	one = "store location updated"
[storelocation_parent_title]
	one = "parent"
[storelocation_entity_title]
	one = "entity"
[storelocation_canstore_title]
	one = "can store"
[storelocation_color_title]
	one = "color"
[storelocation_name_title]
	one = "name"
[storelocation_name_table_header]
	one = "name"
[storelocation_entity_table_header]
	one = "entity"
[storelocation_color_table_header]
	one = "color"
[storelocation_canstore_table_header]
	one = "can store"
[storelocation_parent_table_header]
	one = "parent"
[storelocation_entity_placeholder]
	one = "select an entity"
[storelocation_storelocation_placeholder]
	one = "select an entity first"

[entity_create_title]
	one = "create entity"
[entity_update_title]
	one = "update entity"
[entity_deleted_message]
	one = "entity deleted"
[entity_created_message]
	one = "entity created"
[entity_updated_message]
	one = "entity updated"
[entity_name_table_header]
	one = "name"
[entity_description_table_header]
	one = "description"
[entity_manager_table_header]
	one = "manager(s)"
[entity_manager_placeholder]
	one = "select manager(s)"

[entity_nameexist_validate]
	one = "entity with this name already present" 
[person_emailexist_validate]
	one = "person with this email already present" 
[empiricalformula_validate]
	one = "invalid empirical formula"
[casnumber_validate_wrongcas]
	one = "invalid CAS number"
[casnumber_validate_casspecificity]
	one = "CAS number/specificity pair already exist"
[cenumber_validate]
	one = "invalid EC number"
`)
