package handlers

import (
	"database/sql"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/tbellembois/gochimitheque/models"
)

// TokenSignKey is the JWT token signing key
var TokenSignKey = []byte("secret")
var Decoder = schema.NewDecoder()

// Convertors for sql.Null* types so that they can be
// used with gorilla/schema
func init() {
	SchemaRegisterSQLNulls(Decoder)
}

func SchemaRegisterSQLNulls(d *schema.Decoder) {
	nullString, nullBool, nullInt64, nullFloat64, nullTime := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}, models.NullTime{}

	d.RegisterConverter(nullString, ConvertSQLNullString)
	d.RegisterConverter(nullBool, ConvertSQLNullBool)
	d.RegisterConverter(nullInt64, ConvertSQLNullInt64)
	d.RegisterConverter(nullFloat64, ConvertSQLNullFloat64)
	d.RegisterConverter(nullTime, ConvertSQLNullTime)
}

func ConvertSQLNullTime(value string) reflect.Value {
	v := models.NullTime{}
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
