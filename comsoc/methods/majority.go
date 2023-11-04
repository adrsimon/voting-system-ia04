package methods

import "github.com/adrsimon/voting-system-ia04/comsoc"

func MajoritySWF(p comsoc.Profile, _ ...int64) (count comsoc.Count, err error) {
	count = make(comsoc.Count)
	for _, v := range p[0] {
		count[v] = 0
	}
	for _, v := range p {
		count[v[0]]++
	}
	return count, nil
}

func MajoritySCF(p comsoc.Profile, _ ...int64) (bestAlts []comsoc.Alternative, err error) {
	count, err := MajoritySWF(p)
	return comsoc.MaxCount(count), err
}
