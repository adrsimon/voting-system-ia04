package comsoc

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	bestAlts = make([]Alternative, 0)

	for _, v := range p[0] { // se base sur le fait que le profile soit vérifié
		winner := true
		for _, v2 := range p[0] {
			if v != v2 {
				check, err := IsPrefProfil(v, v2, p)
				if err != nil {
					return bestAlts, err
				}
				if !check {
					winner = false
					break
				}
			}
		}
		if winner {
			bestAlts = append(bestAlts, Alternative(v))
		}
	}
	return bestAlts, nil
}
