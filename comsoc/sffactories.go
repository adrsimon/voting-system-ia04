package comsoc

import "golang.org/x/exp/slices"

func SWFFactory(swf func(p Profile, thresholds ...int64) (Count, error), tiebreak func([]Alternative) (Alternative, error)) func(Profile, ...int64) ([]Alternative, error) {
	f := func(profile Profile, thresholds ...int64) ([]Alternative, error) {
		count, err := swf(profile, thresholds...)
		if err != nil {
			return nil, err
		}
		alts := make([]Alternative, 0)
		for len(count) >= 1 {
			maximum := MaxCount(count)
			for len(maximum) > 1 {
				maxAlt, err := tiebreak(maximum)
				if err != nil {
					return nil, err
				}
				delete(count, maxAlt)
				alts = append(alts, maxAlt)
				idx := slices.Index(maximum, maxAlt)
				maximum = append(maximum[:idx], maximum[idx+1:]...)
			}
			alts = append(alts, maximum[0])
			delete(count, maximum[0])
		}
		return alts, nil
	}
	return f
}

func SCFFactory(scf func(p Profile, thresholds ...int64) ([]Alternative, error), tiebreak func([]Alternative) (Alternative, error)) func(Profile, ...int64) (Alternative, error) {
	f := func(profile Profile, thresholds ...int64) (Alternative, error) {
		alts, err := scf(profile, thresholds...)
		if err != nil {
			return Alternative(0), err
		}
		bestAlt, err := tiebreak(alts)
		if err != nil {
			return Alternative(0), err
		}
		return bestAlt, nil
	}
	return f
}
