package locales

var LOCALES_FR = []byte(`
[test]
	one = "Un test"
	other = "Plusieurs tests"

[created]
	one = "créé"
[modified]
	one = "modifié"

[edit]
	one = "editer"
[delete]
	one = "supprimer"
[save]
	one = "enregistrer"
[close]
	one = "fermer"

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

[logo_information1]
	one = "Logo Chimithèque réalisé par "
[logo_information2]
	one = "Ne pas utiliser ou copier sans sa permission."

[s_custom_name_part_of]
	one = "partie du nom"
[s_name]
	one = "nom"
[s_casnumber]
	one = "CAS"
[s_empiricalformula]
	one = "formule brute"
[s_storage_barecode]
	one = "code barre"
[s_signalword]
	one = "mention d'avertissement"
[s_symbols]
	one = "symbole(s)"
[s_hazardstatements]
	one = "mention(s) de danger H-EUH"
[s_precautionarystatements]
	one = "conseil(s) de prudence P"

[menu_home]
	one = "accueil"
[menu_bookmark]
	one = "mes favoris"
[menu_create_productcard]
	one = "créer fiche produit"
[menu_entity]
	one = "entités"
[menu_storelocation]
	one = "entrepôts"
[menu_people]
	one = "utilisateurs"
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

[storage_quantity_title]
	one = "quantité"
[storage_barecode_title]
	one = "code barre"
[storage_batchnumber_title]
	one = "numéro de lot"
[supplier_label_title]
	one = "fournisseur"
[storage_entrydate_title]
	one = "date d'entrée"
[storage_exitdate_title]
	one = "date de sortie"
[storage_openingdate_title]
	one = "date d'ouverture"
[storage_expirationdate_title]
	one = "date d'expiration"
[storage_comment_title]
	one = "commentaire"

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
[linearformula_label_title]
	one = "formule linéaire"
[hazardstatement_label_title]
	one = "mention(s) de danger H-EUH"
[precautionarystatement_label_title]
	one = "conseil(s) de prudence P"
[classofcompound_label_title]
	one = "famille chimique"
[physicalstate_label_title]
	one = "état physique"
[product_threedformula_title]
	one = "formule 3D"
[product_msds_title]
	one = "FDS"
[product_disposalcomment_title]
	one = "commentaire de destruction"
[product_remark_title]
	one = "remarque"
[product_radioactive_title]
	one = "radioactif"
[product_restricted_title]
	one = "accès restreint"
`)
