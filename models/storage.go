package models

import (
	"strconv"
	"time"
)

// Storage is a product storage in a store location.
type Storage struct {
	StorageID               *int64     `db:"storage_id" json:"storage_id,omitempty" schema:"storage_id" `
	StorageCreationDate     time.Time  `db:"storage_creation_date" json:"storage_creation_date,omitempty" schema:"storage_creation_date"`
	StorageModificationDate time.Time  `db:"storage_modification_date" json:"storage_modification_date,omitempty" schema:"storage_modification_date"`
	StorageEntryDate        *time.Time `db:"storage_entry_date" json:"storage_entry_date,omitempty" schema:"storage_entry_date" `
	StorageExitDate         *time.Time `db:"storage_exit_date" json:"storage_exit_date,omitempty" schema:"storage_exit_date" `
	StorageOpeningDate      *time.Time `db:"storage_opening_date" json:"storage_opening_date,omitempty" schema:"storage_opening_date" `
	StorageExpirationDate   *time.Time `db:"storage_expiration_date" json:"storage_expiration_date,omitempty" schema:"storage_expiration_date" `
	StorageComment          *string    `db:"storage_comment" json:"storage_comment,omitempty" schema:"storage_comment" `
	StorageReference        *string    `db:"storage_reference" json:"storage_reference,omitempty" schema:"storage_reference" `
	StorageBatchNumber      *string    `db:"storage_batch_number" json:"storage_batch_number,omitempty" schema:"storage_batch_number" `
	StorageQuantity         *float64   `db:"storage_quantity" json:"storage_quantity,omitempty" schema:"storage_quantity" `
	// StorageNbItem            int        `db:"-" json:"storage_nbitem,omitempty" schema:"storage_nbitem"`
	// StorageIdenticalBarecode bool    `db:"-" json:"storage_identical_barecode,omitempty" schema:"storage_identical_barecode" `
	StorageBarecode       *string `db:"storage_barecode" json:"storage_barecode,omitempty" schema:"storage_barecode" `
	StorageQRCode         []byte  `db:"storage_qrcode" json:"storage_qrcode,omitempty" schema:"storage_qrcode"`
	StorageToDestroy      bool    `db:"storage_to_destroy" json:"storage_to_destroy,omitempty" schema:"storage_to_destroy" `
	StorageArchive        bool    `db:"storage_archive" json:"storage_archive,omitempty" schema:"storage_archive" `
	StorageConcentration  *int64  `db:"storage_concentration" json:"storage_concentration,omitempty" schema:"storage_concentration" `
	StorageNumberOfBag    *int64  `db:"storage_number_of_bag" json:"storage_number_of_bag,omitempty" schema:"storage_number_of_bag" `
	StorageNumberOfCarton *int64  `db:"storage_number_of_carton" json:"storage_number_of_carton,omitempty" schema:"storage_number_of_carton" `
	Person                `db:"person" json:"person,omitempty" schema:"person"`
	Product               `db:"product" json:"product,omitempty" schema:"product"`
	//Entity                   `db:"entity" json:"entity,omitempty" schema:"entity"`
	StoreLocation     `db:"store_location" json:"store_location,omitempty" schema:"store_location"`
	UnitQuantity      *Unit `db:"unit_quantity" json:"unit_quantity,omitempty" schema:"unit_quantity"`
	UnitConcentration *Unit `db:"unit_concentration" json:"unit_concentration,omitempty" schema:"unit_concentration"`
	*Supplier         `db:"supplier" json:"supplier,omitempty" schema:"supplier"`
	Storage           *Storage   `db:"storage" json:"storage,omitempty" schema:"storage"`       // history reference storage
	Borrowing         *Borrowing `db:"borrowing" json:"borrowing,omitempty" schema:"borrowing"` // not un db but sqlx requires the "db" entry

	// storage history count
	StorageHC int `db:"storage_hc" json:"storage_hc,omitempty" schema:"storage_hc"` // not in db but sqlx requires the "db" entry
}

func (s Storage) StorageToStringSlice() []string {
	ret := make([]string, 0)

	ret = append(ret, strconv.FormatInt(*s.StorageID, 10))
	ret = append(ret, s.Product.Name.NameLabel)
	// ret = append(ret, s.Product.CasNumber.CasNumberLabel.String)
	if s.Product.CasNumber.CasNumberLabel != nil {
		ret = append(ret, *s.Product.CasNumber.CasNumberLabel)
	}
	if s.Product.ProductSpecificity != nil {
		ret = append(ret, *s.Product.ProductSpecificity)
	}

	ret = append(ret, s.StoreLocation.StoreLocationFullPath)

	ret = append(ret, strconv.FormatFloat(*s.StorageQuantity, 'E', -1, 64))

	if s.UnitQuantity.UnitLabel != nil {
		ret = append(ret, *s.UnitQuantity.UnitLabel)
	}

	ret = append(ret, *s.StorageBarecode)

	if s.Supplier.SupplierLabel != nil {
		ret = append(ret, *s.Supplier.SupplierLabel)
	}

	ret = append(ret, s.StorageCreationDate.Format("2006-01-02"))
	ret = append(ret, s.StorageModificationDate.Format("2006-01-02"))
	ret = append(ret, s.StorageEntryDate.Format("2006-01-02"))
	ret = append(ret, s.StorageExitDate.Format("2006-01-02"))
	ret = append(ret, s.StorageOpeningDate.Format("2006-01-02"))
	ret = append(ret, s.StorageExpirationDate.Format("2006-01-02"))

	ret = append(ret, *s.StorageComment)
	ret = append(ret, *s.StorageReference)
	ret = append(ret, *s.StorageBatchNumber)

	ret = append(ret, strconv.FormatBool(s.StorageToDestroy))
	ret = append(ret, strconv.FormatBool(s.StorageArchive))

	return ret
}

// StoragesToCSV returns a file name of the products prs
// exported into CSV.
// func StoragesToCSV(sts []Storage) (string, error) {
// 	var (
// 		err     error
// 		tmpFile *os.File
// 	)

// 	header := []string{
// 		"storage_id",
// 		"product_name",
// 		"product_cas_number",
// 		"product_specificity",
// 		"store_location",
// 		"quantity",
// 		"unit",
// 		"barecode",
// 		"supplier",
// 		"creation_date",
// 		"modification_date",
// 		"entry_date",
// 		"exit_date",
// 		"opening_date",
// 		"expiration_date",
// 		"comment",
// 		"reference",
// 		"batch_number",
// 		"to_destroy?",
// 		"archive?",
// 	}

// 	// create a temp file
// 	if tmpFile, err = os.CreateTemp(os.TempDir(), "chimitheque-"); err != nil {
// 		logger.Log.Error("cannot create temporary file", err)
// 		return "", err
// 	}
// 	// creates a csv writer that uses the io buffer
// 	csvwr := csv.NewWriter(tmpFile)
// 	// write the header
// 	if err = csvwr.Write(header); err != nil {
// 		logger.Log.Error("cannot write header", err)
// 		return "", err
// 	}

// 	for _, s := range sts {
// 		if err = csvwr.Write(s.StorageToStringSlice()); err != nil {
// 			logger.Log.Error("cannot write entry", err)
// 			return "", err
// 		}
// 	}

// 	csvwr.Flush()

// 	return strings.Split(tmpFile.Name(), "chimitheque-")[1], nil
// }
