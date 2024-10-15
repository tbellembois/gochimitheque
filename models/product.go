package models

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/tbellembois/gochimitheque/logger"
)

// Product is a chemical product card.
type Product struct {
	ProductID              int     `db:"product_id" json:"product_id" schema:"product_id"`
	ProductInchi           *string `db:"product_inchi" json:"product_inchi" schema:"product_inchi"`
	ProductInchikey        *string `db:"product_inchikey" json:"product_inchikey" schema:"product_inchikey"`
	ProductCanonicalSmiles *string `db:"product_canonical_smiles" json:"product_canonical_smiles" schema:"product_canonical_smiles"`
	ProductSpecificity     *string `db:"product_specificity" json:"product_specificity" schema:"product_specificity" `
	ProductMSDS            *string `db:"product_msds" json:"product_msds" schema:"product_msds" `
	ProductRestricted      bool    `db:"product_restricted" json:"product_restricted" schema:"product_restricted" `
	ProductRadioactive     bool    `db:"product_radioactive" json:"product_radioactive" schema:"product_radioactive" `
	ProductThreeDFormula   *string `db:"product_threed_formula" json:"product_threed_formula" schema:"product_threed_formula" `
	ProductTwoDFormula     *string `db:"product_twod_formula" json:"product_twod_formula" schema:"product_twod_formula" `
	// ProductMolFormula      sql.NullString `db:"product_molformula" json:"product_molformula" schema:"product_molformula" `
	ProductDisposalComment *string  `db:"product_disposal_comment" json:"product_disposal_comment" schema:"product_disposal_comment" `
	ProductRemark          *string  `db:"product_remark" json:"product_remark" schema:"product_remark" `
	ProductMolecularWeight *float64 `db:"product_molecular_weight" json:"product_molecular_weight" schema:"product_molecular_weight" `
	ProductTemperature     *int64   `db:"product_temperature" json:"product_temperature" schema:"product_temperature" `
	ProductSheet           *string  `db:"product_sheet" json:"product_sheet" schema:"product_sheet" `
	ProductNumberPerCarton *int64   `db:"product_number_per_carton" json:"product_number_per_carton" schema:"product_number_per_carton" `
	ProductNumberPerBag    *int64   `db:"product_number_per_bag" json:"product_number_per_bag" schema:"product_number_per_bag" `
	ProductType            string   `db:"-" json:"product_type" schema:"product_type"`
	EmpiricalFormula       `db:"empirical_formula" json:"empirical_formula" schema:"empirical_formula"`
	LinearFormula          `db:"linear_formula" json:"linear_formula" schema:"linear_formula"`
	PhysicalState          `db:"physical_state" json:"physical_state" schema:"physical_state"`
	SignalWord             `db:"signal_word" json:"signal_word" schema:"signal_word"`
	Person                 `db:"person" json:"person" schema:"person"`
	CasNumber              `db:"cas_number" json:"cas_number" schema:"cas_number"`
	CeNumber               `db:"ce_number" json:"ce_number" schema:"ce_number"`
	Name                   `db:"name" json:"name" schema:"name"`
	ProducerRef            `db:"producer_ref" json:"producer_ref" schema:"producer_ref"`
	Category               `db:"category" json:"category" schema:"category"`
	UnitTemperature        Unit `db:"unit_temperature" json:"unit_temperature" schema:"unit_temperature"`
	UnitMolecularWeight    Unit `db:"unit_molecular_weight" json:"unit_molecular_weight" schema:"unit_molecular_weight"`

	ClassOfCompound         []ClassOfCompound        `db:"-" schema:"classes_of_compound" json:"classes_of_compound"`
	Synonyms                []Name                   `db:"-" schema:"synonyms" json:"synonyms"`
	Symbols                 []Symbol                 `db:"-" schema:"symbols" json:"symbols"`
	HazardStatements        []HazardStatement        `db:"-" schema:"hazard_statements" json:"hazard_statements"`
	PrecautionaryStatements []PrecautionaryStatement `db:"-" schema:"precautionary_statements" json:"precautionary_statements"`
	SupplierRefs            []SupplierRef            `db:"-" json:"supplier_refs" schema:"supplier_refs"`
	Tags                    []Tag                    `db:"-" json:"tags" schema:"tags"`

	Bookmark *Bookmark `db:"bookmark" json:"bookmark" schema:"bookmark"` // not in db but sqlx requires the "db" entry

	// archived storage count in the logged user entity(ies)
	ProductASC int `db:"product_asc" json:"product_asc" schema:"product_asc"` // not in db but sqlx requires the "db" entry
	// total storage count
	ProductTSC int `db:"product_tsc" json:"product_tsc" schema:"product_tsc"` // not in db but sqlx requires the "db" entry
	// storage count in the logged user entity(ies)
	ProductSC int `db:"product_sc" json:"product_sc" schema:"product_sc"` // not in db but sqlx requires the "db" entry
	// storage barecode concatenation
	ProductSL *string `db:"product_sl" json:"product_sl" schema:"product_sl" ` // not in db but sqlx requires the "db" entry
	// hazard statement CMR concatenation
	HazardStatementCMR *string `db:"product_hs_cmr" json:"product_hs_cmr" schema:"product_hs_cmr" ` // not in db but sqlx requires the "db" entry
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

	// ret = append(ret, p.CasNumberLabel.String)
	ret = append(ret, *p.CasNumberLabel)
	ret = append(ret, *p.CeNumberLabel)

	ret = append(ret, *p.ProductSpecificity)
	ret = append(ret, *p.EmpiricalFormulaLabel)
	ret = append(ret, *p.LinearFormulaLabel)
	ret = append(ret, *p.ProductThreeDFormula)

	ret = append(ret, *p.ProductMSDS)

	ret = append(ret, *p.PhysicalStateLabel)

	ret = append(ret, *p.SignalWordLabel)

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

	ret = append(ret, *p.ProductRemark)
	ret = append(ret, *p.ProductDisposalComment)

	ret = append(ret, strconv.FormatBool(p.ProductRestricted))
	ret = append(ret, strconv.FormatBool(p.ProductRadioactive))

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
