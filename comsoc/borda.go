package comsoc

func BordaSWF(p Profile) (count Count, err error) {
	count = make(Count)
	for _, v := range p {
		for j, w := range v {
			count[w] += len(v) - 1 - j
		}
	}
	return count, nil
}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	return maxCount(count), err
}