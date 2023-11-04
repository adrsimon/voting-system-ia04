package methods

import "github.com/adrsimon/voting-system-ia04/comsoc"

func BordaSWF(p comsoc.Profile, _ ...int64) (count comsoc.Count, err error) {
	count = make(comsoc.Count)
	for _, v := range p {
		for j, w := range v {
			count[w] += len(v) - 1 - j
		}
	}
	return count, nil
}

func BordaSCF(p comsoc.Profile, _ ...int64) (bestAlts []comsoc.Alternative, err error) {
	count, err := BordaSWF(p)
	return comsoc.MaxCount(count), err
}
