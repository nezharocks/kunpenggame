package core

// TeamAgent is
type TeamAgent interface {
	Disconnect() error
	Invitation(invitation *Invitation) error
	TeamBattle
}

// TeamBattle is
type TeamBattle interface {
	GetTeamID() int
	GetTeamName() string
	SetTeamID(int)
	SetTeamName(string)
	LegStart(legStart *LegStart) error
	LegEnd(legEnd *LegEnd) error
	Round(round *Round) error
	GameOver(gameOver *GameOver) error
	GetRegCh() chan Registration
	GetActionCh() chan Action
	GetErrCh() chan error
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
}

// Game is
type Game interface {
	NewTeamID() int
	Battle(team TeamAgent)
}
