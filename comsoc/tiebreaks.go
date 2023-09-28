package comsoc

import (
	"errors"
	"slices"
)

// TODO : quand appel a cette fonction, il faut check la composition d'orderedAlts
func TieBreakFactory(orderedAlts []Alternative) (func ([]Alternative) (Alternative, error)) {
	f := func (alts []Alternative) (Alternative, error) {
		if len(alts) == 0 {
			return Alternative(0), errors.New("no alts to check")
		}
		err := checkSubProfile(alts, orderedAlts)
		if err != nil {
			return Alternative(0), err
		}
		for _, v := range orderedAlts {
			if slices.Contains(alts, v) {
				return v, nil
			}
		}
		return Alternative(0), errors.New("unreacheable")
	}
	return f
}