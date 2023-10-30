package agt

import (
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"sync"
	"time"
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

type ServerRest struct {
	sync.Mutex
	id           string
	addr         string
	ballotAgents map[string]ballotAgent
	count        int64
}

type NewBallotRequest struct {
	Rule     string    `json:"rule"`
	Deadline string    `json:"deadline"`
	VoterIds []AgentID `json:"voter-ids"`
	Alts     int64     `json:"#alts"`
	TieBreak []int64   `json:"tie-break"`
}

type ballotAgent struct {
	ballotID   string
	rule       func(profile comsoc.Profile) ([]comsoc.Alternative, error)
	ruleApp    func(profile comsoc.Profile, thresholds []int64) ([]comsoc.Alternative, error)
	deadline   time.Time
	voterID    []AgentID
	profile    comsoc.Profile
	nbrAlt     int64
	tiebreak   []int64
	thresholds []int64
}

type VoteRequest struct {
	VoterID  AgentID              `json:"agent-id"`
	BallotID string               `json:"ballot-id"`
	Prefs    []comsoc.Alternative `json:"prefs"`
	Options  []int64              `json:"options,omitempty"`
}

type NewBallotResponse struct {
	BallotID string `json:"ballot-id"`
}

type ResultRequest struct {
	BallotID string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking,omitempty"`
}
