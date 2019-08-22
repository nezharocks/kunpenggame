package core

// Invitation is
type Invitation struct {
	TeamID int `json:"team_id"`
}

// Message returns the pointer of the generated message of the Invitation object
func (m *Invitation) Message() *Message {
	return &Message{InvitationName, m, false}
}

// Registration is
type Registration struct {
	TeamID   int    `json:"team_id"`
	TeamName string `json:"team_name"`
}

// Message returns the pointer of the generated message of the Registration object
func (m *Registration) Message() *Message {
	return &Message{RegistrationName, m, false}
}

// LegStart is
type LegStart struct {
	Map   *Map    `json:"map"`
	Teams []*Team `json:"teams"`
}

// Message returns the pointer of the generated message of the LegStart object
func (m *LegStart) Message() *Message {
	return &Message{LegStartName, m, false}
}

// LegEnd is
type LegEnd struct {
	Teams []*Team `json:"teams"`
}

// Message returns the pointer of the generated message of the LegEnd object
func (m *LegEnd) Message() *Message {
	return &Message{LegEndName, m, false}
}

// Round is
type Round struct {
	ID      int       `json:"round_id"`
	Mode    string    `json:"mode"`
	Powers  []*Power  `json:"power"`
	Players []*Player `json:"players"`
	Teams   []*Team   `json:"teams"`
}

// Message returns the pointer of the generated message of the Round object
func (m *Round) Message() *Message {
	return &Message{RoundName, m, false}
}

// Action is
type Action struct {
	ID      int             `json:"round_id"`
	Actions []*PlayerAction `json:"actions"`
}

// Message returns the pointer of the generated message of the Action object
func (m *Action) Message() *Message {
	return &Message{ActionName, m, false}
}

// GameOver is
type GameOver struct{}

// Message returns the pointer of the generated message of the GameOver object
func (m *GameOver) Message() *Message {
	return &Message{GameOverName, m, false}
}

// PlayerAction is
type PlayerAction struct {
	Team   int      `json:"team"`
	Player int      `json:"player_id"`
	Move   []string `json:"move"`
}

// LegStartTeam is
type LegStartTeam struct {
	ID      int    `json:"id"`
	Players []int  `json:"players"`
	Force   string `json:"force"`
}

// LegEndTeam is
type LegEndTeam struct {
	ID    int `json:"id"`
	Point int `json:"point"`
}

// RoundTeam is
type RoundTeam struct {
	ID         int `json:"id"`
	Point      int `json:"point"`
	RemainLife int `json:"remain_life"`
}
