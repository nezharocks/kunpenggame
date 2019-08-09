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
	MsgChs      []chan Message
	ErrChs      []chan error
	StopCh      chan struct{}
}

// NewGameAgentImpl is
func NewGameAgentImpl(teamID int, teamName string, serverIP string, serverPort int) *GameAgentImpl {
	return &GameAgentImpl{
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
func (a *GameAgentImpl) OnMessage(msgCh chan Message) {
	a.MsgChs = append(a.MsgChs, msgCh)
}

// OffMessage unbinds message handler by removing Message channel
func (a *GameAgentImpl) OffMessage(msgCh chan Message) {
	for i, ch := range a.MsgChs {
		if ch == msgCh {
			a.MsgChs = append(a.MsgChs[:i], a.MsgChs[i+1:]...)
			break
		}
	}
}

// OnError binds error handler by adding error channel
func (a *GameAgentImpl) OnError(errCh chan error) {
	a.ErrChs = append(a.ErrChs, errCh)
}

// OffError unbinds error handler by removing error channel
func (a *GameAgentImpl) OffError(errCh chan error) {
	for i, ch := range a.ErrChs {
		if ch == errCh {
			a.ErrChs = append(a.ErrChs[:i], a.ErrChs[i+1:]...)
			break
		}
	}
}

// Register is
func (a *GameAgentImpl) Register(registration *Registration) error {
	return a.Wire.Send(registration.Message())
}

// Act is
func (a *GameAgentImpl) Act(action *Action) error {
	return a.Wire.Send(action.Message())
}

func (a *GameAgentImpl) emitMessage(msg Message) {
	for _, handler := range a.MsgChs {
		go func(handler chan Message) {
			handler <- msg
		}(handler)
	}
}

func (a *GameAgentImpl) emitError(err error) {
	for _, handler := range a.ErrChs {
		go func(handler chan error) {
			handler <- err
		}(handler)
	}
}

func (a *GameAgentImpl) receive() {
loop:
	for {
		select {
		case msg := <-a.Wire.MsgCh:
			a.emitMessage(msg)
		case err := <-a.Wire.ErrCh:
			a.emitError(err)
		case <-a.StopCh:
			break loop
		}
	}
}
