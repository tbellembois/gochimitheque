package locales

var LOCALES_FR = []byte(`
[test]
	one = "Un test"
	other = "Plusieurs tests"

[nil]
	one = " "

[created]
	one = "créé"
[modified]
	one = "modifié"
[select_all]
	one = "sélectionner tout"
[none]
	one = "aucun"

[edit]
	one = "editer"
[delete]
	one = "supprimer"
[save]
	one = "enregistrer"
[close]
	one = "fermer"
[list]
	one = "lister"
[create]
	one = "créer"

[required_input]
	one = "champs requis"

[no_empirical_formula]
	one = "pas de formule brute"
[no_cas_number]
	one = "pas de numéro CAS"
[howto_magicalselector]
	one = "comment utiliser le sélecteur magique"

[password]
	one = "mot de passe"
[confirm_password]
	one = "confirmer le mot de passe"
[invalid_password]
	one = "mauvais mot de passe"

[members]
	one = "membres"
[storelocations]
	one = "entrepôts"

[magical_selector]
	one = "selecteur magique"

[nb_duplicate]
	one = "nombre de duplications"

[bt_loadingMessage]
	one = "chargement..."
[bt_recordsPerPage]
	one = "enegistrements par page"
[bt_showingRowsTotal]
	one = "enregistrements total"
[bt_search]
	one = "rechercher"
[bt_noMatches]
	one = "pas de résultat"

[email_placeholder]
	one = "entrez votre email"
[submitlogin_text]
	one = "entrer"
[password_placeholder]
	one = "entrez votre mot de passe"
[resetpassword_text]
	one = "réinitialiser mon mot de passe"
[resetpassword2_text]
	one = "réinitialiser mon mot de passe, je ne suis pas un robot"
[resetpassword_warning_enteremail]
	one = "entrez votre adresse mail dans le formulaire"
[resetpassword_message_mailsentto]
	one = "un mail de réinitialisation a été envoyé à"
[resetpassword_areyourobot]
	one = "êtes vous un robot ?"
[resetpassword_mailbody1]
	one = '''
	Voici votre mot de passe temporaire pour Chimithèque : %s

	Vous pouvez le changer dans l'application.
	'''
[resetpassword_mailsubject1]
	one = "Chimithèque nouveau mot de passe temporaire\r\n"
[resetpassword_mailbody2]
	one = '''
	Cliquez sur ce lien pour réinitialiser votre mot de passe : %s%sreset?token=%s

	Vous recevrez ensuite un nouveau mail avec un mot de passe temporaire.
	'''
[resetpassword_mailsubject2]
	one = "Chimithèque lien de réinitialisation de mot de passe\r\n"
[resetpassword_done]
	one = "Un nouveau mot de passe temporaire a été envoyé à %s"

[createperson_mailsubject]
	one = "Chimithèque nouveau compte\r\n"
[createperson_mailbody]
	one = '''
	Un compte Chimithèque a été créé pour vous.

	Vous pouvez maintenant initialiser votre mot de passe.

	Rendez vous sur la page de connexion %s, entrez votre adresse mail %s et cliquez sur le lien "réinitialiser mon mot de passe ".

	Vous recevrez ensuite un mot de passe temporaire.
	'''

[logo_information1]
	one = "Logo Chimithèque réalisé par "
[logo_information2]
	one = "Ne pas utiliser ou copier sans sa permission."

[welcomeannounce_text_title]
	one = "Texte complémentaire de la page d'accueil"
[welcomeannounce_text_modificationsuccess]
	one = "Texte modifié"

[s_custom_name_part_of]
	one = "partie du nom"
[s_name]
	one = "nom exact"
[s_casnumber]
	one = "CAS"
[s_empiricalformula]
	one = "formule brute"
[s_storage_barecode]
	one = "code barre"
[s_signalword]
	one = "mention d'avertissement"
[s_symbols]
	one = "pictogramme(s)"
[s_hazardstatements]
	one = "mention(s) de danger H-EUH"
[s_precautionarystatements]
	one = "conseil(s) de prudence P"
[s_casnumber_cmr]
	one = "CMR"

[menu_home]
	one = "accueil"
[menu_bookmark]
	one = "mes favoris"
[menu_borrow]
	one = "mes produits empruntés"
[menu_create_productcard]
	one = "créer fiche produit"
[menu_entity]
	one = "entités"
[menu_storelocation]
	one = "entrepôts"
[menu_people]
	one = "utilisateurs"
[menu_welcomeannounce]
	one = "changer le message de page d'accueil"
[menu_password]
	one = "changer mon mot de passe"
[menu_logout]
	one = "déconnexion"

[clearsearch_text]
	one = "effacer le formulaire"
[search_text]
	one = "rechercher"
[advancedsearch_text]
	one = "recherche avancée"

[switchproductview_text]
	one = "vue par produits"
[switchstorageview_text]
	one = "vue par stockages"
[export_text]
	one = "exporter"
[showdeleted_text]
	one = "voir supprimés"
[hidedeleted_text]
	one = "cacher supprimés"
[storeagain_text]
	one = "stocker ce produit"
[totalstock_text]
	one = "afficher le stock total"

[unit_label_title]
	one = "unité"
[supplier_label_title]
	one = "fournisseur"

[storage_create_title]
	one = "créer stockage"
[storage_update_title]
	one = "mise à jour stockage"
[storage_clone]
	one = "cloner"
[storage_borrow]
	one = "emprunter"
[storage_unborrow]
	one = "restituer"
[storage_restore]
	one = "restaurer"
[storage_showhistory]
	one = "voir historique"
[storage_history]
	one = "historique"

[storage_storelocation_title]
	one = "entrepôt"
[storage_quantity_title]
	one = "quantité"
[storage_barecode_title]
	one = "code barre"
[storage_create_barecode_comment]
	one = "si vous laissez ce champs vide, un code barre sera autogénéré"
[storage_batchnumber_title]
	one = "numéro de lot"
[storage_entrydate_title]
	one = "date d'entrée"
[storage_exitdate_title]
	one = "date de sortie"
[storage_openingdate_title]
	one = "date d'ouverture"
[storage_expirationdate_title]
	one = "date d'expiration"
[storage_borrower_title]
	one = "emprunteur"
[storage_comment_title]
	one = "commentaire"
[storage_reference_title]
	one = "référence"
[storage_todestroy_title]
	one = "à détruire"
[storage_product_table_header]
	one = "produit"
[storage_storelocation_table_header]
	one = "entrepôt"
[storage_quantity_table_header]
	one = "quantité"
[storage_barecode_table_header]
	one = "code barre"

[stock_storelocation_title]
	one = "dans cet entrepôt"
[stock_storelocation_sub_title]
	one = "en incluant les sous entrepôts"

[empiricalformula_label_title]
	one = "formule brute"
[cenumber_label_title]
	one = "CE"
[casnumber_label_title]
	one = "CAS"
[casnumber_cmr_title]
	one = "CMR"
[signalword_label_title]
	one = "mention d'avertissement"
[symbol_label_title]
	one = "symbole(s)"
[linearformula_label_title]
	one = "formule linéaire"
[hazardstatement_label_title]
	one = "mention(s) de danger H-EUH"
[precautionarystatement_label_title]
	one = "conseil(s) de prudence P"
[classofcompound_label_title]
	one = "famille(s) chimique(s)"
[physicalstate_label_title]
	one = "état physique"
[name_label_title]
	one = "nom"
[synonym_label_title]
	one = "synonyme(s)"

[product_create_title]
	one = "créer produit"
[product_update_title]
	one = "mettre à jour produit"
[product_threedformula_title]
	one = "formule 3D"
[product_threedformula_mol_title]
	one = "fichier MOL formule 3D"
[product_msds_title]
	one = "FDS"
[product_disposalcomment_title]
	one = "commentaire de destruction"
[product_remark_title]
	one = "remarque"
[product_specificity_title]
	one = "spécificité"
[product_radioactive_title]
	one = "radioactif"
[product_restricted_title]
	one = "accès restreint"
[product_name_table_header]
	one = "nom"
[product_empiricalformula_table_header]
	one = "formule br."
[product_cas_table_header]
	one = "CAS"
[product_specificity_table_header]
	one = "spéc."
[product_cas_placeholder]
	one = "sélectionnez ou entrez un numéro CAS"
[product_ce_placeholder]
	one = "sélectionnez ou entrez un numéro CE"
[product_physicalstate_placeholder]
	one = "sélectionnez un état physique"
[product_signalword_placeholder]
	one = "sélectionnez une mention"
[product_classofcompound_placeholder]
	one = "sélectionnez une ou plusieurs famille(s)"
[product_name_placeholder]
	one = "sélectionnez ou entrez un nom"
[product_synonyms_placeholder]
	one = "sélectionnez ou entrez un ou plusieurs nom(s)"
[product_empiricalformula_placeholder]
	one = "sélectionnez ou entrez une formule"
[product_linearformula_placeholder]
	one = "sélectionnez ou entrez une formule"
[product_symbols_placeholder]
	one = "sélectionnez un ou plusieurs symbole(s)"
[product_hazardstatements_placeholder]
	one = "sélectionnez une ou plusieurs mention(s)"
[product_precautionarystatements_placeholder]
	one = "sélectionnez un ou plusieurs conseil(s)"

[person_create_title]
	one = "créer personne"
[person_update_title]
	one = "mettre à jour personne"
[person_email_title]
	one = "mail"
[person_password_title]
	one = "mot de passe"
[person_entity_title]
	one = "entité(s)"
[person_email_table_header]
	one = "mail"

[storelocation_create_title]
	one = "créer entrepôt"
[storelocation_update_title]
	one = "mettre à jour entrepôt"
[storelocation_deleted_message]
	one = "entrepôt supprimé"
[storelocation_created_message]
	one = "entrepôt créé"
[storelocation_updated_message]
	one = "entrepôt mis à jour"
[storelocation_parent_title]
	one = "parent"
[storelocation_entity_title]
	one = "entité"
[storelocation_canstore_title]
	one = "peut stocker"
[storelocation_color_title]
	one = "couleur"
[storelocation_name_title]
	one = "nom"
[storelocation_name_table_header]
	one = "nom"
[storelocation_entity_table_header]
	one = "entité"
[storelocation_color_table_header]
	one = "couleur"
[storelocation_canstore_table_header]
	one = "peut stocker"
[storelocation_parent_table_header]
	one = "parent"

[entity_create_title]
	one = "créer une entité"
[entity_update_title]
	one = "mettre à jour entité"
[entity_deleted_message]
	one = "entité supprimée"
[entity_created_message]
	one = "entité crée"
[entity_updated_message]
	one = "entité mise à jour"
[entity_name_table_header]
	one = "nom"
[entity_description_table_header]
	one = "description"
[entity_manager_table_header]
	one = "responsable(s)"
	
[entity_nameexist_validate]
	one = "une entité avec ce nom existe déjà"
[person_emailexist_validate]
	one = "une personne avec cet email existe déjà"
[empiricalformula_validate]
	one = "formule brute invalide"
[casnumber_validate_wrongcas]
	one = "numéro CAS invalide"
[casnumber_validate_casspecificity]
	one = "le couple numéro CAS/spécificité existe déjà"
[cenumber_validate]
	one = "numéro CE invalide"
`)
