package models

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tbellembois/gochimitheque/logger"
)

type MyNullTime sql.NullTime
type MyTime struct{ time.Time }

func (t *MyNullTime) Scan(src interface{}) error {
	var (
		tm    time.Time
		valid bool
	)

	switch src.(type) {
	case int64:
		value := src.(int64)
		tm = time.Unix(value, 0)
		valid = false
	case nil:

	default:
		return errors.New("Invalid type for ExpireTimestamp")
	}

	*t = MyNullTime{Valid: valid, Time: tm}
	return nil
}

func (t *MyTime) Scan(src interface{}) error {
	var (
		tm time.Time
	)

	switch src.(type) {
	case int64:
		value := src.(int64)
		tm = time.Unix(value, 0)
	case nil:

	default:
		return errors.New("Invalid type for ExpireTimestamp")
	}

	*t = MyTime{tm}
	return nil
}

// Storage is a product storage in a store location.
type Storage struct {
	StorageID                sql.NullInt64   `db:"storage_id" json:"storage_id" schema:"storage_id" `
	StorageCreationDate      MyTime          `db:"storage_creation_date" json:"storage_creation_date" schema:"storage_creation_date"`
	StorageModificationDate  MyTime          `db:"storage_modification_date" json:"storage_modification_date" schema:"storage_modification_date"`
	StorageEntryDate         MyNullTime      `db:"storage_entry_date" json:"storage_entry_date" schema:"storage_entry_date" `
	StorageExitDate          MyNullTime      `db:"storage_exit_date" json:"storage_exit_date" schema:"storage_exit_date" `
	StorageOpeningDate       MyNullTime      `db:"storage_opening_date" json:"storage_opening_date" schema:"storage_opening_date" `
	StorageExpirationDate    MyNullTime      `db:"storage_expiration_date" json:"storage_expiration_date" schema:"storage_expiration_date" `
	StorageComment           sql.NullString  `db:"storage_comment" json:"storage_comment" schema:"storage_comment" `
	StorageReference         sql.NullString  `db:"storage_reference" json:"storage_reference" schema:"storage_reference" `
	StorageBatchNumber       sql.NullString  `db:"storage_batch_number" json:"storage_batch_number" schema:"storage_batch_number" `
	StorageQuantity          sql.NullFloat64 `db:"storage_quantity" json:"storage_quantity" schema:"storage_quantity" `
	StorageNbItem            int             `db:"-" json:"storage_nbitem" schema:"storage_nbitem"`
	StorageIdenticalBarecode sql.NullBool    `db:"-" json:"storage_identicalbarecode" schema:"storage_identicalbarecode" `
	StorageBarecode          sql.NullString  `db:"storage_barecode" json:"storage_barecode" schema:"storage_barecode" `
	StorageQRCode            []byte          `db:"storage_qrcode" json:"storage_qrcode" schema:"storage_qrcode"`
	StorageToDestroy         sql.NullBool    `db:"storage_to_destroy" json:"storage_to_destroy" schema:"storage_to_destroy" `
	StorageArchive           sql.NullBool    `db:"storage_archive" json:"storage_archive" schema:"storage_archive" `
	StorageConcentration     sql.NullInt64   `db:"storage_concentration" json:"storage_concentration" schema:"storage_concentration" `
	StorageNumberOfUnit      sql.NullInt64   `db:"storage_number_of_unit" json:"storage_number_of_unit" schema:"storage_number_of_unit" `
	StorageNumberOfBag       sql.NullInt64   `db:"storage_number_of_bag" json:"storage_number_of_bag" schema:"storage_number_of_bag" `
	StorageNumberOfCarton    sql.NullInt64   `db:"storage_number_of_carton" json:"storage_number_of_carton" schema:"storage_number_of_carton" `
	Person                   `db:"person" json:"person" schema:"person"`
	Product                  `db:"product" json:"product" schema:"product"`
	StoreLocation            `db:"store_location" json:"store_location" schema:"store_location"`
	UnitQuantity             Unit `db:"unit_quantity" json:"unit_quantity" schema:"unit_quantity"`
	UnitConcentration        Unit `db:"unit_concentration" json:"unit_concentration" schema:"unit_concentration"`
	Supplier                 `db:"supplier" json:"supplier" schema:"supplier"`
	Storage                  *Storage   `db:"storage" json:"storage" schema:"storage"`       // history reference storage
	Borrowing                *Borrowing `db:"borrowing" json:"borrowing" schema:"borrowing"` // not un db but sqlx requires the "db" entry

	// storage history count
	StorageHC int `db:"storage_hc" json:"storage_hc" schema:"storage_hc"` // not in db but sqlx requires the "db" entry
}

func (s Storage) StorageToStringSlice() []string {
	ret := make([]string, 0)

	ret = append(ret, strconv.FormatInt(s.StorageID.Int64, 10))
	ret = append(ret, s.Product.Name.NameLabel)
	// ret = append(ret, s.Product.CasNumber.CasNumberLabel.String)
	ret = append(ret, *s.Product.CasNumber.CasNumberLabel)
	ret = append(ret, *s.Product.ProductSpecificity)

	ret = append(ret, s.StoreLocation.StoreLocationFullPath)

	ret = append(ret, strconv.FormatFloat(s.StorageQuantity.Float64, 'E', -1, 64))
	ret = append(ret, *s.UnitQuantity.UnitLabel)
	ret = append(ret, s.StorageBarecode.String)
	ret = append(ret, *s.Supplier.SupplierLabel)

	ret = append(ret, s.StorageCreationDate.Format("2006-01-02"))
	ret = append(ret, s.StorageModificationDate.Format("2006-01-02"))
	ret = append(ret, s.StorageEntryDate.Time.Format("2006-01-02"))
	ret = append(ret, s.StorageExitDate.Time.Format("2006-01-02"))
	ret = append(ret, s.StorageOpeningDate.Time.Format("2006-01-02"))
	ret = append(ret, s.StorageExpirationDate.Time.Format("2006-01-02"))

	ret = append(ret, s.StorageComment.String)
	ret = append(ret, s.StorageReference.String)
	ret = append(ret, s.StorageBatchNumber.String)

	ret = append(ret, strconv.FormatBool(s.StorageToDestroy.Bool))
	ret = append(ret, strconv.FormatBool(s.StorageArchive.Bool))

	return ret
}

// StoragesToCSV returns a file name of the products prs
// exported into CSV.
func StoragesToCSV(sts []Storage) (string, error) {
	var (
		err     error
		tmpFile *os.File
	)

	header := []string{
		"storage_id",
		"product_name",
		"product_cas_number",
		"product_specificity",
		"store_location",
		"quantity",
		"unit",
		"barecode",
		"supplier",
		"creation_date",
		"modification_date",
		"entry_date",
		"exit_date",
		"opening_date",
		"expiration_date",
		"comment",
		"reference",
		"batch_number",
		"to_destroy?",
		"archive?",
	}

	// create a temp file
	if tmpFile, err = os.CreateTemp(os.TempDir(), "chimitheque-"); err != nil {
		logger.Log.Error("cannot create temporary file", err)
		return "", err
	}
	// creates a csv writer that uses the io buffer
	csvwr := csv.NewWriter(tmpFile)
	// write the header
	if err = csvwr.Write(header); err != nil {
		logger.Log.Error("cannot write header", err)
		return "", err
	}

	for _, s := range sts {
		if err = csvwr.Write(s.StorageToStringSlice()); err != nil {
			logger.Log.Error("cannot write entry", err)
			return "", err
		}
	}

	csvwr.Flush()

	return strings.Split(tmpFile.Name(), "chimitheque-")[1], nil
}
