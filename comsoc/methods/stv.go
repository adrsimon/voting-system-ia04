package methods

import "github.com/adrsimon/voting-system-ia04/comsoc"

func STVSWF(p comsoc.Profile, tiebreak ...int64) (count comsoc.Count, err error) {
	count = make(comsoc.Count)
	nbVote := len(p[0])
	for i := 1; i <= nbVote; i++ {
		countmp, err := MajoritySWF(p)
		if err != nil {
			return nil, err
		}

		worstAlt := comsoc.MinCount(countmp)
		if len(countmp) == 1 {
			count[worstAlt[0]] = i
		} else {
			if len(worstAlt) > 1 {
				k := 0
				for j, v := range tiebreak {
					for _, v2 := range worstAlt {
						if comsoc.Alternative(v) == v2 {
							if j > k {
								k = j
							}
						}
					}
				}
				p = comsoc.deleteAlternative(p, tiebreak[k])
				count[comsoc.Alternative(tiebreak[k])] = i
			} else {
				p = comsoc.deleteAlternative(p, int64(worstAlt[0]))
				count[worstAlt[0]] = i
			}
		}
	}
	return count, nil
}

func STVSCF(p comsoc.Profile, tiebreak ...int64) (bestAlts []comsoc.Alternative, err error) {
	count, err := STVSWF(p, tiebreak...)
	return comsoc.MaxCount(count), err
}
