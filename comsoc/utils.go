package comsoc

import (
	"errors"
	"slices"
)

func rank(alt Alternative, prefs []Alternative) int {
	for i, v := range prefs {
		if alt == v {
			return i
		}
	}
	return -1
}

func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	return rank(alt1, prefs) < rank(alt2, prefs)
}

func maxCount(count Count) (bestAlts []Alternative) {
	max := 0
	for i, v := range count {
		if v > max {
			bestAlts = make([]Alternative, 0)
			bestAlts = append(bestAlts, i)
		} else if v == max {
			bestAlts = append(bestAlts, i)
		}
	}
	return
}

func checkProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) != len(alts) {
		return errors.New("alts and prefs doesn't have the same size")
	}

	prefsCheck := make([]Alternative, 0)
	for _, v := range prefs {
		if slices.Contains(prefsCheck, v) {
			return errors.New("prefs contains several times the same alternative")
		} else {
			prefsCheck = append(prefsCheck, v)
			if !slices.Contains(alts, v) {
				return errors.New("one of prefs values is not a member of alts")
			}
		}
	}
	return nil
}

func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	for _, v := range prefs {
		err := checkProfile(v, alts)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkSubProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) > len(alts) {
		return errors.New("prefs isn't a subset of alts")
	}

	prefsCheck := make([]Alternative, 0)
	for _, v := range prefs {
		if slices.Contains(prefsCheck, v) {
			return errors.New("prefs contains several times the same alternative")
		} else {
			prefsCheck = append(prefsCheck, v)
			if !slices.Contains(alts, v) {
				return errors.New("one of prefs values is not a member of alts")
			}
		}
	}
	return nil
}
