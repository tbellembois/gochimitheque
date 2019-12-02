package locales

var LOCALES_EN = []byte(`
[test]
	one = "One test"
	other = "Several tests"

[nil]
	one = " "

[created]
	one = "created"
[modified]
	one = "modified"
[select_all]
	one = "select all"
[none]
	one = "none"

[edit]
	one = "edit"
[delete]
	one = "delete"
[save]
	one = "save"
[close]
	one = "close"
[list]
	one = "list"
[create]
	one = "create"

[required_input]
	one = "required input"

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

[members]
	one = "members"
[storelocations]
	one = "storelocations"

[magical_selector]
	one = "magical selector"

[nb_duplicate]
	one = "nb of duplicates"

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
[resetpassword_warning_enteremail]
	one = "enter your email in the login form"
[resetpassword_message_mailsentto]
	one = "a reinitialization link has been sent to"
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
	Click on this link to reinitialize your password: %s%sreset?token=%s

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

[s_custom_name_part_of]
	one = "part of name"
[s_name]
	one = "exact name"
[s_casnumber]
	one = "CAS"
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

[menu_home]
	one = "home"
[menu_bookmark]
	one = "my bookmarks"
[menu_scanqr]
	one = "scan a QR code"
[menu_borrow]
	one = "my borrowed products"
[menu_create_productcard]
	one = "create product card"
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

[clearsearch_text]
	one = "clear search form"
[search_text]
	one = "search"
[advancedsearch_text]
	one = "advanced search"

[switchproductview_text]
	one = "switch to product view"
[switchstorageview_text]
	one = "switch to storage view"
[export_text]
	one = "export"
[showdeleted_text]
	one = "show deleted"
[hidedeleted_text]
	one = "hide deleted"
[storeagain_text]
	one = "store this product"
[totalstock_text]
	one = "show total stock"

[unit_label_title]
	one = "unit"
[supplier_label_title]
	one = "supplier"

[storage_create_title]
	one = "create storage"
[storage_update_title]
	one = "update storage"
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

[storage_storelocation_title]
	one = "store location"
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

[stock_storelocation_title]
	one = "in this store location"
[stock_storelocation_sub_title]
	one = "including children store locations"

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

[product_create_title]
	one = "create product"
[product_update_title]
	one = "update product"
[product_threedformula_title]
	one = "3D formula"
[product_threedformula_mol_title]
	one = "3D formula MOL file"
[product_msds_title]
	one = "MSDS"
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
	one = "select a physical state"
[product_signalword_placeholder]
	one = "select a signal word"
[product_classofcompound_placeholder]
	one = "select or enter class(es) of compound"
[product_name_placeholder]
	one = "select or enter a name"
[product_synonyms_placeholder]
	one = "select or enter name(s)"
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

[person_create_title]
	one = "create person"
[person_update_title]
	one = "update person"
[person_email_title]
	one = "email"
[person_password_title]
	one = "password"
[person_entity_title]
	one = "entity(ies)"
[person_email_table_header]
	one = "email"

[storelocation_create_title]
	one = "create store location"
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
