package agt

import (
	"fmt"

	comsoc "github.com/adrsimon/voting-system-ia04/comsoc"
)

func NewAgent(id AgentID, name string, prefs []comsoc.Alternative, c chan []comsoc.Alternative) *Agent {
	return &Agent{id, name, prefs, c}
}

func (ag1 Agent) Equal(ag2 Agent) bool {
	if ag1.ID != ag2.ID || ag1.Name != ag2.Name || len(ag1.Prefs) != len(ag2.Prefs) || ag1.c != ag2.c {
		return false
	}

	for i := range ag1.Prefs {
		if ag1.Prefs[i] != ag2.Prefs[i] {
			return false
		}
	}

	return true
}

func (ag1 Agent) DeepEqual(ag2 Agent) bool {
	return &ag1 == &ag2
}

func (ag Agent) Clone() *Agent {
	return NewAgent(ag.ID, ag.Name, ag.Prefs, ag.c)
}

func (ag Agent) String() string {
	return fmt.Sprintf("ID : %s, Name : %s, Preferences : %v", ag.ID, ag.Name, ag.Prefs)
}

func (ag Agent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	for _, v := range ag.Prefs {
		if v == a {
			return true;
		}
		if v == b {
			return false;
		}
	}
	return false;
}

func (ag Agent) Start() {
	ag.c <- ag.Prefs;
}
