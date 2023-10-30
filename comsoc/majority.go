package comsoc

func MajoritySWF(p Profile) (count Count, err error) {
	count = make(Count)
	for _, v := range p {
		count[v[0]]++
	}
	return count, nil
}

func MajoritySCF(p Profile, _ ...int64) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)
	return maxCount(count), err
}
