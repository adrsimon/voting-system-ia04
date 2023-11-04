// version 2.0.0

package tests

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"github.com/adrsimon/voting-system-ia04/comsoc/methods"
	"testing"
)

func TestBordaSWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := methods.BordaSWF(prefs)
	if res[1] != 4 {
		t.Errorf("error, result for 1 should be 4, %d computed", res[1])
	}
	if res[2] != 3 {
		t.Errorf("error, result for 2 should be 3, %d computed", res[2])
	}
	if res[3] != 2 {
		t.Errorf("error, result for 3 should be 2, %d computed", res[3])
	}
}

func TestBordaSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := methods.BordaSCF(prefs)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestMajoritySWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := methods.MajoritySWF(prefs)

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 0 {
		t.Errorf("error, result for 2 should be 0, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestMajoritySCF(t *testing.T) {
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
}

func TestApprovalSWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 3, 2},
		{2, 3, 1},
	}
	thresholds := []int64{2, 1, 2}

	res, _ := methods.ApprovalSWF(prefs, thresholds...)

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 2 {
		t.Errorf("error, result for 2 should be 2, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestApprovalSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
	}
	thresholds := []int64{2, 1, 2}

	res, err := methods.ApprovalSCF(prefs, thresholds...)

	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

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

	res1, _ := methods.CondorcetWinner(prefs1)
	res2, _ := methods.CondorcetWinner(prefs2)

	if len(res1) == 0 || res1[0] != 1 {
		t.Errorf("error, 1 should be the only best alternative for prefs1")
	}
	if len(res2) != 0 {
		t.Errorf("no best alternative for prefs2")
	}
}
