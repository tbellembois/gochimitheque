ATTACH DATABASE './storage.db' as 'Y';

PRAGMA foreign_keys=off;

INSERT INTO bookmark (
	bookmark_id,
	person,
	product
)
SELECT bookmark_id,
	person,
	product
FROM Y.bookmark;

INSERT INTO borrowing (
	borrowing_id,
	borrowing_comment,
	person,
	borrower,
	storage
)
SELECT borrowing_id,
	borrowing_comment,
	person,
	borrower,
	storage
FROM Y.borrowing;

INSERT INTO cas_number (
	cas_number_id,
	cas_number_label,
	cas_number_cmr
)
SELECT casnumber_id,
	casnumber_label,
	casnumber_cmr
FROM Y.casnumber;

INSERT INTO ce_number (
	ce_number_id,
	ce_number_label
)
SELECT cenumber_id,
	cenumber_label
FROM Y.cenumber;

INSERT INTO category (
	category_id,
	category_label
)
SELECT category_id,
	category_label
FROM Y.category;

INSERT INTO class_of_compound (
	class_of_compound_id,
	class_of_compound_label
)
SELECT classofcompound_id,
	classofcompound_label
FROM Y.classofcompound;

INSERT INTO empirical_formula (
	empirical_formula_id,
	empirical_formula_label
)
SELECT empiricalformula_id,
	empiricalformula_label
FROM Y.empiricalformula;

INSERT INTO linear_formula (
	linear_formula_id,
	linear_formula_label
)
SELECT linearformula_id,
	linearformula_label
FROM Y.linearformula;

INSERT INTO entity (
	entity_id,
	entity_name,
	entity_description
)
SELECT entity_id,
	entity_name,
	entity_description
FROM Y.entity;

INSERT INTO name (
	name_id,
	name_label
)
SELECT name_id,
	name_label
FROM Y.name;

INSERT into entitypeople (
	entitypeople_entity_id,
	entitypeople_person_id
)
SELECT entitypeople_entity_id, entitypeople_person_id
FROM Y.entitypeople;

INSERT INTO hazard_statement (
	hazard_statement_id,
	hazard_statement_label,
	hazard_statement_reference,
	hazard_statement_cmr
)
SELECT hazardstatement_id,
	hazardstatement_label,
	hazardstatement_reference,
	hazardstatement_cmr
FROM Y.hazardstatement;

INSERT INTO permission (
	permission_id,
	person,
	permission_name,
	permission_item,
	permission_entity
)
SELECT permission_id,
	person,
	permission_perm_name,
	permission_item_name,
	permission_entity_id
FROM Y.permission;

DELETE FROM permission WHERE permission_item = 'people';

-- Step 1: Try to insert normally
INSERT OR IGNORE INTO person (person_id, person_email)
SELECT person_id, LOWER(person_email)
FROM Y.person;

-- Step 2: Insert conflicting rows with modified email
INSERT INTO person (person_id, person_email)
SELECT
    person_id,
    LOWER(person_email) || '_' || person_id
FROM Y.person
WHERE LOWER(person_email) IN (
    SELECT LOWER(person_email)
    FROM Y.person
    GROUP BY LOWER(person_email)
    HAVING COUNT(*) > 1
)
AND person_id NOT IN (
    SELECT person_id FROM person
);

INSERT INTO personentities (
	personentities_person_id,
	personentities_entity_id
)
SELECT personentities_person_id, personentities_entity_id
FROM Y.personentities;

INSERT into physical_state (
	physical_state_id,
	physical_state_label
)
SELECT physicalstate_id,
	physicalstate_label
FROM Y.physicalstate;

INSERT INTO precautionary_statement (
	precautionary_statement_id,
	precautionary_statement_label,
	precautionary_statement_reference
)
SELECT precautionarystatement_id,
	precautionarystatement_label,
	precautionarystatement_reference
FROM Y.precautionarystatement;

-- Step 1: Try to insert normally
INSERT OR IGNORE INTO producer (producer_id, producer_label)
SELECT producer_id, producer_label
FROM Y.producer;

-- Step 2: Insert conflicting rows with modified label
INSERT INTO producer (producer_id, producer_label)
SELECT
    producer_id,
    producer_label || '_' || producer_id
FROM Y.producer
WHERE producer_label IN (
    SELECT producer_label
    FROM Y.producer
    GROUP BY producer_label
    HAVING COUNT(*) > 1
)
AND producer_id NOT IN (
    SELECT producer_id FROM producer
);

INSERT INTO producer_ref (
	producer_ref_id,
	producer_ref_label,
	producer
)
SELECT producerref_id,
	producerref_label,
	producer
FROM Y.producerref;

INSERT into product (
	product_id,
	product_type,
	product_specificity,
	product_msds,
	product_restricted,
	product_radioactive,
	product_threed_formula,
	product_twod_formula,
	product_disposal_comment,
	product_remark,
	product_qrcode,
	product_sheet,
	product_concentration,
	product_temperature,
	cas_number,
	ce_number,
	person,
	empirical_formula,
	linear_formula,
	physical_state,
	signal_word,
	name,
	producer_ref,
	unit_temperature,
	category,
	product_number_per_carton,
	product_number_per_bag
)
SELECT product_id,
	random(),
	product_specificity,
	product_msds,
	product_restricted,
	product_radioactive,
	product_threedformula,
	product_twodformula,
	product_disposalcomment,
	product_remark,
	product_qrcode,
	product_sheet,
	product_concentration,
	product_temperature,
	casnumber,
	cenumber,
	person,
	empiricalformula,
	linearformula,
	physicalstate,
	signalword,
	name,
	producerref,
	unit_temperature,
	category,
	product_number_per_carton,
	product_number_per_bag
FROM Y.product;

UPDATE product SET product_type = 'cons' WHERE (product_number_per_carton IS NOT NULL AND product_number_per_carton != 0);
UPDATE product SET product_type = 'bio' WHERE (producer_ref IS NOT NULL AND (product_number_per_carton IS NULL OR product_number_per_carton == 0));
UPDATE product SET product_type = 'chem' WHERE (producer_ref IS NULL AND (product_number_per_carton IS NULL OR product_number_per_carton == 0));

INSERT INTO productclassesofcompounds (
	productclassesofcompounds_product_id,
	productclassesofcompounds_class_of_compound_id
)
SELECT productclassofcompound_product_id,
productclassofcompound_classofcompound_id
FROM Y.productclassofcompound;

INSERT INTO producthazardstatements (
	producthazardstatements_product_id,
	producthazardstatements_hazard_statement_id
)
SELECT producthazardstatements_product_id,
producthazardstatements_hazardstatement_id
FROM Y.producthazardstatements;

INSERT INTO productprecautionarystatements (
	productprecautionarystatements_product_id,
	productprecautionarystatements_precautionary_statement_id
)
SELECT productprecautionarystatements_product_id,
productprecautionarystatements_precautionarystatement_id
FROM Y.productprecautionarystatements;

INSERT INTO productsupplierrefs (
	productsupplierrefs_product_id,
	productsupplierrefs_supplier_ref_id
)
SELECT productsupplierrefs_product_id,
productsupplierrefs_supplierref_id
FROM Y.productsupplierrefs;

INSERT INTO productsymbols (
	productsymbols_product_id,
	productsymbols_symbol_id
)
SELECT productsymbols_product_id,
productsymbols_symbol_id
FROM Y.productsymbols;

INSERT INTO productsynonyms (
	productsynonyms_product_id,
	productsynonyms_name_id
)
SELECT productsynonyms_product_id,
productsynonyms_name_id
FROM Y.productsynonyms;

INSERT INTO producttags (
	producttags_product_id,
	producttags_tag_id
)
SELECT producttags_product_id, producttags_tag_id
FROM Y.producttags;

INSERT INTO signal_word (
	signal_word_id,
	signal_word_label
)
SELECT signalword_id,
	signalword_label
FROM Y.signalword;

INSERT INTO storage (
	storage_id,
	storage_creation_date,
	storage_modification_date,
	storage_entry_date,
	storage_exit_date,
	storage_opening_date,
	storage_expiration_date,
	storage_quantity,
	storage_barecode,
	storage_comment,
	storage_reference,
	storage_batch_number,
	storage_to_destroy,
	storage_archive,
	storage_qrcode,
	storage_concentration,
	storage_number_of_bag,
	storage_number_of_carton,
	person,
	product,
	store_location,
	unit_concentration,
	unit_quantity,
	supplier,
	storage
)
SELECT storage_id,
CAST(unixepoch(storage_creationdate) AS INTEGER) AS storage_creation_date,
CAST(unixepoch(storage_modificationdate) AS INTEGER) AS storage_modification_date,
CAST(unixepoch(storage_entrydate) AS INTEGER) AS storage_entry_date,
CAST(unixepoch(storage_exitdate) AS INTEGER) AS storage_exit_date,
CAST(unixepoch(storage_openingdate) AS INTEGER) AS storage_opening_date,
CAST(unixepoch(storage_expirationdate) AS INTEGER) AS storage_expiration_date,
storage_quantity,
storage_barecode,
storage_comment,
storage_reference,
storage_batchnumber,
storage_todestroy,
storage_archive,
storage_qrcode,
storage_concentration,
storage_number_of_bag,
storage_number_of_carton,
person,
product,
storelocation,
unit_concentration,
unit_quantity,
supplier,
storage
FROM Y.storage;

UPDATE storage SET storage_to_destroy = 0 WHERE storage_to_destroy is NULL;

INSERT INTO store_location (
	store_location_id,
	store_location_name,
	store_location_color,
	store_location_can_store,
	store_location_full_path,
	entity,
	store_location
)
SELECT storelocation_id,
	storelocation_name,
	storelocation_color,
	storelocation_canstore,
	storelocation_fullpath,
	entity,
	storelocation
FROM Y.storelocation;

-- Step 1: Try to insert normally
INSERT OR IGNORE INTO supplier (supplier_id, supplier_label)
SELECT supplier_id, supplier_label
FROM Y.supplier;

-- Step 2: Insert conflicting rows with modified label
INSERT INTO supplier (supplier_id, supplier_label)
SELECT
    supplier_id,
    supplier_label || '_' || supplier_id
FROM Y.supplier
WHERE supplier_label IN (
    SELECT supplier_label
    FROM Y.supplier
    GROUP BY supplier_label
    HAVING COUNT(*) > 1
)
AND supplier_id NOT IN (
    SELECT supplier_id FROM supplier
);

INSERT INTO supplier_ref (
	supplier_ref_id,
	supplier_ref_label,
	supplier
)
SELECT supplierref_id,
	supplierref_label,
	supplier
FROM Y.supplierref;

INSERT INTO symbol (symbol_label) VALUES ('GHS01'), ('GHS02'), ('GHS03'), ('GHS04'), ('GHS05'), ('GHS06'), ('GHS07'), ('GHS08'), ('GHS09');

INSERT INTO tag (
	tag_id,
	tag_label
)
SELECT tag_id,
	tag_label
FROM Y.tag;

INSERT INTO unit (
	unit_id,
	unit_label,
	unit_multiplier,
	unit_type,
	unit
)
SELECT unit_id,
	unit_label,
	unit_multiplier,
	unit_type,
	unit
FROM Y.unit;

INSERT INTO unit (unit_label, unit_multiplier, unit_type) VALUES ('g/mol', 1, 'molecular_weight');

INSERT INTO welcome_announce (
	welcome_announce_id,
	welcome_announce_text
)
SELECT welcomeannounce_id,
	welcomeannounce_text
FROM Y.welcomeannounce;

PRAGMA user_version=10;
PRAGMA foreign_keys=on;
