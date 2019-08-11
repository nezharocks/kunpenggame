package core

import (
	"fmt"
	"log"
	"time"
)

// TeamImpl is
type TeamImpl struct {
	ID     int
	Name   string
	Battle *Battle
}

// NewTeamImpl creates a TeamImpl instance
func NewTeamImpl(name string) *TeamImpl {
	return &TeamImpl{
		Name: name,
	}
}

// GetID is
func (t *TeamImpl) GetID() int {
	return t.ID
}

// GetName is
func (t *TeamImpl) GetName() string {
	return t.Name
}

// SetID is
func (t *TeamImpl) SetID(id int) {
	t.ID = id
}

// GameStart is
func (t *TeamImpl) GameStart() {
	t.Battle = NewBattle(t.ID, time.Now())

	// todo init the game battle's algorithm impl here.
	log.Printf("team %q's game battle is started", fmt.Sprintf("%v:%v", t.ID, t.Name))
}

// LegStart is
func (t *TeamImpl) LegStart(legStart *LegStart) error {
	t.Battle.NewLeg(legStart)
	return nil
}

// LegEnd is
func (t *TeamImpl) LegEnd(legEnd *LegEnd) error {
	return t.Battle.EndLeg(legEnd)
}

// Round is
func (t *TeamImpl) Round(round *Round) (*Action, error) {
	// todo
	return &Action{}, nil
}

// GameOver is
func (t *TeamImpl) GameOver(gameOver *GameOver) error {
	// todo
	log.Printf("team %q's game battle is over", fmt.Sprintf("%v:%v", t.ID, t.Name))
	return nil
}
