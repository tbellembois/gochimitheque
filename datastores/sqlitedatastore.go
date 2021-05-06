package datastores

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/data"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	. "github.com/tbellembois/gochimitheque/models"
)

// SQLiteDataStore implements the Datastore interface
// to store data in SQLite3
type SQLiteDataStore struct {
	*sqlx.DB
}

var (
	regex = func(re, s string) bool {
		m, _ := regexp.MatchString(re, s)
		return m
	}
)

func init() {
	sql.Register("sqlite3_with_go_func",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("regexp", regex, true)
			},
		})
}

// GetWelcomeAnnounce returns the welcome announce
func (db *SQLiteDataStore) GetWelcomeAnnounce() (WelcomeAnnounce, error) {
	var (
		wa   WelcomeAnnounce
		sqlr string
		err  error
	)

	sqlr = `SELECT welcomeannounce.welcomeannounce_id, welcomeannounce.welcomeannounce_text
	FROM welcomeannounce LIMIT 1`
	if err = db.Get(&wa, sqlr); err != nil {
		return WelcomeAnnounce{}, err
	}

	logger.Log.WithFields(logrus.Fields{"wa": wa}).Debug("GetWelcomeAnnounce")
	return wa, nil
}

// UpdateWelcomeAnnounce updates the main page announce
func (db *SQLiteDataStore) UpdateWelcomeAnnounce(w WelcomeAnnounce) error {
	var (
		sqlr string
		tx   *sqlx.Tx
		err  error
	)

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	// updating person
	sqlr = `UPDATE welcomeannounce SET welcomeannounce_text = ?
	WHERE welcomeannounce_id = (SELECT welcomeannounce_id FROM welcomeannounce LIMIT 1)`
	if _, err = tx.Exec(sqlr, w.WelcomeAnnounceText); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	return nil
}

// NewSQLiteDBstore returns a database connection to the given dataSourceName
// ie. a path to the sqlite database file
func NewSQLiteDBstore(dataSourceName string) (*SQLiteDataStore, error) {
	var (
		db  *sqlx.DB
		err error
	)

	logger.Log.WithFields(logrus.Fields{"dbdriver": "sqlite3", "dataSourceName": dataSourceName}).Debug("NewDBstore")
	if db, err = sqlx.Connect("sqlite3_with_go_func", dataSourceName+"?_journal=wal&_fk=1"); err != nil {
		return &SQLiteDataStore{}, err
	}

	return &SQLiteDataStore{db}, nil
}

// ToCasbinJSONAdapter returns a JSON as a slice of bytes
// following the format: https://github.com/casbin/json-adapter#policy-json
func (db *SQLiteDataStore) ToCasbinJSONAdapter() ([]byte, error) {
	var (
		ps   []Permission
		js   []CasbinJSON
		err  error
		res  []byte
		sqlr string
	)

	sqlr = `SELECT person AS "person.person_id", permission_perm_name, permission_item_name, permission_entity_id 
	FROM permission`
	if err = db.Select(&ps, sqlr); err != nil {
		return nil, err
	}

	for _, p := range ps {
		js = append(js, models.CasbinJSON{
			PType: "p",
			V0:    strconv.Itoa(p.Person.PersonID),
			V1:    p.PermissionPermName,
			V2:    p.PermissionItemName,
			V3:    strconv.Itoa(p.PermissionEntityID),
		})
	}

	if res, err = json.Marshal(js); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateDatabase creates the database tables
func (db *SQLiteDataStore) CreateDatabase() error {
	var (
		err         error
		c           int
		userVersion int
		r           *csv.Reader
		records     [][]string
	)

	// tables creation
	logger.Log.Info("  creating sqlite tables")
	if _, err = db.Exec(schema); err != nil {
		return err
	}

	// shema migration
	if err = db.Get(&userVersion, `PRAGMA user_version`); err != nil {
		return err
	}
	logger.Log.Info(fmt.Sprintf("  user_version:%d", userVersion))

	nextVersion := userVersion + 1
	for _, version := range versionToMigration[userVersion:] {

		logger.Log.Infof("  upgrading version to %d ", nextVersion)
		if _, err = db.Exec(version); err != nil {
			return err
		}
		nextVersion++

	}

	// welcome announce
	if err = db.Get(&c, `SELECT count(*) FROM welcomeannounce`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting welcome announce")
		if _, err = db.Exec(inswelcomeannounce); err != nil {
			return err
		}
	}

	// symbols
	if err = db.Get(&c, `SELECT count(*) FROM symbol`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting symbols")
		if _, err = db.Exec(inssymbol); err != nil {
			return err
		}
	}

	// signal words
	if err = db.Get(&c, `SELECT count(*) FROM signalword`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting signal words")
		if _, err = db.Exec(inssignalword); err != nil {
			return err
		}
	}

	// cas numbers
	if err = db.Get(&c, `SELECT count(*) FROM casnumber`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting CMRs")
		r = csv.NewReader(strings.NewReader(data.CMR_CAS))
		r.Comma = ','
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO casnumber (casnumber_label, casnumber_cmr) VALUES (?, ?)`, record[0], record[1]); err != nil {
				return err
			}
		}
	}

	// tags
	if err = db.Get(&c, `SELECT count(*) FROM tag`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting tags")
		r = csv.NewReader(strings.NewReader(data.TAG))
		r.Comma = ','
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO tag (tag_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// categories
	if err = db.Get(&c, `SELECT count(*) FROM category`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting categories")
		r = csv.NewReader(strings.NewReader(data.CATEGORY))
		r.Comma = ';'
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO category (category_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// suppliers
	if err = db.Get(&c, `SELECT count(*) FROM supplier`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting suppliers")
		r = csv.NewReader(strings.NewReader(data.SUPPLIER))
		r.Comma = ','
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO supplier (supplier_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// producers
	if err = db.Get(&c, `SELECT count(*) FROM producer`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting producers")
		r = csv.NewReader(strings.NewReader(data.PRODUCER))
		r.Comma = ','
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO producer (producer_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// hazard statements
	if err = db.Get(&c, `SELECT count(*) FROM hazardstatement`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting hazard statements")
		r = csv.NewReader(strings.NewReader(data.HAZARDSTATEMENT))
		r.Comma = '\t'
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference, hazardstatement_cmr) VALUES (?, ?, ?)`, record[0], record[1], record[2]); err != nil {
				return err
			}
		}
	}

	// precautionary statements
	if err = db.Get(&c, `SELECT count(*) FROM precautionarystatement`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting precautionary statements")
		r = csv.NewReader(strings.NewReader(data.PRECAUTIONARYSTATEMENT))
		r.Comma = '\t'
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO precautionarystatement (precautionarystatement_label, precautionarystatement_reference) VALUES (?, ?)`, record[0], record[1]); err != nil {
				return err
			}
		}
	}

	// inserting default admin
	var admin *Person
	if err = db.Get(&c, `SELECT count(*) FROM person`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting admin user")
		admin = &Person{PersonEmail: "admin@chimitheque.fr", Permissions: []*Permission{{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: -1}}}
		var insertId int64
		insertId, _ = db.CreatePerson(*admin)
		admin.PersonPassword = "chimitheque"
		admin.PersonID = int(insertId)
		if err = db.UpdatePersonPassword(*admin); err != nil {
			return err
		}
	}

	// inserting sample entity
	if err = db.Get(&c, `SELECT count(*) FROM entity`); err != nil {
		return err
	}
	if c == 0 {
		logger.Log.Info("  inserting sample entity")
		sentity := Entity{EntityName: "sample entity", EntityDescription: "you can delete me, I am just a sample entity", Managers: []*Person{admin}}
		if _, err = db.CreateEntity(sentity); err != nil {
			return err
		}
	}

	// tables creation
	logger.Log.Info("  vacuuming database")
	if _, err = db.Exec("VACUUM;"); err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDataStore) Maintenance() {

	var (
		err  error
		sqlr string
		tx   *sql.Tx
	)

	//
	// Cleaning up casnumber labels duplicates.
	//
	if tx, err = db.Begin(); err != nil {
		logger.Log.Error(err)
		return
	}

	var casNumbers []CasNumber
	sqlr = `SELECT casnumber_id, casnumber_label FROM casnumber;`
	if err = db.Select(&casNumbers, sqlr); err != nil {
		logger.Log.Error(err)
		return
	}

	for _, casNumber := range casNumbers {

		if strings.HasPrefix(casNumber.CasNumberLabel.String, " ") || strings.HasSuffix(casNumber.CasNumberLabel.String, " ") {
			logger.Log.Infof("casnumber %s contains spaces", casNumber.CasNumberLabel.String)

			trimmedLabel := strings.Trim(casNumber.CasNumberLabel.String, " ")

			// Checking if the trimmed label already exists.
			var existCasNumber CasNumber
			sqlr = `SELECT casnumber_id, casnumber_label FROM casnumber WHERE casnumber_label=?;`
			if err = db.Get(&existCasNumber, sqlr, trimmedLabel); err != nil {
				switch err {
				case sql.ErrNoRows:
					// Just fixing the label.
					logger.Log.Info("  - fixing it")
					sqlr = `UPDATE casnumber SET casnumber_label=? WHERE casnumber_id=?;`
					if _, err = tx.Exec(sqlr, trimmedLabel, casNumber.CasNumberID); err != nil {
						logger.Log.Error(err)
						if errr := tx.Rollback(); errr != nil {
							logger.Log.Error(err)
							return
						}
						return
					}
					continue
				default:
					logger.Log.Error(err)
					return
				}
			}

			// Updating products with the found casnumber.
			logger.Log.Infof("  - correct cas number found, replacing it: %d -> %d", existCasNumber.CasNumberID.Int64, casNumber.CasNumberID.Int64)
			sqlr = `UPDATE product SET casnumber=? WHERE casnumber=?;`
			if _, err = tx.Exec(sqlr, existCasNumber.CasNumberID, casNumber.CasNumberID); err != nil {
				logger.Log.Error(err)
				if errr := tx.Rollback(); errr != nil {
					logger.Log.Error(err)
					return
				}
				return
			}

			// Deleting the wrong cas number.
			logger.Log.Info("  - deleting it")
			sqlr = `DELETE FROM casnumber WHERE casnumber_id=?;`
			if _, err = tx.Exec(sqlr, casNumber.CasNumberID); err != nil {
				logger.Log.Error(err)
				if errr := tx.Rollback(); errr != nil {
					logger.Log.Error(err)
					return
				}
				return
			}

		}

	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error(err)
		if errr := tx.Rollback(); errr != nil {
			logger.Log.Error(errr)
			return
		}
	}

}

// Import import data from another Chimithèque instance
func (db *SQLiteDataStore) Import(url string) error {

	type r struct {
		Rows  []Product `json:"rows"`
		Total int       `json:"total"`
	}

	var (
		err         error
		httpresp    *http.Response
		bodyresp    r
		admin       Person
		notimported int
	)

	logger.Log.Info("- gathering remote products from " + url + "/e/products")
	if httpresp, err = http.Get(url + "/e/products"); err != nil {
		logger.Log.Error("can not get remote products " + err.Error())
	}
	defer httpresp.Body.Close()

	logger.Log.Info("- decoding response")
	if err = json.NewDecoder(httpresp.Body).Decode(&bodyresp); err != nil {
		logger.Log.Error("can not decode remote response " + err.Error())
	}
	logger.Log.Info(fmt.Sprintf("  found %d products", bodyresp.Total))

	logger.Log.Info("- retrieving default admin")
	if admin, err = db.GetPersonByEmail("admin@chimitheque.fr"); err != nil {
		logger.Log.Error("can not get default admin " + err.Error())
		os.Exit(1)
	}

	logger.Log.Info("- starting import")
	for _, p := range bodyresp.Rows {

		// cas number already exist ?
		if p.CasNumberID.Valid {
			var casnumber CasNumber
			if casnumber, err = db.GetProductsCasNumberByLabel(p.CasNumberLabel.String); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product cas number " + err.Error())
					os.Exit(1)
				}
			}
			// new cas number
			if casnumber == (CasNumber{}) {
				// setting cas number id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.CasNumber.CasNumberID = sql.NullInt64{Valid: true, Int64: -1}
			} else {
				// do not insert products with existing cas number
				notimported++
				continue
			}
		}

		// ce number already exist ?
		if p.CeNumberID.Valid {
			var cenumber CeNumber
			if cenumber, err = db.GetProductsCeNumberByLabel(p.CeNumberLabel.String); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product ce number " + err.Error())
					os.Exit(1)
				}
			}
			// new ce number
			if cenumber == (CeNumber{}) {
				// setting ce number id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.CeNumber.CeNumberID = sql.NullInt64{Valid: true, Int64: -1}
			} else {
				p.CeNumber = cenumber
			}
		}

		// empirical formula already exist ?
		if p.EmpiricalFormula.EmpiricalFormulaID.Valid {
			var eformula EmpiricalFormula
			if eformula, err = db.GetProductsEmpiricalFormulaByLabel(p.EmpiricalFormula.EmpiricalFormulaLabel.String); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product empirical formula " + err.Error())
					os.Exit(1)
				}
			}
			// new empirical formula
			if eformula == (EmpiricalFormula{}) {
				// setting empirical formula id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.EmpiricalFormula.EmpiricalFormulaID = sql.NullInt64{Valid: true, Int64: -1}
			} else {
				p.EmpiricalFormula = eformula
			}
		}

		// linear formula already exist ?
		if p.LinearFormula.LinearFormulaID.Valid {
			var lformula LinearFormula
			if lformula, err = db.GetProductsLinearFormulaByLabel(p.LinearFormula.LinearFormulaLabel.String); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product linear formula " + err.Error())
					os.Exit(1)
				}
			}
			// new linear formula
			if lformula == (LinearFormula{}) {
				// setting linear formula id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.LinearFormula.LinearFormulaID = sql.NullInt64{Valid: true, Int64: -1}
			} else {
				p.LinearFormula = lformula
			}
		}

		// physical state already exist ?
		if p.PhysicalState.PhysicalStateID.Valid {
			var physicalstate PhysicalState
			if physicalstate, err = db.GetProductsPhysicalStateByLabel(p.PhysicalState.PhysicalStateLabel.String); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product physical state " + err.Error())
					os.Exit(1)
				}
			}
			// new physical state
			if physicalstate == (PhysicalState{}) {
				// setting physical state id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.PhysicalState.PhysicalStateID = sql.NullInt64{Valid: true, Int64: -1}
			} else {
				p.PhysicalState = physicalstate
			}
		}

		// signal word already exist ?
		if p.SignalWord.SignalWordID.Valid {
			var signalword SignalWord
			if signalword, err = db.GetProductsSignalWordByLabel(p.SignalWord.SignalWordLabel.String); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product signal word " + err.Error())
					os.Exit(1)
				}
			}
			// new signal word
			if signalword == (SignalWord{}) {
				// setting signal word id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.SignalWord.SignalWordID = sql.NullInt64{Valid: true, Int64: -1}
			} else {
				p.SignalWord = signalword
			}
		}

		// name already exist ?
		var name Name
		if name, err = db.GetProductsNameByLabel(p.Name.NameLabel); err != nil {
			if err != sql.ErrNoRows {
				logger.Log.Error("can not get product name " + err.Error())
				os.Exit(1)
			}
		}
		// new name
		if name == (Name{}) {
			// setting name id to -1 for the CreateProduct method
			// to automatically insert it into the db
			p.Name.NameID = -1
		} else {
			p.Name = name
		}

		// synonyms
		var (
			processedSyn map[string]string
			newSyn       []Name
			ok           bool
		)
		// duplicate names map
		processedSyn = make(map[string]string)
		processedSyn[p.Name.NameLabel] = ""
		for _, syn := range p.Synonyms {
			// duplicates hunting
			if _, ok = processedSyn[syn.NameLabel]; ok {
				logger.Log.Debug("leaving duplicate synonym " + syn.NameLabel)
				continue
			}

			processedSyn[syn.NameLabel] = ""

			// synonym already exist ?
			var syn2 Name
			if syn2, err = db.GetProductsNameByLabel(syn.NameLabel); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product synonym " + err.Error())
					os.Exit(1)
				}
			}
			// new synonym
			if syn2 == (Name{}) {
				// setting synonym id to -1 for the CreateProduct method
				// to automatically insert it into the db
				newSyn = append(newSyn, Name{NameID: -1, NameLabel: syn.NameLabel})
			} else {
				newSyn = append(newSyn, syn2)
			}
		}
		p.Synonyms = newSyn

		// classes of compounds
		for i, coc := range p.ClassOfCompound {
			// class of compounds already exist ?
			var coc2 ClassOfCompound
			if coc2, err = db.GetProductsClassOfCompoundByLabel(coc.ClassOfCompoundLabel); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product class of compounds " + err.Error())
					os.Exit(1)
				}
			}
			// new class of compounds
			if coc2 == (ClassOfCompound{}) {
				// setting class of compounds id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.ClassOfCompound[i].ClassOfCompoundID = -1
			} else {
				p.ClassOfCompound[i] = coc2
			}
		}

		// symbols
		for i, sym := range p.Symbols {
			// symbols already exist ?
			var sym2 Symbol
			if sym2, err = db.GetProductsSymbolByLabel(sym.SymbolLabel); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product symbol " + err.Error())
					os.Exit(1)
				}
			}
			// new symbol
			if sym2 == (Symbol{}) {
				// setting symbol id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.Symbols[i].SymbolID = -1
			} else {
				p.Symbols[i] = sym2
			}
		}

		// hazard statements
		for i, hs := range p.HazardStatements {
			// hazard statement already exist ?
			var hs2 HazardStatement
			if hs2, err = db.GetProductsHazardStatementByReference(hs.HazardStatementReference); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product hazard statement " + err.Error())
					os.Exit(1)
				}
			}
			// new hazard statement
			if hs2 == (HazardStatement{}) {
				// setting hazard statement id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.HazardStatements[i].HazardStatementID = -1
			} else {
				p.HazardStatements[i] = hs2
			}
		}

		// precautionnary statements
		for i, ps := range p.PrecautionaryStatements {
			// precautionary statement already exist ?
			var ps2 PrecautionaryStatement
			if ps2, err = db.GetProductsPrecautionaryStatementByReference(ps.PrecautionaryStatementReference); err != nil {
				if err != sql.ErrNoRows {
					logger.Log.Error("can not get product precautionary statement " + err.Error())
					os.Exit(1)
				}
			}
			// new precautionary statement
			if ps2 == (PrecautionaryStatement{}) {
				// setting precautionary statement id to -1 for the CreateProduct method
				// to automatically insert it into the db
				p.PrecautionaryStatements[i].PrecautionaryStatementID = -1
			} else {
				p.PrecautionaryStatements[i] = ps2
			}
		}

		// setting default admin as creator
		p.Person = admin

		// finally creating the product
		if _, err = db.CreateProduct(p); err != nil {
			logger.Log.Error("can not create product " + err.Error())
			os.Exit(1)
		}

	}

	logger.Log.Info(fmt.Sprintf("%d products not imported (duplicates)", notimported))

	return nil
}

// CSVToMap takes a reader and returns an array of dictionaries, using the header row as the keys
// credit: https://gist.github.com/drernie/5684f9def5bee832ebc50cabb46c377a
func CSVToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		logger.Log.Debug(fmt.Sprintf("record: %s", record))
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error(err)
			return nil
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}

// ImportV1 import data from CSV
func (db *SQLiteDataStore) ImportV1(dir string) error {

	var (
		csvFile *os.File
		//csvReader *csv.Reader
		csvMap []map[string]string
		err    error
		res    sql.Result
		lastid int64
		c      int      // count result
		tx     *sqlx.Tx // db transaction
		sqlr   string   // sql request

		zerocasnumberid        int
		zeroempiricalformulaid int
		zeropersonid           int // admin id
		zerohsid               string
		zeropsid               string

		// ids mappings
		// O:old N:new R:reverse
		mONperson        map[string]string   // oldid <> newid map for user table
		mONsupplier      map[string]string   // oldid <> newid map for supplier table
		mONunit          map[string]string   // oldid <> newid map for unit table
		mONentity        map[string]string   // oldid <> newid map for entity table
		mONstorelocation map[string]string   // oldid <> newid map for storelocation table
		mOOentitypeople  map[string][]string // managers, oldentityid <> oldpersonid
		mRNNcasnumber    map[string]string   // newlabel <> newid
		mRNNcenumber     map[string]string   // newlabel <> newid

		mONproduct                map[string]string // oldid <> newid map for product table
		mONclassofcompound        map[string]string // oldid <> newid map for classofcompound table
		mONempiricalformula       map[string]string // oldid <> newid map for empiricalformula table
		mONlinearformula          map[string]string // oldid <> newid map for linearformula table
		mONname                   map[string]string // oldid <> newid map for name table
		mONphysicalstate          map[string]string // oldid <> newid map for physicalstate table
		mONhazardstatement        map[string]string // oldid <> newid map for hazardstatement table
		mONprecautionarystatement map[string]string // oldid <> newid map for precautionarystatement table
		mONsymbol                 map[string]string // oldid <> newid map for symbol table
		mONsignalword             map[string]string // oldid <> newid map for signalword table

	)

	// init maps
	mONproduct = make(map[string]string)
	mONperson = make(map[string]string)
	mONunit = make(map[string]string)
	mONsupplier = make(map[string]string)
	mONentity = make(map[string]string)
	mONstorelocation = make(map[string]string)
	mOOentitypeople = make(map[string][]string)
	mRNNcasnumber = make(map[string]string)
	mRNNcenumber = make(map[string]string)
	mONclassofcompound = make(map[string]string)
	mONempiricalformula = make(map[string]string)
	mONlinearformula = make(map[string]string)
	mONname = make(map[string]string)
	mONphysicalstate = make(map[string]string)
	mONhazardstatement = make(map[string]string)
	mONprecautionarystatement = make(map[string]string)
	mONsymbol = make(map[string]string)
	mONsignalword = make(map[string]string)

	// number regex
	rnumber := regexp.MustCompile("([0-9]+)")

	// checking tables empty
	if err = db.Get(&c, `SELECT count(*) FROM product`); err != nil {
		return err
	}
	if c != 0 {
		panic("person product not empty - can not import")
	}

	// beginning transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	//
	// entity
	//
	logger.Log.Info("- importing entity")
	rentityName := regexp.MustCompile("user_[0-9]+|root_entity|all_entity")
	if csvFile, err = os.Open(path.Join(dir, "entity.csv")); err != nil {
		return (err)
	}

	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		role := k["role"]
		description := k["description"]
		manager := k["manager"]

		// finding web2py like manager ids
		ms := rnumber.FindAllString(manager, -1)
		for _, m := range ms {
			// leaving hardcoded zeros
			if m != "0" {
				mOOentitypeople[id] = append(mOOentitypeople[id], m)
				logger.Log.Debug("entity with old id " + id + " has manager with old id " + m)
			}
		}

		// leaving web2py specific entries
		if !rentityName.MatchString(role) {
			logger.Log.Debug("  " + role)
			sqlr = `INSERT INTO entity(entity_name, entity_description) VALUES (?, ?)`
			if res, err = tx.Exec(sqlr, role, description); err != nil {
				logger.Log.Error("error importing entity " + role)
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// populating the map
			mONentity[id] = strconv.FormatInt(lastid, 10)
			logger.Log.Debug("entity with old id " + id + " has new  id " + strconv.FormatInt(lastid, 10))
		}
	}

	//
	// storelocation
	//
	logger.Log.Info("- importing store locations")
	if csvFile, err = os.Open(path.Join(dir, "store_location.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		entity := k["entity"]
		parent := k["parent"]
		canStore := k["can_store"]
		color := k["color"]

		newentity := mONentity[entity]
		newparent := sql.NullString{}
		np := mONstorelocation[parent]
		if np != "" {
			newparent = sql.NullString{Valid: true, String: np}
		}
		newcanStore := false
		if canStore == "T" {
			newcanStore = true
		}
		logger.Log.Debug("storelocation " + label + ", entity:" + newentity + ", parent:" + newparent.String)
		sqlr = `INSERT INTO storelocation(storelocation_name, storelocation_color, storelocation_canstore, storelocation_fullpath, entity, storelocation) VALUES (?, ?, ?, ?, ?, ?)`
		if res, err = tx.Exec(sqlr, label, color, newcanStore, "", newentity, newparent); err != nil {
			logger.Log.Error("error importing storelocation " + label)
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the map
		mONstorelocation[id] = strconv.FormatInt(lastid, 10)
	}

	//
	// person
	//
	logger.Log.Info("- importing user")
	if csvFile, err = os.Open(path.Join(dir, "person.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		email := k["email"]
		password := k["password"]

		sqlr = `INSERT INTO person(person_email, person_password) VALUES (?, ?)`
		if res, err = tx.Exec(sqlr, email, password); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the map
		mONperson[id] = strconv.FormatInt(lastid, 10)
	}

	//
	// permissions
	//
	logger.Log.Info("- initializing default permissions (r products)")
	for _, newpid := range mONperson {
		sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) VALUES (?, ?, ?, ?)`
		if _, err = tx.Exec(sqlr, newpid, "r", "products", -1); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	//
	// managers
	//
	logger.Log.Info("- importing managers")
	for oldentityid, oldmanagerids := range mOOentitypeople {
		for _, oldmanagerid := range oldmanagerids {
			newentityid := mONentity[oldentityid]
			newmanagerid := mONperson[oldmanagerid]
			// silently missing entities with no managers
			if newmanagerid != "" {
				sqlr = `INSERT INTO entitypeople(entitypeople_entity_id, entitypeople_person_id) VALUES (?, ?)`
				if _, err = tx.Exec(sqlr, newentityid, newmanagerid); err != nil {
					if errr := tx.Rollback(); errr != nil {
						return errr
					}
					return err
				}
				logger.Log.Debug("person "+newmanagerid+", permission_perm_name: all permission_item_name: all", " permission_entity_id:"+newentityid)
				sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) VALUES (?, ?, ?, ?)`
				if _, err = tx.Exec(sqlr, newmanagerid, "all", "all", newentityid); err != nil {
					if errr := tx.Rollback(); errr != nil {
						return errr
					}
					return err
				}
			}
		}
	}

	//
	// membership
	//
	logger.Log.Info("- importing membership")
	if csvFile, err = os.Open(path.Join(dir, "membership.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		userId := k["user_id"]
		groupId := k["group_id"]
		newuserId := mONperson[userId]
		newgroupId := mONentity[groupId]

		if newuserId != "" && newgroupId != "" {
			sqlr = `INSERT INTO personentities(personentities_person_id, personentities_entity_id) VALUES (?, ?)`
			if _, err = tx.Exec(sqlr, newuserId, newgroupId); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) VALUES (?, ?, ?, ?)`
			if _, err = tx.Exec(sqlr, newuserId, "r", "entities", newgroupId); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
		}
	}

	//
	// class of compounds
	//
	logger.Log.Info("- importing classes of compounds")
	if csvFile, err = os.Open(path.Join(dir, "class_of_compounds.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]

		sqlr = `INSERT INTO classofcompound(classofcompound_id, classofcompound_label) VALUES (?, ?)`
		if res, err = tx.Exec(sqlr, id, label); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the map
		mONclassofcompound[id] = strconv.FormatInt(lastid, 10)
	}

	//
	// empirical formula
	//
	logger.Log.Info("- importing empirical formulas")
	if csvFile, err = os.Open(path.Join(dir, "empirical_formula.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		if label == "----" {
			continue
		}

		sqlr = `INSERT INTO empiricalformula(empiricalformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, label); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the map
		mONempiricalformula[id] = strconv.FormatInt(lastid, 10)
	}

	//
	// linear formula
	//
	logger.Log.Info("- importing linear formulas")
	if csvFile, err = os.Open(path.Join(dir, "linear_formula.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		if label == "----" {
			continue
		}

		sqlr = `INSERT INTO linearformula(linearformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, label); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the map
		mONlinearformula[id] = strconv.FormatInt(lastid, 10)
	}

	//
	// name
	//
	logger.Log.Info("- importing product names")
	if csvFile, err = os.Open(path.Join(dir, "name.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		label = strings.Replace(label, "@", "_", -1)

		logger.Log.Debug("label:" + label)
		sqlr = `INSERT INTO name(name_id, name_label) VALUES (?, ?)`
		if res, err = tx.Exec(sqlr, id, label); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the maps
		mONname[id] = strconv.FormatInt(lastid, 10)
	}

	//
	// physical states
	//
	logger.Log.Info("- importing product physical states")
	if csvFile, err = os.Open(path.Join(dir, "physical_state.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]

		sqlr = `INSERT INTO physicalstate(physicalstate_id, physicalstate_label) VALUES (?, ?)`
		if res, err = tx.Exec(sqlr, id, label); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the map
		mONphysicalstate[id] = strconv.FormatInt(lastid, 10)
	}

	//
	// cas numbers
	//
	logger.Log.Info("- extracting and importing cas numbers from products")
	logger.Log.Info("  gathering existing CMR cas numbers")
	var (
		rows     *sql.Rows
		casid    string
		caslabel string
	)
	if rows, err = tx.Query(`SELECT casnumber_id, casnumber_label FROM casnumber`); err != nil {
		logger.Log.Error("error gathering existing CMR cas numbers")
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	for rows.Next() {
		err := rows.Scan(&casid, &caslabel)
		if err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		mRNNcasnumber[caslabel] = casid
	}
	if csvFile, err = os.Open(path.Join(dir, "product.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {

		casnumber := k["cas_number"]
		logger.Log.Debug(fmt.Sprintf("casnumber: %s", casnumber))
		if _, ok := mRNNcasnumber[casnumber]; !ok {
			sqlr = `INSERT INTO casnumber(casnumber_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, casnumber); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// populating the map
			mRNNcasnumber[casnumber] = strconv.FormatInt(lastid, 10)
		}
	}

	//
	// ce numbers
	//
	logger.Log.Info("- extracting and importing ce numbers from products")
	if csvFile, err = os.Open(path.Join(dir, "product.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {

		cenumber := k["ce_number"]
		if cenumber != "" {
			if _, ok := mRNNcenumber[cenumber]; !ok {
				sqlr = `INSERT INTO cenumber(cenumber_label) VALUES (?)`
				if res, err = tx.Exec(sqlr, cenumber); err != nil {
					if errr := tx.Rollback(); errr != nil {
						return errr
					}
					return err
				}
				// getting the last inserted id
				if lastid, err = res.LastInsertId(); err != nil {
					if errr := tx.Rollback(); errr != nil {
						return errr
					}
					return err
				}
				// populating the map
				mRNNcenumber[cenumber] = strconv.FormatInt(lastid, 10)
			}
		}
	}

	//
	// supplier
	//
	logger.Log.Info("- importing storage suppliers")
	if csvFile, err = os.Open(path.Join(dir, "supplier.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		// csvReader = csv.NewReader(bufio.NewReader(csvFile))
		// i = 0
		// for {
		// 	line, error := csvReader.Read()

		// 	// skip header
		// 	if i == 0 {
		// 		i++
		// 		continue
		// 	}

		// 	if error == io.EOF {
		// 		break
		// 	} else if error != nil {
		// 		if errr := tx.Rollback(); errr != nil {
		// 			return errr
		// 		}
		// 		return err
		// 	}
		// 	id := line[0]
		// 	label := line[1]

		logger.Log.Debug("label:" + label)
		sqlr = `INSERT INTO supplier(supplier_id, supplier_label) VALUES (?, ?)`
		if res, err = tx.Exec(sqlr, id, label); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// populating the maps
		mONsupplier[id] = strconv.FormatInt(lastid, 10)
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	//
	// products
	//
	logger.Log.Info("- importing products")
	logger.Log.Info("  retrieving zero empirical id")
	if err = db.Get(&zeroempiricalformulaid, `SELECT empiricalformula_id FROM empiricalformula WHERE empiricalformula_label = "XXXX"`); err != nil {
		logger.Log.Error("error retrieving zero empirical id")
		return err
	}
	logger.Log.Info("  retrieving zero casnumber id")
	if err = db.Get(&zerocasnumberid, `SELECT casnumber_id FROM casnumber WHERE casnumber_label = "0000-00-0"`); err != nil {
		logger.Log.Error("error retrieving zero casnumber id")
		return err
	}
	logger.Log.Info("  retrieving default admin id")
	if err = db.Get(&zeropersonid, `SELECT person_id FROM person WHERE person_email = "admin@chimitheque.fr"`); err != nil {
		logger.Log.Error("error retrieving default admin id")
		return err
	}
	logger.Log.Info("  gathering hazardstatement ids")
	if csvFile, err = os.Open(path.Join(dir, "hazard_statement.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		reference := k["reference"]
		// csvReader = csv.NewReader(bufio.NewReader(csvFile))
		// i = 0
		// for {
		// 	line, error := csvReader.Read()

		// 	// skip header
		// 	if i == 0 {
		// 		i++
		// 		continue
		// 	}

		// 	if error == io.EOF {
		// 		break
		// 	} else if error != nil {
		// 		return err
		// 	}
		// 	id := line[0]
		// 	reference := line[2]
		if reference == "----" {
			zerohsid = id
			continue
		}
		// finding new id
		var nid int
		if err = db.Get(&nid, `SELECT hazardstatement_id FROM hazardstatement WHERE hazardstatement_reference = ?`, reference); err != nil {
			logger.Log.Info("no hazardstatement id for " + reference + " inserting a new one")
			var (
				res   sql.Result
				nid64 int64
			)
			if res, err = tx.Exec(`INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES (?, ?)`, id, reference); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if nid64, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			nid = int(nid64)
		}
		mONhazardstatement[id] = strconv.Itoa(nid)
	}
	logger.Log.Info("  gathering precautionarystatement ids")
	if csvFile, err = os.Open(path.Join(dir, "precautionary_statement.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		reference := k["reference"]
		// csvReader = csv.NewReader(bufio.NewReader(csvFile))
		// i = 0
		// for {
		// 	line, error := csvReader.Read()

		// 	// skip header
		// 	if i == 0 {
		// 		i++
		// 		continue
		// 	}

		// 	if error == io.EOF {
		// 		break
		// 	} else if error != nil {
		// 		return err
		// 	}
		// 	id := line[0]
		// 	reference := line[2]
		if reference == "----" {
			zeropsid = id
			continue
		}
		// finding new id
		var nid int
		if err = db.Get(&nid, `SELECT precautionarystatement_id FROM precautionarystatement WHERE precautionarystatement_reference = ?`, reference); err != nil {
			logger.Log.Info("no precautionarystatement id for " + reference + " inserting a new one")
			var (
				res   sql.Result
				nid64 int64
			)
			if res, err = tx.Exec(`INSERT INTO precautionarystatement (precautionarystatement_label, precautionarystatement_reference) VALUES (?, ?)`, id, reference); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if nid64, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			nid = int(nid64)
		}
		mONprecautionarystatement[id] = strconv.Itoa(nid)
	}
	logger.Log.Info("  gathering symbol ids")
	if csvFile, err = os.Open(path.Join(dir, "symbol.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		// csvReader = csv.NewReader(bufio.NewReader(csvFile))
		// i = 0
		// for {
		// 	line, error := csvReader.Read()

		// 	// skip header
		// 	if i == 0 {
		// 		i++
		// 		continue
		// 	}

		// 	if error == io.EOF {
		// 		break
		// 	} else if error != nil {
		// 		return err
		// 	}
		// 	id := line[0]
		// 	label := line[1]
		// finding new id
		var nid int
		if err = db.Get(&nid, `SELECT symbol_id FROM symbol WHERE symbol_label = ?`, label); err != nil {
			logger.Log.Error("error gathering symbol id for " + label)
			return err
		}
		mONsymbol[id] = strconv.Itoa(nid)
	}
	logger.Log.Info("  gathering signalword ids")
	if csvFile, err = os.Open(path.Join(dir, "signal_word.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		// csvReader = csv.NewReader(bufio.NewReader(csvFile))
		// i = 0
		// for {
		// 	line, error := csvReader.Read()

		// 	// skip header
		// 	if i == 0 {
		// 		i++
		// 		continue
		// 	}

		// 	if error == io.EOF {
		// 		break
		// 	} else if error != nil {
		// 		return err
		// 	}
		// 	id := line[0]
		// 	label := line[1]
		// finding new id
		var nid int
		if err = db.Get(&nid, `SELECT signalword_id FROM signalword WHERE signalword_label = ?`, label); err != nil {
			logger.Log.Error("error gathering signalword id for " + label)
			return err
		}
		mONsignalword[id] = strconv.Itoa(nid)
	}

	if csvFile, err = os.Open(path.Join(dir, "product.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		cenumber := k["ce_number"]
		person := k["person"]
		name := k["name"]
		synonym := k["synonym"]
		restricted := k["restricted_access"]
		specificity := k["specificity"]
		tdformula := k["tdformula"]
		empiricalformula := k["empirical_formula"]
		linearformula := k["linear_formula"]
		msds := k["msds"]
		physicalstate := k["physical_state"]
		coc := k["class_of_compounds"]
		symbol := k["symbol"]
		signalword := k["signal_word"]
		hazardstatement := k["hazard_statement"]
		precautionarystatement := k["precautionary_statement"]
		disposalcomment := k["disposal_comment"]
		remark := k["remark"]
		archive := k["archive"]
		casnumber := k["cas_number"]
		isradio := k["is_radio"]

		newperson := mONperson[person]
		if newperson == "" {
			newperson = strconv.Itoa(zeropersonid)
		}
		newname := mONname[name]
		newrestricted := false
		if restricted == "T" {
			newrestricted = true
		}
		newspecificity := specificity
		newtdformula := tdformula
		newempiricalformula := mONempiricalformula[empiricalformula]
		if newempiricalformula == "" {
			newempiricalformula = strconv.Itoa(zeroempiricalformulaid)
		}
		newlinearformula := sql.NullInt64{}
		if mONlinearformula[linearformula] != "" {
			i, e := strconv.ParseInt(mONlinearformula[linearformula], 10, 64)
			if e != nil {
				logger.Log.Error("error converting linearformula id for " + mONlinearformula[linearformula])
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			newlinearformula = sql.NullInt64{Valid: true, Int64: i}
		}
		newmsds := msds
		newphysicalstate := sql.NullInt64{}
		if mONphysicalstate[physicalstate] != "" {
			i, e := strconv.ParseInt(mONphysicalstate[physicalstate], 10, 64)
			if e != nil {
				logger.Log.Error("error converting physicalstate id for " + mONphysicalstate[physicalstate])
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			newphysicalstate = sql.NullInt64{Valid: true, Int64: i}
		}
		newsignalword := sql.NullInt64{}
		if mONsignalword[signalword] != "" {
			i, e := strconv.ParseInt(mONsignalword[signalword], 10, 64)
			if e != nil {
				logger.Log.Error("error converting signalword id for " + mONsignalword[signalword])
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			newsignalword = sql.NullInt64{Valid: true, Int64: i}
		}
		newdisposalcomment := disposalcomment
		newremark := remark
		newarchive := false
		if archive == "T" {
			newarchive = true
		}
		newcasnumber := mRNNcasnumber[casnumber]
		if newcasnumber == "" {
			newcasnumber = strconv.Itoa(zerocasnumberid)
		}
		newcenumber := mRNNcenumber[cenumber]
		newisradio := false
		if isradio == "T" {
			newisradio = true
		}

		// do not import archived cards
		if !newarchive {
			reqValues := "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?"
			reqArgs := []interface{}{
				newspecificity,
				newmsds,
				newrestricted,
				newisradio,
				newtdformula,
				newdisposalcomment,
				newremark,
				newempiricalformula,
				newlinearformula,
				newphysicalstate,
				newsignalword,
				newperson,
				newcasnumber,
				newname,
			}
			sqlr = `INSERT INTO product (product_specificity, 
                product_msds, 
                product_restricted, 
                product_radioactive, 
                product_threedformula, 
                product_disposalcomment, 
                product_remark,
                empiricalformula,
                linearformula,
                physicalstate,
                signalword,
                person,
				casnumber,
				name`
			if newcenumber != "" {
				sqlr += ",cenumber"
				reqValues += ",?"
				reqArgs = append(reqArgs, newcenumber)
			}
			sqlr += `) VALUES (` + reqValues + `)`

			logger.Log.Debug(fmt.Sprintf(`newperson: %s,
			newname: %s,
			newrestricted: %t,
			newspecificity: %s,
			newtdformula: %s,
			newempiricalformula: %s,
			newlinearformula: %v,
			newmsds: %s,
			newphysicalstate: %v,
			newsignalword: %v,
			newdisposalcomment: %s,
			newremark: %s,
			newarchive: %t,
			casnumber: %s,
			newcasnumber: %s,
			newcenumber: %s,
			newisradio: %t
			`, newperson,
				newname,
				newrestricted,
				newspecificity,
				newtdformula,
				newempiricalformula,
				newlinearformula,
				newmsds,
				newphysicalstate,
				newsignalword,
				newdisposalcomment,
				newremark,
				newarchive,
				casnumber,
				newcasnumber,
				newcenumber,
				newisradio))

			if res, err = tx.Exec(sqlr, reqArgs...); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				logger.Log.Error("error importing product")
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// populating the map
			mONproduct[id] = strconv.FormatInt(lastid, 10)

			// coc
			cocs := rnumber.FindAllString(coc, -1)
			for _, c := range cocs {
				sqlr = `INSERT INTO productclassofcompound (productclassofcompound_product_id, productclassofcompound_classofcompound_id) VALUES (?,?)`
				if _, err = tx.Exec(sqlr, lastid, mONclassofcompound[c]); err != nil {
					// not leaving on errors
					logger.Log.Debug("non fatal error importing product class of compounds with id " + c + ": " + err.Error())
				}
			}
			// synonym
			syns := rnumber.FindAllString(synonym, -1)
			for _, s := range syns {
				if s == "0" {
					continue
				}
				// leaving hardcoded zeros
				sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
				if _, err = tx.Exec(sqlr, lastid, mONname[s]); err != nil {
					// not leaving on errors
					logger.Log.Debug("non fatal error importing product synonym with id " + s + ": " + err.Error())
				}
			}
			// symbol
			symbols := rnumber.FindAllString(symbol, -1)
			for _, s := range symbols {
				sqlr = `INSERT INTO productsymbols (productsymbols_product_id, productsymbols_symbol_id) VALUES (?,?)`
				if _, err = tx.Exec(sqlr, lastid, mONsymbol[s]); err != nil {
					// not leaving on errors
					logger.Log.Error("error importing product symbol with id " + s + ": " + err.Error())
					if errr := tx.Rollback(); errr != nil {
						return errr
					}
					return err
				}
			}
			// hs
			hss := rnumber.FindAllString(hazardstatement, -1)
			for _, s := range hss {
				if s == zerohsid {
					continue
				}
				sqlr = `INSERT INTO producthazardstatements (producthazardstatements_product_id, producthazardstatements_hazardstatement_id) VALUES (?,?)`
				if _, err = tx.Exec(sqlr, lastid, mONhazardstatement[s]); err != nil {
					// not leaving on errors
					logger.Log.Error("error importing product hazardstatement with id " + s + ": " + err.Error())
					if errr := tx.Rollback(); errr != nil {
						return errr
					}
					return err
				}
			}
			// ps
			pss := rnumber.FindAllString(precautionarystatement, -1)
			for _, s := range pss {
				if s == zeropsid {
					continue
				}
				sqlr = `INSERT INTO productprecautionarystatements (productprecautionarystatements_product_id, productprecautionarystatements_precautionarystatement_id) VALUES (?,?)`
				if _, err = tx.Exec(sqlr, lastid, mONprecautionarystatement[s]); err != nil {
					// not leaving on errors
					logger.Log.Error("error importing product precautionarystatement with id " + s + ": " + err.Error())
					if errr := tx.Rollback(); errr != nil {
						return errr
					}
					return err
				}
			}
		}

	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	//
	// storages
	//
	logger.Log.Info("- importing storages")
	logger.Log.Info("  gathering unit ids")
	if csvFile, err = os.Open(path.Join(dir, "unit.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		id := k["id"]
		label := k["label"]
		// csvReader = csv.NewReader(bufio.NewReader(csvFile))
		// i = 0
		// for {
		// 	line, error := csvReader.Read()

		// 	// skip header
		// 	if i == 0 {
		// 		i++
		// 		continue
		// 	}

		// 	if error == io.EOF {
		// 		break
		// 	} else if error != nil {
		// 		return err
		// 	}
		// 	id := line[0]
		// 	label := line[1]
		// uppercase liter
		label = strings.Replace(label, "l", "L", -1)
		// finding new id
		var nid int
		if err = db.Get(&nid, `SELECT unit_id FROM unit WHERE unit_label = ?`, label); err != nil {
			logger.Log.Error("error gathering unit id for " + label)
			return err
		}
		mONunit[id] = strconv.Itoa(nid)
	}

	if csvFile, err = os.Open(path.Join(dir, "storage.csv")); err != nil {
		return (err)
	}
	csvMap = CSVToMap(bufio.NewReader(csvFile))
	for _, k := range csvMap {
		oldid := k["id"]
		product := k["product"]
		person := k["person"]
		storeLocation := k["store_location"]
		unit := k["unit"]
		entrydate := k["entry_datetime"]
		exitdate := k["exit_datetime"]
		comment := k["comment"]
		barecode := k["barecode"]
		reference := k["reference"]
		batchNumber := k["batch_number"]
		supplier := k["supplier"]
		archive := k["archive"]
		creationdate := k["creation_datetime"]
		volumeWeight := k["volume_weight"]
		openingdate := k["opening_datetime"]
		toDestroy := k["to_destroy"]
		expirationdate := k["expiration_datetime"]

		logger.Log.Debug(logger.Log.WithFields(logrus.Fields{
			"oldid":         oldid,
			"product":       product,
			"person":        person,
			"storeLocation": storeLocation,
			"unit":          unit,
			"entrydate":     entrydate,
			"exitdate":      exitdate,
			"supplier":      supplier,
		}))

		newproduct := mONproduct[product]
		newperson := mONperson[person]
		if newperson == "" {
			newperson = strconv.Itoa(zeropersonid)
		}
		newstoreLocation := mONstorelocation[storeLocation]
		newunit := mONunit[unit]
		var newentrydate *time.Time
		if entrydate != "" {
			newentrydate = &time.Time{}
			*newentrydate, _ = time.Parse("2006-01-02 15:04:05", entrydate)
		}
		var newexitdate *time.Time
		if exitdate != "" {
			newexitdate = &time.Time{}
			*newexitdate, _ = time.Parse("2006-01-02 15:04:05", exitdate)
		}
		newcomment := comment
		newbarecode := barecode
		newreference := reference
		newbatchNumber := batchNumber
		newsupplier := mONsupplier[supplier]
		newarchive := false
		if archive == "T" {
			newarchive = true
		}
		newstorageCreationdate := time.Now()
		if creationdate != "" {
			newstorageCreationdate, _ = time.Parse("2006-01-02 15:04:05", creationdate)
		}
		newvolumeWeight := volumeWeight
		if newvolumeWeight == "" {
			newvolumeWeight = "1"
		}
		var newopeningdate *time.Time
		if openingdate != "" {
			newopeningdate = &time.Time{}
			*newopeningdate, _ = time.Parse("2006-01-02 15:04:05", openingdate)
		}
		newtoDestroy := false
		if toDestroy == "T" {
			newtoDestroy = true
		}
		var newexpirationdate *time.Time
		if expirationdate != "" {
			newexpirationdate = &time.Time{}
			*newexpirationdate, _ = time.Parse("2006-01-02 15:04:05", expirationdate)
		}

		// do not import archived cards
		if !newarchive {
			reqValues := "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?"
			reqArgs := []interface{}{
				newstorageCreationdate,
				newstorageCreationdate,
				newcomment,
				newreference,
				newbatchNumber,
				newvolumeWeight,
				newbarecode,
				newtoDestroy,
				newperson,
				newproduct,
				newstoreLocation,
			}
			sqlr = `INSERT INTO storage (storage_creationdate, 
                storage_modificationdate, 
                storage_comment, 
                storage_reference, 
                storage_batchnumber, 
                storage_quantity, 
                storage_barecode,
                storage_todestroy,
                person,
                product,
				storelocation`
			if newunit != "" {
				sqlr += ",unit"
				reqValues += ",?"
				reqArgs = append(reqArgs, newunit)
			}
			if newsupplier != "" {
				sqlr += ",supplier"
				reqValues += ",?"
				reqArgs = append(reqArgs, newsupplier)
			}
			if newentrydate != nil {
				sqlr += ",storage_entrydate"
				reqValues += ",?"
				reqArgs = append(reqArgs, newentrydate)
			}
			if newexitdate != nil {
				sqlr += ",storage_exitdate"
				reqValues += ",?"
				reqArgs = append(reqArgs, newexitdate)
			}
			if newopeningdate != nil {
				sqlr += ",storage_openingdate"
				reqValues += ",?"
				reqArgs = append(reqArgs, newopeningdate)
			}
			if newexpirationdate != nil {
				sqlr += ",storage_expirationdate"
				reqValues += ",?"
				reqArgs = append(reqArgs, newexpirationdate)
			}

			sqlr += `) VALUES (` + reqValues + `)`

			logger.Log.Debug(logger.Log.WithFields(logrus.Fields{
				"newstorageCreationdate": newstorageCreationdate,
				"newcomment":             newcomment,
				"newreference":           newreference,
				"newbatchNumber":         newbatchNumber,
				"newvolumeWeight":        newvolumeWeight,
				"newbarecode":            newbarecode,
				"newtoDestroy":           newtoDestroy,
				"newperson":              newperson,
				"newproduct":             newproduct,
				"newstoreLocation":       newstoreLocation,
				"newunit":                newunit,
				"newsupplier":            newsupplier,
				"newentrydate":           newentrydate,
				"newexitdate":            newexitdate,
				"newopeningdate":         newopeningdate,
				"newexpirationdate":      newexpirationdate,
			}))

			if _, err = tx.Exec(sqlr, reqArgs...); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	logger.Log.Info("- updating store locations full path")
	var sls []StoreLocation
	if err = db.Select(&sls, ` SELECT s.storelocation_id AS "storelocation_id", 
        s.storelocation_name AS "storelocation_name", 
        s.storelocation_canstore, 
        s.storelocation_color,
        storelocation.storelocation_id AS "storelocation.storelocation_id",
        storelocation.storelocation_name AS "storelocation.storelocation_name"
        FROM storelocation AS s
        LEFT JOIN storelocation on s.storelocation = storelocation.storelocation_id`); err != nil {
		return err
	}

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	for _, sl := range sls {
		logger.Log.Debug("  " + sl.StoreLocationName.String)
		sl.StoreLocationFullPath = db.buildFullPath(sl, tx)
		sqlr = `UPDATE storelocation SET storelocation_fullpath = ? WHERE storelocation_id = ?`
		if _, err = tx.Exec(sqlr, sl.StoreLocationFullPath, sl.StoreLocationID.Int64); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	return nil
}
