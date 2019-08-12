package core

import (
	"fmt"
	"log"
	"time"
)

// TeamImpl is
type TeamImpl struct {
	ID       int
	Name     string
	RegCh    chan Registration
	ActionCh chan Action
	ErrCh    chan error
	Battle   *Battle
}

// NewTeamImpl creates a TeamImpl instance
func NewTeamImpl(name string) *TeamImpl {
	return &TeamImpl{
		Name:     name,
		RegCh:    make(chan Registration, 1),
		ActionCh: make(chan Action, 1),
		ErrCh:    make(chan error, 10),
	}
}

// GetTeamID is
func (t *TeamImpl) GetTeamID() int {
	return t.ID
}

// GetTeamName is
func (t *TeamImpl) GetTeamName() string {
	return t.Name
}

// SetTeamID is
func (t *TeamImpl) SetTeamID(id int) {
	t.ID = id
}

// SetTeamName is
func (t *TeamImpl) SetTeamName(name string) {
	t.Name = name
}

func (t *TeamImpl) ensureBattle() {
	if t.Battle == nil {
		t.Battle = NewBattle(t.ID, time.Now())
		// todo init the game battle's algorithm impl here.
		log.Printf("team %q's game battle is started", fmt.Sprintf("%v:%v", t.ID, t.Name))
	}
}

// LegStart is
func (t *TeamImpl) LegStart(legStart *LegStart) error {
	t.ensureBattle()
	t.Battle.StartLeg(legStart)
	return nil
}

// LegEnd is
func (t *TeamImpl) LegEnd(legEnd *LegEnd) error {
	return t.Battle.EndLeg(legEnd)
}

// Round is
func (t *TeamImpl) Round(round *Round) error {
	err := t.Battle.AddRound(round)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		action, err := t.Battle.CalcAction()
		if err != nil {
			log.Println(err)
		}
		t.ActionCh <- *action
	}()

	return nil
}

// GameOver is
func (t *TeamImpl) GameOver(gameOver *GameOver) error {
	// todo
	log.Printf("team %q's game battle is over", fmt.Sprintf("%v:%v", t.ID, t.Name))
	return nil
}

// GetRegCh is
func (t *TeamImpl) GetRegCh() chan Registration {
	return t.RegCh
}

// GetActionCh is
func (t *TeamImpl) GetActionCh() chan Action {
	return t.ActionCh
}

// GetErrCh is
func (t *TeamImpl) GetErrCh() chan error {
	return t.ErrCh
}
