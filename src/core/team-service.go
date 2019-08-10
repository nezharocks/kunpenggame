package core

import (
	"log"
	"time"
)

// TeamService is
type TeamService struct {
	Team      TeamStrategy
	GameAgent GameAgent
	ErrCh     chan error
	StopCh    chan struct{}
}

// NewTeamService creates a TeamService instance
func NewTeamService(team TeamStrategy, game GameAgent) *TeamService {
	return &TeamService{
		Team:      team,
		GameAgent: game,
		ErrCh:     make(chan error, 10),
		StopCh:    make(chan struct{}, 10),
	}
}

// Start is
func (t *TeamService) Start() error {
	err := t.GameAgent.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	go t.handleInternalMessages()

	t.waitInvitationFor(time.Second * 5)

	// register the team
	err = t.GameAgent.Registration(&Registration{t.Team.GetID(), t.Team.GetName()})
	if err != nil {
		log.Println(err) // todo
		t.Stop()
	}

	// start a game battle
	go t.Team.GameStart()

	// Wait and process the incoming messsages from game server
	go t.handleExternalMessages()

	<-t.StopCh
	return nil
}

// Stop is
func (t *TeamService) Stop() {
	err := t.Team.GameOver(&GameOver{})
	if err != nil {
		log.Println(err) // todo
	}
	err = t.GameAgent.Disconnect()
	if err != nil {
		log.Println(err) // todo
	}
	t.StopCh <- struct{}{}
	t.StopCh <- struct{}{}
}

func (t *TeamService) handleInternalMessages() {
loop:
	for {
		select {
		case err := <-t.ErrCh:
			log.Println(err) // todo
		case <-t.StopCh:
			break loop
		}
	}
}

func (t *TeamService) handleExternalMessages() {
loop:
	for {
		select {
		case legStart := <-t.GameAgent.GetLegStartCh():
			go t.onLegStart(&legStart)
		case legEnd := <-t.GameAgent.GetLegEndCh():
			go t.onLegEnd(&legEnd)
		case round := <-t.GameAgent.GetRoundCh():
			go t.onRound(&round)
		case gameOver := <-t.GameAgent.GetGameOverCh():
			go t.onGameOver(&gameOver)
		case err := <-t.GameAgent.GetErrorCh():
			t.ErrCh <- err
		case <-time.After(time.Second * 10): // server timeout
			t.Stop()
		case <-t.StopCh:
			break loop
		}
	}
}

// wait invitation for some time to set team ID
func (t *TeamService) waitInvitationFor(d time.Duration) {
	if NetworkMode {
		select {
		case inv := <-t.GameAgent.GetInvitationCh():
			t.Team.SetID(inv.TeamID)
		case <-time.After(d):
			log.Printf("team service - ignore waiting invitation for %v\n", d)
		}
	}
}

func (t *TeamService) onLegStart(legStart *LegStart) {
	if err := t.Team.LegStart(legStart); err != nil {
		t.ErrCh <- err
	}
}

func (t *TeamService) onLegEnd(legEnd *LegEnd) {
	if err := t.Team.LegEnd(legEnd); err != nil {
		t.ErrCh <- err
	}
}

func (t *TeamService) onRound(round *Round) {
	action, err := t.Team.Round(round)
	if err != nil {
		t.ErrCh <- err
		return
	}
	err = t.GameAgent.Action(action)
	if err != nil {
		t.ErrCh <- err
	}
}

func (t *TeamService) onGameOver(gameOver *GameOver) {
	if err := t.Team.GameOver(gameOver); err != nil {
		t.ErrCh <- err
	}
}
