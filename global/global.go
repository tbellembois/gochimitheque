package global

import (
	"database/sql"
	"database/sql/driver"
	"github.com/gorilla/schema"
	"reflect"
	"time"
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

var (
	// TokenSignKey is the JWT token signing key
	TokenSignKey       []byte
	Decoder            *schema.Decoder
	ProxyPath          string // application proxy path if behind a proxy
	ProxyURL           string // application url if behind a proxy
	MailServerAddress  string
	MailServerPassword string
	MailServerSender   string
	MailServerUser     string
	MailServerPort     string
)

// Convertors for sql.Null* types so that they can be
// used with gorilla/schema
func init() {
	TokenSignKey = []byte("secret")
	Decoder = schema.NewDecoder()
	SchemaRegisterSQLNulls(Decoder)
}

func SchemaRegisterSQLNulls(d *schema.Decoder) {
	nullString, nullBool, nullInt64, nullFloat64, nullTime := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}, NullTime{}

	d.RegisterConverter(nullString, ConvertSQLNullString)
	d.RegisterConverter(nullBool, ConvertSQLNullBool)
	d.RegisterConverter(nullInt64, ConvertSQLNullInt64)
	d.RegisterConverter(nullFloat64, ConvertSQLNullFloat64)
	d.RegisterConverter(nullTime, ConvertSQLNullTime)
}

func ConvertSQLNullTime(value string) reflect.Value {
	v := NullTime{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

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
