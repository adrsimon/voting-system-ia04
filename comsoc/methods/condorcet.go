package methods

import "github.com/adrsimon/voting-system-ia04/comsoc"

func CondorcetWinner(p comsoc.Profile, _ ...int64) (bestAlts []comsoc.Alternative, err error) {
	bestAlts = make([]comsoc.Alternative, 0)

	for _, v := range p[0] { // se base sur le fait que le profile soit vérifié
		winner := true
		for _, v2 := range p[0] {
			if v != v2 {
				check, err := comsoc.IsPrefProfil(v, v2, p)
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
			bestAlts = append(bestAlts, comsoc.Alternative(v))
		}
	}
	return bestAlts, nil
}
