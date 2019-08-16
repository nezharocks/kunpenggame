package core

import (
	"log"
)

// GameImpl is
type GameImpl struct {
	Name  string
	Judge *Judge
}

// NewGameImpl creates a GameImpl instance
func NewGameImpl(name string) *GameImpl {
	judge := NewFirstGame(defaultMapData, defaultVision, defaultWidth, defaultHeight)
	if err := judge.Init(); err != nil {
		log.Println(err)
	}
	return &GameImpl{
		Name:  name,
		Judge: judge,
	}
}

// NewTeamID is
func (g *GameImpl) NewTeamID() int {
	return g.Judge.NewTeamID()
}

// Battle is
func (g *GameImpl) Battle(guest TeamAgent) {
	log.Printf("the battle of AI vs. %v/%v is starting...", guest.GetTeamID(), guest.GetTeamName())
	ai := NewTeamImpl("ai_team")
	ai.SetTeamID(g.Judge.NewTeamID())
	battle := g.Judge.NewBattle(guest, ai)
	battle.Run()
	log.Println("Battle is ended")
}
