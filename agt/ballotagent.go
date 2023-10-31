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
	return &ServerRest{id: addr, addr: addr, ballotAgents: make(map[string]ballotAgent), count: 0}
}

func newBallotAgent(ballotID string, rule func(comsoc.Profile, ...int64) (comsoc.Alternative, error), deadline time.Time, voterID []AgentID, profile comsoc.Profile, nbrAlt int64, tiebreak []comsoc.Alternative, thresholds []int64) *ballotAgent {
	return &ballotAgent{ballotID: ballotID, rule: rule, deadline: deadline, voterID: voterID, profile: profile, nbrAlt: nbrAlt, tiebreak: tiebreak, thresholds: thresholds}
}

func (vs *ServerRest) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := fmt.Fprintf(w, "method %q not allowed", r.Method)
		if err != nil {
			return false
		}
		return false
	}
	return true
}

func (ba *ballotAgent) removeVoter(agID AgentID) {
	for i, v := range ba.voterID {
		if v == agID {
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

	tieB := make([]comsoc.Alternative, 0)
	if req.TieBreak == nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		for _, v := range req.TieBreak {
			tieB = append(tieB, v)
		}
	}

	ballotID := fmt.Sprintf("ballot-%d", vs.count)
	ba := *newBallotAgent(ballotID, nil, end, req.VoterIds, make(comsoc.Profile, 0), req.Alts, req.TieBreak, make([]int64, 0))

	switch req.Rule {
	case "majority", "borda", "approval", "stv", "copeland":
		ba.rule = comsoc.SCFFactory(SCFMap[req.Rule], comsoc.TieBreakFactory(tieB))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Printf("starting a new session based on the %s rule\n", req.Rule)
	vs.ballotAgents[ba.ballotID] = ba
	w.WriteHeader(http.StatusOK)
	buf.Reset()
	resp, err := json.Marshal(NewBallotResponse{ballotID})
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
	ba.ballotID = ""
	for _, b := range vs.ballotAgents {
		if b.ballotID == req.BallotID {
			ba = b
		}
	}

	if ba.ballotID == "" { // pas de ballotID => 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ba.deadline.Before(time.Now()) { // deadline dépassée => 503
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if !slices.Contains(ba.voterID, req.VoterID) { // pas autorisé à voter => 403
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ba.profile = append(ba.profile, req.Prefs)
	ba.thresholds = append(ba.thresholds, req.Options...)
	ba.removeVoter(req.VoterID)
	vs.ballotAgents[ba.ballotID] = ba
	fmt.Printf("%s has voted for %s, with preferences %v, and approval threshold %v \n", req.VoterID, req.BallotID, req.Prefs, req.Options)
	w.WriteHeader(http.StatusOK)
}

func (vs *ServerRest) result(w http.ResponseWriter, r *http.Request) {
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

	req := ResultRequest{}
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	ba := ballotAgent{}
	ba.ballotID = ""
	for _, b := range vs.ballotAgents {
		if b.ballotID == req.BallotID {
			ba = b
		}
	}

	if ba.ballotID == "" { // pas de ballotID => 404
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if ba.deadline.After(time.Now()) { // résultats pas disponibles => 425
		w.WriteHeader(http.StatusTooEarly)
		return
	}

	winner, err := ba.rule(ba.profile, ba.thresholds...)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	buf.Reset()
	resp, err := json.Marshal(ResultResponse{Winner: winner, Ranking: nil})
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

func (vs *ServerRest) methods(w http.ResponseWriter, r *http.Request) {
	if !vs.checkMethod("GET", w, r) {
		return
	}

	vs.Lock()
	defer vs.Unlock()

	methods := make([]string, 0)
	for v := range SCFMap {
		methods = append(methods, v)
	}

	buf := new(bytes.Buffer)
	resp, err := json.Marshal(MethodsResponse{methods})
	if err != nil {
		return
	}
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
	mux.HandleFunc("/vote", vs.vote)
	mux.HandleFunc("/result", vs.result)
	mux.HandleFunc("/methods", vs.methods)

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
