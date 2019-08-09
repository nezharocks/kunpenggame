package core

import (
	"fmt"
	"log"
	"time"
)

const defaultBufferSize = 1024 * 10
const defaultConnRetries = 30

// TeamService is
type TeamService struct {
	Team      TeamStrategy
	GameAgent GameAgent
	MsgCh     chan Message
	ErrCh     chan error
	StopCh    chan struct{}
}

// NewTeamService creates a TeamService instance
func NewTeamService(team TeamStrategy, game GameAgent) *TeamService {
	return &TeamService{
		Team:      team,
		GameAgent: game,
		MsgCh:     make(chan Message, 10),
		ErrCh:     make(chan error, 10),
		StopCh:    make(chan struct{}, 1),
	}
}

// Start is
func (t *TeamService) Start() error {
	// t.GameAgent = NewGameAgentImpl(t.Team.ID, t.Team.Name, t.ServerIP, t.ServerPort)
	// err := t.GameAgent.Connect()
	// if err != nil {
	// 	log.Println(err) // todo
	// 	return err
	// }

	// wait for invitation to set team ID
	if NetworkMode {
		select {
		case inv := <-t.GameAgent.GetInvitationCh():
			t.Team.SetID(inv.TeamID)
		case <-time.After(time.Second * 5):
			log.Println("team service - ignore waiting for invitation")
		}
	}

	// register the team
	err := t.GameAgent.Registration(&Registration{t.Team.GetID(), t.Team.GetName()})
	if err != nil {
		log.Println(err) // todo
		err = t.GameAgent.Disconnect()
		if err != nil {
			log.Println(err) // todo
		}
		return err
	}

	// start a game battle
	go t.Team.GameStart()

	// Wait and process the incoming messsages from game server
	go t.Handle()

	return nil
}

// Stop is
func (t *TeamService) Stop() {
	t.StopCh <- struct{}{}
}

// Handle is
func (t *TeamService) Handle() {
loop:
	for {
		select {
		case legStart := <-t.GameAgent.GetLegStartCh():
			go t.Team.LegStart(&legStart)
		case legEnd := <-t.GameAgent.GetLegEndCh():
			go t.Team.LegEnd(&legEnd)
		case round := <-t.GameAgent.GetRoundCh():
			go func() {
				action, err := t.Team.Round(&round)
				if err != nil {
					t.ErrCh <- err
					return
				}
				t.GameAgent.Action(action)
			}()
		case gameOver := <-t.GameAgent.GetGameOverCh():
			go t.Team.GameOver(&gameOver)
			break
		case err := <-t.GameAgent.GetErrorCh():
			fmt.Printf("ERROR: %v", err) // todo
		case <-t.StopCh:
			break loop
		}
	}

	// case msg := <-t.MsgCh:
	// 	switch msg.Name {
	// 	case LegStartName:
	// 		legStart, err := msg.LegStart()
	// 		if err != nil {
	// 			t.ErrCh <- err
	// 			continue loop
	// 		}
	// 		if err := t.Team.LegStart(legStart); err != nil {
	// 			t.ErrCh <- err
	// 		}
	// 	case LegEndName:
	// 		legEnd, err := msg.LegEnd()
	// 		if err != nil {
	// 			t.ErrCh <- err
	// 			continue loop
	// 		}
	// 		if err := t.Team.LegEnd(legEnd); err != nil {
	// 			t.ErrCh <- err
	// 		}
	// 	case RoundName:
	// 		round, err := msg.Round()
	// 		if err != nil {
	// 			t.ErrCh <- err
	// 			continue loop
	// 		}
	// 		action, err := t.Team.Round(round)
	// 		if err != nil {
	// 			t.ErrCh <- err
	// 			continue loop
	// 		}
	// 		err = t.GameAgent.Action(action)
	// 		if err != nil {
	// 			t.ErrCh <- err
	// 		}
	// 	case GameOverName:
	// 		gameOver, err := msg.GameOver()
	// 		if err != nil {
	// 			t.ErrCh <- err
	// 		}
	// 		err = t.Team.GameOver(gameOver)
	// 		if err != nil {
	// 			t.ErrCh <- err
	// 		}
	// 		go t.Stop()
	// 	default:
	// 		log.Printf("message error - unknown message %q with msg data:\n%v", msg.Name, msg.String())
	// 	}
	// case err := <-t.ErrCh:
	// 	log.Println(err) // todo
	// case <-t.StopCh:
	// 	break loop
	// }
	// }

}
