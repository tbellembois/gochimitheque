package datastores

var versionToMigration = []string{migrationOne, migrationTwo, migrationThree, migrationFour, migrationFive, migrationSix, migrationSeven, migrationEight, migrationNine, migrationTen}

var migrationOne = `BEGIN TRANSACTION;

ALTER TABLE hazardstatement ADD hazardstatement_cmr string;

UPDATE hazardstatement SET hazardstatement_cmr='M1' WHERE hazardstatement_reference='H340';
UPDATE hazardstatement SET hazardstatement_cmr='M2' WHERE hazardstatement_reference='H341';
UPDATE hazardstatement SET hazardstatement_cmr='C1' WHERE hazardstatement_reference='H350';
UPDATE hazardstatement SET hazardstatement_cmr='C1' WHERE hazardstatement_reference='H350i';
UPDATE hazardstatement SET hazardstatement_cmr='C2' WHERE hazardstatement_reference='H351';
UPDATE hazardstatement SET hazardstatement_cmr='R1' WHERE hazardstatement_reference='H360';
UPDATE hazardstatement SET hazardstatement_cmr='R1' WHERE hazardstatement_reference='H360F';
UPDATE hazardstatement SET hazardstatement_cmr='R1' WHERE hazardstatement_reference='H360D';
UPDATE hazardstatement SET hazardstatement_cmr='R1' WHERE hazardstatement_reference='H360Fd';
UPDATE hazardstatement SET hazardstatement_cmr='R1' WHERE hazardstatement_reference='H360Df';
UPDATE hazardstatement SET hazardstatement_cmr='R1' WHERE hazardstatement_reference='H360FD';
UPDATE hazardstatement SET hazardstatement_cmr='R2' WHERE hazardstatement_reference='H361';
UPDATE hazardstatement SET hazardstatement_cmr='R2' WHERE hazardstatement_reference='H361f';
UPDATE hazardstatement SET hazardstatement_cmr='R2' WHERE hazardstatement_reference='H361d';
UPDATE hazardstatement SET hazardstatement_cmr='R2' WHERE hazardstatement_reference='H361fd';
UPDATE hazardstatement SET hazardstatement_cmr='L' WHERE hazardstatement_reference='H362';

PRAGMA user_version=1;
COMMIT;
`

var migrationTwo = `BEGIN TRANSACTION;
		
DELETE FROM permission WHERE permission_item_name='storelocations';
DELETE FROM permission WHERE permission_id IN (SELECT p1.permission_id FROM permission p1 INNER JOIN permission p2 WHERE p1.person=p2.person AND p1.permission_perm_name="r" AND p2.permission_perm_name="w" AND p1.permission_item_name=p2.permission_item_name AND p1.permission_entity_id=p2.permission_entity_id);

PRAGMA user_version=2;
COMMIT;
`

var migrationThree = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS new_unit (
	unit_id integer PRIMARY KEY,
	unit_label string UNIQUE NOT NULL,
	unit_multiplier integer NOT NULL default 1,
	unit_type string,
	unit integer,
	FOREIGN KEY(unit) references unit(unit_id));

INSERT into new_unit (
	unit_id,
	unit_label,
	unit_multiplier,
	unit
)
SELECT unit_id,
	unit_label,
	unit_multiplier,
	unit
FROM unit;

DROP table unit;
ALTER TABLE new_unit RENAME TO unit; 

INSERT OR IGNORE INTO unit (unit_label) VALUES ("L"), ("mL"), ("µL");
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="L"), unit_multiplier=0.001 WHERE unit_label="mL";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="L"), unit_multiplier=0.000001 WHERE unit_label="µL";
UPDATE unit SET unit_multiplier=1 WHERE unit_label="L";

INSERT OR IGNORE INTO unit (unit_label) VALUES ("kg"), ("g"), ("mg"), ("µg");
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="g"), unit_multiplier=1000 WHERE unit_label="kg";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="g"), unit_multiplier=0.001 WHERE unit_label="mg";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="g"), unit_multiplier=0.000001 WHERE unit_label="µg";
UPDATE unit SET unit_multiplier=1 WHERE unit_label="g";

INSERT OR IGNORE INTO unit (unit_label) VALUES ("m"), ("dm"), ("cm");
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="m"), unit_multiplier=10 WHERE unit_label="dm";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="m"), unit_multiplier=100 WHERE unit_label="cm";
UPDATE unit SET unit_multiplier=1 WHERE unit_label="m";

INSERT OR IGNORE INTO unit (unit_label) VALUES ("°K"), ("°F"), ("°C");
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="°K") WHERE unit_label="°F";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="°K") WHERE unit_label="°C";

INSERT OR IGNORE INTO unit (unit_label) VALUES ("nM"), ("µM"), ("mM");
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="mM") WHERE unit_label="µM";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="mM") WHERE unit_label="mM";

INSERT OR IGNORE INTO unit (unit_label) VALUES ("ng/L"), ("µg/L"), ("mg/L"), ("g/L");
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="g/L") WHERE unit_label="ng/L";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="g/L") WHERE unit_label="µg/L";
UPDATE unit SET unit=(SELECT unit_id FROM unit WHERE unit_label="g/L") WHERE unit_label="mg/L";

UPDATE unit SET unit_type="quantity" WHERE unit_label="L";
UPDATE unit SET unit_type="quantity" WHERE unit=(SELECT unit_id FROM unit WHERE unit_label="L");
UPDATE unit SET unit_type="quantity" WHERE unit_label="g";
UPDATE unit SET unit_type="quantity" WHERE unit=(SELECT unit_id FROM unit WHERE unit_label="g");
UPDATE unit SET unit_type="quantity" WHERE unit_label="m";
UPDATE unit SET unit_type="quantity" WHERE unit=(SELECT unit_id FROM unit WHERE unit_label="m");
UPDATE unit SET unit_type="temperature" WHERE unit_label="°K";
UPDATE unit SET unit_type="temperature" WHERE unit=(SELECT unit_id FROM unit WHERE unit_label="°K");
UPDATE unit SET unit_type="concentration" WHERE unit_type IS NULL;

CREATE TABLE IF NOT EXISTS tag (
	tag_id integer PRIMARY KEY,
	tag_label string NOT NULL);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tag_label ON tag(tag_label);

CREATE TABLE IF NOT EXISTS category (
	category_id integer PRIMARY KEY,
	category_label string NOT NULL);
CREATE UNIQUE INDEX IF NOT EXISTS idx_category_label ON category(category_label);

CREATE TABLE IF NOT EXISTS producer (
	producer_id integer PRIMARY KEY,
	producer_label string NOT NULL);
CREATE UNIQUE INDEX IF NOT EXISTS idx_producer_label ON producer(producer_label);

CREATE TABLE IF NOT EXISTS producerref (
	producerref_id integer PRIMARY KEY,
	producerref_label string NOT NULL,
	producer integer,
	FOREIGN KEY(producer) references producer(producer_id));
CREATE UNIQUE INDEX IF NOT EXISTS idx_producerref_label ON producerref(producerref_label);

CREATE TABLE IF NOT EXISTS supplierref (
	supplierref_id integer PRIMARY KEY,
	supplierref_label string NOT NULL,
	supplier integer,
	FOREIGN KEY(supplier) references supplier(supplier_id));
CREATE UNIQUE INDEX IF NOT EXISTS idx_supplierref_label ON supplierref(supplierref_label);

CREATE TABLE IF NOT EXISTS productsupplierrefs (
	productsupplierrefs_product_id integer NOT NULL,
	productsupplierrefs_supplierref_id integer NOT NULL,
	PRIMARY KEY(productsupplierrefs_product_id, productsupplierrefs_supplierref_id),
	FOREIGN KEY(productsupplierrefs_product_id) references product(product_id),
	FOREIGN KEY(productsupplierrefs_supplierref_id) references supplierref(supplierref_id));
CREATE UNIQUE INDEX IF NOT EXISTS idx_productsupplierrefs ON productsupplierrefs(productsupplierrefs_product_id, productsupplierrefs_supplierref_id);

CREATE TABLE IF NOT EXISTS producttags (
	producttags_product_id integer NOT NULL,
	producttags_tag_id integer NOT NULL,
	PRIMARY KEY(producttags_product_id, producttags_tag_id),
	FOREIGN KEY(producttags_product_id) references product(product_id),
	FOREIGN KEY(producttags_tag_id) references tag(tag_id));
CREATE UNIQUE INDEX IF NOT EXISTS idx_producttags ON producttags(producttags_product_id, producttags_tag_id);

CREATE TABLE IF NOT EXISTS new_storage (
	storage_id integer PRIMARY KEY,
	storage_creationdate datetime NOT NULL,
	storage_modificationdate datetime NOT NULL,
	storage_entrydate datetime,
	storage_exitdate datetime,
	storage_openingdate datetime,
	storage_expirationdate datetime,
	storage_quantity float,
	storage_barecode string,
	storage_comment string,
	storage_reference string,
	storage_batchnumber string,
	storage_todestroy boolean default 0,
	storage_archive boolean default 0,
	storage_qrcode blob,
	storage_concentration integer,
	storage_number_of_unit integer,
	storage_number_of_bag integer,
	storage_number_of_carton integer,
	person integer NOT NULL,
	product integer NOT NULL,
	storelocation integer NOT NULL,
	unit_concentration integer,
	unit_quantity integer,
	supplier integer,
	storage integer,
	FOREIGN KEY(unit_concentration) references unit(unit_id),
	FOREIGN KEY(storage) references storage(storage_id),
	FOREIGN KEY(unit_quantity) references unit(unit_id),
	FOREIGN KEY(supplier) references supplier(supplier_id),
	FOREIGN KEY(person) references person(person_id),
	FOREIGN KEY(product) references product(product_id),
	FOREIGN KEY(storelocation) references storelocation(storelocation_id));

CREATE TABLE IF NOT EXISTS new_product (
	product_id integer PRIMARY KEY,
	product_specificity string,
	product_msds string,
	product_restricted boolean default 0,
	product_radioactive boolean default 0,
	product_threedformula string,
	product_twodformula string,
	product_molformula blob,
	product_disposalcomment string,
	product_remark string,
	product_qrcode string,
	product_sheet string,
	product_concentration integer,
	product_temperature integer,
	product_number_per_carton integer,
	product_number_per_bag integer,
	casnumber integer,
	cenumber integer,
	person integer NOT NULL,
	empiricalformula integer,
	linearformula integer,
	physicalstate integer,
	signalword integer,
	name integer NOT NULL,
	producerref integer,
	unit_temperature integer,
	category integer,
	FOREIGN KEY(unit_temperature) references unit(unit_id),
	FOREIGN KEY(producerref) references producerref(producerref_id),
	FOREIGN KEY(category) references category(category_id),
	FOREIGN KEY(casnumber) references casnumber(casnumber_id),
	FOREIGN KEY(cenumber) references cenumber(cenumber_id),
	FOREIGN KEY(person) references person(person_id),
	FOREIGN KEY(empiricalformula) references empiricalformula(empiricalformula_id),
	FOREIGN KEY(linearformula) references linearformula(linearformula_id),
	FOREIGN KEY(physicalstate) references physicalstate(physicalstate_id),
	FOREIGN KEY(signalword) references signalword(signalword_id),
	FOREIGN KEY(name) references name(name_id));

INSERT INTO new_product (
	product_id,
	product_specificity,
	product_msds,
	product_restricted,
	product_radioactive,
	product_threedformula,
	product_molformula,
	product_disposalcomment,
	product_remark,
	product_qrcode,
	casnumber,
	cenumber,
	person,
	empiricalformula,
	linearformula,
	physicalstate,
	signalword,
	name
)
SELECT product_id,
	product_specificity,
	product_msds,
	product_restricted,
	product_radioactive,
	product_threed_formula,
	product_molformula,
	product_disposal_comment,
	product_remark,
	product_qrcode,
	casnumber,
	cenumber,
	person,
	empiricalformula,
	linearformula,
	physicalstate,
	signalword,
	name
FROM product;

INSERT INTO new_storage (
	storage_id,
	storage_creationdate,
	storage_modificationdate,
	storage_entrydate,
	storage_exitdate,
	storage_openingdate,
	storage_expirationdate,
	storage_quantity,
	storage_barecode,
	storage_comment,
	storage_reference,
	storage_batchnumber,
	storage_todestroy,
	storage_archive,
	storage_qrcode,
	person,
	product,
	storelocation,
	unit_quantity,
	supplier,
	storage
)
SELECT storage_id,
	storage_creationdate,
	storage_modificationdate,
	storage_entrydate,
	storage_exitdate,
	storage_openingdate,
	storage_expirationdate,
	storage_quantity,
	storage_barecode,
	storage_comment,
	storage_reference,
	storage_batchnumber,
	storage_todestroy,
	storage_archive,
	storage_qrcode,
	person,
	product,
	storelocation,
	unit,
	supplier,
	storage
FROM storage;

DROP TABLE product;
ALTER TABLE new_product RENAME TO product; 

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_casnumber ON product(product_id, casnumber);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_cenumber ON product(product_id, cenumber);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_empiricalformula ON product(product_id, empiricalformula);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_name ON product(product_id, name);

DROP TABLE storage;
ALTER TABLE new_storage RENAME TO storage; 

CREATE UNIQUE INDEX IF NOT EXISTS idx_storage_product ON storage(storage_id, product);
CREATE UNIQUE INDEX IF NOT EXISTS idx_storage_storelocation ON storage(storage_id, storelocation);
CREATE UNIQUE INDEX IF NOT EXISTS idx_storage_storelocation_product ON storage(storage_id, storelocation, product);


UPDATE product SET empiricalformula=null WHERE empiricalformula=(SELECT empiricalformula_id FROM empiricalformula WHERE empiricalformula_label="XXXX");
DELETE FROM empiricalformula where empiricalformula_label="XXXX";

UPDATE product SET casnumber=null WHERE casnumber=(SELECT casnumber_id FROM casnumber WHERE casnumber_label="0000-00-0");
DELETE FROM casnumber where casnumber_label="0000-00-0";

CREATE INDEX "idx_permission_person" ON "permission" (
	"person" ASC
);
CREATE INDEX "idx_permission_perm_name" ON "permission" (
	"permission_perm_name"	ASC
);
CREATE INDEX "idx_permission_item_name" ON "permission" (
	"permission_item_name"	ASC
);
CREATE INDEX "idx_permission_entity_id" ON "permission" (
	"permission_entity_id"	ASC
);

PRAGMA user_version=3;
COMMIT;
PRAGMA foreign_keys=on;
`

var migrationFour = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS new_storage (
	storage_id integer PRIMARY KEY,
	storage_creationdate datetime NOT NULL,
	storage_modificationdate datetime NOT NULL,
	storage_entrydate datetime,
	storage_exitdate datetime,
	storage_openingdate datetime,
	storage_expirationdate datetime,
	storage_quantity float,
	storage_barecode text,
	storage_comment text,
	storage_reference text,
	storage_batchnumber text,
	storage_todestroy boolean default 0,
	storage_archive boolean default 0,
	storage_qrcode blob,
	storage_concentration integer,
	storage_number_of_unit integer,
	storage_number_of_bag integer,
	storage_number_of_carton integer,
	person integer NOT NULL,
	product integer NOT NULL,
	storelocation integer NOT NULL,
	unit_concentration integer,
	unit_quantity integer,
	supplier integer,
	storage integer,
	FOREIGN KEY(unit_concentration) references unit(unit_id),
	FOREIGN KEY(storage) references storage(storage_id),
	FOREIGN KEY(unit_quantity) references unit(unit_id),
	FOREIGN KEY(supplier) references supplier(supplier_id),
	FOREIGN KEY(person) references person(person_id),
	FOREIGN KEY(product) references product(product_id),
	FOREIGN KEY(storelocation) references storelocation(storelocation_id));

INSERT INTO new_storage (
	storage_id,
	storage_creationdate,
	storage_modificationdate,
	storage_entrydate,
	storage_exitdate,
	storage_openingdate,
	storage_expirationdate,
	storage_quantity,
	storage_barecode,
	storage_comment,
	storage_reference,
	storage_batchnumber,
	storage_todestroy,
	storage_archive,
	storage_qrcode,
	storage_concentration,
	storage_number_of_unit,
	storage_number_of_bag,
	storage_number_of_carton,
	person,
	product,
	storelocation,
	unit_concentration,
	unit_quantity,
	supplier,
	storage
)
SELECT storage_id,
	storage_creationdate,
	storage_modificationdate,
	storage_entrydate,
	storage_exitdate,
	storage_openingdate,
	storage_expirationdate,
	storage_quantity,
	storage_barecode,
	storage_comment,
	storage_reference,
	storage_batchnumber,
	storage_todestroy,
	storage_archive,
	storage_qrcode,
	storage_concentration,
	storage_number_of_unit,
	storage_number_of_bag,
	storage_number_of_carton,
	person,
	product,
	storelocation,
	unit_concentration,
	unit_quantity,
	supplier,
	storage
FROM storage;

DROP TABLE storage;
ALTER TABLE new_storage RENAME TO storage; 

PRAGMA user_version=4;
COMMIT;
PRAGMA foreign_keys=on;
`

var migrationFive = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

INSERT INTO unit (unit_label, unit_multiplier, unit_type) 
VALUES ("%", 1, "concentration"), ("X", 1, "concentration");

PRAGMA user_version=5;
COMMIT;
PRAGMA foreign_keys=on;
`

var migrationSix = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS new_person(
	person_id integer PRIMARY KEY,
	person_email string NOT NULL,
	person_password string NOT NULL,
	person_aeskey string NOT NULL);

INSERT INTO new_person(
	person_id,
	person_email,
	person_password,
	person_aeskey
)
SELECT person_id,
	person_email,
	person_password,
	"32byteskeytobechangedbygocode+++"
FROM person;

DROP TABLE person;
ALTER TABLE new_person RENAME TO person;

CREATE UNIQUE INDEX IF NOT EXISTS idx_person ON person(person_email);

PRAGMA user_version=6;
COMMIT;
PRAGMA foreign_keys=on;`

var migrationSeven = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS entityldapgroups (
	entityldapgroups_entity_id integer NOT NULL,
	entityldapgroups_ldapgroup string NOT NULL,
	PRIMARY KEY(entityldapgroups_entity_id, entityldapgroups_ldapgroup),
	FOREIGN KEY(entityldapgroups_entity_id) references entity(entity_id));
CREATE UNIQUE INDEX IF NOT EXISTS idx_entityldapgroups ON entityldapgroups(entityldapgroups_entity_id, entityldapgroups_ldapgroup);

PRAGMA user_version=7;
COMMIT;
PRAGMA foreign_keys=on;`

var migrationEight = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.May cause endocrine disruption in humans.', 'EUH380');
INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.Suspected of causing endocrine disruption in humans.', 'EUH381');
INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.May cause endocrine disruption in the environment.', 'EUH430');
INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.Suspected of causing endocrine disruption in the environment.', 'EUH431');
INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.Accumulates in the environment and living organisms including in humans.', 'EUH440');
INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.Strongly accumulates in the environment and living organisms including in humans.', 'EUH441');
INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.Can cause long-lasting and diffuse contamination of water resources.', 'EUH450');
INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES ('.Can cause very long-lasting and diffuse contamination of water resources.', 'EUH451');

PRAGMA user_version=8;
COMMIT;
PRAGMA foreign_keys=on;`

var migrationNine = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

ALTER TABLE person DROP COLUMN person_password;
ALTER TABLE person DROP COLUMN person_aeskey;

DROP TABLE entityldapgroups;

PRAGMA user_version=9;
COMMIT;
PRAGMA foreign_keys=on;`

var migrationTen = `PRAGMA foreign_keys=off;

BEGIN TRANSACTION;

CREATE TABLE bookmark_new (
	bookmark_id	INTEGER,
	person	INTEGER NOT NULL,
	product	INTEGER NOT NULL,
	FOREIGN KEY(person) REFERENCES person(person_id),
	FOREIGN KEY(product) REFERENCES product(product_id),
	PRIMARY KEY(bookmark_id)
) STRICT;
INSERT INTO bookmark_new (
	bookmark_id,
	person,
	product
)
SELECT bookmark_id,
	person,
	product
FROM bookmark;
DROP TABLE bookmark;
ALTER TABLE bookmark_new RENAME TO bookmark; 

CREATE TABLE borrowing_new (
	borrowing_id	INTEGER,
	borrowing_comment	TEXT,
	person	INTEGER NOT NULL,
	borrower	INTEGER NOT NULL,
	storage	INTEGER NOT NULL UNIQUE,
	FOREIGN KEY(person) REFERENCES person(person_id),
	FOREIGN KEY(storage) REFERENCES storage(storage_id),
	FOREIGN KEY(borrower) REFERENCES person(person_id),
	PRIMARY KEY(borrowing_id)
) STRICT;
INSERT INTO borrowing_new (
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
FROM borrowing;
DROP TABLE borrowing;
ALTER TABLE borrowing_new RENAME TO borrowing; 

CREATE TABLE cas_number_new (
	cas_number_id	INTEGER,
	cas_number_label	TEXT NOT NULL UNIQUE,
	cas_number_cmr	TEXT,
	PRIMARY KEY(cas_number_id)
) STRICT;
INSERT INTO cas_number_new (
	cas_number_id,
	cas_number_label,
	cas_number_cmr
)
SELECT casnumber_id,
	casnumber_label,
	casnumber_cmr
FROM casnumber;
DROP TABLE casnumber;
ALTER TABLE cas_number_new RENAME TO cas_number;

CREATE TABLE category_new (
	category_id	INTEGER,
	category_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(category_id)
) STRICT;
INSERT INTO category_new (
	category_id,
	category_label
)	
SELECT category_id,
	category_label
FROM category;
DROP TABLE category;
ALTER TABLE category_new RENAME TO category; 

CREATE TABLE ce_number_new (
	ce_number_id	INTEGER,
	ce_number_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(ce_number_id)
) STRICT;
INSERT INTO ce_number_new (
	ce_number_id,
	ce_number_label
)
SELECT cenumber_id,
	cenumber_label 
FROM cenumber;
DROP TABLE cenumber;
ALTER TABLE ce_number_new RENAME TO ce_number;

CREATE TABLE class_of_compound_new (
	class_of_compound_id	INTEGER,
	class_of_compound_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(class_of_compound_id)
) STRICT;
INSERT INTO class_of_compound_new (
	class_of_compound_id,
	class_of_compound_label
)
SELECT classofcompound_id,
	classofcompound_label
FROM classofcompound;
DROP TABLE classofcompound;
ALTER TABLE class_of_compound_new RENAME TO class_of_compound;

CREATE TABLE empirical_formula_new (
	empirical_formula_id	INTEGER,
	empirical_formula_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(empirical_formula_id)
) STRICT;
INSERT INTO empirical_formula_new (
	empirical_formula_id,
	empirical_formula_label
)
SELECT empiricalformula_id,
	empiricalformula_label
FROM empiricalformula;
DROP TABLE empiricalformula;
ALTER TABLE empirical_formula_new RENAME TO empirical_formula;

CREATE TABLE entity_new (
	entity_id	INTEGER,
	entity_name	TEXT NOT NULL UNIQUE,
	entity_description	TEXT,
	PRIMARY KEY(entity_id)
) STRICT;
INSERT INTO entity_new (
	entity_id,
	entity_name,
	entity_description
)
SELECT entity_id,
	entity_name,
	entity_description
FROM entity;
DROP TABLE entity;
ALTER TABLE entity_new RENAME TO entity; 

CREATE TABLE hazard_statement_new (
	hazard_statement_id	INTEGER,
	hazard_statement_label	TEXT NOT NULL,
	hazard_statement_reference	TEXT NOT NULL UNIQUE,
	hazard_statement_cmr	TEXT,
	PRIMARY KEY(hazard_statement_id)
) STRICT;
INSERT INTO hazard_statement_new (
	hazard_statement_id,
	hazard_statement_label,
	hazard_statement_reference,
	hazard_statement_cmr
)
SELECT hazardstatement_id,
	hazardstatement_label,
	hazardstatement_reference,
	hazardstatement_cmr
FROM hazardstatement;
DROP TABLE hazardstatement;
ALTER TABLE hazard_statement_new RENAME TO hazard_statement;

CREATE TABLE linear_formula_new (
	linear_formula_id	INTEGER,
	linear_formula_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(linear_formula_id)
) STRICT;
INSERT INTO linear_formula_new (
	linear_formula_id,
	linear_formula_label
)
SELECT linearformula_id,
	linearformula_label
FROM linearformula;
DROP TABLE linearformula;
ALTER TABLE linear_formula_new RENAME TO linear_formula;

CREATE TABLE name_new (
	name_id	INTEGER,
	name_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(name_id)
) STRICT;
INSERT INTO name_new (
	name_id,
	name_label
)
SELECT name_id,
	name_label
FROM name;
DROP TABLE name;
ALTER TABLE name_new RENAME TO name; 

CREATE TABLE permission_new (
	permission_id	INTEGER,
	person	INTEGER NOT NULL,
	permission_perm_name	TEXT NOT NULL,
	permission_item_name	TEXT NOT NULL,
	permission_entity_id	INTEGER,
	FOREIGN KEY(person) REFERENCES person(person_id),
	PRIMARY KEY(permission_id)
) STRICT;
INSERT INTO permission_new (
	permission_id,
	person,
	permission_perm_name,
	permission_item_name,
	permission_entity_id
)
SELECT permission_id,
	person,
	permission_perm_name,
	permission_item_name,
	permission_entity_id
FROM permission;
DROP TABLE permission;
ALTER TABLE permission_new RENAME TO permission; 

CREATE TABLE person_new (
	person_id	INTEGER,
	person_email	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(person_id)
) STRICT;
INSERT INTO person_new (
	person_id,
	person_email
)
SELECT person_id,
	person_email
FROM person;
DROP TABLE person;
ALTER TABLE person_new RENAME TO person; 

CREATE TABLE physical_state_new (
	physical_state_id	INTEGER,
	physical_state_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(physical_state_id)
) STRICT;
INSERT into physical_state_new (
	physical_state_id,
	physical_state_label
)
SELECT physicalstate_id,
	physicalstate_label
FROM physicalstate;
DROP TABLE physicalstate;
ALTER TABLE physical_state_new RENAME TO physical_state;

CREATE TABLE precautionary_statement_new (
	precautionary_statement_id	INTEGER,
	precautionary_statement_label	TEXT NOT NULL,
	precautionary_statement_reference	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(precautionary_statement_id)
) STRICT;
INSERT INTO precautionary_statement_new (
	precautionary_statement_id,
	precautionary_statement_label,
	precautionary_statement_reference
)
SELECT precautionarystatement_id,
	precautionarystatement_label,
	precautionarystatement_reference
FROM precautionarystatement;
DROP TABLE precautionarystatement;
ALTER TABLE precautionary_statement_new RENAME TO precautionary_statement;

CREATE TABLE producer_new (
	producer_id	INTEGER,
	producer_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(producer_id)
) STRICT;
INSERT INTO producer_new (
	producer_id,
	producer_label
)
SELECT producer_id,
	producer_label
FROM producer;
DROP TABLE producer;
ALTER TABLE producer_new RENAME TO producer; 

CREATE TABLE producer_ref_new (
	producer_ref_id	INTEGER,
	producer_ref_label	TEXT NOT NULL,
	producer	INTEGER,
	FOREIGN KEY(producer) REFERENCES producer(producer_id),
	PRIMARY KEY(producer_ref_id)
) STRICT;
INSERT INTO producer_ref_new (
	producer_ref_id,
	producer_ref_label,
	producer
)
SELECT producerref_id,
	producerref_label,
	producer
FROM producerref;
DROP TABLE producerref;
ALTER TABLE producer_ref_new RENAME TO producer_ref;

CREATE TABLE product_new (
	product_id	INTEGER,
	product_inchi TEXT,
	product_inchikey TEXT,
	product_canonical_smiles TEXT,
	product_specificity	TEXT,
	product_msds	TEXT,
	product_restricted	INTEGER DEFAULT 0,
	product_radioactive	INTEGER DEFAULT 0,
	product_threed_formula	TEXT,
	product_twod_formula	TEXT,
	product_disposal_comment	TEXT,
	product_remark	TEXT,
	product_qrcode	TEXT,
	product_sheet	TEXT,
	product_concentration	REAL,
	product_temperature	REAL,
	product_molecular_weight REAL,
	cas_number	INTEGER,
	ce_number	INTEGER,
	person	INTEGER NOT NULL,
	empirical_formula	INTEGER,
	linear_formula	INTEGER,
	physical_state	INTEGER,
	signal_word	INTEGER,
	name	INTEGER NOT NULL,
	producer_ref	INTEGER,
	unit_molecular_weight INTEGER,
	unit_temperature	INTEGER,
	category	INTEGER,
	product_number_per_carton	INTEGER,
	product_number_per_bag	INTEGER,
	FOREIGN KEY(person) REFERENCES person(person_id),
	FOREIGN KEY(empirical_formula) REFERENCES empirical_formula(empirical_formula_id),
	FOREIGN KEY(linear_formula) REFERENCES linear_formula(linear_formula_id),
	FOREIGN KEY(cas_number) REFERENCES cas_number(cas_number_id),
	FOREIGN KEY(ce_number) REFERENCES ce_number(ce_number_id),
	FOREIGN KEY(producer_ref) REFERENCES producer_ref(producer_ref_id),
	FOREIGN KEY(category) REFERENCES category(category_id),
	PRIMARY KEY(product_id),
	FOREIGN KEY(unit_temperature) REFERENCES unit(unit_id),
	FOREIGN KEY(unit_molecular_weight) REFERENCES unit(unit_id),
	FOREIGN KEY(physical_state) REFERENCES physical_state(physical_state_id),
	FOREIGN KEY(signal_word) REFERENCES signal_word(signal_word_id),
	FOREIGN KEY(name) REFERENCES name(name_id)
) STRICT;
INSERT into product_new (
	product_id,
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
FROM product;
DROP TABLE product;
ALTER TABLE product_new RENAME TO product; 

CREATE TABLE signal_word_new (
	signal_word_id	INTEGER,
	signal_word_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(signal_word_id)
) STRICT;
INSERT INTO signal_word_new (
	signal_word_id,
	signal_word_label
)
SELECT signalword_id,
	signalword_label
FROM signalword;
DROP TABLE signalword;
ALTER TABLE signal_word_new RENAME TO signal_word;

CREATE TABLE storage_new (
	storage_id	INTEGER,
	storage_creation_date	INTEGER NOT NULL DEFAULT current_timestamp,
	storage_modification_date	INTEGER NOT NULL DEFAULT current_timestamp,
	storage_entry_date	INTEGER,
	storage_exit_date	INTEGER,
	storage_opening_date	INTEGER,
	storage_expiration_date	INTEGER,
	storage_quantity	REAL,
	storage_barecode	TEXT,
	storage_comment	TEXT,
	storage_reference	TEXT,
	storage_batch_number	TEXT,
	storage_to_destroy	INTEGER DEFAULT 0,
	storage_archive	INTEGER DEFAULT 0,
	storage_qrcode	BLOB,
	storage_concentration	REAL,
	storage_number_of_unit	INTEGER,
	storage_number_of_bag	INTEGER,
	storage_number_of_carton	INTEGER,
	person	INTEGER NOT NULL,
	product	INTEGER NOT NULL,
	store_location	INTEGER NOT NULL,
	unit_concentration	REAL,
	unit_quantity	REAL,
	supplier	INTEGER,
	storage	INTEGER,
	FOREIGN KEY(unit_concentration) REFERENCES unit(unit_id),
	FOREIGN KEY(storage) REFERENCES storage(storage_id),
	FOREIGN KEY(unit_quantity) REFERENCES unit(unit_id),
	FOREIGN KEY(supplier) REFERENCES supplier(supplier_id),
	FOREIGN KEY(person) REFERENCES person(person_id),
	FOREIGN KEY(product) REFERENCES product(product_id),
	FOREIGN KEY(store_location) REFERENCES store_location(store_location_id),
	PRIMARY KEY(storage_id)
) STRICT;
INSERT INTO storage_new (
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
	storage_number_of_unit,
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
storage_number_of_unit,
storage_number_of_bag,
storage_number_of_carton,
person,
product,
storelocation,
unit_concentration,
unit_quantity,
supplier,
storage
FROM storage;
DROP TABLE storage;
ALTER TABLE storage_new RENAME TO storage; 

-- CREATE TRIGGER insert_storage_modification_date_Trigger
-- AFTER INSERT ON storage
-- BEGIN
-- UPDATE storage SET storage_modification_date = current_timestamp WHERE storage_id = NEW.storage_id;
-- END;

CREATE TABLE store_location_new (
	store_location_id	INTEGER,
	store_location_name	TEXT NOT NULL,
	store_location_color	TEXT,
	store_location_can_store	INTEGER DEFAULT 0,
	store_location_full_path	TEXT,
	entity	INTEGER NOT NULL,
	store_location	INTEGER,
	FOREIGN KEY(store_location) REFERENCES store_location(store_location_id),
	FOREIGN KEY(entity) REFERENCES entity(entity_id),
	PRIMARY KEY(store_location_id)
) STRICT;
INSERT INTO store_location_new (
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
FROM storelocation;
DROP TABLE storelocation;
ALTER TABLE store_location_new RENAME TO store_location;

CREATE TABLE supplier_new (
	supplier_id	INTEGER,
	supplier_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(supplier_id)
) STRICT;
INSERT INTO supplier_new (
	supplier_id,
	supplier_label
)
SELECT supplier_id,
	supplier_label
FROM supplier;
DROP TABLE supplier;
ALTER TABLE supplier_new RENAME TO supplier;

CREATE TABLE supplier_ref_new (
	supplier_ref_id	INTEGER,
	supplier_ref_label	TEXT NOT NULL,
	supplier	INTEGER,
	FOREIGN KEY(supplier) REFERENCES supplier(supplier_id),
	PRIMARY KEY(supplier_ref_id)
) STRICT;
INSERT INTO supplier_ref_new (
	supplier_ref_id,
	supplier_ref_label,
	supplier
)
SELECT supplierref_id,
	supplierref_label,
	supplier
FROM supplierref;
DROP TABLE supplierref;
ALTER TABLE supplier_ref_new RENAME TO supplier_ref;

CREATE TABLE symbol_new (
	symbol_id	INTEGER,
	symbol_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(symbol_id)
) STRICT;

INSERT INTO symbol_new (symbol_label) VALUES ("GHS01"), ("GHS02"), ("GHS03"), ("GHS04"), ("GHS05"), ("GHS06"), ("GHS07"), ("GHS08"), ("GHS09");

DROP TABLE symbol;
ALTER TABLE symbol_new RENAME TO symbol;

CREATE TABLE tag_new (
	tag_id	INTEGER,
	tag_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(tag_id)
) STRICT;
INSERT INTO tag_new (
	tag_id,
	tag_label
)
SELECT tag_id,
	tag_label
FROM tag;
DROP TABLE tag;
ALTER TABLE tag_new RENAME TO tag;

CREATE TABLE unit_new (
	unit_id	INTEGER,
	unit_label	TEXT NOT NULL UNIQUE,
	unit_multiplier	REAL NOT NULL DEFAULT 1,
	unit_type	TEXT,
	unit	INTEGER,
	FOREIGN KEY(unit) REFERENCES unit(unit_id),
	PRIMARY KEY(unit_id)
) STRICT;
INSERT INTO unit_new (
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
FROM unit;
DROP TABLE unit;
ALTER TABLE unit_new RENAME TO unit;

CREATE TABLE welcome_announce_new (
	welcome_announce_id	INTEGER,
	welcome_announce_text	TEXT,
	PRIMARY KEY(welcome_announce_id)
) STRICT;
INSERT INTO welcome_announce_new (
	welcome_announce_id,
	welcome_announce_text
)
SELECT welcomeannounce_id,
	welcomeannounce_text
FROM welcomeannounce;
DROP TABLE welcomeannounce;
ALTER TABLE welcome_announce_new RENAME TO welcome_announce;

CREATE TABLE productclassesofcompounds (
	productclassesofcompounds_product_id INTEGER NOT NULL,
	productclassesofcompounds_class_of_compound_id INTEGER NOT NULL,
	PRIMARY KEY (productclassesofcompounds_product_id, productclassesofcompounds_class_of_compound_id),
	FOREIGN KEY (productclassesofcompounds_product_id) REFERENCES product(product_id),
	FOREIGN KEY (productclassesofcompounds_class_of_compound_id) REFERENCES class_of_compound(class_of_compound_id)
) STRICT;
INSERT INTO productclassesofcompounds (
	productclassesofcompounds_product_id,
	productclassesofcompounds_class_of_compound_id
)
SELECT productclassofcompound_product_id,
	productclassofcompound_classofcompound_id
FROM productclassofcompound;
DROP TABLE productclassofcompound;

CREATE TABLE producthazardstatements_new (
	producthazardstatements_product_id INTEGER NOT NULL,
	producthazardstatements_hazard_statement_id INTEGER NOT NULL,
	PRIMARY KEY (producthazardstatements_product_id, producthazardstatements_hazard_statement_id),
	FOREIGN KEY (producthazardstatements_product_id) REFERENCES product(product_id),
	FOREIGN KEY (producthazardstatements_hazard_statement_id) REFERENCES hazard_statement(hazard_statement_id)
) STRICT;
INSERT INTO producthazardstatements_new (
	producthazardstatements_product_id,
	producthazardstatements_hazard_statement_id
)
SELECT producthazardstatements_product_id,
producthazardstatements_hazardstatement_id
FROM producthazardstatements;
DROP TABLE producthazardstatements;
ALTER TABLE producthazardstatements_new RENAME TO producthazardstatements;

CREATE TABLE productprecautionarystatements_new (
productprecautionarystatements_product_id INTEGER NOT NULL,
productprecautionarystatements_precautionary_statement_id INTEGER NOT NULL,
PRIMARY KEY (productprecautionarystatements_product_id, productprecautionarystatements_precautionary_statement_id),
FOREIGN KEY (productprecautionarystatements_product_id) REFERENCES product(product_id),
FOREIGN KEY (productprecautionarystatements_precautionary_statement_id) REFERENCES precautionary_statement(precautionary_statement_id)
) STRICT;
INSERT INTO productprecautionarystatements_new (
productprecautionarystatements_product_id,
productprecautionarystatements_precautionary_statement_id
)
SELECT productprecautionarystatements_product_id,
productprecautionarystatements_precautionarystatement_id
FROM productprecautionarystatements;
DROP TABLE productprecautionarystatements;
ALTER TABLE productprecautionarystatements_new RENAME TO productprecautionarystatements;

CREATE TABLE productsupplierrefs_new (
productsupplierrefs_product_id INTEGER NOT NULL,
productsupplierrefs_supplier_ref_id INTEGER NOT NULL,
PRIMARY KEY (productsupplierrefs_product_id, productsupplierrefs_supplier_ref_id),
FOREIGN KEY (productsupplierrefs_product_id) REFERENCES product(product_id),
FOREIGN KEY (productsupplierrefs_supplier_ref_id) REFERENCES supplier_ref(supplier_ref_id)
) STRICT;
INSERT INTO productsupplierrefs_new (
productsupplierrefs_product_id,
productsupplierrefs_supplier_ref_id
)
SELECT productsupplierrefs_product_id,
productsupplierrefs_supplierref_id
FROM productsupplierrefs;
DROP TABLE productsupplierrefs;
ALTER TABLE productsupplierrefs_new RENAME TO productsupplierrefs;

DROP INDEX IF EXISTS idx_producerref_label;
DROP INDEX IF EXISTS idx_supplierref_label;

DROP TABLE IF EXISTS captcha;
DROP TABLE IF EXISTS entityldapgroups;

DROP INDEX IF EXISTS idx_product_casnumber;
DROP INDEX IF EXISTS idx_product_cenumber;
DROP INDEX IF EXISTS idx_product_empiricalformula;

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_casnumber ON product(product_id, cas_number);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_cenumber ON product(product_id, ce_number);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_empiricalformula ON product(product_id, empirical_formula);

INSERT INTO unit (unit_label, unit_multiplier, unit_type) VALUES ("g/mol", 1, "molecular_weight");

COMMIT;

PRAGMA user_version=10;
PRAGMA foreign_keys=on;`
