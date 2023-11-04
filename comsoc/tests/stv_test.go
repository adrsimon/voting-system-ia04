package tests

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"github.com/adrsimon/voting-system-ia04/comsoc/methods"
	"testing"
)

func TestSTVSWF(t *testing.T) { // pensez  a expliquer valeur du count readme
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := methods.STVSWF(prefs, []int64{1, 2, 3}...)

	if res[1] != 3 {
		t.Errorf("error, result for 1 should be 3, %d computed", res[1])
	}
	if res[2] != 1 {
		t.Errorf("error, result for 2 should be 1, %d computed", res[2])
	}
	if res[3] != 2 {
		t.Errorf("error, result for 3 should be 2, %d computed", res[3])
	}
	prefs2 := [][]comsoc.Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
	}

	res2, _ := methods.STVSWF(prefs2, []int64{4, 2, 3, 1}...)

	if res2[1] != 3 {
		t.Errorf("error, result for 1 should be 2, %d computed", res2[1])
	}
	if res2[2] != 1 {
		t.Errorf("error, result for 2 should be 0, %d computed", res2[2])
	}
	if res2[3] != 4 {
		t.Errorf("error, result for 3 should be 4, %d computed", res2[3])
	}
	if res2[4] != 2 {
		t.Errorf("error, result for 4 should be 0, %d computed", res2[4])
	}
}

func TestSTVSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := methods.MajoritySCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
	prefs2 := [][]comsoc.Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
	}

	res2, err := methods.STVSCF(prefs2, []int64{4, 2, 3, 1}...)
	if err != nil {
		t.Error(err)
	}

	if len(res2) != 1 || res2[0] != 3 {
		t.Errorf("error, 3 should be the only best Alternative")
	}
}
