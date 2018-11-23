package main

import (
	"testing"

	"github.com/tbellembois/gochimitheque/utils"
)

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
