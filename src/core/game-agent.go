package core

import (
	"fmt"
	"log"
	"net"
	"time"
)

// GameAgentImpl is
type GameAgentImpl struct {
	TeamID      int
	TeamName    string
	ServerIP    string
	ServerPort  int
	Conn        net.Conn
	Connected   bool
	DialRetries int
	BufferSize  int
	Wire        *Wire
	// MsgChs       []chan Message
	// ErrChs       []chan error
	InvitationCh chan Invitation
	LegStartCh   chan LegStart
	LegEndCh     chan LegEnd
	RoundCh      chan Round
	GameOverCh   chan GameOver
	ErrorCh      chan error
	StopCh       chan struct{}
}

// NewGameAgentImpl is
func NewGameAgentImpl(teamID int, teamName string, serverIP string, serverPort int) *GameAgentImpl {
	return &GameAgentImpl{
		TeamID:       teamID,
		TeamName:     teamName,
		ServerIP:     serverIP,
		ServerPort:   serverPort,
		DialRetries:  defaultConnRetries,
		BufferSize:   defaultBufferSize,
		InvitationCh: make(chan Invitation, 1),
		LegStartCh:   make(chan LegStart, 1),
		LegEndCh:     make(chan LegEnd, 1),
		RoundCh:      make(chan Round, 10),
		GameOverCh:   make(chan GameOver, 1),
		ErrorCh:      make(chan error, 10),
		StopCh:       make(chan struct{}, 1),
	}
}

// Connect is
func (a *GameAgentImpl) Connect() (err error) {
	address := fmt.Sprintf("%s:%d", a.ServerIP, a.ServerPort)
	teamDesc := fmt.Sprintf("%v:%v", a.TeamID, a.TeamName)
	retries := a.DialRetries
	if retries <= 0 {
		retries = defaultConnRetries
	}

	log.Printf("team (%v) client is connecting to game server@%v", teamDesc, address)
	for i := 1; i <= retries; i++ {
		a.Conn, err = net.DialTimeout("tcp4", address, time.Second*1)
		if err != nil {
			log.Printf("client dial error - try %vth time to dial %v, error: %v\n", i, address, err)
		} else {
			log.Printf("team (%v) client is connected to game server@%v", teamDesc, address)
			a.Wire = NewWire(a.Conn, a.BufferSize)
			a.Connected = true
			go a.receive()
			return nil
		}
	}
	return fmt.Errorf("team (%v) client is connected to game server@%v, error: %v", teamDesc, address, err)
}

// Disconnect is
func (a *GameAgentImpl) Disconnect() (err error) {
	address := fmt.Sprintf("%s:%d", a.ServerIP, a.ServerPort)
	teamDesc := fmt.Sprintf("%v:%v", a.TeamID, a.TeamName)
	if !a.Connected {
		log.Printf("team (%v) client is not connected to game server@%v, no need to disconnect", teamDesc, address)
		return nil
	}

	if err = a.Conn.Close(); err != nil {
		log.Printf("team (%v) client fails to disconnect game server@%v", teamDesc, address)
		return err
	}
	a.Connected = false
	a.StopCh <- struct{}{}
	log.Printf("team (%v) client is disconnected to game server@%v", teamDesc, address)
	return nil
}

// OnMessage binds message handler by adding Message channel
// func (a *GameAgentImpl) OnMessage(msgCh chan Message) {
// 	a.MsgChs = append(a.MsgChs, msgCh)
// }

// OffMessage unbinds message handler by removing Message channel
// func (a *GameAgentImpl) OffMessage(msgCh chan Message) {
// 	for i, ch := range a.MsgChs {
// 		if ch == msgCh {
// 			a.MsgChs = append(a.MsgChs[:i], a.MsgChs[i+1:]...)
// 			break
// 		}
// 	}
// }

// OnError binds error handler by adding error channel
// func (a *GameAgentImpl) OnError(errCh chan error) {
// 	a.ErrChs = append(a.ErrChs, errCh)
// }

// OffError unbinds error handler by removing error channel
// func (a *GameAgentImpl) OffError(errCh chan error) {
// 	for i, ch := range a.ErrChs {
// 		if ch == errCh {
// 			a.ErrChs = append(a.ErrChs[:i], a.ErrChs[i+1:]...)
// 			break
// 		}
// 	}
// }

// Registration is
func (a *GameAgentImpl) Registration(registration *Registration) error {
	return a.Wire.Send(registration.Message())
}

// Action is
func (a *GameAgentImpl) Action(action *Action) error {
	return a.Wire.Send(action.Message())
}

// func (a *GameAgentImpl) emitMessage(msg Message) {
// 	for _, handler := range a.MsgChs {
// 		go func(handler chan Message) {
// 			handler <- msg
// 		}(handler)
// 	}
// }

// func (a *GameAgentImpl) emitError(err error) {
// 	for _, handler := range a.ErrChs {
// 		go func(handler chan error) {
// 			handler <- err
// 		}(handler)
// 	}
// }

func (a *GameAgentImpl) receive() {
loop:
	for {
		select {
		case msg := <-a.Wire.MsgCh:
			switch msg.Name {
			case InvitationName:
				inv, err := msg.Invitation()
				if err != nil {
					a.ErrorCh <- err
					continue loop
				}
				go func() {
					a.InvitationCh <- *inv
				}()
			case LegStartName:
				start, err := msg.LegStart()
				if err != nil {
					a.ErrorCh <- err
					continue loop
				}
				a.LegStartCh <- *start
			case LegEndName:
				end, err := msg.LegEnd()
				if err != nil {
					a.ErrorCh <- err
					continue loop
				}
				a.LegEndCh <- *end
			case RoundName:
				round, err := msg.Round()
				if err != nil {
					a.ErrorCh <- err
					continue loop
				}
				go func() {
					a.RoundCh <- *round
				}()
			case GameOverName:
				gameOver, err := msg.GameOver()
				if err != nil {
					a.ErrorCh <- err
					continue loop
				}
				a.GameOverCh <- *gameOver
			}
		case err := <-a.Wire.ErrCh:
			a.ErrorCh <- err
		case <-a.StopCh:
			break loop
		}
	}
}

// GetInvitationCh is
func (a *GameAgentImpl) GetInvitationCh() chan Invitation {
	return a.InvitationCh
}

// GetLegStartCh is
func (a *GameAgentImpl) GetLegStartCh() chan LegStart {
	return a.LegStartCh
}

// GetLegEndCh is
func (a *GameAgentImpl) GetLegEndCh() chan LegEnd {
	return a.LegEndCh
}

// GetRoundCh is
func (a *GameAgentImpl) GetRoundCh() chan Round {
	return a.RoundCh
}

// GetGameOverCh is
func (a *GameAgentImpl) GetGameOverCh() chan GameOver {
	return a.GameOverCh
}

// GetGameOverCh is
func (a *GameAgentImpl) GetErrorCh() chan error {
	return a.ErrorCh
}
