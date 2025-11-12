BEGIN TRANSACTION;

DROP TABLE IF EXISTS "bookmark";
CREATE TABLE "bookmark" (
	"bookmark_id"	INTEGER,
	"person"	INTEGER NOT NULL,
	"product"	INTEGER NOT NULL,
	PRIMARY KEY("bookmark_id"),
	FOREIGN KEY("person") REFERENCES "person"("person_id") ON DELETE CASCADE,
	FOREIGN KEY("product") REFERENCES "product"("product_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "borrowing";
CREATE TABLE "borrowing" (
	"borrowing_id"	INTEGER,
	"borrowing_comment"	TEXT,
	"person"	INTEGER NOT NULL,
	"borrower"	INTEGER NOT NULL,
	"storage"	INTEGER NOT NULL UNIQUE,
	PRIMARY KEY("borrowing_id"),
	FOREIGN KEY("borrower") REFERENCES "person"("person_id") ON DELETE CASCADE,
	FOREIGN KEY("person") REFERENCES "person"("person_id") ON DELETE CASCADE,
	FOREIGN KEY("storage") REFERENCES "storage"("storage_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "cas_number";
CREATE TABLE "cas_number" (
	"cas_number_id"	INTEGER,
	"cas_number_label"	TEXT NOT NULL UNIQUE,
	"cas_number_cmr"	TEXT,
	PRIMARY KEY("cas_number_id")
) STRICT;

DROP TABLE IF EXISTS "ce_number";
CREATE TABLE "ce_number" (
	"ce_number_id"	INTEGER,
	"ce_number_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("ce_number_id")
) STRICT;

DROP TABLE IF EXISTS "category";
CREATE TABLE "category" (
	"category_id"	INTEGER,
	"category_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("category_id")
) STRICT;

DROP TABLE IF EXISTS "class_of_compound";
CREATE TABLE "class_of_compound" (
	"class_of_compound_id"	INTEGER,
	"class_of_compound_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("class_of_compound_id")
) STRICT;

DROP TABLE IF EXISTS "empirical_formula";
CREATE TABLE "empirical_formula" (
	"empirical_formula_id"	INTEGER,
	"empirical_formula_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("empirical_formula_id")
) STRICT;

DROP TABLE IF EXISTS "linear_formula";
CREATE TABLE "linear_formula" (
	"linear_formula_id"	INTEGER,
	"linear_formula_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("linear_formula_id")
) STRICT;

DROP TABLE IF EXISTS "hazard_statement";
CREATE TABLE "hazard_statement" (
	"hazard_statement_id"	INTEGER,
	"hazard_statement_label"	TEXT NOT NULL,
	"hazard_statement_reference"	TEXT NOT NULL UNIQUE,
	"hazard_statement_cmr"	TEXT,
	PRIMARY KEY("hazard_statement_id")
) STRICT;

DROP TABLE IF EXISTS "precautionary_statement";
CREATE TABLE "precautionary_statement" (
	"precautionary_statement_id"	INTEGER,
	"precautionary_statement_label"	TEXT NOT NULL,
	"precautionary_statement_reference"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("precautionary_statement_id")
) STRICT;

DROP TABLE IF EXISTS "physical_state";
CREATE TABLE "physical_state" (
	"physical_state_id"	INTEGER,
	"physical_state_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("physical_state_id")
) STRICT;

DROP TABLE IF EXISTS "signal_word";
CREATE TABLE "signal_word" (
	"signal_word_id"	INTEGER,
	"signal_word_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("signal_word_id")
) STRICT;

DROP TABLE IF EXISTS "symbol";
CREATE TABLE "symbol" (
	"symbol_id"	INTEGER,
	"symbol_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("symbol_id")
) STRICT;

DROP TABLE IF EXISTS "tag";
CREATE TABLE "tag" (
	"tag_id"	INTEGER,
	"tag_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("tag_id")
) STRICT;

DROP TABLE IF EXISTS "producer";
CREATE TABLE "producer" (
	"producer_id"	INTEGER,
	"producer_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("producer_id")
) STRICT;

DROP TABLE IF EXISTS "producer_ref";
CREATE TABLE "producer_ref" (
	"producer_ref_id"	INTEGER,
	"producer_ref_label"	TEXT NOT NULL,
	"producer"	INTEGER,
	PRIMARY KEY("producer_ref_id"),
	FOREIGN KEY("producer") REFERENCES "producer"("producer_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "supplier";
CREATE TABLE "supplier" (
	"supplier_id"	INTEGER,
	"supplier_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("supplier_id")
) STRICT;

DROP TABLE IF EXISTS "supplier_ref";
CREATE TABLE "supplier_ref" (
	"supplier_ref_id"	INTEGER,
	"supplier_ref_label"	TEXT NOT NULL,
	"supplier"	INTEGER,
	PRIMARY KEY("supplier_ref_id"),
	FOREIGN KEY("supplier") REFERENCES "supplier"("supplier_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "name";
CREATE TABLE "name" (
	"name_id"	INTEGER,
	"name_label"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("name_id")
) STRICT;

DROP TABLE IF EXISTS "unit";
CREATE TABLE "unit" (
	"unit_id"	INTEGER,
	"unit_label"	TEXT NOT NULL UNIQUE,
	"unit_multiplier"	REAL NOT NULL DEFAULT 1,
	"unit_type"	TEXT,
	"unit"	INTEGER,
	PRIMARY KEY("unit_id"),
	FOREIGN KEY("unit") REFERENCES "unit"("unit_id")
) STRICT;

DROP TABLE IF EXISTS "permission";
CREATE TABLE "permission" (
	"person"	INTEGER NOT NULL,
	"permission_name"	TEXT NOT NULL,
	"permission_item"	TEXT NOT NULL,
	"permission_entity"	INTEGER NOT NULL,
	-- PRIMARY KEY("permission_id"),
	PRIMARY KEY("person", "permission_name", "permission_item", "permission_entity"),
	FOREIGN KEY("person") REFERENCES "person"("person_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "entity";
CREATE TABLE "entity" (
	"entity_id"	INTEGER,
	"entity_name"	TEXT NOT NULL UNIQUE,
	"entity_description"	TEXT,
	PRIMARY KEY("entity_id")
) STRICT;

DROP TABLE IF EXISTS "person";
CREATE TABLE "person" (
	"person_id"	INTEGER,
	"person_email"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("person_id")
) STRICT;

DROP TABLE IF EXISTS "product";
CREATE TABLE "product" (
	"product_id"	INTEGER,
	"product_type"	TEXT NOT NULL,
	"product_inchi"	TEXT,
	"product_inchikey"	TEXT,
	"product_canonical_smiles"	TEXT,
	"product_specificity"	TEXT,
	"product_msds"	TEXT,
	"product_restricted"	INTEGER DEFAULT 0,
	"product_radioactive"	INTEGER DEFAULT 0,
	"product_threed_formula"	TEXT,
	"product_twod_formula"	TEXT,
	"product_disposal_comment"	TEXT,
	"product_remark"	TEXT,
	"product_qrcode"	TEXT,
	"product_sheet"	TEXT,
	"product_concentration"	REAL,
	"product_temperature"	REAL,
	"product_molecular_weight"	REAL,
	"cas_number"	INTEGER,
	"ce_number"	INTEGER,
	"person"	INTEGER NOT NULL DEFAULT 1,
	"empirical_formula"	INTEGER,
	"linear_formula"	INTEGER,
	"physical_state"	INTEGER,
	"signal_word"	INTEGER,
	"name"	INTEGER NOT NULL,
	"producer_ref"	INTEGER,
	"unit_molecular_weight"	INTEGER,
	"unit_temperature"	INTEGER,
	"category"	INTEGER,
	"product_number_per_carton"	INTEGER,
	"product_number_per_bag"	INTEGER,
	PRIMARY KEY("product_id"),
	FOREIGN KEY("cas_number") REFERENCES "cas_number"("cas_number_id"),
	FOREIGN KEY("category") REFERENCES "category"("category_id"),
	FOREIGN KEY("ce_number") REFERENCES "ce_number"("ce_number_id"),
	FOREIGN KEY("empirical_formula") REFERENCES "empirical_formula"("empirical_formula_id"),
	FOREIGN KEY("linear_formula") REFERENCES "linear_formula"("linear_formula_id"),
	FOREIGN KEY("name") REFERENCES "name"("name_id"),
	FOREIGN KEY("person") REFERENCES "person"("person_id") ON DELETE SET DEFAULT,
	FOREIGN KEY("physical_state") REFERENCES "physical_state"("physical_state_id"),
	FOREIGN KEY("producer_ref") REFERENCES "producer_ref"("producer_ref_id"),
	FOREIGN KEY("signal_word") REFERENCES "signal_word"("signal_word_id"),
	FOREIGN KEY("unit_molecular_weight") REFERENCES "unit"("unit_id"),
	FOREIGN KEY("unit_temperature") REFERENCES "unit"("unit_id")
) STRICT;

DROP TABLE IF EXISTS "store_location";
CREATE TABLE "store_location" (
	"store_location_id"	INTEGER,
	"store_location_name"	TEXT NOT NULL,
	"store_location_color"	TEXT,
	"store_location_can_store"	INTEGER DEFAULT 0,
	"store_location_full_path"	TEXT,
	"entity"	INTEGER NOT NULL,
	"store_location"	INTEGER,
	PRIMARY KEY("store_location_id"),
	FOREIGN KEY("entity") REFERENCES "entity"("entity_id"),
	FOREIGN KEY("store_location") REFERENCES "store_location"("store_location_id")
) STRICT;

DROP TABLE IF EXISTS "storage";
CREATE TABLE "storage" (
	"storage_id"	INTEGER,
	"storage_creation_date"	INTEGER NOT NULL DEFAULT current_timestamp,
	"storage_modification_date"	INTEGER NOT NULL DEFAULT current_timestamp,
	"storage_entry_date"	INTEGER,
	"storage_exit_date"	INTEGER,
	"storage_opening_date"	INTEGER,
	"storage_expiration_date"	INTEGER,
	"storage_quantity"	REAL,
	"storage_barecode"	TEXT,
	"storage_comment"	TEXT,
	"storage_reference"	TEXT,
	"storage_batch_number"	TEXT,
	"storage_to_destroy"	INTEGER DEFAULT 0,
	"storage_archive"	INTEGER DEFAULT 0,
	"storage_qrcode"	BLOB,
	"storage_concentration"	REAL,
	"storage_number_of_bag"	INTEGER,
	"storage_number_of_carton"	INTEGER,
	"person"	INTEGER NOT NULL DEFAULT 1,
	"product"	INTEGER NOT NULL,
	"store_location"	INTEGER NOT NULL,
	"unit_concentration"	REAL,
	"unit_quantity"	REAL,
	"supplier"	INTEGER,
	"storage"	INTEGER,
	PRIMARY KEY("storage_id"),
	FOREIGN KEY("person") REFERENCES "person"("person_id") ON DELETE SET DEFAULT,
	FOREIGN KEY("product") REFERENCES "product"("product_id"),
	FOREIGN KEY("storage") REFERENCES "storage"("storage_id"),
	FOREIGN KEY("store_location") REFERENCES "store_location"("store_location_id"),
	FOREIGN KEY("supplier") REFERENCES "supplier"("supplier_id"),
	FOREIGN KEY("unit_concentration") REFERENCES "unit"("unit_id"),
	FOREIGN KEY("unit_quantity") REFERENCES "unit"("unit_id")
) STRICT;

DROP TABLE IF EXISTS "welcome_announce";
CREATE TABLE "welcome_announce" (
	"welcome_announce_id"	INTEGER,
	"welcome_announce_text"	TEXT,
	PRIMARY KEY("welcome_announce_id")
) STRICT;

DROP TABLE IF EXISTS "productclassesofcompounds";
CREATE TABLE "productclassesofcompounds" (
	"productclassesofcompounds_product_id"	INTEGER NOT NULL,
	"productclassesofcompounds_class_of_compound_id"	INTEGER NOT NULL,
	PRIMARY KEY("productclassesofcompounds_product_id","productclassesofcompounds_class_of_compound_id"),
	FOREIGN KEY("productclassesofcompounds_class_of_compound_id") REFERENCES "class_of_compound"("class_of_compound_id") ON DELETE CASCADE,
	FOREIGN KEY("productclassesofcompounds_product_id") REFERENCES "product"("product_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "producthazardstatements";
CREATE TABLE "producthazardstatements" (
	"producthazardstatements_product_id"	INTEGER NOT NULL,
	"producthazardstatements_hazard_statement_id"	INTEGER NOT NULL,
	PRIMARY KEY("producthazardstatements_product_id","producthazardstatements_hazard_statement_id"),
	FOREIGN KEY("producthazardstatements_hazard_statement_id") REFERENCES "hazard_statement"("hazard_statement_id") ON DELETE CASCADE,
	FOREIGN KEY("producthazardstatements_product_id") REFERENCES "product"("product_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "productprecautionarystatements";
CREATE TABLE "productprecautionarystatements" (
	"productprecautionarystatements_product_id"	INTEGER NOT NULL,
	"productprecautionarystatements_precautionary_statement_id"	INTEGER NOT NULL,
	PRIMARY KEY("productprecautionarystatements_product_id","productprecautionarystatements_precautionary_statement_id"),
	FOREIGN KEY("productprecautionarystatements_precautionary_statement_id") REFERENCES "precautionary_statement"("precautionary_statement_id") ON DELETE CASCADE,
	FOREIGN KEY("productprecautionarystatements_product_id") REFERENCES "product"("product_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "productsupplierrefs";
CREATE TABLE "productsupplierrefs" (
	"productsupplierrefs_product_id"	INTEGER NOT NULL,
	"productsupplierrefs_supplier_ref_id"	INTEGER NOT NULL,
	PRIMARY KEY("productsupplierrefs_product_id","productsupplierrefs_supplier_ref_id"),
	FOREIGN KEY("productsupplierrefs_product_id") REFERENCES "product"("product_id") ON DELETE CASCADE,
	FOREIGN KEY("productsupplierrefs_supplier_ref_id") REFERENCES "supplier_ref"("supplier_ref_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "productsymbols";
CREATE TABLE "productsymbols" (
	"productsymbols_product_id"	integer NOT NULL,
	"productsymbols_symbol_id"	integer NOT NULL,
	PRIMARY KEY("productsymbols_product_id","productsymbols_symbol_id"),
	FOREIGN KEY("productsymbols_product_id") REFERENCES "product"("product_id") ON DELETE CASCADE,
	FOREIGN KEY("productsymbols_symbol_id") REFERENCES "symbol"("symbol_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "productsynonyms";
CREATE TABLE "productsynonyms" (
	"productsynonyms_product_id"	integer NOT NULL,
	"productsynonyms_name_id"	integer NOT NULL,
	PRIMARY KEY("productsynonyms_product_id","productsynonyms_name_id"),
	FOREIGN KEY("productsynonyms_name_id") REFERENCES "name"("name_id") ON DELETE CASCADE,
	FOREIGN KEY("productsynonyms_product_id") REFERENCES "product"("product_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "producttags";
CREATE TABLE "producttags" (
	"producttags_product_id"	integer NOT NULL,
	"producttags_tag_id"	integer NOT NULL,
	PRIMARY KEY("producttags_product_id","producttags_tag_id"),
	FOREIGN KEY("producttags_product_id") REFERENCES "product"("product_id") ON DELETE CASCADE,
	FOREIGN KEY("producttags_tag_id") REFERENCES "tag"("tag_id") ON DELETE CASCADE
) STRICT;

DROP TABLE IF EXISTS "entitypeople";
CREATE TABLE "entitypeople" (
	"entitypeople_entity_id"	integer NOT NULL,
	"entitypeople_person_id"	integer NOT NULL,
	PRIMARY KEY("entitypeople_entity_id","entitypeople_person_id"),
	FOREIGN KEY("entitypeople_entity_id") REFERENCES "entity"("entity_id") ON DELETE CASCADE,
	FOREIGN KEY("entitypeople_person_id") REFERENCES "person"("person_id") ON DELETE RESTRICT
) STRICT;

DROP TABLE IF EXISTS "personentities";
CREATE TABLE "personentities" (
	"personentities_person_id"	integer NOT NULL,
	"personentities_entity_id"	integer NOT NULL,
	PRIMARY KEY("personentities_person_id","personentities_entity_id"),
	FOREIGN KEY("personentities_entity_id") REFERENCES "entity"("entity_id") ON DELETE CASCADE,
	FOREIGN KEY("personentities_person_id") REFERENCES "person"("person_id") ON DELETE CASCADE
) STRICT;

DROP INDEX IF EXISTS "idx_entitypeople";
CREATE UNIQUE INDEX "idx_entitypeople" ON "entitypeople" (
	"entitypeople_entity_id",
	"entitypeople_person_id"
);
DROP INDEX IF EXISTS "idx_personentities";
CREATE UNIQUE INDEX "idx_personentities" ON "personentities" (
	"personentities_person_id",
	"personentities_entity_id"
);
DROP INDEX IF EXISTS "idx_product_casnumber";
CREATE UNIQUE INDEX "idx_product_casnumber" ON "product" (
	"product_id",
	"cas_number"
);
DROP INDEX IF EXISTS "idx_product_cenumber";
CREATE UNIQUE INDEX "idx_product_cenumber" ON "product" (
	"product_id",
	"ce_number"
);
DROP INDEX IF EXISTS "idx_product_empiricalformula";
CREATE UNIQUE INDEX "idx_product_empiricalformula" ON "product" (
	"product_id",
	"empirical_formula"
);
DROP INDEX IF EXISTS "idx_productsymbols";
CREATE UNIQUE INDEX "idx_productsymbols" ON "productsymbols" (
	"productsymbols_product_id",
	"productsymbols_symbol_id"
);
DROP INDEX IF EXISTS "idx_productsynonyms";
CREATE UNIQUE INDEX "idx_productsynonyms" ON "productsynonyms" (
	"productsynonyms_product_id",
	"productsynonyms_name_id"
);
DROP INDEX IF EXISTS "idx_producttags";
CREATE UNIQUE INDEX "idx_producttags" ON "producttags" (
	"producttags_product_id",
	"producttags_tag_id"
);
COMMIT;
