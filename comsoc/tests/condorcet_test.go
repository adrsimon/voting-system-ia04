package tests

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"github.com/adrsimon/voting-system-ia04/comsoc/methods"
	"testing"
)

func TestCondorcetWinner(t *testing.T) {
	prefs1 := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	prefs2 := [][]comsoc.Alternative{
		{1, 2, 3},
		{2, 3, 1},
		{3, 1, 2},
	}

	prefs3 := [][]comsoc.Alternative{
		{1, 3, 2},
		{3, 2, 1},
		{3, 2, 1},
	}

	res1, _ := methods.CondorcetWinner(prefs1)
	res2, _ := methods.CondorcetWinner(prefs2)
	res3, _ := methods.CondorcetWinner(prefs3)

	if len(res1) == 0 || res1[0] != 1 {
		t.Errorf("error, 1 should be the only best alternative for prefs1")
	}
	if len(res2) != 0 {
		t.Errorf("no best alternative for prefs2")
	}
	if len(res3) != 1 && res3[0] != 3 {
		t.Errorf("error, winner should be 3, %d computed", res3[0])
	}
}
