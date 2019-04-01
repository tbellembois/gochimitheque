package locales

var LOCALES_EN = []byte(`
[test]
	one = "One test"
	other = "Several tests"

[created]
	one = "created"
[modified]
	one = "modified"

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

[members]
	one = "members"
[storelocations]
	one = "storelocations"

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

[s_custom_name_part_of]
	one = "part of name"
[s_name]
	one = "name"
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
[menu_create_productcard]
	one = "create product card"
[menu_entity]
	one = "entities"
[menu_storelocation]
	one = "store locations"
[menu_people]
	one = "people"
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

[storage_quantity_title]
	one = "quantity"
[storage_barecode_title]
	one = "barecode"
[storage_batchnumber_title]
	one = "batch number"
[supplier_label_title]
	one = "supplier"
[storage_entrydate_title]
	one = "entry date"
[storage_exitdate_title]
	one = "exit date"
[storage_openingdate_title]
	one = "opening date"
[storage_expirationdate_title]
	one = "expiration date"
[storage_comment_title]
	one = "comment"

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

[stock_storelocation_title]
	one = "in this store location"
[stock_storelocation_sub_title]
	one = "including children store locations"

[empiricalformula_label_title]
	one = "empirical formula"
[cenumber_label_title]
	one = "CE"
[casnumber_label_title]
	one = "CAS"
[casnumber_cmr_title]
	one = "CMR"
[signalword_label_title]
	one = "signal word"
[linearformula_label_title]
	one = "liner formula"
[hazardstatement_label_title]
	one = "hazard statement(s)"
[precautionarystatement_label_title]
	one = "precautionary statement(s)"
[classofcompound_label_title]
	one = "class of compounds"
[physicalstate_label_title]
	one = "physical state"
[product_threedformula_title]
	one = "3D formula"
[product_msds_title]
	one = "MSDS"
[product_disposalcomment_title]
	one = "disposal comment"
[product_remark_title]
	one = "remark"
[product_radioactive_title]
	one = "radioactive"
[product_restricted_title]
	one = "restricted access"

[entity_create_title]
	one = "create entity"
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
`)
