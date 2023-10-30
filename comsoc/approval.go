package comsoc

import "fmt"

func ApprovalSWF(p Profile, thresholds []int64) (count Count, err error) {
	count = make(Count)
	for i, v := range p {
		for j := int64(0); j < thresholds[i]; j++ {
			count[v[j]]++
		}
		fmt.Println(count)
	}
	return count, nil
}

func ApprovalSCF(p Profile, thresholds ...int64) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	return maxCount(count), err
}
