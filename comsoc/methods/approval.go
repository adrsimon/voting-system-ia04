package methods

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
)

func ApprovalSWF(p comsoc.Profile, thresholds ...int64) (count comsoc.Count, err error) {
	count = make(comsoc.Count)
	for i, v := range p {
		for j := int64(0); j < thresholds[i]; j++ {
			count[v[j]]++
		}
	}
	return count, nil
}

func ApprovalSCF(p comsoc.Profile, thresholds ...int64) (bestAlts []comsoc.Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds...)
	return comsoc.MaxCount(count), err
}
