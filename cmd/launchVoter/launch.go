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
	alts := make([]comsoc.Alternative, 5)
	for i := 0; i < 5; i++ {
		alts[i] = comsoc.Alternative(i)
	}

	nbVoters := 10
	agents := make(map[agt.AgentID]agt.Agent, nbVoters)
	for i := 0; i < nbVoters; i++ {
		id := agt.AgentID(fmt.Sprintf("agent-%d", i))
		threshold := make([]int64, 0)
		threshold = append(threshold, int64(rand.Intn(len(alts)+1)))
		agents[id] = *agt.NewAgent(id, randomPreferences(alts), threshold)

	}

	// variables nécessaires à la création d'un vote
	deadline := time.Now().Add(3 * time.Second).Format(time.RFC3339)
	ids := make([]agt.AgentID, 0, len(agents))
	for k := range agents {
		ids = append(ids, k)
	}

	// on récupère le premier agent qui se chargera de créer les sessions et de récupérer les résultats
	organizer := agents["agent-1"]
	tb := alts
	ballotID, err := organizer.StartSession("Borda", deadline, ids, int64(len(alts)), tb)
	if err != nil {
		return
	}

	// on fait voter tous les agents
	for _, ag := range agents {
		ag.Vote(ballotID)
	}

	// on récupère les résultats
	time.Sleep(5 * time.Second)
	organizer.GetResults(ballotID)
}
