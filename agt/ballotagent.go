package agt

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"github.com/adrsimon/voting-system-ia04/comsoc/methods"
	"log"
	"net/http"
	"slices"
	"time"
)

func NewServerRest(addr string) *ServerRest {
	return &ServerRest{id: addr, addr: addr, ballotAgents: make(map[string]*ballotAgent), count: 0}
}

func newBallotAgent(ballotID string, rulename string, rule func(comsoc.Profile, ...int64) ([]comsoc.Alternative, error), deadline time.Time, voterID []AgentID, profile comsoc.Profile, alts []comsoc.Alternative, tiebreak []comsoc.Alternative, thresholds []int64) *ballotAgent {
	return &ballotAgent{ballotID: ballotID, rulename: rulename, rule: rule, deadline: deadline, voterID: voterID, profile: profile, alternatives: alts, tiebreak: tiebreak, thresholds: thresholds}
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

func (ba *ballotAgent) vote(req VoteRequest, c chan int) {
	ba.Lock()
	defer ba.Unlock()
	if ba.deadline.Before(time.Now()) { // missed deadline => 503
		c <- http.StatusServiceUnavailable
		return
	} else if !slices.Contains(ba.voterID, req.VoterID) { // unauthorized => 403
		c <- http.StatusForbidden
		return
	} else if comsoc.CheckProfile(req.Prefs, ba.alternatives) != nil { // bad vote => 400
		c <- http.StatusBadRequest
		return
	} else if ba.rulename == "approval" && (len(req.Options) != 1 || req.Options[0] <= 0 || req.Options[0] > int64(len(ba.alternatives))) { // missing option => 400
		c <- http.StatusBadRequest
		return
	} else {
		ba.profile = append(ba.profile, req.Prefs)
		ba.thresholds = append(ba.thresholds, req.Options...)
		ba.removeVoter(req.VoterID)
		c <- http.StatusOK
	}
}

func (vs *ServerRest) newBallot(w http.ResponseWriter, r *http.Request) {
	if !vs.checkMethod("POST", w, r) {
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := NewBallotRequest{}
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Alts <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if time.Now().After(end) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vs.Lock()
	defer vs.Unlock()
	ballotID := fmt.Sprintf("ballot-%d", vs.count)
	ba := *newBallotAgent(ballotID, req.Rule, nil, end, req.VoterIds, make(comsoc.Profile, 0), make([]comsoc.Alternative, 0), make([]comsoc.Alternative, 0), make([]int64, 0))

	for i := int64(0); i < req.Alts; i++ {
		ba.alternatives = append(ba.alternatives, comsoc.Alternative(i))
	}

	tieB := make([]comsoc.Alternative, 0)
	if req.TieBreak == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		for _, v := range req.TieBreak {
			tieB = append(tieB, v)
		}
	}
	ba.tiebreak = tieB

	err = comsoc.CheckProfile(tieB, ba.alternatives)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Rule {
	case "majority", "borda", "approval", "stv", "copeland":
		ba.rule = comsoc.SWFFactory(SWFMap[req.Rule], comsoc.TieBreakFactory(tieB))
	case "condorcet":
		ba.rule = methods.CondorcetWinner
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("starting a new session based on the %s rule\n", req.Rule)
	vs.ballotAgents[ba.ballotID] = &ba
	vs.count++
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
	ballotIndex := ""
	for i, b := range vs.ballotAgents {
		if b.ballotID == req.BallotID {
			ballotIndex = i
		}
	}

	if ballotIndex == "" { // ballot not existing => 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := make(chan int)
	go vs.ballotAgents[ballotIndex].vote(req, c)
	val := <-c
	if val != http.StatusOK {
		w.WriteHeader(val)
		return
	} else {
		w.WriteHeader(val)
		fmt.Printf("%s has voted for %s, with preferences %v, and options %v \n", req.VoterID, req.BallotID, req.Prefs, req.Options)
	}
}

func (vs *ServerRest) result(w http.ResponseWriter, r *http.Request) {
	if !vs.checkMethod("POST", w, r) {
		return
	}

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
			ba = *b
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

	ranking, err := make([]comsoc.Alternative, 0), error(nil)
	if ba.rulename == "approval" {
		ranking, err = ba.rule(ba.profile, ba.thresholds...)
	} else if ba.rulename == "stv" {
		tiebreak := make([]int64, 0)
		for _, v := range ba.tiebreak {
			tiebreak = append(tiebreak, int64(v))
		}
		ranking, err = ba.rule(ba.profile, tiebreak...)
	} else {
		ranking, err = ba.rule(ba.profile)
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	winner := comsoc.Alternative(-1)
	if len(ranking) > 0 {
		winner = ranking[0]
	}
	buf.Reset()
	resp, err := json.Marshal(ResultResponse{Winner: winner, Ranking: ranking})
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

	methods := make([]string, 0)
	for v := range SWFMap {
		methods = append(methods, v)
	}
	methods = append(methods, "condorcet")

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
