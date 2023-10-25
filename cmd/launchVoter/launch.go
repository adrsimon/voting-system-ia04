package main

import (
	"github.com/adrsimon/voting-system-ia04/agt"
	"github.com/adrsimon/voting-system-ia04/comsoc"
)

func main() {
	ag := agt.NewAgent("1", []comsoc.Alternative{1, 2, 3}, []int64{})

	voter := make([]agt.AgentID, 3)
	voter[0] = "1"
	voter[1] = "2"
	voter[2] = "3"
	tb := make([]int64, 3)
	tb[0] = 1
	tb[1] = 2
	tb[2] = 3

	ballotID, err := ag.StartSession("majority", "2024-10-02T10:00:00-05:00", voter, 3, tb)
	if err != nil {
		return
	}

	ag.Vote(ballotID) // success
	// TODO : cant vote two times
}
