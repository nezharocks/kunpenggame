package core

// TeamAgent is
type TeamAgent interface {
	LegStart(legStart *LegStart) error
	LegEnd(legEnd *LegEnd) error
	Round(round *Round) error
	GameOver(gameOver *GameOver) error
	GetTeamID() int
	GetTeamName() string
	GetRegCh() chan Registration
	GetActionCh() chan Action
	GetErrCh() chan error
	Disconnect()
}

// TeamStrategy is
type TeamStrategy interface {
	GetRegistration() *Registration
	LegStart(legStart *LegStart) error
	LegEnd(legEnd *LegEnd) error
	Round(round *Round) (*Action, error)
	GameOver(gameOver *GameOver) error
}

// GameAgent is
type GameAgent interface {
	Register(registration *Registration) error
	Act(action *Action) error
	OnMessage(msgCh chan Message)
	OffMessage(msgCh chan Message)
	OnError(errCh chan error)
	OffError(errCh chan error)
}

// GameStrategy is
type GameStrategy interface {
	NewTeamID() int
	Battle(team TeamAgent)
}
