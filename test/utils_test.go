package main

import (
	"github.com/tbellembois/gochimitheque/utils"
	"testing"
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
	f := "ClNaC2"
	if sortedf, err = utils.SortSimpleFormula(f); err != nil {
		t.Errorf("%s is not a valid formula: %v", f, err)
	}
	if f != sortedf {
		t.Errorf("%s was not sorted", f)
	}
}
