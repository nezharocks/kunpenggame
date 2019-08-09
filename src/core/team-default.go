package core

// TeamSimple is
type TeamSimple struct {
	ID   int
	Name string
}

// NewTeamSimple creates a TeamSimple instance
func NewTeamSimple(id int, name string) *TeamSimple {
	return &TeamSimple{
		ID:   id,
		Name: name,
	}
}

// GetRegistration is
func (t *TeamSimple) GetRegistration() *Registration {
	return &Registration{t.ID, t.Name}
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
func (t *TeamSimple) GameOver() error {
	// todo
	return nil
}
