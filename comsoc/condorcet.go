package comsoc

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	bestAlts = make([]Alternative, 1)

	for i := 0; i < len(p[0]); i++ { // se base sur le fait que le profile soit vérifié
		winner := true
		for j := 0; i < len(p[0]); i++ {
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
