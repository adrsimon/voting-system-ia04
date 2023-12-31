package tests

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"github.com/adrsimon/voting-system-ia04/comsoc/methods"
	"testing"
)

func TestTiebreakFactory(t *testing.T) { // the lower alt wins
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
		{3, 2, 1},
	}

	res, err := methods.MajoritySCF(prefs)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 2 {
		t.Errorf("error, there should be a draw")
	}

	tiebreak := comsoc.TieBreakFactory([]comsoc.Alternative{1, 2, 3})

	res2, err2 := tiebreak(res)
	if err2 != nil {
		t.Error(err2)
	}
	if res2 != 1 {
		t.Errorf("error, winner should be 1, %d computed", res)
	}

}

func TestSWFFactory(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{3, 2, 1},
		{3, 2, 1},
	}

	tiebreak := comsoc.TieBreakFactory([]comsoc.Alternative{1, 2, 3})
	swf := comsoc.SWFFactory(methods.BordaSWF, tiebreak)
	res, err := swf(prefs)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 3 || res[0] != 3 {
		t.Errorf("error, winner should be 3, %d computed", res[0])
	}
	if res[1] != 1 {
		t.Errorf("error, second should be 1, %d computed", res[1])
	}
	if res[2] != 2 {
		t.Errorf("error, winner should be 2, %d computed", res[2])
	}
	tiebreak2 := comsoc.TieBreakFactory([]comsoc.Alternative{3, 2, 1})
	swf2 := comsoc.SWFFactory(methods.BordaSWF, tiebreak2)
	res2, err2 := swf2(prefs)
	if err2 != nil {
		t.Error(err2)
	}
	if len(res2) != 3 || res2[0] != 3 {
		t.Errorf("error, winner should be 3, %d computed", res2[0])
	}
	if res2[1] != 2 {
		t.Errorf("error, second should be 2, %d computed", res2[1])
	}
	if res2[2] != 1 {
		t.Errorf("error, winner should be 1, %d computed", res2[2])
	}

}

func TestScfFactory(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{3, 2, 1},
		{3, 2, 1},
	}

	tiebreak := comsoc.TieBreakFactory([]comsoc.Alternative{1, 2, 3})
	scf := comsoc.SCFFactory(methods.BordaSCF, tiebreak)
	res, err := scf(prefs)

	if err != nil {
		t.Error(err)
	}
	if res != 3 {
		t.Errorf("error, winner should be 3, %d computed", res)
	}

	prefs2 := [][]comsoc.Alternative{
		{1, 3, 2},
		{3, 2, 1},
		{2, 1, 3},
	}
	scf2 := comsoc.SCFFactory(methods.BordaSCF, tiebreak)
	res2, err := scf2(prefs2)

	if err != nil {
		t.Error(err)
	}
	if res2 != 1 {
		t.Errorf("error, winner should be 1, %d computed", res)
	}
}
