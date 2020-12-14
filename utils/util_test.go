package utils_test

import (
	"testing"

	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/utils"
)

func TestIsCasNumber(t *testing.T) {
	c := "1825-62-3"
	if !handlers.IsCasNumber(c) {
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
	f := "C8H14O6"
	if sortedf, err = utils.SortEmpiricalFormula(f); err != nil {
		t.Errorf("%s is not a valid formula: %v", f, err)
	}
	if sortedf != "C8H14O6" {
		t.Errorf("%s was not sorted - output: %s", f, sortedf)
	}
}

func TestLinearToEmpiricalFormula(t *testing.T) {
	var (
		sortedf string
		err     error
	)
	f := "NaCl"
	if sortedf, err = utils.SortEmpiricalFormula(f); err != nil {
		t.Errorf("%s is not a valid formula: %v", f, err)
	}
	if sortedf != "ClNa" {
		t.Errorf("%s was not converted - output: %s", f, sortedf)
	}
}
