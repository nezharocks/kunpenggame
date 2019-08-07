package core

// TeamAgent is
type TeamAgent interface {
	LegStart(legStart *LegStart) error
	LegEnd(legEnd *LegEnd) error
	Round(round *Round) (*Action, error)
	GameOver() error
}

// ClientStrategy is
type ClientStrategy interface {
	// Register(registration *Registration) error
	TeamAgent
}

// CenterAgent is
type CenterAgent interface {
	Register(registration *Registration) error
	Act(action *Action) error
}

// CenterStrategy is
type CenterStrategy interface {
	NewTeamID() int
	CenterAgent
	Battle(team TeamAgent) (*Team, error)
}