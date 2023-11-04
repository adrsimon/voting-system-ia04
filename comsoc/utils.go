package comsoc

import (
	"errors"
	"slices"
)

func Rank(alt Alternative, prefs []Alternative) int {
	for i, v := range prefs {
		if alt == v {
			return i
		}
	}
	return -1
}

func IsPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	return Rank(alt1, prefs) < Rank(alt2, prefs)
}

func IsPrefProfil(alt1, alt2 Alternative, p Profile) (bool, error) {
	count := 0
	for _, alts := range p {
		if IsPref(alt1, alt2, alts) {
			count++
		} else {
			count--
		}
	}
	if count == 0 {
		return false, errors.New("pas de préférence, erreur relation d'ordre nécessaire")
	} else if count >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func MaxCount(count Count) (bestAlts []Alternative) {
	maximum := -(2 << 31)
	for i, v := range count {
		if v > maximum {
			bestAlts = make([]Alternative, 0)
			bestAlts = append(bestAlts, i)
			maximum = v
		} else if v == maximum {
			bestAlts = append(bestAlts, i)
		}
	}
	return
}

func MinCount(count Count) (worstAlts []Alternative) {
	minimum := 2 << 31
	for i, v := range count {
		if v < minimum {
			worstAlts = make([]Alternative, 0)
			worstAlts = append(worstAlts, i)
			minimum = v
		} else if v == minimum {
			worstAlts = append(worstAlts, i)
		}
	}
	return
}

func DeleteAlternative(p Profile, i int64) Profile {
	for j, v := range p {
		rank := Rank(Alternative(i), v)
		p[j] = append(v[:rank], v[rank+1:]...)
	}
	return p
}

func CheckProfile(prefs []Alternative, alts []Alternative) error {
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

func CheckProfileAlternative(prefs Profile, alts []Alternative) error {
	for _, v := range prefs {
		err := CheckProfile(v, alts)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckSubProfile(prefs []Alternative, alts []Alternative) error {
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
