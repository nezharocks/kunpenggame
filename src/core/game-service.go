package core

import (
	"log"
	"net"
	"time"
)

// GameService is
type GameService struct {
	Game     Game
	Listener net.Listener
}

// NewGameService creates a GameService instance
func NewGameService(game Game, ln net.Listener) *GameService {
	return &GameService{
		Game:     game,
		Listener: ln,
	}
}

// Serve is
func (s *GameService) Serve() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Println(err) // todo logging
			continue
		}

		go func() {
			wire := NewWire(conn, defaultBufferSize)
			go wire.Receive()
			teamAgent := NewTeamAgentImpl(wire)

			if NetworkMode {
				inv := &Invitation{s.Game.NewTeamID()}
				err := teamAgent.Invitation(inv)
				if err != nil {
					log.Println(err)
					teamAgent.Disconnect()
					return
				}
			}

		reg_wait:
			for {
				select {
				case <-teamAgent.GetRegCh():
					break reg_wait
				case <-time.After(time.Second * 10):
					break reg_wait
				case err := <-teamAgent.GetErrCh():
					log.Println(err)
				}
			}
			if teamAgent.GetTeamID() == 0 { // team is not registered, so exit
				teamAgent.Disconnect()
				return
			}
			s.Game.Battle(teamAgent)
		}()
	}
}
