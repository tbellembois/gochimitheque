package zmqclient

type Filter struct {
	Search  string `json:"search"`
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
