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

CREATE TABLE casnumber_new (
	casnumber_id	INTEGER,
	casnumber_label	TEXT NOT NULL UNIQUE,
	casnumber_cmr	TEXT,
	PRIMARY KEY(casnumber_id)
) STRICT;
INSERT INTO casnumber_new (
	casnumber_id,
	casnumber_label,
	casnumber_cmr
)
SELECT casnumber_id,
	casnumber_label,
	casnumber_cmr
FROM casnumber;
DROP TABLE casnumber;
ALTER TABLE casnumber_new RENAME TO casnumber; 

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

CREATE TABLE cenumber_new (
	cenumber_id	INTEGER,
	cenumber_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(cenumber_id)
) STRICT;
INSERT INTO cenumber_new (
	cenumber_id,
	cenumber_label
)
SELECT cenumber_id,
	cenumber_label 
FROM cenumber;
DROP TABLE cenumber;
ALTER TABLE cenumber_new RENAME TO cenumber; 

CREATE TABLE classofcompound_new (
	classofcompound_id	INTEGER,
	classofcompound_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(classofcompound_id)
) STRICT;
INSERT INTO classofcompound_new (
	classofcompound_id,
	classofcompound_label
)
SELECT classofcompound_id,
	classofcompound_label
FROM classofcompound;
DROP TABLE classofcompound;
ALTER TABLE classofcompound_new RENAME TO classofcompound; 

CREATE TABLE empiricalformula_new (
	empiricalformula_id	INTEGER,
	empiricalformula_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(empiricalformula_id)
) STRICT;
INSERT INTO empiricalformula_new (
	empiricalformula_id,
	empiricalformula_label
)
SELECT empiricalformula_id,
	empiricalformula_label
FROM empiricalformula;
DROP TABLE empiricalformula;
ALTER TABLE empiricalformula_new RENAME TO empiricalformula; 

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

CREATE TABLE hazardstatement_new (
	hazardstatement_id	INTEGER,
	hazardstatement_label	TEXT NOT NULL,
	hazardstatement_reference	TEXT NOT NULL UNIQUE,
	hazardstatement_cmr	TEXT,
	PRIMARY KEY(hazardstatement_id)
) STRICT;
INSERT INTO hazardstatement_new (
	hazardstatement_id,
	hazardstatement_label,
	hazardstatement_reference,
	hazardstatement_cmr
)
SELECT hazardstatement_id,
	hazardstatement_label,
	hazardstatement_reference,
	hazardstatement_cmr
FROM hazardstatement;
DROP TABLE hazardstatement;
ALTER TABLE hazardstatement_new RENAME TO hazardstatement; 

CREATE TABLE linearformula_new (
	linearformula_id	INTEGER,
	linearformula_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(linearformula_id)
) STRICT;
INSERT INTO linearformula_new (
	linearformula_id,
	linearformula_label
)
SELECT linearformula_id,
	linearformula_label
FROM linearformula;
DROP TABLE linearformula;
ALTER TABLE linearformula_new RENAME TO linearformula; 

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

CREATE TABLE physicalstate_new (
	physicalstate_id	INTEGER,
	physicalstate_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(physicalstate_id)
) STRICT;
INSERT into physicalstate_new (
	physicalstate_id,
	physicalstate_label
)
SELECT physicalstate_id,
	physicalstate_label
FROM physicalstate;
DROP TABLE physicalstate;
ALTER TABLE physicalstate_new RENAME TO physicalstate; 

CREATE TABLE precautionarystatement_new (
	precautionarystatement_id	INTEGER,
	precautionarystatement_label	TEXT NOT NULL,
	precautionarystatement_reference	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(precautionarystatement_id)
) STRICT;
INSERT INTO precautionarystatement_new (
	precautionarystatement_id,
	precautionarystatement_label,
	precautionarystatement_reference
)
SELECT precautionarystatement_id,
	precautionarystatement_label,
	precautionarystatement_reference
FROM precautionarystatement;
DROP TABLE precautionarystatement;
ALTER TABLE precautionarystatement_new RENAME TO precautionarystatement; 

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

CREATE TABLE producerref_new (
	producerref_id	INTEGER,
	producerref_label	TEXT NOT NULL,
	producer	INTEGER,
	FOREIGN KEY(producer) REFERENCES producer(producer_id),
	PRIMARY KEY(producerref_id)
) STRICT;
INSERT INTO producerref_new (
	producerref_id,
	producerref_label,
	producer
)
SELECT producerref_id,
	producerref_label,
	producer
FROM producerref;
DROP TABLE producerref;
ALTER TABLE producerref_new RENAME TO producerref; 

CREATE TABLE product_new (
	product_id	INTEGER,
	product_inchi TEXT,
	product_inchikey TEXT,
	product_canonical_smiles TEXT,
	product_specificity	TEXT,
	product_msds	TEXT,
	product_restricted	INTEGER DEFAULT 0,
	product_radioactive	INTEGER DEFAULT 0,
	product_threedformula	TEXT,
	product_twodformula	TEXT,
	product_disposalcomment	TEXT,
	product_remark	TEXT,
	product_qrcode	TEXT,
	product_sheet	TEXT,
	product_concentration	REAL,
	product_temperature	REAL,
	product_molecularweight REAL,
	casnumber	INTEGER,
	cenumber	INTEGER,
	person	INTEGER NOT NULL,
	empiricalformula	INTEGER,
	linearformula	INTEGER,
	physicalstate	INTEGER,
	signalword	INTEGER,
	name	INTEGER NOT NULL,
	producerref	INTEGER,
	unit_molecularweight INTEGER,
	unit_temperature	INTEGER,
	category	INTEGER,
	product_number_per_carton	INTEGER,
	product_number_per_bag	INTEGER,
	FOREIGN KEY(person) REFERENCES person(person_id),
	FOREIGN KEY(empiricalformula) REFERENCES empiricalformula(empiricalformula_id),
	FOREIGN KEY(linearformula) REFERENCES linearformula(linearformula_id),
	FOREIGN KEY(casnumber) REFERENCES casnumber(casnumber_id),
	FOREIGN KEY(cenumber) REFERENCES cenumber(cenumber_id),
	FOREIGN KEY(producerref) REFERENCES producerref(producerref_id),
	FOREIGN KEY(category) REFERENCES category(category_id),
	PRIMARY KEY(product_id),
	FOREIGN KEY(unit_temperature) REFERENCES unit(unit_id),
	FOREIGN KEY(unit_molecularweight) REFERENCES unit(unit_id),
	FOREIGN KEY(physicalstate) REFERENCES physicalstate(physicalstate_id),
	FOREIGN KEY(signalword) REFERENCES signalword(signalword_id),
	FOREIGN KEY(name) REFERENCES name(name_id)
) STRICT;
INSERT into product_new (
	product_id,
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

CREATE TABLE signalword_new (
	signalword_id	INTEGER,
	signalword_label	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(signalword_id)
) STRICT;
INSERT INTO signalword_new (
	signalword_id,
	signalword_label
)
SELECT signalword_id,
	signalword_label
FROM signalword;
DROP TABLE signalword;
ALTER TABLE signalword_new RENAME TO signalword; 

CREATE TABLE storage_new (
	storage_id	INTEGER,
	storage_creationdate	INTEGER NOT NULL,
	storage_modificationdate	INTEGER NOT NULL,
	storage_entrydate	INTEGER,
	storage_exitdate	INTEGER,
	storage_openingdate	INTEGER,
	storage_expirationdate	INTEGER,
	storage_quantity	REAL,
	storage_barecode	TEXT,
	storage_comment	TEXT,
	storage_reference	TEXT,
	storage_batchnumber	TEXT,
	storage_todestroy	INTEGER DEFAULT 0,
	storage_archive	INTEGER DEFAULT 0,
	storage_qrcode	BLOB,
	storage_concentration	REAL,
	storage_number_of_unit	INTEGER,
	storage_number_of_bag	INTEGER,
	storage_number_of_carton	INTEGER,
	person	INTEGER NOT NULL,
	product	INTEGER NOT NULL,
	storelocation	INTEGER NOT NULL,
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
	FOREIGN KEY(storelocation) REFERENCES storelocation(storelocation_id),
	PRIMARY KEY(storage_id)
) STRICT;
INSERT INTO storage_new (
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
CAST(unixepoch(storage_creationdate) AS INTEGER) AS storage_creationdate,
CAST(unixepoch(storage_modificationdate) AS INTEGER) AS storage_modificationdate,
CAST(unixepoch(storage_entrydate) AS INTEGER) AS storage_entrydate,
CAST(unixepoch(storage_exitdate) AS INTEGER) AS storage_exitdate,
CAST(unixepoch(storage_openingdate) AS INTEGER) AS storage_openingdate,
CAST(unixepoch(storage_expirationdate) AS INTEGER) AS storage_expirationdate,
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

CREATE TABLE storelocation_new (
	storelocation_id	INTEGER,
	storelocation_name	TEXT NOT NULL,
	storelocation_color	TEXT,
	storelocation_canstore	INTEGER DEFAULT 0,
	storelocation_fullpath	TEXT,
	entity	INTEGER NOT NULL,
	storelocation	INTEGER,
	FOREIGN KEY(storelocation) REFERENCES storelocation(storelocation_id),
	FOREIGN KEY(entity) REFERENCES entity(entity_id),
	PRIMARY KEY(storelocation_id)
) STRICT;
INSERT INTO storelocation_new (
	storelocation_id,
	storelocation_name,
	storelocation_color,
	storelocation_canstore,
	storelocation_fullpath,
	entity,
	storelocation
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
ALTER TABLE storelocation_new RENAME TO storelocation;

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

CREATE TABLE supplierref_new (
	supplierref_id	INTEGER,
	supplierref_label	TEXT NOT NULL,
	supplier	INTEGER,
	FOREIGN KEY(supplier) REFERENCES supplier(supplier_id),
	PRIMARY KEY(supplierref_id)
) STRICT;
INSERT INTO supplierref_new (
	supplierref_id,
	supplierref_label,
	supplier
)
SELECT supplierref_id,
	supplierref_label,
	supplier
FROM supplierref;
DROP TABLE supplierref;
ALTER TABLE supplierref_new RENAME TO supplierref;

CREATE TABLE symbol_new (
	symbol_id	INTEGER,
	symbol_label	TEXT NOT NULL UNIQUE,
	symbol_image	TEXT,
	PRIMARY KEY(symbol_id)
) STRICT;

INSERT INTO symbol_new (symbol_label, symbol_image) VALUES 
	("GHS01", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAInSURBVFiFzdi9b45RGMfxz11KNalS0napgVaE0HhLxEtpVEIQxKAbiUhMJCoisTyriYTBLCZGsRCrP8JkkLB4iVlyDL1EU9refZ7raZ3kJPfrub75/a7zWpVSpJSqaoBSGintlVJarzQKJWojo81sqDS4XKUSlcu3LwkuFyoRLh8qCa49UAlw6VDoRUercOlK4T7WtapcNlQHHmfYmgm1CVO4MOPZBow31V4SVDde4i22xrMeXEfVVK5myB4gh/AQVwPqSljbEe/XYVfd9lOgIvAt3AnAS1iBczgV168wVTdOClSAPcMwzmIg4EbRP+u7behZKF6r9q3BTTzFC1wLO49iD/owHioex2nswGpsnC9uU1BYhUE8R8EH3As1DuIYtmAnDsT9SZwPJScxMp8o9RKRtQHSFUk8jBHcxpPIr95QqC+svIxHGKiVDrM4VqpRSik/qqoaxTecwSe8CUWO4Dve4W6o9xFf8Bl9VVV1RgfoDLXfl1J+LhR0bp+nVRjGZoxhLw7jRNhzIwAKXmMCD/AVDVxsRq3ayY/1GEK/6RF+u+k5cTAUGJoxVk1ionaPnjf568HtD6h9GJunY3RjN7qahfobrEYP9Xv0brUuaoCt+VO7oeYGaydcS5N4u+BSlj3ZcKkLxSy4tiytW4Vr62ak2SBLsn1bbLAl3fDWDbosRwQLBV/WQ5W5IP6LY6h/w6VA5YAl2jez1lrBLlhKaaiqP9cJ5Rf+De5Q3HyidwAAAABJRU5ErkJggg=="),
	("GHS02", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAIvSURBVFiFzdgxbI5BGMDx36uNJsJApFtFIh2QEIkBtYiYJFKsrDaJqYOofhKRMFhsFhNNMBgkTAaD0KUiBomN1EpYBHWGnvi8vu9r7/2eti65vHfv3T3P/57nuXvfuyqlJCRVVQuk1AqRl1LqP9NKpJxbETKjocLgYi0VaLl49wXBxUIFwsVDBcGFQWEAg1FwYZbCGMajLBfmPkzgUZRbw2IKFzGPrRFw/bpvD/bn8jUkXM719f3A9eu+k3iXA/92Bnub2yYx1NgDfbrvXIYZx8dcThjBExxvPOmGltqLIzmuEt63QSVczc+z/2whSw2ThpbajS+4UgOq59O4gYFSuGaByWb8zKvwN8RXXKiBPc7PLaWx3ARqY37O1CBe5/cvO1huVy+ZnfSX7y9MYxRTNeX32lZj+/sXWNfVnV3g1tT/aJeQ5vAGp3L9eXbjTFv7NzzM9VncSSnNF2lp4MqjNYvcxwEcy+0HcQg32/q8Kndl+YrcgM9Z4YdsrZ21PtvxHT9yv1vNgr8cbiIrnMUmbKu177PwVZjLgKPNt4sCOKzF0ww32aF9CA+yxSZKoTqDlVnucI6lMxhpg76OuxhrKr8oIENyXx/xxQKTE/hUkIdLJ1tlRd3TwtF/KtcuSalVVdUwdvQe+Fd6ljhfl9NzRKT5I8cvq/B+xi3vzFfk+FaqbEUPvEtVuipXBIspX9VLlW4Q/8U1VGe4EKgYsED3tefBgt271y7dUlV/ygHpF8bRglXiwx7BAAAAAElFTkSuQmCC"),
	("GHS03", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAJhSURBVFiFzdjNq41RFMfxz/FSXK4iUV4GUjJRMpG8lQnJUF0mTAy8xkApimOEMDA2oSgD5R9gcE0obxO6JkooBroJhRttg7Nu94lzzz37OdvLrt16zrPXs9d3/9Ze+5zzNFJKirRGowlSahaZL6XUe6eZSNGbJeYsDVUMrqxSBZUrn75CcMWgsARTSsEVUwq7sLWUcmXS1wK7iOvd+pcDa5++qVgU1/fxKq4n9wrXa/qW4Bb6MIKE2diPmb3A1VYq7MqAuRQ2YTdeY6CXtNZVaj4uYG0FaLR/D3sc0+vC1d3okwJgsAI0iB+Vz5dxe1TdXLg6UHPCvg2AT2E34VobBaflzD8+2AQPRYqu4kUEPh1KzcKOuPck7CMcQF92nOyVtCquqsg8PI5C2IyHWBFjn8NuzM5Mdu7ZGcGO4k1U5EgF9CNO4QuG4t6x7ALLrhY2RLB9uBMAJ7Ea63A+CuMVlobvidzqzz9fmFtR5jvWtPHZHj4Xww5MNO+vHJNktpTSezxAP26klO618bkZah4JRe/mxslOZSiyLZQ43MHnTPicy1Wr1uavBH6Hsx3Gr+ADZudC1TouKoFv4CX624wtwDBO1oH6HSwDDsvxTetrZ2Hl/jKtg3UYs+pAtQfLg1uldcqPhH2qVanPsL4u1Phg3ayIPdiLg3hu7IAdwqEY24vFtRZdew/wtQLTqW+ptYc7gnWYLPbS8i76jFyo7sBqTFri+T86eS/P/dmV/5W/b7nB/uof3m6D/pNXBBMF/6cvVcaD+C9eQ7WHKwJVBqxg+qp9SvYvy3YtpaZGY+y6QPsJlPiFVobY9AkAAAAASUVORK5CYII="),
	("GHS04", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAFtSURBVFiFzdixLgRBGADgb4UoBYXKA3gAHYU3EKWH0FKuAk8gGm+hvU4oJAqlaDVqiUSCUbjEJe5ys7v/7NnkT66Ynf/bmZ3bf6ZKKQm5qqoGKdUh/aWUugd1Ig2jjugzGhWGix2pwJGLn74gXCwqEBePCsKVQQXgyqE63lcW1eH+8qiW/fSDatFff6iG/faLatB//6jMPKEoLGEX53jEHTaxj0sMcvN1QmEBWzjGLT6QxsQX7nGUO3KNUdjAAa7wOgEyGjdYazqt+auEQzxnQEbjBdtt3rn5nCq3qqp1nGJuStM3XGMwjIc0fKrGV9YKYQXv/o7Ip58X/AQ7WIxaofnLl7Mh5gkX2MNyqb+NrEYjuNXOkMx8jRr3hRoP6wPX6pNUGtfpI14KF1L2RONCC8UoXJHSuiuu6GakbZJetm9Nk/W64c1NOpMjgmnJZ3qoMgnxL46hxuNCUDGwwOkbjawKNqParFXV7++A6xtDLLIHRMAuWAAAAABJRU5ErkJggg=="),
	("GHS05", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAI9SURBVFiFzdhPiI1RGMfxz2FmgYVo1CyYhVJISuPPDmFGNwtW2EzKn4TEELKaa2NNdiyUlWxtUGYptiytrRSNIk3GsbinceO+7vvOPe8dp56677k9v+d7nnPOczonxBhlaSE0QYzNLHoxxt6NZiQma+bQzA2VDS5vpjJmLv/0ZYLLC5URLj9UJrh6oDLA1QfVo1+9UD341w+1QJ2exLAG+3AKq9DAGUxgJ1YsFG6gy9k3lb5uibEZQhjHNnxFxEe8w1ucwxfcxxw24GAIYTkCBmOMTSFIulNCKDxbO4N1gEq/P6SsfMdw6hvFDB5hGS5iLYZSRr/hE6bRAikDV2X6sC4B7MZQ13XSAlufII5iBGOllknFNTWGq6W3PJO4ngYygRu4WSoJJTO1N62TXWnqtmK0BNgR3E4DmsQdjJRJRtlMNXAX53EBz7G6BNgAjmEH9hT6dIhfagtjCV5gC67hUgmoA3iNl3iDe7iSQAe7wRWWixBCAyfauj5jo9bOPBRCeNz239kY40yb72lsx5MEth9PY4zvQwgr8aMo7nwrTCWXtWpVGRv+I1vH8SxlaRxLsblKIQ9J6K/aFXiITV1H1mrTMcbZ9o4QwmGc1CrGr/BTa83N4kGMca5T3PmattAjI4uVKhf9hqtUYPsFV6YS9OJcF9S/weqAq6CXVSynTi2iOfxrFe/Fr96R9+X6VjVYXy+8ZYMuyhNBt+CL+qhSBPFfPEN1hssClQcs4/S1W/GFt0r7fVcsvMBWbb8AgnCJLinP5ycAAAAASUVORK5CYII="),
	("GHS06", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAK6SURBVFiFzdhPiFdVFMDxzxsnpTIFNQpJZjEoVEwQyoRYRP82JUQySAYuZITEDBeCuvPnInXhzGAaRTMtDNyMBElUi9m4ECVGI6FVMOuilUG4COO0eBfmj2/m997vd0e7cOG9y7vnfN8555577i0iQpZWFC0Q0coiLyK677SCSL2VQ2ZuqGxweS2V0XL53ZcJLi9URrj8UJnglgcqA1zXUHgc7+EsvsQxbOwWrluoAXyDHXgkjfXhC7zbDVw3UIO4gFWVgjmIjzqF6xTqhWSVNxcVzAq0cKgTuI4CF6cxgr4lhTOG850siMZQSeEr+B29NcBONZVfDba0pXrmPH+Ll7EX7+OxNP56GtuBy1iXxot5P9Lu5xtADWICn6AHn2M3Av9gc/puOo19jXNp7FlM4kJtfQ3cdx39eAPncS4p3IPXUrD34xkcwUZ8hW24hJXJtUO1Flhtn/M9Pk6JdBc+wx8pLfRgP+5gLV7ET7iRIA+keXtxpk74FFEu6RNl2eikRSrQoijeSop/w3Cyzq94CbeTpQaSu99JMXUN2zGOX7ATP0TE3xUK5nH0VEFUtYiYiojpiPgrIkZxFE/jJp7HFA4rF8R3CXImIvZFxPWIuBsRk5VQiyhsnpVL1w0pk+wVPKeMvQ8T7EVlDI5hSzt51a4sFd1nyoUuLYriOJ5Srrgn8LPSnffQi1X4E5twFa9iBquxRunaf/FjREzNEVytt9YKKS00jrH0vi8J2zDnm0ct2AmU8TehjL+VuIUPmqWL9nCFsqxZj5G27pmdt1a52W/FcGcJtj3cIEbxZF2wNO9tfIrVtdNTk4DM0rvaxJcLLkvZkxsua6GYC25ZSutu4Zb1MNKpkgdyfGuq7IEeeOsqfShXBO2UP9RLlcUg/hfXUNVwWaDygGV039zeW10+NmwRLUUx+5yh/QdzLVcJBJ5ddQAAAABJRU5ErkJggg=="),
	("GHS07", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAF/SURBVFiFzdghT8NAGMbx/y1g0HMLQWDxJDgCkk+AQhE8BlkDGL4JwczyAXBkCY4vQAhBgiI7xJpwK1fW9nm6ccklNe/7/u7eW9c2xBixjBAKAGIsLPlijPqEIkIsZ+HI6UbZcN6dMu6cv30mnBdlxPlRJlw/KAOuP5QYJycHBsARMAYmwCUwVHHyioELIFbmvbpzchuAmwzsSW2rfDaA0wzsTs4rrwz2M7BruRPyWYBRBnYin10FVcIC8FGB7TWJ/fPXrqAS3KQCGzaOr6kro0rYbYJ6bxufq+/5w4WrBPbQJUfVMZAfgWfjuea6+zC1cht4Y7Zjx55WGg5/idsEDhyoeZiIA7aAXQfqN6wjDjgHvoApi+76Det0CsrAXpi/j+0oqHpYSxzwmKCmwEjNK213AjsEXoFP4MyyWPUsJLh1YMOBagbrkNQR32tyJa7flS/l9a1tsaW+8DYtupJPBIuKr/SjSh3iX3yGyuMsKA/M2L50rpmeNgtC+Lk2jG/Rx4o589viKwAAAABJRU5ErkJggg=="),
	("GHS08", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAKJSURBVFiFzdixix1FGADw38QDUbQQqyNewiEYCxNN4ECjkuIaFYLNQQpTCBLkzD8gJsVLOqPGJiBYeGB1SeFhYXEhmCIpTgIpRDAEosWhckXU6AknCUyKN49bn+/tze6OnAsDuzA789v5vvnevg0xRkWOEHogxl6R8WKM3Ru9SEytV2LM0qhiuLIrVXDlyoevEK4sqiCuGApTeBo7SuA6o7AX3yKmdh3PdsV1RU3hzwpq0NbxZBdc1/AtjEAN2lddcq41KsG+qYH9jYfb4lqjEuxWDSxissl4ebCMQbBSg/qp7bjjYbk3c7oGttjpoduiEuz5GthbndKkLaqCu5Qgd/FdOv8FD3XaWF1QCfZawqzhWDo/mnv/2FLUEfUA3kuYL/FyOn87GzZm/rbhm8JH+LmSU6t4pHL9A87gmTa4Hf9+px19hBAmQgjzIYQVnMIGHqt0+QKvVq6ncRAzIYSlEMLZEMKu3PmyQ4kTldWI+lX/MD7GZ5jHJH7ESbyJs/qbYnDP79idF8rc7Tu6mK7hHexJfQLewE7jfxWO5Cd/zvblypiJ/sBzQ32vjun7T1hWudjqCdiH20OTrOJz7MJTqd8reAk3RqCu4fGsCDWqLezHr2mSdcxgFss4nkK4iPcxh8v6myTiazyandONCx8H8Ck+wM3KapxP+Ta43tDfqXM4hweb1MzGha+CuzMUpoWU+MPh+yR3g+XD6nEvpMQfAJaxNIT6EKEpKg9Wj3vR5jv/b/r1a4B6t81KNYPV4w7hL3yfcu+e6ivPf/pnZGvcLC7iAl7vimoOq8dN44kSqHawnEm35RPBVpNv60eVcYj/xWeo0bgiqDKwguGrtgkljhh7Qtg8L3DcB497IINNg8B2AAAAAElFTkSuQmCC"),
	("GHS09", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAI8SURBVFiFzdi7a1VBEMDhb0EjYhFfaCEICopBi4haWUiqqGgQW1OKaBpNIQEVvREi2lhaWCiC/4CSImgv2IgQELSSdKKlMUKEtTgr3jxuPI+9iQMD5+y9u/NjdmbO7IYYoywSQgvE2MqyXoyxudKKxKStHGvmhsoGl9dTGT2Xf/syweWFygiXHyoTXHegMsDVgprjLk5hS4wRenLD1fYUhjGans9jF/rRlwOuMhS2Yz924F0CO4EBjOTyXK0YwUVcwAcEHMJr9OaKudqLYBA/cQx7cT+NB1xrCld5cvLOGA7gqSIRricv7sED9Df1XB1PbcVOXMZnzOA7pvEcg3UTqjNY1QDlKmLSsRX+14NNlXamLlQyeDJBfenw+xmcxQ3cxunSCVYpINmYatUfHUlg3xaN9+MmjqR5uxPkaNltrZbChcFYUi8tmjuOqbIxt25pT7uifMThtvcBPMS84kvwKY2fizE+hhDCBkXGzuBZaUsN4qtPUVQj3uOJIhmGFSXlOO4loEmsr5KhtYIfvRjCUTxS1K8reKnIwAkLt/UHhqqUjUblIkFuxi18TRCv8KsNah5v8UbqRqqVi5JwOKhI/ym8wKzlg3+xzmGyrJ2QjC2U4ox4J72NS2fFEMI27CsdwEtlNsY43Wn9BVIlILNoo494t+CytD254bI2irngutJaN4Xr6mGkrpFVOb5VNbaqB96yRtfkiuBfxtf0UqUTxH9xDbU8XBaoPGAZt69dq3awy0uMLSH8fc4gvwFyuYuihNiCxwAAAABJRU5ErkJggg==");

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

CREATE TABLE welcomeannounce_new (
	welcomeannounce_id	INTEGER,
	welcomeannounce_text	TEXT,
	PRIMARY KEY(welcomeannounce_id)
) STRICT;
INSERT INTO welcomeannounce_new (
	welcomeannounce_id,
	welcomeannounce_text
)
SELECT welcomeannounce_id,
	welcomeannounce_text
FROM welcomeannounce;
DROP TABLE welcomeannounce;
ALTER TABLE welcomeannounce_new RENAME TO welcomeannounce;

DROP INDEX IF EXISTS idx_producerref_label;
DROP INDEX IF EXISTS idx_supplierref_label;

DROP TABLE IF EXISTS captcha;
DROP TABLE IF EXISTS entityldapgroups;

INSERT INTO unit (unit_label, unit_multiplier, unit_type) VALUES ("g/mol", 1, "molecularweight");

COMMIT;

PRAGMA user_version=10;
PRAGMA foreign_keys=on;`
