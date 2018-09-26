package main

import (
	"github.com/tbellembois/gochimitheque/utils"
	"testing"
)

func TestIsCasNumber(t *testing.T) {
	c := "7732-18-5"
	if !utils.IsCasNumber(c) {
		t.Errorf("%s is a valid cas number", c)
	}
}
