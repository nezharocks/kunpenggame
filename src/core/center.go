package core

// Center is
type Center struct {
	ID       int
	Name     string
	Strategy CenterStrategy
}

// NewCenter creates a Center instance
func NewCenter(id int, name string, strategy CenterStrategy) *Center {
	return &Center{}
}

// NewTeamID is
func (c *Center) NewTeamID() int {
	return c.Strategy.NewTeamID()
}

// Register is
func (c *Center) Register(registration *Registration) error {
	return c.Strategy.Register(registration)
}

// Act is
func (c *Center) Act(action *Action) error {
	return c.Strategy.Act(action)
}

// Battle is
func (c *Center) Battle(team TeamAgent) (*Team, error) {
	return c.Strategy.Battle(team)
}
