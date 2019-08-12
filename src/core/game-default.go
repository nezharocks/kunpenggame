package core

import "log"

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
func (g *GameImpl) Battle(team1 TeamAgent) {
	log.Printf("the battle of AI vs. %v/%v is starting...", team1.GetTeamID(), team1.GetTeamName())
	team2 := NewTeamImpl("ai_team")
	team2.SetID(g.NewTeamID())
	team2.GameStart()

	// legStart := &legStart{}

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
