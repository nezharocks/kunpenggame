package core

// TeamAgent is
type TeamAgent interface {
	Invitation(invitation *Invitation) error
	LegStart(legStart *LegStart) error
	LegEnd(legEnd *LegEnd) error
	Round(round *Round) error
	GameOver(gameOver *GameOver) error
	GetTeamID() int
	GetTeamName() string
	GetRegCh() chan Registration
	GetActionCh() chan Action
	GetErrCh() chan error
	Disconnect() error
}

// TeamStrategy is
type TeamStrategy interface {
	GetID() int
	GetName() string
	SetID(int)
	GameStart()
	LegStart(legStart *LegStart) error
	LegEnd(legEnd *LegEnd) error
	Round(round *Round) (*Action, error)
	GameOver(gameOver *GameOver) error
}

// GameAgent is
type GameAgent interface {
	Connect() error
	Disconnect() error
	Registration(registration *Registration) error
	Action(action *Action) error
	GetInvitationCh() chan Invitation
	GetLegStartCh() chan LegStart
	GetLegEndCh() chan LegEnd
	GetRoundCh() chan Round
	GetGameOverCh() chan GameOver
	GetErrorCh() chan error

	// OnMessage(msgCh chan Message)
	// OffMessage(msgCh chan Message)
	// OnError(errCh chan error)
	// OffError(errCh chan error)
}

// GameStrategy is
type GameStrategy interface {
	NewTeamID() int
	Battle(team TeamAgent)
}
