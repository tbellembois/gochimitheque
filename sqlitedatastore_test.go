package main

import (
	"database/sql"
	"testing"

	
	"github.com/tbellembois/gochimitheque/models"
)

var (
	err       error
	dbname    = "./storage.db"
	datastore models.Datastore
)

func init() {
	if datastore, err = models.NewSQLiteDBstore(dbname); err != nil {
		global.Log.Fatal(err)
	}
}

func BenchmarkComputeStockStorelocation(b *testing.B) {

	p := models.Product{ProductID: 2407}
	s := models.StoreLocation{StoreLocationID: sql.NullInt64{Int64: 281}}
	u := models.Unit{UnitID: sql.NullInt64{Int64: 1}}
	for n := 0; n < b.N; n++ {
		datastore.ComputeStockStorelocation(p, &s, u)
	}
}
