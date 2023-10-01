package comsoc

func SWFFactory(swf func(p Profile) (Count, error), tiebreak func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	f := func(profile Profile) ([]Alternative, error) {
		count, err := swf(profile)
		if err != nil {
			return nil, err
		}
		alts := make([]Alternative, len(count))
		for len(count) >= 1 {
			max := maxCount(count)
			for len(max) > 1 {
				maxAlt, err := tiebreak(max)
				if err != nil {
					return nil, err
				}
				delete(count, maxAlt)
				max = append(max[:maxAlt], max[maxAlt+1:]...)
			}
			alts = append(alts, max[0])
			delete(count, max[0])
		}
		return alts, nil
	}
	return f
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tiebreak func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	f := func(profile Profile) (Alternative, error) {
		alts, err := scf(profile)
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
