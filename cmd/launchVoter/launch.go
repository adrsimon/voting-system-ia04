package main

import (
	"github.com/adrsimon/voting-system-ia04/agt"
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"time"
)

func main() {
	ag := agt.NewAgent("1", []comsoc.Alternative{1, 2, 3}, make(map[string][]int64))
	ag2 := agt.NewAgent("2", []comsoc.Alternative{2, 1, 3}, make(map[string][]int64))
	ag3 := agt.NewAgent("3", []comsoc.Alternative{1, 3, 2}, make(map[string][]int64))
	voter := make([]agt.AgentID, 3)
	voter[0] = "1"
	voter[1] = "2"
	voter[2] = "3"
	tb := make([]int64, 3)
	tb[0] = 1
	tb[1] = 2
	tb[2] = 3

	ballotID, err := ag.StartSession("Majority", "2023-10-28T12:04:00+02:00", voter, 3, tb)
	if err != nil {
		return
	}

	go ag.Vote(ballotID)  // success
	go ag2.Vote(ballotID) // failure
	go ag3.Vote(ballotID)
	time.Sleep(1 * time.Minute)
	ag.GetResults(ballotID)
	// TODO : cant vote two times
}
