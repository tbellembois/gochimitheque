package zmqclient

// product type
type PubchemProduct struct {
	Name                *string   `json:"name"`
	Inchi               *string   `json:"inchi"`
	InchiKey            *string   `json:"inchi_key"`
	CanonicalSmiles     *string   `json:"canonical_smiles"`
	MolecularFormula    *string   `json:"molecular_formula"`
	Cas                 *string   `json:"cas"`
	Ec                  *string   `json:"ec"`
	MolecularWeight     *string   `json:"molecular_weight"`
	MolecularWeightUnit *string   `json:"molecular_weight_unit"`
	Synonyms            *[]string `json:"synonyms"`
	Symbols             *[]string `json:"symbols"`
	Signal              *[]string `json:"signal"`
	Hs                  *[]string `json:"hs"`
	Ps                  *[]string `json:"ps"`
	Twodpicture         *string   `json:"twodpicture"` // base64 encoded png
}

// compound types
type PropValue struct {
	Ival   *int     `json:"ival"`
	Fval   *float64 `json:"fval"`
	Binary *string  `json:"binary"`
	Sval   *string  `json:"sval"`
}

type PropURN struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

type Prop struct {
	URN   PropURN   `json:"urn"`
	Value PropValue `json:"value"`
}

type ID struct {
	ID CID `json:"id"`
}

type CID struct {
	CID int `json:"cid"`
}

type Section struct {
	TOCHeading  *string        `json:"TOCHeading"`
	TOCID       *int           `json:"TOCID"`
	Description string         `json:"Description"`
	URL         string         `json:"URL"`
	Section     *[]Section     `json:"Section"`
	Information *[]Information `json:"Information"`
}

type Information struct {
	ReferenceNumber int      `json:"ReferenceNumber"`
	Name            string   `json:"Name"`
	Description     string   `json:"Description"`
	Reference       []string `json:"Reference"`
	LicenseNote     []string `json:"LicenseNote"`
	LicenseURL      []string `json:"LicenseURL"`
	Value           Value    `json:"Value"`
}

type Value struct {
	Number               []float64           `json:"Number"`
	DateISO8601          []string            `json:"DateISO8601"`
	Boolean              []bool              `json:"Boolean"`
	Binary               []string            `json:"Binary"`
	BinaryToStore        []string            `json:"BinaryToStore"`
	ExternalDataURL      []string            `json:"ExternalDataURL"`
	ExternalTableName    string              `json:"ExternalTableName"`
	Unit                 string              `json:"Unit"`
	MimeType             string              `json:"MimeType"`
	ExternalTableNumRows int                 `json:"ExternalTableNumRows"`
	StringWithMarkup     *[]StringWithMarkup `json:"StringWithMarkup"`
}

type Markup struct {
	Start  float64 `json:"Start"`
	Length float64 `json:"Length"`
	URL    string  `json:"URL"`
	Type   string  `json:"Type"`
	Extra  string  `json:"Extra"`
}

type StringWithMarkup struct {
	String string   `json:"String"`
	Markup []Markup `json:"Markup"`
}

type Record struct {
	Record RecordContent `json:"Record"`
}

type RecordContent struct {
	RecordType        string        `json:"RecordType"`
	RecordNumber      int           `json:"RecordNumber"`
	RecordAccession   string        `json:"RecordAccession"`
	RecordTitle       string        `json:"RecordTitle"`
	RecordExternalURL string        `json:"RecordExternalURL"`
	Section           []Section     `json:"Section"`
	Information       []Information `json:"Information"`
}

type PCCompound struct {
	ID     ID     `json:"id"`
	Props  []Prop `json:"props"`
	Record Record `json:"record"`
}

type Compounds struct {
	PCCompounds []PCCompound `json:"PC_Compounds"`
	Record      Record       `json:"record"`
	Base64Png   string       `json:"base64_png"`
}

// autocomplete types
type DictionnaryTerms struct {
	Compound []string `json:"compound"`
}
type PubchemAutocomplete struct {
	Total           uint64           `json:"total"`
	DictionaryTerms DictionnaryTerms `json:"dictionary_terms"`
}

type RequestFilter struct {
	Search  string `json:"search"`
	Id      uint64 `json:"id"`
	OrderBy string `json:"order_by"`
	Order   string `json:"order"`
	Offset  uint64 `json:"offset"`
	Limit   uint64 `json:"limit"`

	Bookmark                bool   `json:"bookmark"`
	Borrowing               bool   `json:"borrowing"`
	CasNumber               int    `json:"cas_number"`
	CasNumberCmr            bool   `json:"ce_number"`
	Category                int    `json:"category"`
	CustomNamePartOf        string `json:"custom_name_part_of"`
	EmpiricalFormula        int    `json:"empirical_formula"`
	Entity                  int    `json:"entity"`
	EntityName              string `json:"entity_name"`
	HazardStatements        []int  `json:"hazard_statements"`
	History                 bool   `json:"history"`
	Ids                     []int  `json:"storages"`
	Name                    int    `json:"name"`
	Permission              string `json:"permission"`
	PrecautionaryStatements []int  `json:"precautionary_statements"`
	Producer                int    `json:"producer"`
	ProducerRef             int    `json:"producer_ref"`
	Product                 int    `json:"product"`
	ProductSpecificity      string `json:"product_specificity"`
	ShowBio                 bool   `json:"show_bio"`
	ShowChem                bool   `json:"show_chem"`
	ShowConsu               bool   `json:"show_consu"`
	SignalWord              int    `json:"signal_word"`
	Storage                 int    `json:"storage"`
	StorageArchive          bool   `json:"storage_archive"`
	StorageBarecode         string `json:"storage_barecode"`
	StorageBatchNumber      string `json:"storage_batch_number"`
	StorageToDestroy        bool   `json:"storage_to_destroy"`
	Storelocation           int    `json:"store_location"`
	StoreLocationCanStore   bool   `json:"store_location_can_store"`
	Supplier                int    `json:"supplier"`
	Symbols                 []int  `json:"symbols"`
	Tags                    []int  `json:"tags"`
	UnitType                string `json:"unit_type"`
}
