package comsoc

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	bestAlts = make([]Alternative, 0)

	for i := 1; i <= len(p[0]); i++ { // se base sur le fait que le profile soit vérifié
		winner := true
		for j := 1; j <= len(p[0]); j++ {
			if i != j {
				check, err := isPrefProfil(Alternative(i), Alternative(j), p)
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
			bestAlts = append(bestAlts, Alternative(i))
		}
	}
	return bestAlts, nil
}
