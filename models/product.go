package models

import (
	"database/sql"
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/tbellembois/gochimitheque/logger"
)

// Product is a chemical product card.
type Product struct {
	ProductID              int            `db:"product_id" json:"product_id" schema:"product_id"`
	ProductInchi           sql.NullString `db:"product_inchi" json:"product_inchi" schema:"product_inchi"`
	ProductInchikey        sql.NullString `db:"product_inchikey" json:"product_inchikey" schema:"product_inchikey"`
	ProductCanonicalSmiles sql.NullString `db:"product_canonical_smiles" json:"product_canonical_smiles" schema:"product_canonical_smiles"`
	ProductSpecificity     sql.NullString `db:"product_specificity" json:"product_specificity" schema:"product_specificity" `
	ProductMSDS            sql.NullString `db:"product_msds" json:"product_msds" schema:"product_msds" `
	ProductRestricted      sql.NullBool   `db:"product_restricted" json:"product_restricted" schema:"product_restricted" `
	ProductRadioactive     sql.NullBool   `db:"product_radioactive" json:"product_radioactive" schema:"product_radioactive" `
	ProductThreeDFormula   sql.NullString `db:"product_threedformula" json:"product_threedformula" schema:"product_threedformula" `
	ProductTwoDFormula     sql.NullString `db:"product_twodformula" json:"product_twodformula" schema:"product_twodformula" `
	// ProductMolFormula      sql.NullString `db:"product_molformula" json:"product_molformula" schema:"product_molformula" `
	ProductDisposalComment sql.NullString  `db:"product_disposalcomment" json:"product_disposalcomment" schema:"product_disposalcomment" `
	ProductRemark          sql.NullString  `db:"product_remark" json:"product_remark" schema:"product_remark" `
	ProductMolecularWeight sql.NullFloat64 `db:"product_molecularweight" json:"product_molecularweight" schema:"product_molecularweight" `
	ProductTemperature     sql.NullInt64   `db:"product_temperature" json:"product_temperature" schema:"product_temperature" `
	ProductSheet           sql.NullString  `db:"product_sheet" json:"product_sheet" schema:"product_sheet" `
	ProductNumberPerCarton sql.NullInt64   `db:"product_number_per_carton" json:"product_number_per_carton" schema:"product_number_per_carton" `
	ProductNumberPerBag    sql.NullInt64   `db:"product_number_per_bag" json:"product_number_per_bag" schema:"product_number_per_bag" `
	ProductType            string          `db:"-" json:"product_type" schema:"product_type"`
	EmpiricalFormula       `db:"empiricalformula" json:"empiricalformula" schema:"empiricalformula"`
	LinearFormula          `db:"linearformula" json:"linearformula" schema:"linearformula"`
	PhysicalState          `db:"physicalstate" json:"physicalstate" schema:"physicalstate"`
	SignalWord             `db:"signalword" json:"signalword" schema:"signalword"`
	Person                 `db:"person" json:"person" schema:"person"`
	CasNumber              `db:"casnumber" json:"casnumber" schema:"casnumber"`
	CeNumber               `db:"cenumber" json:"cenumber" schema:"cenumber"`
	Name                   `db:"name" json:"name" schema:"name"`
	ProducerRef            `db:"producerref" json:"producerref" schema:"producerref"`
	Category               `db:"category" json:"category" schema:"category"`
	UnitTemperature        Unit `db:"unit_temperature" json:"unit_temperature" schema:"unit_temperature"`
	UnitMolecularWeight    Unit `db:"unit_molecularweight" json:"unit_molecularweight" schema:"unit_molecularweight"`

	ClassOfCompound         []ClassOfCompound        `db:"-" schema:"classofcompound" json:"classofcompound"`
	Synonyms                []Name                   `db:"-" schema:"synonyms" json:"synonyms"`
	Symbols                 []Symbol                 `db:"-" schema:"symbols" json:"symbols"`
	HazardStatements        []HazardStatement        `db:"-" schema:"hazardstatements" json:"hazardstatements"`
	PrecautionaryStatements []PrecautionaryStatement `db:"-" schema:"precautionarystatements" json:"precautionarystatements"`
	SupplierRefs            []SupplierRef            `db:"-" json:"supplierrefs" schema:"supplierrefs"`
	Tags                    []Tag                    `db:"-" json:"tags" schema:"tags"`

	Bookmark *Bookmark `db:"bookmark" json:"bookmark" schema:"bookmark"` // not in db but sqlx requires the "db" entry

	// archived storage count in the logged user entity(ies)
	ProductASC int `db:"product_asc" json:"product_asc" schema:"product_asc"` // not in db but sqlx requires the "db" entry
	// total storage count
	ProductTSC int `db:"product_tsc" json:"product_tsc" schema:"product_tsc"` // not in db but sqlx requires the "db" entry
	// storage count in the logged user entity(ies)
	ProductSC int `db:"product_sc" json:"product_sc" schema:"product_sc"` // not in db but sqlx requires the "db" entry
	// storage barecode concatenation
	ProductSL sql.NullString `db:"product_sl" json:"product_sl" schema:"product_sl" ` // not in db but sqlx requires the "db" entry
	// hazard statement CMR concatenation
	HazardStatementCMR sql.NullString `db:"hazardstatement_cmr" json:"hazardstatement_cmr" schema:"hazardstatement_cmr" ` // not in db but sqlx requires the "db" entry
}

func (p Product) ProductToStringSlice() []string {
	ret := make([]string, 0)

	ret = append(ret, strconv.Itoa(p.ProductID))

	ret = append(ret, p.NameLabel)
	syn := ""

	for _, s := range p.Synonyms {
		syn += "|" + s.NameLabel
	}

	ret = append(ret, syn)

	ret = append(ret, p.CasNumberLabel.String)
	ret = append(ret, p.CeNumberLabel.String)

	ret = append(ret, p.ProductSpecificity.String)
	ret = append(ret, p.EmpiricalFormulaLabel.String)
	ret = append(ret, p.LinearFormulaLabel.String)
	ret = append(ret, p.ProductThreeDFormula.String)

	ret = append(ret, p.ProductMSDS.String)

	ret = append(ret, p.PhysicalStateLabel.String)

	ret = append(ret, p.SignalWordLabel.String)

	coc := ""

	for _, c := range p.ClassOfCompound {
		coc += "|" + c.ClassOfCompoundLabel
	}

	ret = append(ret, coc)
	sym := ""

	for _, s := range p.Symbols {
		sym += "|" + s.SymbolLabel
	}

	ret = append(ret, sym)
	hs := ""

	for _, h := range p.HazardStatements {
		hs += "|" + h.HazardStatementReference
	}

	ret = append(ret, hs)
	ps := ""

	for _, p := range p.PrecautionaryStatements {
		ps += "|" + p.PrecautionaryStatementReference
	}

	ret = append(ret, ps)

	ret = append(ret, p.ProductRemark.String)
	ret = append(ret, p.ProductDisposalComment.String)

	ret = append(ret, strconv.FormatBool(p.ProductRestricted.Bool))
	ret = append(ret, strconv.FormatBool(p.ProductRadioactive.Bool))

	return ret
}

// ProductsToCSV returns a file name of the products prs
// exported into CSV.
func ProductsToCSV(prs []Product) string {
	header := []string{
		"product_id",
		"product_name",
		"product_synonyms",
		"product_cas",
		"product_ce",
		"product_specificity",
		"empirical_formula",
		"linear_formula",
		"3D_formula",
		"MSDS",
		"class_of_compounds",
		"physical_state",
		"signal_word",
		"symbols",
		"hazard_statements",
		"precautionary_statements",
		"remark",
		"disposal_comment",
		"restricted?",
		"radioactive?",
	}

	// create a temp file
	tmpFile, err := os.CreateTemp(os.TempDir(), "chimitheque-")
	if err != nil {
		logger.Log.Error("cannot create temporary file", err)
	}
	// creates a csv writer that uses the io buffer
	csvwr := csv.NewWriter(tmpFile)
	// write the header
	_ = csvwr.Write(header)

	for _, p := range prs {
		_ = csvwr.Write(p.ProductToStringSlice())
	}

	csvwr.Flush()

	return strings.Split(tmpFile.Name(), "chimitheque-")[1]
}
