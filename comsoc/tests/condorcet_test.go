package tests

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"github.com/adrsimon/voting-system-ia04/comsoc/methods"
	"testing"
)

func TestCondorcet(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{3, 2, 1},
		{3, 2, 1},
	}

	bests, err := methods.CondorcetWinner(prefs)
	if err != nil {
		t.Error(err)
	}
	if len(bests) != 1 && bests[0] != 3 {
		t.Errorf("error, winner should be 3, %d computed", bests[0])
	}

	// implémenter cas d'erreur dans condorcet (propagation d'erreur de isBestProfile)
}
