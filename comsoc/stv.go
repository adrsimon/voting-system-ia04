package comsoc

func STVSWF(p Profile) (count Count, err error) {
	count = make(Count)
	for _, v := range p {
		count[v[0]]++
	}
	return count, nil
}

func STVSCF(p Profile, _ ...int64) (bestAlts []Alternative, err error) {
	for {
		count, err := STVSWF(p)
		if err != nil {
			return nil, err
		}

		if len(count) == 1 {
			return maxCount(count), nil
		}

		worstAlts := minCount(count)
		if len(worstAlts) == len(count) || len(worstAlts) == 0 {
			return maxCount(count), nil
		}

		for _, worstAlt := range worstAlts {
			for i, voter := range p {
				for j, alt := range voter {
					if alt == worstAlt {
						p[i] = append(voter[:j], voter[j+1:]...)
					}
				}
			}
			delete(count, worstAlt)
		}
	}
}
