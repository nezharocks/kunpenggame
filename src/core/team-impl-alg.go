package core

import (
	"fmt"
	"log"
	"time"
)

// TeamAlg is
type TeamAlg struct {
	ID       int
	Name     string
	RegCh    chan Registration
	ActionCh chan Action
	ErrCh    chan error
	Battle   *Battle
}

// NewTeamAlg creates a TeamAlg instance
func NewTeamAlg(name string) *TeamAlg {
	return &TeamAlg{
		Name:     name,
		RegCh:    make(chan Registration, 1),
		ActionCh: make(chan Action, 1),
		ErrCh:    make(chan error, 10),
	}
}

// GetTeamID is
func (t *TeamAlg) GetTeamID() int {
	return t.ID
}

// GetTeamName is
func (t *TeamAlg) GetTeamName() string {
	return t.Name
}

// SetTeamID is
func (t *TeamAlg) SetTeamID(id int) {
	t.ID = id
}

// SetTeamName is
func (t *TeamAlg) SetTeamName(name string) {
	t.Name = name
}

func (t *TeamAlg) ensureBattle() {
	if t.Battle == nil {
		t.Battle = NewBattle(t.ID, time.Now())
		// todo init the game battle's algorithm impl here.
		log.Printf("team %q's game battle is started", fmt.Sprintf("%v:%v", t.ID, t.Name))
	}
}

// LegStart is
func (t *TeamAlg) LegStart(legStart *LegStart) error {
	t.ensureBattle()
	t.Battle.StartLeg(legStart)
	return nil
}

// LegEnd is
func (t *TeamAlg) LegEnd(legEnd *LegEnd) error {
	return t.Battle.EndLeg(legEnd)
}

// Round is
func (t *TeamAlg) Round(round *Round) error {
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
func (t *TeamAlg) GameOver(gameOver *GameOver) error {
	// todo
	log.Printf("team %q's game battle is over", fmt.Sprintf("%v:%v", t.ID, t.Name))
	return nil
}

// GetRegCh is
func (t *TeamAlg) GetRegCh() chan Registration {
	return t.RegCh
}

// GetActionCh is
func (t *TeamAlg) GetActionCh() chan Action {
	return t.ActionCh
}

// GetErrCh is
func (t *TeamAlg) GetErrCh() chan error {
	return t.ErrCh
}
