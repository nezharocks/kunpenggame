package core

import "log"

// TeamAgentImpl is
type TeamAgentImpl struct {
	ID       int
	Name     string
	Wire     *Wire
	RegCh    chan Registration
	ActionCh chan Action
	ErrCh    chan error
	StopCh   chan struct{}
}

// NewTeamAgentImpl is
func NewTeamAgentImpl(wire *Wire) *TeamAgentImpl {
	teamAgent := &TeamAgentImpl{
		Wire:     wire,
		RegCh:    make(chan Registration, 1),
		ActionCh: make(chan Action, 1),
		ErrCh:    make(chan error, 10),
		StopCh:   make(chan struct{}, 1),
	}
	go teamAgent.receive()
	return teamAgent
}

func (a *TeamAgentImpl) receive() {
loop:
	for {
		select {
		case msg := <-a.Wire.MsgCh:
			switch msg.Name {
			case RegistrationName:
				reg, err := msg.Registration()
				if err != nil {
					a.ErrCh <- err
					continue loop
				}
				a.ID = reg.TeamID
				a.Name = reg.TeamName
				go func() {
					a.RegCh <- *reg
				}()
			case ActionName:
				action, err := msg.Action()
				if err != nil {
					a.ErrCh <- err
					continue loop
				}
				go func() {
					a.ActionCh <- *action
				}()
			default:
				log.Printf("message error - unknown message %q with msg data:\n%v", msg.Name, msg.String())
			}
		case err := <-a.Wire.ErrCh:
			a.ErrCh <- err
		case <-a.StopCh:
			break loop
		}
	}
}

// GetTeamID is
func (a *TeamAgentImpl) GetTeamID() int {
	return a.ID
}

// GetTeamName is
func (a *TeamAgentImpl) GetTeamName() string {
	return a.Name
}

// SetTeamID is
func (a *TeamAgentImpl) SetTeamID(id int) {
}

// SetTeamName is
func (a *TeamAgentImpl) SetTeamName(name string) {
}

// GetRegCh is
func (a *TeamAgentImpl) GetRegCh() chan Registration {
	return a.RegCh
}

// GetActionCh is
func (a *TeamAgentImpl) GetActionCh() chan Action {
	return a.ActionCh
}

// GetErrCh is
func (a *TeamAgentImpl) GetErrCh() chan error {
	return a.ErrCh
}

// Invitation sends the Invitation message to the team
func (a *TeamAgentImpl) Invitation(invitation *Invitation) error {
	return a.Wire.Send(invitation.Message())
}

// LegStart sends the LegStart message to the team
func (a *TeamAgentImpl) LegStart(legStart *LegStart) error {
	return a.Wire.Send(legStart.Message())
}

// LegEnd sends the LegEnd message to the team
func (a *TeamAgentImpl) LegEnd(legEnd *LegEnd) error {
	return a.Wire.Send(legEnd.Message())
}

// Round sends the Round message to the team
func (a *TeamAgentImpl) Round(round *Round) error {
	return a.Wire.Send(round.Message())
}

// GameOver sends the GameOver message to the team
func (a *TeamAgentImpl) GameOver(gameOver *GameOver) error {
	return a.Wire.Send(gameOver.Message())
}

// Disconnect is
func (a *TeamAgentImpl) Disconnect() error {
	// todo
	return nil
}
