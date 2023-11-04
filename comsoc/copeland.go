package comsoc

func CopelandSWF(p Profile, _ ...int64) (count Count, err error) {
	count = make(Count)
	for _, v1 := range p[0] {
		count[v1] = 0
		for _, v2 := range p[0] {
			cmp := 0
			if v1 != v2 {
				for _, v := range p {
					if IsPref(v1, v2, v) {
						cmp++
					} else {
						cmp--
					}
				}
			}
			if cmp > 0 {
				count[v1]++
			} else if cmp < 0 {
				count[v1]--
			}
		}
	}
	return count, nil
}

func CopelandSCF(p Profile, _ ...int64) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)
	return MaxCount(count), err
}
