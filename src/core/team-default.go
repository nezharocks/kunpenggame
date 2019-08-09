package core

// TeamSimple is
type TeamSimple struct {
	ID   int
	Name string
}

// NewTeamSimple creates a TeamSimple instance
func NewTeamSimple(name string) *TeamSimple {
	return &TeamSimple{
		Name: name,
	}
}

// GetID is
func (t *TeamSimple) GetID() int {
	return t.ID
}

// GetName is
func (t *TeamSimple) GetName() string {
	return t.Name
}

// SetID is
func (t *TeamSimple) SetID(id int) {
	t.ID = id
}

// GameStart is
func (t *TeamSimple) GameStart() {
	// todo
}

// LegStart is
func (t *TeamSimple) LegStart(legStart *LegStart) error {
	// todo
	return nil
}

// LegEnd is
func (t *TeamSimple) LegEnd(legEnd *LegEnd) error {
	// todo
	return nil
}

// Round is
func (t *TeamSimple) Round(round *Round) (*Action, error) {
	// todo
	return nil, nil
}

// GameOver is
func (t *TeamSimple) GameOver(gameOver *GameOver) error {
	// todo
	return nil
}
