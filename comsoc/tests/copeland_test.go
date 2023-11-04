// version 2.0.0

package tests

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"github.com/adrsimon/voting-system-ia04/comsoc/methods"
	"testing"
)

func TestCopelandSWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := methods.CopelandSWF(prefs)
	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 0 {
		t.Errorf("error, result for 2 should be 0, %d computed", res[2])
	}
	if res[3] != -2 {
		t.Errorf("error, result for 3 should be -2, %d computed", res[3])
	}
}

func TestCopelandSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := methods.CopelandSCF(prefs)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative, %d found", res[0])
	}
}
