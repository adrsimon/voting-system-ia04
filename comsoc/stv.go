package comsoc

func STVSWF(p Profile, tiebreak ...int64) (count Count, err error) {
	count = make(Count)
	nbVote := len(p[0])
	for i := 1; i <= nbVote; i++ {
		countmp, err := MajoritySWF(p)
		if err != nil {
			return nil, err
		}

		worstAlt := MinCount(countmp)
		if len(countmp) == 1 {
			count[worstAlt[0]] = i
		} else {
			if len(worstAlt) > 1 {
				k := 0
				for j, v := range tiebreak {
					for _, v2 := range worstAlt {
						if Alternative(v) == v2 {
							if j > k {
								k = j
							}
						}
					}
				}
				p = deleteAlternative(p, tiebreak[k])
				count[Alternative(tiebreak[k])] = i
			} else {
				p = deleteAlternative(p, int64(worstAlt[0]))
				count[worstAlt[0]] = i
			}
		}
	}
	return count, nil
}

func STVSCF(p Profile, tiebreak ...int64) (bestAlts []Alternative, err error) {
	count, err := STVSWF(p, tiebreak...)
	return MaxCount(count), err
}
