package core

import (
	"log"
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

	/*
	 * bind center-message handler and error handler
	 */
	t.GameAgent.OnMessage(t.MsgCh)
	t.GameAgent.OnError(t.ErrCh)

	/*
	 * register team
	 */
	reg := t.Team.GetRegistration()
	err := t.GameAgent.Register(reg)
	if err != nil {
		log.Println(err) // todo
		return err
	}

	/*
	 * handle the coming messages
	 */
	go t.Handle()

	return nil
}

// Stop is
func (t *TeamService) Stop() {
	t.StopCh <- struct{}{}
	t.GameAgent.OffMessage(t.MsgCh)
	t.GameAgent.OffError(t.ErrCh)
}

// Handle is
func (t *TeamService) Handle() {
loop:
	for {
		select {
		case msg := <-t.MsgCh:
			switch msg.Name {
			case LegStartName:
				legStart, err := msg.LegStart()
				if err != nil {
					t.ErrCh <- err
					continue loop
				}
				if err := t.Team.LegStart(legStart); err != nil {
					t.ErrCh <- err
				}
			case LegEndName:
				legEnd, err := msg.LegEnd()
				if err != nil {
					t.ErrCh <- err
					continue loop
				}
				if err := t.Team.LegEnd(legEnd); err != nil {
					t.ErrCh <- err
				}
			case RoundName:
				round, err := msg.Round()
				if err != nil {
					t.ErrCh <- err
					continue loop
				}
				action, err := t.Team.Round(round)
				if err != nil {
					t.ErrCh <- err
					continue loop
				}
				err = t.GameAgent.Act(action)
				if err != nil {
					t.ErrCh <- err
				}
			case GameOverName:
				gameOver, err := msg.GameOver()
				if err != nil {
					t.ErrCh <- err
				}
				err = t.Team.GameOver(gameOver)
				if err != nil {
					t.ErrCh <- err
				}
				go t.Stop()
			default:
				log.Printf("message error - unknown message %q with msg data:\n%v", msg.Name, msg.String())
			}
		case err := <-t.ErrCh:
			log.Println(err) // todo
		case <-t.StopCh:
			break loop
		}
	}
	// select {
	// case legStart := <-t.GameAgent.LegStartCh:
	// 	if err := t.Team.LegStart(&legStart); err != nil {
	// 		t.GameAgent.ErrCh <- err
	// 	}
	// case legEnd := <-t.GameAgent.LegEndCh:
	// 	if err := t.Team.LegEnd(&legEnd); err != nil {
	// 		t.GameAgent.ErrCh <- err
	// 	}
	// case round := <-t.GameAgent.RoundCh:
	// 	action, err := t.Team.Round(&round)
	// 	if err != nil {
	// 		t.GameAgent.ErrCh <- err
	// 		break
	// 	}
	// 	if err := t.GameAgent.Act(action); err != nil {
	// 		t.GameAgent.ErrCh <- err
	// 	}
	// case <-t.GameAgent.GameOverCh:
	// 	if err := t.Team.GameOver(); err != nil {
	// 		t.GameAgent.ErrCh <- err
	// 	}
	// 	break
	// case err := <-t.GameAgent.ErrCh:
	// 	fmt.Printf("ERROR: %v", err) // todo
	// }
	// }
}
