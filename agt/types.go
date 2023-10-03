package agt

import comsoc "github.com/adrsimon/voting-system-ia04/comsoc"

type AgentID string

type Agent struct {
	ID    AgentID
	Name  string
	Prefs []comsoc.Alternative
	c     chan []comsoc.Alternative
}

type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a comsoc.Alternative, b comsoc.Alternative) bool
	Start()
}
