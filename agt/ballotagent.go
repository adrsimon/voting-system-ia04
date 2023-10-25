package agt

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"log"
	"net/http"
	"sync"
	"time"
)

type ServerRest struct {
	sync.Mutex
	id           string
	addr         string
	ballotAgents []ballotAgent
	count        int64
}

func NewServerRest(addr string) *ServerRest {
	return &ServerRest{id: addr, addr: addr, ballotAgents: make([]ballotAgent, 0), count: 0}
}

type ballotAgent struct {
	ballotID int64
	rule     string
	deadline time.Time
	voterID  []string
	profile  comsoc.Profile
	nbrAlt   int64
	tiebreak []int64
}

func newBallotAgent(ballotID int64, rule string, deadline time.Time, voterID []string, profile comsoc.Profile, nbrAlt int64, tiebreak []int64) *ballotAgent {
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
	vs.count++
	w.WriteHeader(http.StatusOK)
	buf.Reset()
	resp, err := json.Marshal(Response{req.Rule, vs.count})
	err = binary.Write(buf, binary.LittleEndian, resp)
	if err != nil {
		return
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return
	}
}

func (vs *ServerRest) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", vs.newBallot)

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
