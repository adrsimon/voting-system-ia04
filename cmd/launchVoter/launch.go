package main

import (
	"fmt"
	"github.com/adrsimon/voting-system-ia04/agt"
	"github.com/adrsimon/voting-system-ia04/comsoc"
	"math/rand"
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
	alts := make([]comsoc.Alternative, 5)
	for i := 0; i < 5; i++ {
		alts[i] = comsoc.Alternative(i)
	}

	nbVoters := 15
	agents := make(map[agt.AgentID]agt.Agent, nbVoters)
	for i := 0; i < nbVoters; i++ {
		id := agt.AgentID(fmt.Sprintf("agent-%d", i))
		threshold := make([]int64, 0)
		threshold = append(threshold, 1+int64(rand.Intn(len(alts))))
		agents[id] = *agt.NewAgent(id, randomPreferences(alts), threshold)
		fmt.Println(agents[id])
		fmt.Println(threshold)
	}

	badVoterId := agt.AgentID("agent-11")
	threshold := make([]int64, 0)
	threshold = append(threshold, int64(rand.Intn(len(alts))))
	badPreferences := make([]comsoc.Alternative, 0)
	badPreferences = append(badPreferences, comsoc.Alternative(12))
	agents[badVoterId] = *agt.NewAgent(badVoterId, badPreferences, threshold)

	// variables nécessaires à la création d'un vote
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)
	ids := make([]agt.AgentID, 0, len(agents))
	for k := range agents {
		ids = append(ids, k)
	}

	// on récupère le premier agent qui se chargera de créer les sessions et de récupérer les résultats
	organizer := agents["agent-1"]
	tb := alts
	ballotID, err := organizer.StartSession("majority", deadline, ids, int64(len(alts)), tb)
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
