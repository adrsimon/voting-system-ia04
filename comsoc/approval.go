package comsoc

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	count = make(Count)
	for i, v := range p {
		for j := 0; j < thresholds[i]; j++ {
			count[v[j]]++
		}
	}
	return count, nil
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	return maxCount(count), err
}
