package locales

var LOCALES_FR = []byte(`
[test]
	one = "Un test"
	other = "Plusieurs tests"

[nil]
	one = " "

[project_leader]
	one = "chef de project"
[project_site]
	one = "site web"
[project_support]
	one = "support/rapport d'erreurs"
[project_license]
	one = "license"
[project_version]
	one = "version"

[wasm_loading]
	one = "Chargement d'un module logiciel."
[wasm_loaded]
	one = "Module logiciel chargé."

[created]
	one = "créé"
[modified]
	one = "modifié"
[selected]
	one = "sélectionné"
[select]
	one = "sélectionner"
[select_all]
	one = "sélectionner tout"
[none]
	one = "aucun"
[nocamera]
	one = "pas de caméra détectée"

[confirm]
	one = "confirmer"
[edit]
	one = "editer"
[delete]
	one = "supprimer"
[archive]
	one = "archiver"
[save]
	one = "enregistrer"
[close]
	one = "fermer"
[list]
	one = "lister"
[create]
	one = "créer"
[check_all]
	one = "sélectionner tout"
[import]
	one = "importer"
[replace]
	one = "remplacer"

[required_input]
	one = "champs requis"
[error_occurred]
	one = "une erreur est survenue"
[no_result]
	one = "pas de résultat"
[no_item]
	one = "pas d'élément"
[active_filter]
	one = "filtre(s) actif"
[no_filter]
	one = "pas de filtre"
[remove_filter]
	one = "supprimer le filtre"

[empirical_formula_convert]
	one = "convertir en formule brute"
[empirical_formula_convert_failed]
	one = "impossible de trouver la formule brute"
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
[invalid_email]
	one = "adresse email invalide"
[not_same_password]
	one = "vous n'avez pas saisi le même mot de passe"

[members]
	one = "membres"
[store_locations]
	one = "entrepôts"

[magical_selector]
	one = "selecteur magique"

[nb_duplicate]
	one = "nombre d'éléments (bouteilles, boites...)"
[nb_duplicate_comment]
	one = "créera une fiche stockage par élément avec un code barre différent (sauf si \"code barre identique\" est coché)"
[identical_barecode]
	one = "code barre identique"
[identical_barecode_comment]
	one = "génère le même code barre pour chaque fiche de stockage - scanner le qrcode d'une fiche stockage retournera aussi les stockages avec un code barre identique"

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

[login_openid]
	one = "Connexion OpenID"
[login_openid_text]
	one = "Cliquez sur 'Je m'identifie' ci-dessous"
[login_openid_link]
	one = "Je m'identifie"
[login_forbidden]
	one = "Votre compte n'a pas été encore validé par un gestionnaire d'entité"

[qrcode_login]
	one = "connexion par QRCode"
[email_placeholder]
	one = "entrez votre email"
[submitlogin_text]
	one = "entrer"
[password_placeholder]
	one = "entrez votre mot de passe"
[resetpassword_alert_ldap]
	one = "vous utilisez un compte institutionnel qui ne peut être changé dans Chimithèque"
[resetpassword_text]
	one = "réinitialiser mon mot de passe"
[resetpassword2_text]
	one = "réinitialiser mon mot de passe, je ne suis pas un robot"
[resetpassword_message_mailsentto]
	one = "un mail de réinitialisation a été envoyé à %s"
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
	Cliquez sur ce lien pour réinitialiser votre mot de passe : %sreset?token=%s

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
	one = "logo Chimithèque réalisé par "
[logo_information2]
	one = "Ne pas utiliser ou copier sans sa permission."

[welcome_announce_text_title]
	one = "Texte complémentaire de la page d'accueil"
[welcome_announce_text_modificationsuccess]
	one = "Texte modifié"

[s_tags]
	one = "tag(s)"
[s_category]
	one = "categorie"
[s_entity]
	one = "entité"
[s_store_location]
	one = "entrepôt"
[s_custom_name_part_of]
	one = "partie du nom"
[s_producer_ref]
	one = "numéro référence fabriquant"
[s_name]
	one = "nom ou synonyme exact"
[s_cas_number]
	one = "numéro CAS"
[s_empirical_formula]
	one = "formule brute"
[s_storage_barecode]
	one = "code barre"
[s_signal_word]
	one = "mention d'avertissement"
[s_symbols]
	one = "pictogramme(s)"
[s_hazard_statements]
	one = "mention(s) de danger H-EUH"
[s_precautionary_statements]
	one = "conseil(s) de prudence P"
[s_cas_number_cmr]
	one = "CMR"
[s_borrowing]
	one = "stockages empruntés"
[s_storage_to_destroy]
	one = "stockages à détruire"

[menu_home]
	one = "accueil"
[menu_bookmark]
	one = "favoris"
[menu_scanqr]
	one = "scanner"
[menu_borrow]
	one = "mes produits empruntés"
[menu_create_productcard]
	one = "créer une fiche"
[menu_entity]
	one = "entités"
[menu_store_location]
	one = "entrepôts"
[menu_people]
	one = "utilisateurs"
[menu_welcome_announce]
	one = "message de page d'accueil"
[menu_password]
	one = "changer mon mot de passe"
[menu_qrcode]
	one = "gérer mon accès par QRCode"
[menu_logout]
	one = "déconnexion"
[menu_about]
	one = "à propos"
[menu_account]
	one = "mon compte"
[menu_management]
	one = "gestion"
[menu_settings]
	one = "paramètres"
[menu_pubchem]
	one = "pubchem"

[clearsearch_text]
	one = "supprimer tous les filtres"
[search_text]
	one = "rechercher"
[advancedsearch_text]
	one = "recherche avancée"

[chemical_product]
	one = "produit chimique"
[biological_product]
	one = "réactif biologique"
[consumable_product]
	one = "consommable"
[all_product]
	one = "tous les produits"

[switchproductview_text]
	one = "vue par produits"
[switchstorageview_text]
	one = "vue par stockages"
[export_text]
	one = "exporter"
[download_export]
	one = "télécharger l'export"
[export_progress]
	one = "export en cours -  cette opération peut être longue"
[export_done]
	one = "export effectué"
[showdeleted_text]
	one = "voir archives"
[hidedeleted_text]
	one = "cacher archives"
[storeagain_text]
	one = "stocker ce produit"
[totalstock_text]
	one = "calculer le stock total"

[unit_label_title]
	one = "unité"
[supplier_label_title]
	one = "fournisseur"
[supplier_ref_label_title]
	one = "numéro référence fournisseur"

[add_producer_title]
	one = "ajouter un fabriquant à la liste"
[producer_added]
	one = "fabriquant ajouté"
[add_supplier_title]
	one = "ajouter un fournisseur à la liste"
[supplier_added]
	one = "fournisseur ajouté"

[store]
	one = "stocker"
[storages]
	one = "stockages"
[storage]
	one = "stockage"
[archives]
	one = "archives"
[ostorages]
	one = "disponibilité"
[storage_create_title]
	one = "stocker un produit"
[storage_update_title]
	one = "mise à jour d'un stockage"
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
[storage_restored_message]
	one = "stockage restauré"
[storage_trashed_message]
	one = "stockage mis à la corbeille"
[storage_deleted_message]
	one = "storage supprimé"
[storage_borrow_updated]
	one = "emprunt mis à jour"
[storage_created_message]
	one = "stockage créé"
[storage_updated_message]
	one = "stockage mis à jour"

[storage_store_location_title]
	one = "entrepôt"
[storage_concentration_title]
	one = "concentration"
[storage_quantity_title]
	one = "quantité"
[storage_barecode_title]
	one = "code barre"
[storage_create_barecode_comment]
	one = "si vous laissez ce champs vide, un code barre sera autogénéré"
[storage_batch_number_title]
	one = "lot"
[storage_entry_date_title]
	one = "date d'entrée"
[storage_exit_date_title]
	one = "date de sortie"
[storage_opening_date_title]
	one = "date d'ouverture"
[storage_expiration_date_title]
	one = "date d'expiration"
[storage_borrower_title]
	one = "emprunteur"
[storage_comment_title]
	one = "commentaire"
[storage_reference_title]
	one = "référence"
[storage_to_destroy_title]
	one = "à détruire"
[storage_product_table_header]
	one = "produit"
[storage_store_location_table_header]
	one = "entrepôt"
[storage_quantity_table_header]
	one = "quantité"
[storage_barecode_table_header]
	one = "code barre"
[storage_store_location_placeholder]
	one = "selectionnez un entrepôt"
[storage_borrower_placeholder]
	one = "selectionnez un emprunteur"
[storage_supplier_placeholder]
	one = "selectionnez ou entrez un fournisseur"
[storage_print_qrcode]
	one = "imprimer le qrcode"
[storage_number_of_unit]
	one = "nombre d'unité(s)"
[storage_number_of_bag]
	one = "nombre de sac(s)"
[storage_number_of_bag_comment]
	one = "seulement si le nombre d'unités par sachet est défini pour le produit"
[storage_number_of_carton]
	one = "nombre de carton(s)"
[storage_number_of_carton_comment]
	one = "seulement si le nombre d'unités par carton est défini pour le produit"
[storage_one_number_required]
	one = "au moins un des nombres requis"

[stock_store_location_title]
	one = "dans cet entrepôt"
[stock_store_location_sub_title]
	one = "avec les sous entrepôts"

[empirical_formula_label_title]
	one = "formule brute"
[ce_number_label_title]
	one = "CE"
[cas_number_label_title]
	one = "CAS"
[cas_number_cmr_title]
	one = "CMR"
[signal_word_label_title]
	one = "mention d'avertissement"
[symbol_label_title]
	one = "symbole(s)"
[linear_formula_label_title]
	one = "formule linéaire"
[hazard_statement_label_title]
	one = "mention(s) de danger H-EUH"
[precautionary_statement_label_title]
	one = "conseil(s) de prudence P"
[class_of_compound_label_title]
	one = "famille(s) chimique(s)"
[physical_state_label_title]
	one = "état physique"
[name_label_title]
	one = "nom"
[synonym_label_title]
	one = "synonyme(s)"

[search_pubchem_title]
	one = "rechercher par nom (ou une partie)"

[restricted]
	one = "accès restreint"
[bookmark]
	one = "ajouter aux favoris"
[unbookmark]
	one = "retirer des favoris"
[product_create_title]
	one = "créer une fiche produit"
[product_update_title]
	one = "mettre à jour produit"
[product_inchi_title]
	one = "code inchi"
[product_smiles_title]
	one = "SMILES"
[product_molecular_weight_title]
	one = "poids moléculaire"
[product_inchi_key_title]
	one = "clé inchi"
[product_threed_formula_title]
	one = "formule 3D"
[product_twod_formula_title]
	one = "image molécule"
[product_threed_formula_mol_title]
	one = "fichier MOL formule 3D"
[product_msds_title]
	one = "lien FDS"
[product_sheet_title]
	one = "fiche produit fabriquant"
[product_temperature_title]
	one = "température de stockage préconisée"
[product_number_per_carton_title]
	one = "nombre d'unités par carton"
[product_number_per_bag_title]
	one = "nombre d'unités par sachet"
[producer_label_title]
	one = "fabriquant"
[producer_ref_label_title]
	one = "numéro référence fabriquant"
[producer_ref_create_needproducer]
	one = "pour créer une nouvelle référence sélectionnez un fabriquant d'abord"
[supplier_ref_create_needsupplier]
	one = "pour créer une nouvelle référence sélectionnez un fournisseur d'abord"
[category_label_title]
	one = "catégorie"
[tag_label_title]
	one = "tag(s)"
[product_disposal_comment_title]
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
[product_empirical_formula_table_header]
	one = "formule br."
[product_cas_table_header]
	one = "CAS"
[product_specificity_table_header]
	one = "spéc."
[product_cas_placeholder]
	one = "sélectionnez ou entrez un numéro CAS"
[product_ce_placeholder]
	one = "sélectionnez ou entrez un numéro CE"
[product_physical_state_placeholder]
	one = "sélectionnez ou entrez un état physique"
[product_signal_word_placeholder]
	one = "sélectionnez une mention"
[product_class_of_compound_placeholder]
	one = "sélectionnez ou entrez une ou plusieurs famille(s)"
[product_name_placeholder]
	one = "sélectionnez ou entrez un nom"
[product_synonyms_placeholder]
	one = "sélectionnez ou entrez un ou plusieurs nom(s)"
[producer_ref_placeholder]
	one = "sélectionnez ou entrez une référence"
[product_empirical_formula_placeholder]
	one = "sélectionnez ou entrez une formule"
[product_linear_formula_placeholder]
	one = "sélectionnez ou entrez une formule"
[product_symbols_placeholder]
	one = "sélectionnez un ou plusieurs symbole(s)"
[product_hazard_statements_placeholder]
	one = "sélectionnez une ou plusieurs mention(s)"
[product_precautionary_statements_placeholder]
	one = "sélectionnez un ou plusieurs conseil(s)"
[product_producer_placeholder]
	one = "sélectionnez un fabriquant"
[product_producer_ref_placeholder]
	one = "sélectionnez ou entrez une référence fabriquant"
[product_supplier_placeholder]
	one = "sélectionnez un fournisseur"
[product_supplier_ref_placeholder]
	one = "sélectionnez ou entrez une ou plusieurs référence(s) fournisseur"
[product_tag_placeholder]
	one = "sélectionnez ou entrez un ou plusieurs tag(s)"
[product_category_placeholder]
	one = "sélectionnez ou entrez une catégorie"
[product_unit_placeholder]
	one = "sélectionnez une unité"
[product_deleted_message]
	one = "produit supprimé"
[product_updated_message]
	one = "produit mis à jour"
[product_created_message]
	one = "produit créé"
[product_flammable]
	one = "inflammable"

[person_create_title]
	one = "créer une personne"
[person_update_title]
	one = "mettre à jour personne"
[person_deleted_message]
	one = "personne supprimée"
[person_email_title]
	one = "mail"
[person_password_title]
	one = "mot de passe"
[person_entity_title]
	one = "entité(s)"
[person_permission_title]
	one = "permissions"
[person_email_table_header]
	one = "mail"
[person_can_not_remove_entity_manager]
	one = "cette entité ne peut pas être supprimée, l'utilisateur est un de ses gestionnaires"
[person_created_message]
	one = "personne crée"
[person_updated_message]
	one = "personne mise à jour"
[person_password_updated_message]
	one = "mot de passe mis à jour"
[person_entity_placeholder]
	one = "sélectionnez une ou plusieurs entité(s)"
[person_select_all_none_storage]
	one = "sélectionner tous les 'aucune permission'"
[person_select_all_r_storage]
	one = "sélectionner tous les 'voir seulement'"
[person_select_all_rw_storage]
	one = "sélectionner tous les 'voir, modifier, créer et supprimer'"
[person_show_password]
  	one = "afficher le champs mot de passe"
[person_qrcode_regenerate]
	one = "générer un nouveau QRcode"

[permission_product]
	one = "produits"
[permission_rproduct]
	one = "produits restreints"
[permission_storages]
	one = "stockages"
[permission_none]
	one = "aucune permission"
[permission_read]
	one = "voir seulement"
[permission_crud]
	one = "voir, modifier, créer et supprimer"

[store_location_create_title]
	one = "créer un entrepôt"
[store_location_update_title]
	one = "mettre à jour entrepôt"
[store_location_deleted_message]
	one = "entrepôt supprimé"
[store_location_created_message]
	one = "entrepôt créé"
[store_location_updated_message]
	one = "entrepôt mis à jour"
[store_location_parent_title]
	one = "parent"
[store_location_entity_title]
	one = "entité"
[store_location_can_store_title]
	one = "peut stocker"
[store_location_color_title]
	one = "couleur"
[store_location_name_title]
	one = "nom"
[store_location_name_table_header]
	one = "nom"
[store_location_entity_table_header]
	one = "entité"
[store_location_color_table_header]
	one = "couleur"
[store_location_can_store_table_header]
	one = "peut stocker"
[store_location_parent_table_header]
	one = "parent"
[store_location_entity_placeholder]
	one = "sélectionnez une entité"
[store_location_store_location_placeholder]
	one = "sélectionnez une entité d'abord"

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
	one = "gestionnaires(s)"
[entity_manager_placeholder]
	one = "sélectionnez un ou plusieurs gestionnaires"
[entity_ldap_group]
	one = "groupe(s) LDAP pour enrollement automatique"
[entity_ldap_group_placeholder]
	one = "sélectionnez un ou plusieurs groupes LDAP"

[entity_nameexist_validate]
	one = "une entité avec ce nom existe déjà"
[person_emailexist_validate]
	one = "une personne avec cet email existe déjà"
[empirical_formula_validate]
	one = "formule brute invalide"
[cas_number_validate_wrongcas]
	one = "numéro CAS invalide"
[cas_number_validate_casspecificity]
	one = "le couple numéro CAS/spécificité existe déjà"
[ce_number_validate]
	one = "numéro CE invalide"
`)
