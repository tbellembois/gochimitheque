package globals

import (
	"bytes"
	"database/sql"
	"errors"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/casbin/casbin/v2"
	"github.com/gorilla/schema"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/locales"
	"golang.org/x/text/language"
)

const (
	// LetterBytes is the list of letters
	LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// AtomByLength is a string slice sorter.
type AtomByLength []string

func (s AtomByLength) Len() int           { return len(s) }
func (s AtomByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s AtomByLength) Less(i, j int) bool { return len(s[i]) > len(s[j]) }

var (
	// Atoms is the list of existing atoms
	Atoms = map[string]string{
		"H":  "hydrogen",
		"He": "helium",
		"Li": "lithium",
		"Be": "berylium",
		"B":  "boron",
		"C":  "carbon",
		"N":  "nitrogen",
		"O":  "oxygen",
		"F":  "fluorine",
		"Ne": "neon",
		"Na": "sodium",
		"Mg": "magnesium",
		"Al": "aluminium",
		"Si": "silicon",
		"P":  "phosphorus",
		"S":  "sulfure",
		"Cl": "chlorine",
		"Ar": "argon",
		"K":  "potassium",
		"Ca": "calcium",
		"Sc": "scandium",
		"Ti": "titanium",
		"V":  "vanadium",
		"Cr": "chromium",
		"Mn": "manganese",
		"Fe": "iron",
		"Co": "cobalt",
		"Ni": "nickel",
		"Cu": "copper",
		"Zn": "zinc",
		"Ga": "gallium",
		"Ge": "germanium",
		"As": "arsenic",
		"Se": "sefeniuo",
		"Br": "bromine",
		"Kr": "krypton",
		"Rb": "rubidium",
		"Sr": "strontium",
		"Y":  "yltrium",
		"Zr": "zirconium",
		"Nb": "niobium",
		"Mo": "molybdenum",
		"Tc": "technetium",
		"Ru": "ruthenium",
		"Rh": "rhodium",
		"Pd": "palladium",
		"Ag": "silver",
		"Cd": "cadmium",
		"In": "indium",
		"Sn": "tin",
		"Sb": "antimony",
		"Te": "tellurium",
		"I":  "iodine",
		"Xe": "xenon",
		"Cs": "caesium",
		"Ba": "barium",
		"Hf": "hafnium",
		"Ta": "tantalum",
		"W":  "tungsten",
		"Re": "rhenium",
		"Os": "osmium",
		"Ir": "iridium",
		"Pt": "platinium",
		"Au": "gold",
		"Hg": "mercury",
		"Tl": "thallium",
		"Pb": "lead",
		"Bi": "bismuth",
		"Po": "polonium",
		"At": "astatine",
		"Rn": "radon",
		"Fr": "francium",
		"Ra": "radium",
		"Rf": "rutherfordium",
		"Db": "dubnium",
		"Sg": "seaborgium",
		"Bh": "bohrium",
		"Hs": "hassium",
		"Mt": "meitnerium",
		"Ds": "darmstadtium",
		"Rg": "roentgenium",
		"Cn": "copemicium",
		"La": "lanthanum",
		"Ce": "cerium",
		"Pr": "praseodymium",
		"Nd": "neodymium",
		"Pm": "promethium",
		"Sm": "samarium",
		"Eu": "europium",
		"Gd": "gadolinium",
		"Tb": "terbium",
		"Dy": "dysprosium",
		"Ho": "holmium",
		"Er": "erbium",
		"Tm": "thulium",
		"Yb": "ytterbium",
		"Lu": "lutetium",
		"Ac": "actinium",
		"Th": "thorium",
		"Pa": "protactinium",
		"U":  "uranium",
		"Np": "neptunium",
		"Pu": "plutonium",
		"Am": "americium",
		"Cm": "curium",
		"Bk": "berkelium",
		"Cf": "californium",
		"Es": "einsteinium",
		"Fm": "fermium",
		"Md": "mendelevium",
		"No": "nobelium",
		"Lr": "lawrencium",
		"D":  "deuterium",
	}

	// FormulaRe is the regex matching a chemical formula
	FormulaRe *regexp.Regexp

	// BasicMolRe is the regex matching a chemical formula (atoms and numbers only)
	BasicMolRe *regexp.Regexp

	// OneGroupMolRe is a (AYZ)n molecule like regex
	OneGroupMolRe *regexp.Regexp

	// Log is the general application logger
	Log *logrus.Logger
	// LogInternal is the application logger used to log fatal errors
	LogInternal *logrus.Logger
	// TokenSignKey is the JWT token signing key
	TokenSignKey []byte
	// Enforcer is the casbin enforcer
	Enforcer *casbin.Enforcer
	// JSONAdapterData is the casbin JSON data source
	JSONAdapterData []byte
	// Decoder is the form<>struct gorilla decoder
	Decoder *schema.Decoder
	// ProxyPath is the application proxy path if behind a proxy
	// "/"" by default
	ProxyPath string
	// ProxyURL is application base url
	// "http://localhost:8081" by default
	ProxyURL string
	// ApplicationFullURL is application full url
	// "http://localhost:8081" by default
	// "ProxyURL + ProxyPath" if behind a proxy
	ApplicationFullURL string
	// MailServerAddress is the SMTP server address
	// such as smtp.univ.fr
	MailServerAddress string
	// MailServerSender is the username used
	// to send mails
	MailServerSender string
	// MailServerPort is the SMTP server port
	MailServerPort string
	// MailServerUseTLS specify if a TLS SMTP connection
	// should be used
	MailServerUseTLS bool
	// MailServerTLSSkipVerify bypass the SMTP TLS verification
	MailServerTLSSkipVerify bool
	// Bundle is the i18n configuration bundle
	Bundle *i18n.Bundle
	// Localizer is the i18n translator
	Localizer *i18n.Localizer
	// BuildID is a compile time variable
	BuildID string
	// DisableCache disables the views cache
	DisableCache bool

	err error
)

// Convertors for sql.Null* types so that they can be
// used with gorilla/schema
func init() {
	// generate JWT signing key
	if TokenSignKey, err = GenSymmetricKey(64); err != nil {
		panic(err)
	}

	Decoder = schema.NewDecoder()
	SchemaRegisterSQLNulls(Decoder)

	// load translations
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	Bundle.MustParseMessageFileBytes(locales.LOCALES_EN, "en.toml")
	Bundle.MustParseMessageFileBytes(locales.LOCALES_FR, "fr.toml")

	Localizer = i18n.NewLocalizer(Bundle)

	// general formula regex
	FormulaRe = regexp.MustCompile(`[A-Za-z0-9,\^]+`)

	sortedAtoms := make([]string, 0, len(Atoms))
	for k := range Atoms {
		sortedAtoms = append(sortedAtoms, k)
	}
	// the atom must be sorted by decreasing size
	// to match first Cl before C for example
	sort.Sort(AtomByLength(sortedAtoms))

	// building the basic molecule regex
	// (atom1|atom2|...)([1-9]*)
	var buf bytes.Buffer
	buf.WriteString("(")
	for _, a := range sortedAtoms {
		buf.WriteString(a)
		buf.WriteString("|")
	}
	// removing the last |
	buf.Truncate(buf.Len() - 1)
	buf.WriteString(")")
	buf.WriteString("([1-9]+)*")
	BasicMolRe = regexp.MustCompilePOSIX(buf.String())

	// building the one group molecule regex
	buf.Reset()
	buf.WriteString("(?:\\(|\\[)")
	buf.WriteString("((?:[")
	for _, a := range sortedAtoms {
		buf.WriteString(a)
		buf.WriteString("|")
	}
	// removing the last |
	buf.Truncate(buf.Len() - 1)
	buf.WriteString("]+[1-9]*)+)")
	buf.WriteString("(?:\\)|\\])")
	buf.WriteString("([1-9]*)")

	OneGroupMolRe = regexp.MustCompile(buf.String())
}

// SchemaRegisterSQLNulls registers the custom null type to the application
func SchemaRegisterSQLNulls(d *schema.Decoder) {
	nullString, nullBool, nullInt64, nullFloat64, nullTime := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}, sql.NullTime{}

	d.RegisterConverter(nullString, ConvertSQLNullString)
	d.RegisterConverter(nullBool, ConvertSQLNullBool)
	d.RegisterConverter(nullInt64, ConvertSQLNullInt64)
	d.RegisterConverter(nullFloat64, ConvertSQLNullFloat64)
	d.RegisterConverter(nullTime, ConvertSQLNullTime)
}

// ConvertSQLNullTime converts a string into a NullTime
func ConvertSQLNullTime(value string) reflect.Value {
	var (
		e error
		t time.Time
	)
	if t, e = time.Parse("2006-01-02", string(value)); e != nil {
		Log.Error(e.Error())
		return reflect.Value{}
	}

	v := sql.NullTime{}
	if err := v.Scan(t); err != nil {
		Log.Error(e.Error())
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

// ConvertSQLNullString converts a string into a NullString
func ConvertSQLNullString(value string) reflect.Value {
	v := sql.NullString{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

// ConvertSQLNullBool converts a string into a NullBool
func ConvertSQLNullBool(value string) reflect.Value {
	v := sql.NullBool{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

// ConvertSQLNullInt64 converts a string into a NullInt64
func ConvertSQLNullInt64(value string) reflect.Value {
	v := sql.NullInt64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

// ConvertSQLNullFloat64 converts a string into a NullFloat64
func ConvertSQLNullFloat64(value string) reflect.Value {
	v := sql.NullFloat64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

// GenSymmetricKey generates a key for the JWT encryption
// https://github.com/northbright/Notes/blob/master/jwt/generate_hmac_secret_key_for_jwt.md
func GenSymmetricKey(bits int) (k []byte, err error) {
	if bits <= 0 || bits%8 != 0 {
		return nil, errors.New("key size error")
	}

	size := bits / 8
	k = make([]byte, size)
	if _, err = rand.Read(k); err != nil {
		return nil, err
	}

	return k, nil
}
