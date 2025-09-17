package datastores

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	"github.com/doug-martin/goqu/v9"

	// register sqlite3 driver.
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

// DeleteStorage deletes the storages with the given id.
func (db *SQLiteDataStore) DeleteStorage(id int) error {
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeleteStorage")

	var (
		sqlr string
		err  error
	)

	// Delete history first.
	sqlr = `DELETE FROM storage
	WHERE storage = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `DELETE FROM storage
	WHERE storage_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	return nil
}

// ArchiveStorage archives the storages with the given id.
func (db *SQLiteDataStore) ArchiveStorage(id int) error {
	var (
		sqlr string
		err  error
	)

	sqlr = `UPDATE storage SET storage_archive = true
	WHERE storage_id = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `UPDATE storage SET storage_archive = true
	WHERE storage.storage = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	return nil
}

// RestoreStorage restores (unarchive) the storages with the given id.
func (db *SQLiteDataStore) RestoreStorage(id int) error {
	var (
		sqlr string
		err  error
	)

	sqlr = `UPDATE storage SET storage_archive = false
	WHERE storage_id = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `UPDATE storage SET storage_archive = false
	WHERE storage.storage = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	return nil
}

// CreateStorage creates a new storage.
func (db *SQLiteDataStore) CreateUpdateStorage(s models.Storage, itemNumber int, update bool) (lastInsertID int64, err error) {
	var (
		tx           *sql.Tx
		sqlr         string
		res          sql.Result
		args         []interface{}
		prefix       string
		major, minor string
	)

	logger.Log.WithFields(logrus.Fields{"s": fmt.Sprintf("%+v", s)}).Debug("CreateUpdateStorage")

	// Default major.
	major = strconv.Itoa(s.ProductID)

	dialect := goqu.Dialect("sqlite3")
	tableStorage := goqu.T("storage")

	if tx, err = db.Begin(); err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			logger.Log.Error(err)
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Log.Error(rbErr)
				err = rbErr

				return
			}

			return
		}

		err = tx.Commit()
	}()

	if update {
		// create an history of the storage
		sqlr = `INSERT into storage (storage_creation_date,
		storage_modification_date,
		storage_entry_date,
		storage_exit_date,
		storage_opening_date,
		storage_expiration_date,
		storage_comment,
		storage_reference,
		storage_batch_number,
		storage_quantity,
		storage_barecode,
		storage_to_destroy,
		storage_archive,
		storage_concentration,
		storage_number_of_bag,
		storage_number_of_carton,
		person,
		product,
		store_location,
		unit_quantity,
		unit_concentration,
		supplier,
		storage) select storage_creation_date,
				storage_modification_date,
				storage_entry_date,
				storage_exit_date,
				storage_opening_date,
				storage_expiration_date,
				storage_comment,
				storage_reference,
				storage_batch_number,
				storage_quantity,
				storage_barecode,
				storage_to_destroy,
				storage_archive,
				storage_concentration,
				storage_number_of_bag,
				storage_number_of_carton,
				person,
				product,
				store_location,
				unit_quantity,
				unit_concentration,
				supplier,
				? FROM storage WHERE storage_id = ?`
		if _, err = tx.Exec(sqlr, s.StorageID, s.StorageID); err != nil {
			logger.Log.Error("error creating storage history")
			return
		}
	}

	// Generating barecode if empty.
	if !update {
		if s.StorageBarecode == nil || *s.StorageBarecode == "" {
			//
			// Getting the barecode prefix from the store_location name.
			//
			// regex to detect store locations names starting with [_a-zA-Z] to build barecode prefixes
			prefixRegex := regexp.MustCompile(`^\[(?P<groupone>[_a-zA-Z]{1,5})\].*$`)
			groupNames := prefixRegex.SubexpNames()
			matches := prefixRegex.FindAllStringSubmatch(s.StoreLocationName.String, -1)
			// Building a map of matches.
			matchesMap := map[string]string{}

			logger.Log.WithFields(logrus.Fields{"s.StoreLocationName.String": s.StoreLocationName.String, "matches": matches}).Debug("CreateStorage")

			if len(matches) != 0 {
				for i, j := range matches[0] {
					matchesMap[groupNames[i]] = j
				}
			}

			if len(matchesMap) > 0 {
				prefix = matchesMap["groupone"]
			} else {
				prefix = "_"
			}

			//
			// Getting the storage barecodes matching the regex
			// for the same product in the same entity.
			//
			sqlr := `SELECT storage_barecode FROM storage
		JOIN store_location on storage.store_location = store_location.store_location_id
		WHERE product = ? AND store_location.entity = ? AND regexp('^[_a-zA-Z]{0,5}[0-9]+\.[0-9]+$', '' || storage_barecode || '') = true
		ORDER BY storage_barecode desc`

			var rows *sql.Rows

			if rows, err = tx.Query(sqlr, s.ProductID, s.EntityID); err != nil && err != sql.ErrNoRows {
				logger.Log.Error("error getting storage barecode")
				return
			}

			var (
				count    = 0
				newMinor = 0
			)

			for rows.Next() {
				var barecode string
				if err = rows.Scan(&barecode); err != nil && err != sql.ErrNoRows {
					return
				}

				majorRegex := regexp.MustCompile(`^[_a-zA-Z]{0,5}(?P<groupone>[0-9]+)\.(?P<grouptwo>[0-9]+)$`)
				groupNames = majorRegex.SubexpNames()
				matches = majorRegex.FindAllStringSubmatch(barecode, -1)
				// Building a map of matches.
				matchesMap = map[string]string{}

				if len(matches) != 0 {
					for i, j := range matches[0] {
						matchesMap[groupNames[i]] = j
					}
				}

				if count == 0 {
					// All of the major number are the same.
					// Extracting it ones.
					major = matchesMap["groupone"]
				}

				minor = matchesMap["grouptwo"]

				var iminor int

				if iminor, err = strconv.Atoi(minor); err != nil {
					return 0, err
				}

				if iminor > newMinor {
					newMinor = iminor
				}

				count++
			}

			if !s.StorageIdenticalBarecode || (s.StorageIdenticalBarecode && itemNumber == 1) {
				newMinor++
			}

			minor = strconv.Itoa(newMinor)

			*s.StorageBarecode = prefix + major + "." + minor

			logger.Log.WithFields(logrus.Fields{"s.StorageBarecode.String": s.StorageBarecode}).Debug("CreateStorage")
		}
	}

	// if SupplierID = -1 then it is a new supplier
	if s.Supplier.SupplierID != nil && err == nil && *s.Supplier.SupplierID == -1 {
		sqlr = `INSERT INTO supplier (supplier_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, s.Supplier.SupplierLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the storage SupplierId (SupplierLabel already set)
		*s.Supplier.SupplierID = lastInsertID
	}
	if err != nil {
		logger.Log.Error("supplier error - " + err.Error())
		return
	}

	// finally updating the storage
	insertCols := goqu.Record{}
	if s.StorageComment != nil {
		insertCols["storage_comment"] = s.StorageComment
	} else {
		insertCols["storage_comment"] = nil
	}

	if s.StorageQuantity != nil {
		insertCols["storage_quantity"] = s.StorageQuantity
	} else {
		insertCols["storage_quantity"] = nil
	}

	if s.StorageBarecode != nil {
		insertCols["storage_barecode"] = s.StorageBarecode
	} else {
		insertCols["storage_barecode"] = nil
	}

	if s.UnitQuantity.UnitID != nil {
		insertCols["unit_quantity"] = *s.UnitQuantity.UnitID
	} else {
		insertCols["unit_quantity"] = nil
	}

	if s.Supplier.SupplierID != nil {
		insertCols["supplier"] = *s.SupplierID
	} else {
		insertCols["supplier"] = nil
	}

	if s.StorageEntryDate != nil {
		insertCols["storage_entry_date"] = s.StorageEntryDate.Unix()
	} else {
		insertCols["storage_entry_date"] = nil
	}

	if s.StorageExitDate != nil {
		insertCols["storage_exit_date"] = s.StorageExitDate.Unix()
	} else {
		insertCols["storage_exit_date"] = nil
	}

	if s.StorageOpeningDate != nil {
		insertCols["storage_opening_date"] = s.StorageOpeningDate.Unix()
	} else {
		insertCols["storage_opening_date"] = nil
	}

	if s.StorageExpirationDate != nil {
		insertCols["storage_expiration_date"] = s.StorageExpirationDate.Unix()
	} else {
		insertCols["storage_expiration_date"] = nil
	}

	if s.StorageReference != nil {
		insertCols["storage_reference"] = s.StorageReference
	} else {
		insertCols["storage_reference"] = nil
	}

	if s.StorageBatchNumber != nil {
		insertCols["storage_batch_number"] = s.StorageBatchNumber
	} else {
		insertCols["storage_batch_number"] = nil
	}

	if s.StorageToDestroy {
		insertCols["storage_to_destroy"] = s.StorageToDestroy
	} else {
		insertCols["storage_to_destroy"] = false
	}

	if s.StorageConcentration != nil {
		insertCols["storage_concentration"] = int(*s.StorageConcentration)
	} else {
		insertCols["storage_concentration"] = nil
	}

	if s.StorageNumberOfBag != nil {
		insertCols["storage_number_of_bag"] = int(*s.StorageNumberOfBag)
	} else {
		insertCols["storage_number_of_bag"] = nil
	}

	if s.StorageNumberOfCarton != nil {
		insertCols["storage_number_of_carton"] = int(*s.StorageNumberOfCarton)
	} else {
		insertCols["storage_number_of_carton"] = nil
	}

	if s.UnitConcentration.UnitID != nil {
		insertCols["unit_concentration"] = int(*s.UnitConcentration.UnitID)
	} else {
		insertCols["unit_concentration"] = nil
	}

	insertCols["person"] = s.PersonID
	insertCols["store_location"] = s.StoreLocationID.Int64
	insertCols["product"] = s.ProductID
	insertCols["storage_creation_date"] = s.StorageCreationDate.Unix()
	insertCols["storage_modification_date"] = s.StorageModificationDate.Unix()
	insertCols["storage_archive"] = false

	if update {
		iQuery := dialect.Update(tableStorage).Set(insertCols).Where(goqu.I("storage_id").Eq(s.StorageID))
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			logger.Log.Error("error preparing update storage")
			return
		}
	} else {
		iQuery := dialect.Insert(tableStorage).Rows(insertCols)
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			logger.Log.Error("error preparing create storage")
			return
		}
	}

	// logger.Log.Debug(sqlr)
	// logger.Log.Debug(args)

	if res, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Error("error creating/updating storage")
		return
	}

	// getting the last inserted id
	if !update {
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
	}

	//
	// qrcode
	//
	qr := strconv.FormatInt(lastInsertID, 10)
	if s.StorageQRCode, err = qrcode.Encode(qr, qrcode.Medium, 512); err != nil {
		return
	}

	sqlr = `UPDATE storage SET storage_qrcode=? WHERE storage_id=?`
	if _, err = tx.Exec(sqlr, s.StorageQRCode, lastInsertID); err != nil {
		return
	}

	var storage_id int64 = lastInsertID
	s.StorageID = &storage_id

	logger.Log.WithFields(logrus.Fields{"s": s}).Debug("CreateUpdateStorage")

	return
}

// UpdateAllQRCodes updates the storages QRCodes.
func (db *SQLiteDataStore) UpdateAllQRCodes() error {
	var (
		err  error
		tx   *sqlx.Tx
		sts  []models.Storage
		png  []byte
		sqlr string
	)

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	// retrieving storages
	if err = db.Select(&sts, ` SELECT storage_id
        FROM storage`); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	for _, s := range sts {
		// generating qrcode
		newqrcode := strconv.FormatInt(*s.StorageID, 10)
		logger.Log.Debug("  " + strconv.FormatInt(*s.StorageID, 10) + " " + newqrcode)

		if png, err = qrcode.Encode(newqrcode, qrcode.Medium, 512); err != nil {
			return err
		}

		sqlr = `UPDATE storage
				SET storage_qrcode = ?
				WHERE storage_id = ?`

		if _, err = tx.Exec(sqlr, png, s.StorageID); err != nil {
			logger.Log.Error("error updating storage qrcode")
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
