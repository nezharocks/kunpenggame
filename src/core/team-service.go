package core

import (
	"log"
	"time"
)

// TeamService is
type TeamService struct {
	TeamBattle TeamBattle
	GameAgent  GameAgent
	ErrCh      chan error
	StopCh     chan struct{}
}

// NewTeamService creates a TeamService instance
func NewTeamService(tb TeamBattle, game GameAgent) *TeamService {
	return &TeamService{
		TeamBattle: tb,
		GameAgent:  game,
		ErrCh:      make(chan error, 10),
		StopCh:     make(chan struct{}, 10),
	}
}

// Start is
func (s *TeamService) Start() error {
	err := s.GameAgent.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	go s.handleInternalMessages()

	s.waitInvitationFor(time.Second * 5)

	// register the team
	err = s.GameAgent.Registration(&Registration{s.TeamBattle.GetTeamID(), s.TeamBattle.GetTeamName()})
	if err != nil {
		log.Println(err) // todo
		s.Stop()
	}

	// Wait and process the incoming messsages from game server
	go s.handleExternalMessages()

	<-s.StopCh
	return nil
}

// Stop is
func (s *TeamService) Stop() {
	err := s.TeamBattle.GameOver(&GameOver{})
	if err != nil {
		log.Println(err) // todo
	}
	err = s.GameAgent.Disconnect()
	if err != nil {
		log.Println(err) // todo
	}
	s.StopCh <- struct{}{}
	s.StopCh <- struct{}{}
}

func (s *TeamService) handleInternalMessages() {
loop:
	for {
		select {
		case err := <-s.ErrCh:
			log.Println(err) // todo
		case <-s.StopCh:
			break loop
		}
	}
}

func (s *TeamService) handleExternalMessages() {
loop:
	for {
		select {
		case legStart := <-s.GameAgent.GetLegStartCh():
			s.onLegStart(&legStart)
		case legEnd := <-s.GameAgent.GetLegEndCh():
			s.onLegEnd(&legEnd)
		case round := <-s.GameAgent.GetRoundCh():
			s.onRound(&round)
		case gameOver := <-s.GameAgent.GetGameOverCh():
			s.onGameOver(&gameOver)
		case err := <-s.GameAgent.GetErrorCh():
			s.ErrCh <- err
		case <-time.After(time.Second * ServerTimeout): // server timeout
			s.Stop()
		case <-s.StopCh:
			break loop
		}
	}
}

// wait invitation for some time to set team ID
func (s *TeamService) waitInvitationFor(d time.Duration) {
	if NetworkMode {
		select {
		case inv := <-s.GameAgent.GetInvitationCh():
			s.TeamBattle.SetTeamID(inv.TeamID)
		case <-time.After(d):
			log.Printf("team service - ignore waiting invitation for %v\n", d)
		}
	}
}

func (s *TeamService) onLegStart(legStart *LegStart) {
	if err := s.TeamBattle.LegStart(legStart); err != nil {
		s.ErrCh <- err
	}
}

func (s *TeamService) onLegEnd(legEnd *LegEnd) {
	if err := s.TeamBattle.LegEnd(legEnd); err != nil {
		s.ErrCh <- err
	}
}

func (s *TeamService) onRound(round *Round) {
	err := s.TeamBattle.Round(round)
	if err != nil {
		s.ErrCh <- err
		return
	}
	action := <-s.TeamBattle.GetActionCh()
	err = s.GameAgent.Action(&action)
	if err != nil {
		s.ErrCh <- err
	}
}

func (s *TeamService) onGameOver(gameOver *GameOver) {
	if err := s.TeamBattle.GameOver(gameOver); err != nil {
		s.ErrCh <- err
	}
}
