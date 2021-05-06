package datastores

var versionToMigration = []string{migrationOne, migrationTwo, migrationThree}

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
