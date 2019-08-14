package core

import (
	"fmt"
	"log"
	"time"
)

// GameImpl is
type GameImpl struct {
	Name        string
	teamIDSeq   int
	playerIDSeq int
}

// NewGameImpl creates a GameImpl instance
func NewGameImpl(name string) *GameImpl {
	return &GameImpl{
		Name: name,
	}
}

// NewTeamID is
func (g *GameImpl) NewTeamID() int {
	g.teamIDSeq++
	return g.teamIDSeq
}

// NewPlayerID is
func (g *GameImpl) NewPlayerID() int {
	g.playerIDSeq++
	return g.playerIDSeq
}

// Battle is
func (g *GameImpl) Battle(t1 TeamAgent) {
	log.Printf("the battle of AI vs. %v/%v is starting...", t1.GetTeamID(), t1.GetTeamName())
	t2 := NewTeamImpl("ai_team")
	t2.SetTeamID(g.NewTeamID())

	const (
		vision    = 3
		playerNum = 4
		roundNum  = 300
	)
	m, err := NewMapFromString(map1)
	if err != nil {
		log.Println(err)
		return
	}
	m.Vision = vision
	tt1 := g.newLegTeam(t1.GetTeamID(), playerNum, Beat)
	tt2 := g.newLegTeam(t2.GetTeamID(), playerNum, Think)
	legStart := &LegStart{
		Map:   m,
		Teams: []*Team{tt1, tt2},
	}
	err = t1.LegStart(legStart)
	if err != nil {
		log.Println(err)
		return
	}

	err = t2.LegStart(legStart)
	if err != nil {
		log.Println(err)
		return
	}

	switchRound := roundNum / 2
	mode1 := BeatMode
	// mode2 := mode1.Reverse()
loop:
	for i := 0; i < roundNum; i++ {
		if i == switchRound {
			mode1 = mode1.Reverse()
			// mode2 = mode1.Reverse()
		}
		r1 := &Round{
			ID:   i,
			Mode: string(mode1),
		}
		t1.Round(r1)
		select {
		case action1 := <-t1.GetActionCh():
			fmt.Println(action1)
		case <-time.After(time.Millisecond * 800):
			break loop
		}

	}

	log.Println("Battle is ended")
}

func (g *GameImpl) newLegTeam(teamID, n int, force string) *Team {
	team := &Team{
		ID:      teamID,
		Players: make([]int, n, n),
		Force:   force,
	}

	for i := 0; i < n; i++ {
		team.Players[i] = g.NewPlayerID()

	}
	return team
}
