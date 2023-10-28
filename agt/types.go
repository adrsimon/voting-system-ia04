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
	options map[string][]int64
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
	ballotAgents []ballotAgent
	count        int64
}

type NewBallotRequest struct {
	Rule     string    `json:"rule"`
	Deadline string    `json:"deadline"`
	VoterIds []AgentID `json:"voterIds"`
	Alts     int64     `json:"alts"`
	TieBreak []int64   `json:"tieBreak"`
}

type ballotAgent struct {
	ballotID   int64
	rule       func(profile comsoc.Profile) ([]comsoc.Alternative, error)
	ruleApp    func(profile comsoc.Profile, tresholds []int) ([]comsoc.Alternative, error)
	deadline   time.Time
	voterID    []AgentID
	profile    comsoc.Profile
	nbrAlt     int64
	tiebreak   []int64
	thresholds []int
}

type VoteRequest struct {
	VoterID  AgentID              `json:"agent-id"`
	BallotID int64                `json:"ballot-id"`
	Prefs    []comsoc.Alternative `json:"prefs"`
	Options  map[string][]int64   `json:"options,omitempty"`
}

type NewBallotResponse struct {
	BallotID int64 `json:"ballot-id"`
}

type ResultRequest struct {
	BallotID int64 `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking,omitempty"`
}
