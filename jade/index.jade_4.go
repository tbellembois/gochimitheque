// Code generated by "jade.go"; DO NOT EDIT.

package jade

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	index_4__30 = `</span></a></li></ul></div></nav><div id="test"></div><div id="search" class="row pt-sm-4 pb-sm-4 pl-sm-2 pr-sm-2 mt-sm-2 mb-sm-2 ml-sm-5 mr-sm-5 bg-light border rounded collapse show">`
	index_4__31 = `<div class="col"><div class="row"><div class="col-sm-12">`
	index_4__32 = `</div></div><div class="row"><div class="col-sm-6">`
	index_4__33 = `</div><div class="col-sm-6"></div><div class="col-sm-6">`
	index_4__34 = `</div><div class="col-sm-6">`
	index_4__37 = `</div></div><div class="row collapse" id="advancedsearch"><div class="col-sm-12"><div class="form-row"><div class="form-group col-sm-6">`
	index_4__38 = `</div><div class="form-group col-sm-6">`
	index_4__39 = `</div></div><div class="form-row"><div class="form-group col-sm-6">`
	index_4__42 = `</div></div></div></div><div class="row"><div class="col col-sm-10"><button id="clearsearch" class="btn btn-link mr-sm-2" type="button" onclick="clearsearch();"><span class="mdi mdi-broom mdi-24px iconlabel">`
	index_4__43 = `</span></button><button id="search" class="btn btn-lg btn-link mr-sm-2" type="button" onclick="search();"><span class="mdi mdi-magnify mdi-24px iconlabel">`
	index_4__44 = `</span></button></div><div class="col col-sm-2"><button class="btn btn-link" data-toggle="collapse" href="#advancedsearch" aria-expanded="false"><span class="mdi mdi-magnify-plus-outline mdi-24px iconlabel">`
	index_4__45 = `</span></button></div></div></div></div><div class="row"><div class="col-sm-10"><div id="filter-item"></div></div><div class="col-sm-2 d-flex justify-content-end"><div id="button-store"></div></div></div><div id="toolbar" class="row"><div class="col toggleable"><button class="btn btn-link d-none" id="switchview" type="button" onclick="switchProductStorageView()"><span class="mdi mdi-cube-unfolded mdi-24px iconlabel">`
	index_4__46 = `</span></button></div><div class="col"><button class="btn btn-link" id="export" type="button" onclick="exportAll()"><span class="mdi mdi-content-save mdi-24px iconlabel">`
	index_4__47 = `</span></button></div></div><div id="exportlink" class="modal fade" role="dialog" tabindex="-1" aria-labelledby="exportlinkLabel" aria-hidden="true"><div class="modal-dialog modal-sm" role="document"><div class="modal-content"><div class="modal-body mx-auto" id="exportlink-body"></div><div class="modal-footer"><button class="btn btn-link" type="button" data-dismiss="modal"><span class="mdi mdi-close-box mdi-24px iconlabel">`
	index_4__48 = `</span></button></div></div></div></div><div id="accordion"><div id="list-collapse" class="collapse show" data-parent="#accordion"><header class="row"><div class="col-sm-12"><table id="table" data-toggle="table" data-striped="true" data-search="false" data-toolbar="#toolbar" data-side-pagination="server" data-page-list="[5, 10, 20, 50, 100, 200, 500]" data-pagination="true" data-ajax="getData" data-query-params="queryParams" data-sort-name="name.name_label" data-detail-view="true" data-detail-formatter="detailFormatter" data-row-attributes="rowAttributes"><thead><tr><!--  th(data-field='product_id' data-sortable='true') ID --><th class="th-width-200" data-field="name.name_label" data-sortable="true">name</th><th data-field="empiricalformula.empiricalformula_label" data-sortable="true">empirical formula</th><th data-field="casnumber.casnumber_label" data-sortable="true">CAS</th><th data-field="product_specificity" data-sortable="false" data-formatter="product_specificityFormatter">specificity</th><th data-field="product_sl" data-formatter="product_slFormatter" data-sortable="false"></th><th data-field="operate" data-formatter="operateFormatter" data-events="operateEvents"></th></tr></thead></table></div></header></div><div id="edit-collapse" class="collapse" data-parent="#accordion">`
	index_4__49 = `<form id="product"><input id="index" type="hidden" name="index" value=""/><input id="product_id" type="hidden" name="product_id" value=""/><input id="exactMatchEmpiricalFormula" type="hidden"/><input id="exactMatchlinearFormula" type="hidden"/><input id="exactMatchCasNumber" type="hidden"/><input id="exactMatchCeNumber" type="hidden"/><input id="exactMatchName" type="hidden"/><input id="exactMatchSynonyms" type="hidden"/><input id="exactMatchClassofcompounds" type="hidden"/><input id="exactMatchPhysicalstate" type="hidden"/><div class="form-row"><div class="form-group col-sm-auto"><span class="badge badge-pill badge-danger">&nbsp;</span></div><div class="form-group col-sm-5">`
	index_4__50 = `</div><div class="form-group col-sm-5">`
	index_4__51 = `</div></div><div class="form-row"><div class="form-group col-sm-12">`
	index_4__52 = `</div></div><div class="form-row"><div class="form-group col-sm-auto"><span>magic selector coming soon</span></div></div><div class="form-row"><div class="form-group col-sm-auto"><span class="badge badge-pill badge-danger">&nbsp;</span></div><div class="form-group col-sm-5">`
	index_4__54 = `</div><div class="form-group col-sm-1"><button id="fconverter" class="btn btn-link" type="button" data-toggle="popover" data-content="no result" title="convert linear to empirical formula" onclick="linearToEmpirical();"><i class="material-icons">loop</i></button></div></div><div class="form-row"><div class="form-group col-sm-auto"><span class="badge badge-pill badge-danger">&nbsp;</span></div><div class="form-group col-sm-6">`
	index_4__69 = `</div></div><button id="save" class="btn btn-link" type="button" onclick="saveProduct()"><span class="mdi mdi-content-save mdi-24px iconlabel">`
	index_4__70 = `</span></button><button class="btn btn-link" type="button" onclick="closeEdit();"><span class="mdi mdi-content-save mdi-24px iconlabel">`
	index_4__71 = `</span></button></form></div></div></div><!--  Code generated by go generate; DO NOT EDIT. --><script>    
	var locale_en_advancedsearch_text = "advanced search";
	
	var locale_en_casnumber_cmr_title = "CMR";
	
	var locale_en_casnumber_label_title = "CAS";
	
	var locale_en_cenumber_label_title = "CE";
	
	var locale_en_classofcompound_label_title = "class of compounds";
	
	var locale_en_clearsearch_text = "clear search form";
	
	var locale_en_close = "close";
	
	var locale_en_created = "created";
	
	var locale_en_createperson_mailsubject = "Chimithèque new account\r\n";
	
	var locale_en_delete = "delete";
	
	var locale_en_edit = "edit";
	
	var locale_en_email_placeholder = "enter your email";
	
	var locale_en_empiricalformula_label_title = "empirical formula";
	
	var locale_en_entity_create_title = "create entity";
	
	var locale_en_entity_created_message = "entity created";
	
	var locale_en_entity_deleted_message = "entity deleted";
	
	var locale_en_entity_description_table_header = "description";
	
	var locale_en_entity_manager_table_header = "manager(s)";
	
	var locale_en_entity_name_table_header = "name";
	
	var locale_en_entity_updated_message = "entity updated";
	
	var locale_en_export_text = "export";
	
	var locale_en_hazardstatement_label_title = "hazard statement(s)";
	
	var locale_en_hidedeleted_text = "hide deleted";
	
	var locale_en_linearformula_label_title = "liner formula";
	
	var locale_en_logo_information1 = "Chimithèque logo designed by ";
	
	var locale_en_logo_information2 = "Do not use or copy without her permission.";
	
	var locale_en_members = "members";
	
	var locale_en_menu_bookmark = "my bookmarks";
	
	var locale_en_menu_create_productcard = "create product card";
	
	var locale_en_menu_entity = "entities";
	
	var locale_en_menu_home = "home";
	
	var locale_en_menu_logout = "logout";
	
	var locale_en_menu_password = "change my password";
	
	var locale_en_menu_people = "people";
	
	var locale_en_menu_storelocation = "store locations";
	
	var locale_en_modified = "modified";
	
	var locale_en_password_placeholder = "enter your password";
	
	var locale_en_physicalstate_label_title = "physical state";
	
	var locale_en_precautionarystatement_label_title = "precautionary statement(s)";
	
	var locale_en_product_disposalcomment_title = "disposal comment";
	
	var locale_en_product_msds_title = "MSDS";
	
	var locale_en_product_radioactive_title = "radioactive";
	
	var locale_en_product_remark_title = "remark";
	
	var locale_en_product_restricted_title = "restricted access";
	
	var locale_en_product_threedformula_title = "3D formula";
	
	var locale_en_resetpassword2_text = "reset my password, I am not a robot";
	
	var locale_en_resetpassword_areyourobot = "are you a robot?";
	
	var locale_en_resetpassword_done = "A new temporary password has been sent to %s";
	
	var locale_en_resetpassword_mailsubject1 = "Chimithèque new temporary password\r\n";
	
	var locale_en_resetpassword_mailsubject2 = "Chimithèque password reset link\r\n";
	
	var locale_en_resetpassword_message_mailsentto = "a reinitialization link has been sent to";
	
	var locale_en_resetpassword_text = "reset password";
	
	var locale_en_resetpassword_warning_enteremail = "enter your email in the login form";
	
	var locale_en_s_casnumber = "CAS";
	
	var locale_en_s_casnumber_cmr = "CMR";
	
	var locale_en_s_custom_name_part_of = "part of name";
	
	var locale_en_s_empiricalformula = "emp. formula";
	
	var locale_en_s_hazardstatements = "hazard statement(s)";
	
	var locale_en_s_name = "name";
	
	var locale_en_s_precautionarystatements = "precautionary statement(s)";
	
	var locale_en_s_signalword = "signal word";
	
	var locale_en_s_storage_barecode = "barecode";
	
	var locale_en_s_symbols = "symbol(s)";
	
	var locale_en_save = "save";
	
	var locale_en_search_text = "search";
	
	var locale_en_showdeleted_text = "show deleted";
	
	var locale_en_signalword_label_title = "signal word";
	
	var locale_en_stock_storelocation_sub_title = "including children store locations";
	
	var locale_en_stock_storelocation_title = "in this store location";
	
	var locale_en_storage_barecode_title = "barecode";
	
	var locale_en_storage_batchnumber_title = "batch number";
	
	var locale_en_storage_borrow = "borrow";
	
	var locale_en_storage_clone = "clone";
	
	var locale_en_storage_comment_title = "comment";
	
	var locale_en_storage_entrydate_title = "entry date";
	
	var locale_en_storage_exitdate_title = "exit date";
	
	var locale_en_storage_expirationdate_title = "expiration date";
	
	var locale_en_storage_openingdate_title = "opening date";
	
	var locale_en_storage_quantity_title = "quantity";
	
	var locale_en_storage_restore = "restore";
	
	var locale_en_storage_showhistory = "show history";
	
	var locale_en_storage_unborrow = "unborrow";
	
	var locale_en_storelocations = "storelocations";
	
	var locale_en_submitlogin_text = "enter";
	
	var locale_en_supplier_label_title = "supplier";
	
	var locale_en_switchproductview_text = "switch to product view";
	
	var locale_en_switchstorageview_text = "switch to storage view";
	
	var locale_en_test = "One test";
	
    
	var locale_fr_advancedsearch_text = "recherche avancée";
	
	var locale_fr_casnumber_cmr_title = "CMR";
	
	var locale_fr_casnumber_label_title = "CAS";
	
	var locale_fr_cenumber_label_title = "CE";
	
	var locale_fr_classofcompound_label_title = "famille chimique";
	
	var locale_fr_clearsearch_text = "effacer le formulaire";
	
	var locale_fr_close = "fermer";
	
	var locale_fr_created = "créé";
	
	var locale_fr_createperson_mailsubject = "Chimithèque nouveau compte\r\n";
	
	var locale_fr_delete = "supprimer";
	
	var locale_fr_edit = "editer";
	
	var locale_fr_email_placeholder = "entrez votre email";
	
	var locale_fr_empiricalformula_label_title = "formule brute";
	
	var locale_fr_entity_create_title = "créer une entité";
	
	var locale_fr_entity_created_message = "entité crée";
	
	var locale_fr_entity_deleted_message = "entité supprimée";
	
	var locale_fr_entity_description_table_header = "description";
	
	var locale_fr_entity_manager_table_header = "responsable(s)";
	
	var locale_fr_entity_name_table_header = "nom";
	
	var locale_fr_entity_updated_message = "entité mise à jour";
	
	var locale_fr_export_text = "exporter";
	
	var locale_fr_hazardstatement_label_title = "mention(s) de danger H-EUH";
	
	var locale_fr_hidedeleted_text = "cacher supprimés";
	
	var locale_fr_linearformula_label_title = "formule linéaire";
	
	var locale_fr_logo_information1 = "Logo Chimithèque réalisé par ";
	
	var locale_fr_logo_information2 = "Ne pas utiliser ou copier sans sa permission.";
	
	var locale_fr_members = "membres";
	
	var locale_fr_menu_bookmark = "mes favoris";
	
	var locale_fr_menu_create_productcard = "créer fiche produit";
	
	var locale_fr_menu_entity = "entités";
	
	var locale_fr_menu_home = "accueil";
	
	var locale_fr_menu_logout = "déconnexion";
	
	var locale_fr_menu_password = "changer mon mot de passe";
	
	var locale_fr_menu_people = "utilisateurs";
	
	var locale_fr_menu_storelocation = "entrepôts";
	
	var locale_fr_modified = "modifié";
	
	var locale_fr_password_placeholder = "entrez votre mot de passe";
	
	var locale_fr_physicalstate_label_title = "état physique";
	
	var locale_fr_precautionarystatement_label_title = "conseil(s) de prudence P";
	
	var locale_fr_product_disposalcomment_title = "commentaire de destruction";
	
	var locale_fr_product_msds_title = "FDS";
	
	var locale_fr_product_radioactive_title = "radioactif";
	
	var locale_fr_product_remark_title = "remarque";
	
	var locale_fr_product_restricted_title = "accès restreint";
	
	var locale_fr_product_threedformula_title = "formule 3D";
	
	var locale_fr_resetpassword2_text = "réinitialiser mon mot de passe, je ne suis pas un robot";
	
	var locale_fr_resetpassword_areyourobot = "êtes vous un robot ?";
	
	var locale_fr_resetpassword_done = "Un nouveau mot de passe temporaire a été envoyé à %s";
	
	var locale_fr_resetpassword_mailsubject1 = "Chimithèque nouveau mot de passe temporaire\r\n";
	
	var locale_fr_resetpassword_mailsubject2 = "Chimithèque lien de réinitialisation de mot de passe\r\n";
	
	var locale_fr_resetpassword_message_mailsentto = "un mail de réinitialisation a été envoyé à";
	
	var locale_fr_resetpassword_text = "réinitialiser mon mot de passe";
	
	var locale_fr_resetpassword_warning_enteremail = "entrez votre adresse mail dans le formulaire";
	
	var locale_fr_s_casnumber = "CAS";
	
	var locale_fr_s_custom_name_part_of = "partie du nom";
	
	var locale_fr_s_empiricalformula = "formule brute";
	
	var locale_fr_s_hazardstatements = "mention(s) de danger H-EUH";
	
	var locale_fr_s_name = "nom";
	
	var locale_fr_s_precautionarystatements = "conseil(s) de prudence P";
	
	var locale_fr_s_signalword = "mention d'avertissement";
	
	var locale_fr_s_storage_barecode = "code barre";
	
	var locale_fr_s_symbols = "symbole(s)";
	
	var locale_fr_save = "enregistrer";
	
	var locale_fr_search_text = "rechercher";
	
	var locale_fr_showdeleted_text = "voir supprimés";
	
	var locale_fr_signalword_label_title = "mention d'avertissement";
	
	var locale_fr_storage_barecode_title = "code barre";
	
	var locale_fr_storage_batchnumber_title = "numéro de lot";
	
	var locale_fr_storage_borrow = "emprunter";
	
	var locale_fr_storage_clone = "cloner";
	
	var locale_fr_storage_comment_title = "commentaire";
	
	var locale_fr_storage_entrydate_title = "date d'entrée";
	
	var locale_fr_storage_exitdate_title = "date de sortie";
	
	var locale_fr_storage_expirationdate_title = "date d'expiration";
	
	var locale_fr_storage_openingdate_title = "date d'ouverture";
	
	var locale_fr_storage_quantity_title = "quantité";
	
	var locale_fr_storage_restore = "restaurer";
	
	var locale_fr_storage_showhistory = "voir historique";
	
	var locale_fr_storage_unborrow = "restituer";
	
	var locale_fr_storelocations = "entrepôts";
	
	var locale_fr_submitlogin_text = "entrer";
	
	var locale_fr_supplier_label_title = "fournisseur";
	
	var locale_fr_switchproductview_text = "vue par produits";
	
	var locale_fr_switchstorageview_text = "vue par stockages";
	
	var locale_fr_test = "Un test";
	</script>`
	index_4__89  = `"></script><script src="../js/chim/product_storage.js"></script><script src="../js/chim/product.js"></script></body></html>`
	index_4__90  = `<input id="`
	index_4__91  = `" type="hidden" name="`
	index_4__92  = `" value="`
	index_4__93  = `"/>`
	index_4__114 = `<div class="row"><div class="col-sm-1"><span class="`
	index_4__115 = `"></span></div><div class="col-sm-11"><select style="width: 100% !important;" id="`
	index_4__117 = `"></select></div></div>`
	index_4__118 = `<div class="form-group row"><label class="col-sm-3 col-form-label" for="`
	index_4__120 = `</label><div class="col-sm-9"><input class="form-control" type="text" id="`
	index_4__123 = `"/></div></div>`
	index_4__126 = `</label><div class="col-sm-9"><select class="form-control" style="width: 100% !important;" id="`
	index_4__149 = `"></select></div>`
	index_4__165 = `<div class="form-check"><input class="form-check-input" type="checkbox" id="`
	index_4__166 = `"/><label class="form-check-label" for="`
	index_4__168 = `</label></div>`
	index_4__223 = `</label><input class="form-control" type="file" id="`
	index_4__271 = `</label><textarea class="form-control" type="text" id="`
	index_4__274 = `"></textarea></div>`
)

func Productindex(c ViewContainer, wr io.Writer) {
	buffer := &WriterAsBuffer{wr}

	buffer.WriteString(index__0)
	WriteAll(c.ProxyPath+"css/bootstrap.min.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/bootstrap-table.min.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/select2.min.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/bootstrap-colorpicker.min.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/fontawesome.all.min.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/chimitheque.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/materialdesignicons.min.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/bootstrap-toggle.min.css", true, buffer)
	buffer.WriteString(index__1)
	WriteAll(c.ProxyPath+"css/animate.min.css", true, buffer)
	buffer.WriteString(index__9)
	WriteAll(c.ProxyPath+"js/jquery-3.3.1.min.js", true, buffer)
	buffer.WriteString(index__10)
	WriteAll(c.ProxyPath+"img/logo_chimitheque_small.png", true, buffer)
	buffer.WriteString(index__11)
	WriteAll(c.ProxyPath+"v/products", true, buffer)
	buffer.WriteString(index__12)
	WriteAll(T("menu_home", 1), true, buffer)
	buffer.WriteString(index__13)
	WriteAll(c.ProxyPath+"v/products?bookmark=true", true, buffer)
	buffer.WriteString(index__14)
	WriteAll(T("menu_bookmark", 1), true, buffer)
	buffer.WriteString(index__15)
	WriteAll(c.ProxyPath+"vc/products", true, buffer)
	buffer.WriteString(index__16)
	WriteAll(T("menu_create_productcard", 1), true, buffer)
	buffer.WriteString(index__17)
	WriteAll(T("menu_entity", 1), true, buffer)
	buffer.WriteString(index__18)
	WriteAll(c.ProxyPath+"v/entities", true, buffer)
	buffer.WriteString(index__19)
	WriteAll(c.ProxyPath+"vc/entities", true, buffer)
	buffer.WriteString(index__20)
	WriteAll(T("menu_storelocation", 1), true, buffer)
	buffer.WriteString(index__18)
	WriteAll(c.ProxyPath+"v/storelocations", true, buffer)
	buffer.WriteString(index__22)
	WriteAll(c.ProxyPath+"vc/storelocations", true, buffer)
	buffer.WriteString(index__23)
	WriteAll(T("menu_people", 1), true, buffer)
	buffer.WriteString(index__18)
	WriteAll(c.ProxyPath+"v/people", true, buffer)
	buffer.WriteString(index__25)
	WriteAll(c.ProxyPath+"vc/people", true, buffer)
	buffer.WriteString(index__26)
	WriteAll(c.ProxyPath+"vu/peoplepass", true, buffer)
	buffer.WriteString(index__27)
	WriteAll(T("menu_password", 1), true, buffer)
	buffer.WriteString(index__28)
	WriteAll(c.ProxyPath+"delete-token", true, buffer)
	buffer.WriteString(index__29)
	WriteAll(T("menu_logout", 1), true, buffer)
	buffer.WriteString(index_4__30)

	{
		var (
			name  = "s_entity"
			value = ""
		)

		buffer.WriteString(index_4__90)
		WriteEscString("hidden_"+name, buffer)
		buffer.WriteString(index_4__91)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__92)
		WriteEscString(value, buffer)
		buffer.WriteString(index_4__93)
	}

	{
		var (
			name  = "s_product"
			value = ""
		)

		buffer.WriteString(index_4__90)
		WriteEscString("hidden_"+name, buffer)
		buffer.WriteString(index_4__91)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__92)
		WriteEscString(value, buffer)
		buffer.WriteString(index_4__93)
	}

	{
		var (
			name  = "s_history"
			value = ""
		)

		buffer.WriteString(index_4__90)
		WriteEscString("hidden_"+name, buffer)
		buffer.WriteString(index_4__91)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__92)
		WriteEscString(value, buffer)
		buffer.WriteString(index_4__93)
	}

	{
		var (
			name  = "s_storage"
			value = ""
		)

		buffer.WriteString(index_4__90)
		WriteEscString("hidden_"+name, buffer)
		buffer.WriteString(index_4__91)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__92)
		WriteEscString(value, buffer)
		buffer.WriteString(index_4__93)
	}

	{
		var (
			name  = "s_storage_archive"
			value = ""
		)

		buffer.WriteString(index_4__90)
		WriteEscString("hidden_"+name, buffer)
		buffer.WriteString(index_4__91)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__92)
		WriteEscString(value, buffer)
		buffer.WriteString(index_4__93)
	}

	{
		var (
			name  = "s_bookmark"
			value = ""
		)

		buffer.WriteString(index_4__90)
		WriteEscString("hidden_"+name, buffer)
		buffer.WriteString(index_4__91)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__92)
		WriteEscString(value, buffer)
		buffer.WriteString(index_4__93)
	}

	buffer.WriteString(index_4__31)

	{
		var (
			name = "s_storelocation"
			icon = "docker"
		)

		buffer.WriteString(index_4__114)
		WriteEscString("mdi-"+icon+" mdi mdi-36px", buffer)
		buffer.WriteString(index_4__115)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__117)

	}

	buffer.WriteString(index_4__32)

	{
		var (
			label = T("s_custom_name_part_of", 1)
			name  = "s_custom_name_part_of"
		)

		buffer.WriteString(index_4__118)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_4__120)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__123)

	}

	buffer.WriteString(index_4__33)

	{
		var (
			label = T("s_name", 1)
			name  = "s_name"
		)

		buffer.WriteString(index_4__118)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_4__126)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__117)

	}

	buffer.WriteString(index_4__34)

	{
		var (
			label = T("s_casnumber", 1)
			name  = "s_casnumber"
		)

		buffer.WriteString(index_4__118)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_4__126)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__117)

	}

	buffer.WriteString(index_4__34)

	{
		var (
			label = T("s_empiricalformula", 1)
			name  = "s_empiricalformula"
		)

		buffer.WriteString(index_4__118)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_4__126)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__117)

	}

	buffer.WriteString(index_4__34)

	{
		var (
			label = T("s_storage_barecode", 1)
			name  = "s_storage_barecode"
		)

		buffer.WriteString(index_4__118)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_4__120)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__123)

	}

	buffer.WriteString(index_4__37)

	{
		var (
			label = T("s_signalword", 1)
			name  = "s_signalword"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__38)

	{
		var (
			label = T("s_symbols", 1)
			name  = "s_symbols"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__78)

	}

	buffer.WriteString(index_4__39)

	{
		var (
			label = T("s_hazardstatements", 1)
			name  = "s_hazardstatements"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__78)

	}

	buffer.WriteString(index_4__38)

	{
		var (
			label = T("s_precautionarystatements", 1)
			name  = "s_precautionarystatements"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__78)

	}

	buffer.WriteString(index_4__39)

	{
		var (
			label = T("s_casnumber_cmr", 1)
			name  = "s_casnumber_cmr"
		)

		buffer.WriteString(index_4__165)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__166)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteAll(label, true, buffer)
		buffer.WriteString(index_4__168)
	}

	buffer.WriteString(index_4__42)
	WriteAll(T("clearsearch_text", 1), true, buffer)
	buffer.WriteString(index_4__43)
	WriteAll(T("search_text", 1), true, buffer)
	buffer.WriteString(index_4__44)
	WriteAll(T("advancedsearch_text", 1), true, buffer)
	buffer.WriteString(index_4__45)
	WriteAll(T("switchstorageview_text", 1), true, buffer)
	buffer.WriteString(index_4__46)
	WriteAll(T("export_text", 1), true, buffer)
	buffer.WriteString(index_4__47)
	WriteAll(T("close", 1), true, buffer)
	buffer.WriteString(index_4__48)

	{
		var (
			iconitem   = "border-color"
			iconaction = "tag"
			label      = "update product"
		)

		buffer.WriteString(index_2__58)
		WriteEscString("mdi-"+iconitem+" mdi mdi-48px", buffer)
		buffer.WriteString(index_2__59)
		WriteEscString("mdi-"+iconaction+" mdi mdi-18px", buffer)
		buffer.WriteString(index_2__60)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__61)

	}

	buffer.WriteString(index_4__49)

	{
		var (
			label = "name"
			name  = "name"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__50)

	{
		var (
			label = "synonym(s)"
			name  = "synonyms"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__78)

	}

	buffer.WriteString(index_4__51)

	{
		var (
			label = "specificity"
			name  = "product_specificity"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__64)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__67)
	}

	buffer.WriteString(index_4__52)

	{
		var (
			label = "empirical formula"
			name  = "empiricalformula"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__50)

	{
		var (
			label = "linear formula"
			name  = "linearformula"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__54)

	{
		var (
			label = "cas number"
			name  = "casnumber"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__50)

	{
		var (
			label = "ce number"
			name  = "cenumber"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__51)

	{
		var (
			label = "MSDS link"
			name  = "product_msds"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__64)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__67)
	}

	buffer.WriteString(index_4__39)

	{
		var (
			label = "3D formula link"
			name  = "product_threedformula"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__64)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__67)
	}

	buffer.WriteString(index_4__38)

	{
		var (
			label = "3D formula mol file"
			name  = "product_molformula"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_4__223)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__67)
	}

	{
		var (
			name  = "product_molformula_content"
			value = ""
		)

		buffer.WriteString(index_4__90)
		WriteEscString("hidden_"+name, buffer)
		buffer.WriteString(index_4__91)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__92)
		WriteEscString(value, buffer)
		buffer.WriteString(index_4__93)
	}

	buffer.WriteString(index_4__39)

	{
		var (
			label = "physical state"
			name  = "physicalstate"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__38)

	{
		var (
			label = "class of compound"
			name  = "classofcompound"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__39)

	{
		var (
			label = "signal word"
			name  = "signalword"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__149)

	}

	buffer.WriteString(index_4__38)

	{
		var (
			label = "symbol(s)"
			name  = "symbols"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__78)

	}

	buffer.WriteString(index_4__39)

	{
		var (
			label = "hazard statement(s)"
			name  = "hazardstatements"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__78)

	}

	buffer.WriteString(index_4__38)

	{
		var (
			label = "precautionary statement(s)"
			name  = "precautionarystatements"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_2__76)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__78)

	}

	buffer.WriteString(index_4__51)

	{
		var (
			label = "restricted access"
			name  = "product_restricted"
		)

		buffer.WriteString(index_4__165)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__166)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_4__168)
	}

	buffer.WriteString(index_4__51)

	{
		var (
			label = "radioactive"
			name  = "product_radioactive"
		)

		buffer.WriteString(index_4__165)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__166)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_4__168)
	}

	buffer.WriteString(index_4__51)

	{
		var (
			label = "disposal comment"
			name  = "product_disposalcomment"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_4__271)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__274)

	}

	buffer.WriteString(index_4__51)

	{
		var (
			label = "remark"
			name  = "product_remark"
		)

		buffer.WriteString(index_2__62)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__63)
		WriteEscString(label, buffer)
		buffer.WriteString(index_4__271)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__65)
		WriteEscString(name, buffer)
		buffer.WriteString(index_2__66)
		WriteEscString(name, buffer)
		buffer.WriteString(index_4__274)

	}

	buffer.WriteString(index_4__69)
	WriteAll(T("save", 1), true, buffer)
	buffer.WriteString(index_4__70)
	WriteAll(T("close", 1), true, buffer)
	buffer.WriteString(index_4__71)

	json, _ := json.Marshal(c)

	var out string
	for key, value := range c.URLValues {
		out += fmt.Sprintf("URLValues.set(%s, %s)\n", key, value)
	}

	buffer.WriteString(index__31)
	WriteAll(c.ProxyPath, false, buffer)
	buffer.WriteString(index__32)
	buffer.WriteString(fmt.Sprintf("%s", json))
	buffer.WriteString(index__33)
	buffer.WriteString(out)
	buffer.WriteString(index__34)
	WriteAll(c.ProxyPath+"js/jquery.formautofill.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/jquery.validate.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/jquery.validate.additional-methods.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/select2.full.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/popper.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/bootstrap.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/bootstrap-table.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/bootstrap-confirmation.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/bootstrap-colorpicker.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/bootstrap-toggle.min.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/JSmol.lite.nojq.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/chim/gjs-common.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/chim/chimcommon.js", true, buffer)
	buffer.WriteString(index__35)
	WriteAll(c.ProxyPath+"js/chim/login.js", true, buffer)
	buffer.WriteString(index_4__89)

}
