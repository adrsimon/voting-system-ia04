package agt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	comsoc "github.com/adrsimon/voting-system-ia04/comsoc"
)

type AgentID string

type Agent struct {
	agentId AgentID
	prefs   []comsoc.Alternative
	options []int64
}

type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a comsoc.Alternative, b comsoc.Alternative) bool
	Start()
}

type NewBallotRequest struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	VoterIds []string `json:"voterIds"`
	Alts     int64    `json:"alts"`
	TieBreak []int64  `json:"tieBreak"`
}

type Response struct {
	Rule       string `json:"rule"`
	NbSessions int64  `json:"nbSessions"`
}

func NewAgent(id AgentID, prefs []comsoc.Alternative, opts []int64) *Agent {
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
		if ag1.options[i] != ag2.options[i] {
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

func (ag Agent) TreatResponse(r *http.Response) (Response, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return Response{}, err
	}

	var resp Response

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		fmt.Println("failed unmarshalling")
		return Response{}, err
	}

	return resp, nil
}

func (ag Agent) StartSession(rule string, deadline string, voterIds []string, alts int64, tieBreak []int64) (res int64, err error) {
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
		return
	}

	result, err := ag.TreatResponse(resp)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed treating response")
		return
	}
	fmt.Printf("new session based on the %s method has been started\n", result.Rule)
	fmt.Printf("there are %d sessions ongoing\n", result.NbSessions)
	return
}

func (ag Agent) Vote() {

}

func (ag Agent) GetResults() {

}
