package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// CenterAgentImpl is
type CenterAgentImpl struct {
	TeamID      int
	TeamName    string
	ServerIP    string
	ServerPort  int
	Conn        net.Conn
	Connected   bool
	DialRetries int
	BufferSize  int
	Wire        *Wire
	MsgChs      []chan Message
	ErrChs      []chan error
	StopCh      chan struct{}
}

// NewCenterAgentImpl is
func NewCenterAgentImpl(teamID int, teamName string, serverIP string, serverPort int) *CenterAgentImpl {
	return &CenterAgentImpl{
		TeamID:      teamID,
		TeamName:    teamName,
		ServerIP:    serverIP,
		ServerPort:  serverPort,
		DialRetries: defaultConnRetries,
		BufferSize:  defaultBufferSize,
		StopCh:      make(chan struct{}, 1),
	}
}

// Connect is
func (s *CenterAgentImpl) Connect() (err error) {
	address := fmt.Sprintf("%s:%d", s.ServerIP, s.ServerPort)
	teamDesc := fmt.Sprintf("%v:%v", s.TeamID, s.TeamName)
	retries := s.DialRetries
	if retries <= 0 {
		retries = defaultConnRetries
	}

	log.Printf("team (%v) client is connecting to game server@%v", teamDesc, address)
	for i := 1; i <= retries; i++ {
		s.Conn, err = net.DialTimeout("tcp4", address, time.Second*1)
		if err != nil {
			log.Printf("client dial error - try %vth time to dial %v, error: %v\n", i, address, err)
		} else {
			log.Printf("team (%v) client is connected to game server@%v", teamDesc, address)
			s.Wire = NewWire(s.Conn, s.BufferSize)
			s.Connected = true
			go s.receive()
			return nil
		}
	}
	errMsg := fmt.Sprintf("team (%v) client is connected to game server@%v, error: %v", teamDesc, address, err)
	log.Println(errMsg)
	return errors.New(errMsg)
}

// Disconnect is
func (s *CenterAgentImpl) Disconnect() (err error) {
	address := fmt.Sprintf("%s:%d", s.ServerIP, s.ServerPort)
	teamDesc := fmt.Sprintf("%v:%v", s.TeamID, s.TeamName)
	if !s.Connected {
		log.Printf("team (%v) client is not connected to game server@%v, no need to disconnect", teamDesc, address)
		return nil
	}

	if err = s.Conn.Close(); err != nil {
		log.Printf("team (%v) client fails to disconnect game server@%v", teamDesc, address)
		return err
	}
	s.Connected = false
	s.StopCh <- struct{}{}
	log.Printf("team (%v) client is disconnected to game server@%v", teamDesc, address)
	return nil
}

// OnMessage binds message handler by adding Message channel
func (s *CenterAgentImpl) OnMessage(msgCh chan Message) {
	s.MsgChs = append(s.MsgChs, msgCh)
}

// OffMessage unbinds message handler by removing Message channel
func (s *CenterAgentImpl) OffMessage(msgCh chan Message) {
	for i, ch := range s.MsgChs {
		if ch == msgCh {
			s.MsgChs = append(s.MsgChs[:i], s.MsgChs[i+1:]...)
			break
		}
	}
}

// OnError binds error handler by adding error channel
func (s *CenterAgentImpl) OnError(errCh chan error) {
	s.ErrChs = append(s.ErrChs, errCh)
}

// OffError unbinds error handler by removing error channel
func (s *CenterAgentImpl) OffError(errCh chan error) {
	for i, ch := range s.ErrChs {
		if ch == errCh {
			s.ErrChs = append(s.ErrChs[:i], s.ErrChs[i+1:]...)
			break
		}
	}
}

// Register is
func (s *CenterAgentImpl) Register(registration *Registration) error {
	return s.Wire.Send(registration.Message())
}

// Act is
func (s *CenterAgentImpl) Act(action *Action) error {
	return s.Wire.Send(action.Message())
}

func (s *CenterAgentImpl) emitMessage(msg Message) {
	for _, handler := range s.MsgChs {
		go func(handler chan Message) {
			handler <- msg
		}(handler)
	}
}

func (s *CenterAgentImpl) emitError(err error) {
	for _, handler := range s.ErrChs {
		go func(handler chan error) {
			handler <- err
		}(handler)
	}
}

/*
	switch msg.Name {
	case "leg_start":
		legStart := new(LegStart)
		err := json.Unmarshal(msg.Data.(json.RawMessage), legStart)
		if err != nil {
			s.ErrCh <- err
		} else {
			s.LegStartCh <- *legStart
		}
	case "round":
		round := new(Round)
		err := json.Unmarshal(msg.Data.(json.RawMessage), round)
		if err != nil {
			s.ErrCh <- err
		} else {
			s.RoundCh <- *round
		}
	case "leg_end":
		legEnd := new(LegEnd)
		err := json.Unmarshal(msg.Data.(json.RawMessage), legEnd)
		if err != nil {
			s.ErrCh <- err
		} else {
			s.LegEndCh <- *legEnd
		}
	case "game_over":
		s.GameOverCh <- struct{}{}
	default:
		s.ErrCh <- fmt.Errorf("message error - unknown message %q with msg data:\n%v", msg.Name, msg)
	}
*/

func (s *CenterAgentImpl) receive() {
loop:
	for {
		select {
		case msg := <-s.Wire.MsgCh:
			s.emitMessage(msg)
		case err := <-s.Wire.ErrCh:
			s.emitError(err)
		case <-s.StopCh:
			break loop
		}
	}
}
