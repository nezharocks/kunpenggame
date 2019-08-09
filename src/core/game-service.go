package core

import (
	"log"
	"net"
	"time"
)

// GameService is
type GameService struct {
	Game     GameStrategy
	Listener net.Listener
}

// NewGameService creates a GameService instance
func NewGameService(game GameStrategy, ln net.Listener) *GameService {
	return &GameService{
		Game:     game,
		Listener: ln,
	}
}

// Serve is
func (g *GameService) Serve() {
	for {
		conn, err := g.Listener.Accept()
		if err != nil {
			// todo logging
			log.Println(err)
			continue
		}

		go func() {
			wire := NewWire(conn, defaultBufferSize)
			teamAgent := NewTeamAgentImpl(wire)

			if NetworkMode {
				inv := &Invitation{g.Game.NewTeamID()}
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
			g.Game.Battle(teamAgent)
		}()
	}
}
