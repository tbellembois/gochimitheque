package global

import (
	"database/sql"
	"database/sql/driver"
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
	TokenSignKey            []byte
	Decoder                 *schema.Decoder
	ProxyPath               string // application proxy path if behind a proxy
	ProxyURL                string // application url if behind a proxy
	MailServerAddress       string
	MailServerSender        string
	MailServerPort          string
	MailServerUseTLS        bool
	MailServerTLSSkipVerify bool
	Bundle                  *i18n.Bundle    // i18n bundle
	Localizer               *i18n.Localizer // application i18n localizer
)

// Convertors for sql.Null* types so that they can be
// used with gorilla/schema
func init() {
	TokenSignKey = []byte("secret")
	Decoder = schema.NewDecoder()
	SchemaRegisterSQLNulls(Decoder)

	// load translations
	Bundle = &i18n.Bundle{DefaultLanguage: language.English}
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	Bundle.MustParseMessageFileBytes(locales.LOCALES_EN, "en.toml")
	Bundle.MustParseMessageFileBytes(locales.LOCALES_FR, "fr.toml")

	Localizer = i18n.NewLocalizer(Bundle)
}

func SchemaRegisterSQLNulls(d *schema.Decoder) {
	//nullString, nullBool, nullInt64, nullFloat64, nullTime := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}, NullTime{}
	nullString, nullBool, nullInt64, nullFloat64 := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}

	d.RegisterConverter(nullString, ConvertSQLNullString)
	d.RegisterConverter(nullBool, ConvertSQLNullBool)
	d.RegisterConverter(nullInt64, ConvertSQLNullInt64)
	d.RegisterConverter(nullFloat64, ConvertSQLNullFloat64)
	//d.RegisterConverter(nullTime, ConvertSQLNullTime)
}

// func ConvertSQLNullTime(value string) reflect.Value {
// 	v := NullTime{}
// 	if err := v.Scan(value); err != nil {
// 		return reflect.Value{}
// 	}

// 	return reflect.ValueOf(v)
// }

func ConvertSQLNullString(value string) reflect.Value {
	v := sql.NullString{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func ConvertSQLNullBool(value string) reflect.Value {
	v := sql.NullBool{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func ConvertSQLNullInt64(value string) reflect.Value {
	v := sql.NullInt64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func ConvertSQLNullFloat64(value string) reflect.Value {
	v := sql.NullFloat64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}
