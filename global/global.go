package global

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/schema"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tbellembois/gochimitheque/locales"
	"golang.org/x/text/language"
)

// NullTime represent a nullable time
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// UnmarshalText parse the string date from the view
func (nt *NullTime) UnmarshalText(text []byte) (err error) {
	var e error
	nt.Time, e = time.Parse("2006-01-02", string(text))
	if e == nil {
		nt.Valid = true
	}
	return e
}

var (
	// TokenSignKey is the JWT token signing key
	TokenSignKey []byte
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
	// InternalServerErrorLog error log file
	InternalServerErrorLog *os.File
	// Bundle is the i18n configuration bundle
	Bundle *i18n.Bundle
	// Localizer is the i18n translator
	Localizer *i18n.Localizer
	// BuildID is a compile time variable
	BuildID string

	err error
)

// ChimithequeContextKey is the Go request context
// used in each request
type ChimithequeContextKey string

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
}

// SchemaRegisterSQLNulls registers the custom null type to the application
func SchemaRegisterSQLNulls(d *schema.Decoder) {
	nullString, nullBool, nullInt64, nullFloat64 := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}

	d.RegisterConverter(nullString, ConvertSQLNullString)
	d.RegisterConverter(nullBool, ConvertSQLNullBool)
	d.RegisterConverter(nullInt64, ConvertSQLNullInt64)
	d.RegisterConverter(nullFloat64, ConvertSQLNullFloat64)
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
