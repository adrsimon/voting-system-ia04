package methods

import "github.com/adrsimon/voting-system-ia04/comsoc"

func CopelandSWF(p comsoc.Profile, _ ...int64) (count comsoc.Count, err error) {
	count = make(comsoc.Count)
	for _, v1 := range p[0] {
		count[v1] = 0
		for _, v2 := range p[0] {
			cmp := 0
			if v1 != v2 {
				for _, v := range p {
					if comsoc.IsPref(v1, v2, v) {
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

func CopelandSCF(p comsoc.Profile, _ ...int64) (bestAlts []comsoc.Alternative, err error) {
	count, err := CopelandSWF(p)
	return comsoc.MaxCount(count), err
}
