package comsoc

func CopelandSWF(p Profile, _ ...int64) (count Count, err error) {
	count = make(Count)
	for _, v := range p {
		for i := 0; i < len(v); i++ {
			for j := i + 1; j < len(v); j++ {
				if v[i] != v[j] {
					if IsPref(v[i], v[j], v) {
						count[v[i]]++
					} else {
						count[v[j]]++
					}
				}
			}
		}
	}
	return count, nil
}

func CopelandSCF(p Profile, _ ...int64) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)
	return MaxCount(count), err
}
