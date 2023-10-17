package main

import "github.com/adrsimon/voting-system-ia04/agt"

func main() {
	ag := agt.Agent{}

	voter := make([]string, 3)
	tb := make([]int64, 3)
	tb[0] = 1
	tb[1] = 2
	tb[2] = 3
	voter[0] = "1"
	voter[1] = "2"
	voter[2] = "3"

	_, err := ag.StartSession("majority", "2002-10-02T10:00:00-05:00", voter, 3, tb)
	if err != nil {
		return
	}
}
