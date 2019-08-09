package core

import "log"

// GameImpl is
type GameImpl struct {
	Name  string
	idSeq int
}

// NewGameImpl creates a GameImpl instance
func NewGameImpl(name string) *GameImpl {
	return &GameImpl{
		Name: name,
	}
}

// NewTeamID is
func (c *GameImpl) NewTeamID() int {
	c.idSeq++
	return c.idSeq
}

// Battle is
func (c *GameImpl) Battle(team TeamAgent) {
	log.Printf("AI battle with %v/%v is starting...", team.GetTeamID(), team.GetTeamName())

	log.Println("Battle is ended")
}
