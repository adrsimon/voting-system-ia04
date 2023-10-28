package agt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	comsoc "github.com/adrsimon/voting-system-ia04/comsoc"
)

func NewAgent(id AgentID, prefs []comsoc.Alternative, opts map[string][]int64) *Agent {
	return &Agent{id, prefs, opts}
}

func (ag1 Agent) Equal(ag2 Agent) bool {
	if ag1.agentId != ag2.agentId || len(ag1.prefs) != len(ag2.prefs) || len(ag1.options) != len(ag2.options) {
		return false
	}

	for i := range ag1.prefs {
		if ag1.prefs[i] != ag2.prefs[i] {
			return false
		}
	}

	for i := range ag1.options {
		if reflect.DeepEqual(ag1.options[i], ag2.options[i]) {
			return false
		}
	}

	return true
}

func (ag1 Agent) DeepEqual(ag2 Agent) bool {
	return &ag1 == &ag2
}

func (ag Agent) Clone() *Agent {
	return NewAgent(ag.agentId, ag.prefs, ag.options)
}

func (ag Agent) String() string {
	return fmt.Sprintf("ID : %s, Preferences : %v", ag.agentId, ag.prefs)
}

func (ag Agent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	for _, v := range ag.prefs {
		if v == a {
			return true
		}
		if v == b {
			return false
		}
	}
	return false
}

func (ag Agent) TreatResponse(r *http.Response) (NewBallotResponse, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return NewBallotResponse{}, err
	}

	var resp NewBallotResponse

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		fmt.Println("failed unmarshalling")
		return NewBallotResponse{}, err
	}

	return resp, nil
}

func (ag Agent) StartSession(rule string, deadline string, voterIds []AgentID, alts int64, tieBreak []int64) (res int64, err error) {
	port := 8080
	requestURL := fmt.Sprintf("http://localhost:%d/new_ballot", port)

	session := NewBallotRequest{
		Rule:     rule,
		Deadline: deadline,
		VoterIds: voterIds,
		Alts:     alts,
		TieBreak: tieBreak,
	}

	data, _ := json.Marshal(session)

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		fmt.Println(err)
		return
	}

	result, err := ag.TreatResponse(resp)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed treating response")
		return
	}
	fmt.Printf("new session started with id %d\n", result.BallotID)
	return result.BallotID, nil
}

func (ag Agent) Vote(sessionID int64) {
	port := 8080
	requestURL := fmt.Sprintf("http://localhost:%d/vote", port)

	vote := VoteRequest{
		VoterID:  ag.agentId,
		BallotID: sessionID,
		Prefs:    ag.prefs,
		Options:  ag.options,
	}

	data, _ := json.Marshal(vote)

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		fmt.Println(err)
		return
	}

	fmt.Println("vote has been sent successfully")
	return
}

func (ag Agent) GetResults(sessionID int64) {
	port := 8080
	requestURL := fmt.Sprintf("http://localhost:%d/result", port)

	obj := ResultRequest{sessionID}
	data, _ := json.Marshal(obj)

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		fmt.Println(err)
		return
	}
	buf := new(bytes.Buffer)
	_, err2 := buf.ReadFrom(resp.Body)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	var result ResultResponse
	result.Ranking = make([]comsoc.Alternative, 0)

	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		fmt.Println("failed unmarshalling")
		return
	}
	fmt.Printf("the winner of the vote %d is %d ", sessionID, result.Winner)
	if len(result.Ranking) > 0 {
		fmt.Printf("the ranking of the vote is '%v'", result.Ranking)
	}
}
