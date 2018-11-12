package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/utils"
)

func TestComputeStockStorelocation(t *testing.T) {
	var (
		datastore *models.SQLiteDataStore
		err       error
		m         models.StockMap
	)

	m = make(models.StockMap)

	s := models.StoreLocation{StoreLocationID: sql.NullInt64{Valid: true, Int64: 1}}
	p := models.Product{ProductID: 8}
	u := models.Unit{UnitID: sql.NullInt64{Valid: true, Int64: 1}}

	if datastore, err = models.NewDBstore("/home/thbellem/workspace/workspace_Go/src/github.com/tbellembois/gochimitheque/storage.db"); err != nil {
		log.Panic(err)
	}
	datastore.ComputeStockStorelocation(p, s, u, &m)

}

func TestComputeStockEntity(t *testing.T) {
	var (
		datastore *models.SQLiteDataStore
		err       error
	)

	s := models.Entity{EntityID: 1}
	p := models.Product{ProductID: 8}
	if datastore, err = models.NewDBstore("/home/thbellem/workspace/workspace_Go/src/github.com/tbellembois/gochimitheque/storage.db"); err != nil {
		log.Panic(err)
	}
	datastore.ComputeStockEntity(p, s)
}

func TestIsCasNumber(t *testing.T) {
	c := "7732-18-5"
	if !utils.IsCasNumber(c) {
		t.Errorf("%s is not a valid cas number", c)
	}
}

func TestSortSimpleFormula(t *testing.T) {
	var (
		sortedf string
		err     error
	)
	f := "NaCl2"
	if sortedf, err = utils.SortSimpleFormula(f); err != nil {
		t.Errorf("%s is not a valid formula: %v", f, err)
	}
	if sortedf != "Cl2Na" {
		t.Errorf("%s was not sorted - output: %s", f, sortedf)
	}
}

func TestSortEmpiricalFormula(t *testing.T) {
	var (
		sortedf string
		err     error
	)
	f := "NaCl2"
	if sortedf, err = utils.SortEmpiricalFormula(f); err != nil {
		t.Errorf("%s is not a valid formula: %v", f, err)
	}
	if sortedf != "Cl2Na" {
		t.Errorf("%s was not sorted - output: %s", f, sortedf)
	}
}
