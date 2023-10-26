package agt

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"log"
	"net/http"
	"slices"
	"time"
)

func NewServerRest(addr string) *ServerRest {
	return &ServerRest{id: addr, addr: addr, ballotAgents: make([]ballotAgent, 0), count: 0}
}

func newBallotAgent(ballotID int64, rule string, deadline time.Time, voterID []AgentID, profile comsoc.Profile, nbrAlt int64, tiebreak []int64) *ballotAgent {
	return &ballotAgent{ballotID: ballotID, rule: rule, deadline: deadline, voterID: voterID, profile: profile, nbrAlt: nbrAlt, tiebreak: tiebreak}
}

func (vs *ServerRest) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (ba *ballotAgent) removeVoter(agID AgentID) {
	for i, v := range ba.voterID {
		if v == agID {
			a := ba.voterID[:i]
			b := ba.voterID[i+1:]
			fmt.Println(a, b)
			ba.voterID = append(ba.voterID[:i], ba.voterID[i+1:]...)
		}
	}
}

func (vs *ServerRest) newBallot(w http.ResponseWriter, r *http.Request) {
	if !vs.checkMethod("POST", w, r) {
		return
	}

	vs.Lock()
	defer vs.Unlock()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return
	}

	req := NewBallotRequest{}
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	//faire verif request good
	end, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Printf("starting a new session based on the %s rule\n", req.Rule)
	vs.ballotAgents = append(vs.ballotAgents, *newBallotAgent(vs.count, req.Rule, end, req.VoterIds, make(comsoc.Profile, 0), req.Alts, req.TieBreak))
	w.WriteHeader(http.StatusOK)
	buf.Reset()
	resp, err := json.Marshal(NewBallotResponse{vs.count})
	vs.count++
	err = binary.Write(buf, binary.LittleEndian, resp)
	if err != nil {
		return
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return
	}
}

func (vs *ServerRest) vote(w http.ResponseWriter, r *http.Request) {
	if !vs.checkMethod("POST", w, r) {
		return
	}

	vs.Lock()
	defer vs.Unlock()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return
	}

	req := VoteRequest{}
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	ba := ballotAgent{}
	for _, b := range vs.ballotAgents {
		if b.ballotID == req.BallotID {
			ba = b
		}
	}

	/*if ba.ballotID == 0 { // pas de ballotID => 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/
	if ba.deadline.Before(time.Now()) { // deadline dépassée => 503
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if !slices.Contains(ba.voterID, req.VoterID) { // pas autorisé à voter => 403
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ba.profile = append(ba.profile, req.Prefs)
	ba.removeVoter(req.VoterID)
	vs.ballotAgents[ba.ballotID] = ba
	fmt.Printf("voter n°%s has voted for vote n°%d \n", req.VoterID, req.BallotID)
	w.WriteHeader(http.StatusOK)
}

func (vs *ServerRest) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", vs.newBallot)
	mux.HandleFunc("/vote", vs.vote)

	s := &http.Server{
		Addr:           vs.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Listening on", vs.addr)
	go log.Fatal(s.ListenAndServe())
}
