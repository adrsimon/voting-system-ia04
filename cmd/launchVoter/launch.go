package main

import (
	"fmt"
	"github.com/adrsimon/voting-system-ia04/agt"
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"golang.org/x/exp/rand"
	"time"
)

func randomPreferences(alts []comsoc.Alternative) []comsoc.Alternative {
	availableAlternatives := make([]comsoc.Alternative, len(alts))
	copy(availableAlternatives, alts)
	randomizedAlts := make([]comsoc.Alternative, len(alts))
	for i := 0; i < len(alts); i++ {
		randomIndex := rand.Intn(len(availableAlternatives))
		randomizedAlts[i] = availableAlternatives[randomIndex]
		availableAlternatives = append(availableAlternatives[:randomIndex], availableAlternatives[randomIndex+1:]...)
	}
	return randomizedAlts
}

func main() {
	// on crée 5 alternatives
	alts := make([]comsoc.Alternative, 4)
	for i := 0; i < 4; i++ {
		alts[i] = comsoc.Alternative(i)
	}

	nbVoters := 3
	agents := make(map[agt.AgentID]agt.Agent, nbVoters)
	for i := 0; i < nbVoters; i++ {
		id := agt.AgentID(fmt.Sprintf("agent-%d", i))
		threshold := make([]int64, 0)
		agents[id] = *agt.NewAgent(id, randomPreferences(alts), threshold)

	}

	badVoterId := agt.AgentID("agent-11")
	threshold := make([]int64, 0)
	badPreferences := make([]comsoc.Alternative, 0)
	badPreferences = append(badPreferences, comsoc.Alternative(12))
	agents[badVoterId] = *agt.NewAgent(badVoterId, badPreferences, threshold)

	// variables nécessaires à la création d'un vote
	deadline := time.Now().Add(1 * time.Second).Format(time.RFC3339)
	ids := make([]agt.AgentID, 0, len(agents))
	for k := range agents {
		ids = append(ids, k)
	}

	// on récupère le premier agent qui se chargera de créer les sessions et de récupérer les résultats
	organizer := agents["agent-1"]
	tb := alts
	ballotID, err := organizer.StartSession("copeland", deadline, ids, int64(len(alts)), tb)
	if err != nil {
		return
	}

	// on fait voter tous les agents
	for _, ag := range agents {
		ag.Vote(ballotID)
	}

	// on récupère les résultats
	time.Sleep(2 * time.Second)
	organizer.GetResults(ballotID)
}
