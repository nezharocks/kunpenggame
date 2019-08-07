package core

// Power is
type Power struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Point int `json:"point"`
}

// Meteor is
type Meteor struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Tunnel is
type Tunnel struct {
	Direction string `json:"direction"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
}

// Wormhole is
type Wormhole struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// Map is
type Map struct {
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Vision   int        `json:"vision"`
	Meteor   []Meteor   `json:"meteor"`
	Tunnel   []Tunnel   `json:"tunnel"`
	Wormhole []Wormhole `json:"wormhole"`
}

// Team is ,omitempty
type Team struct {
	ID         int    `json:"id"`
	NAME       string `json:"name,omitempty"`
	Players    []int  `json:"players"`
	Force      string `json:"force"`
	Point      int    `json:"point"`
	RemainLife int    `json:"remain_life"`
}

// Player is one of the member of a team which joins the game as a game player
type Player struct {
	Team  int `json:"team"`
	ID    int `json:"id"`
	Score int `json:"score"`
	Sleep int `json:"sleep"`
	X     int `json:"x"`
	Y     int `json:"y"`
}

// Registration is
type Registration struct {
	TeamID   int    `json:"team_id"`
	TeamName string `json:"team_name"`
}

// LegStart is
type LegStart struct {
	Map   Map    `json:"map"`
	Teams []Team `json:"teams"`
}

// LegEnd is
type LegEnd struct {
	Teams []Team `json:"teams"`
}

// Round is
type Round struct {
	ID      int      `json:"round_id"`
	Mode    string   `json:"mode"`
	Powers  []Power  `json:"power"`
	Players []Player `json:"players"`
	Teams   []Team   `json:"teams"`
}

// Action is
type Action struct {
	ID      int            `json:"round_id"`
	Actions []PlayerAction `json:"actions"`
}

// PlayerAction is
type PlayerAction struct {
	Team   int      `json:"team"`
	Player int      `json:"player_id"`
	Move   []string `json:"move"`
}

// APIMsg is
type APIMsg struct {
	Name string     `json:"msg_name"`
	Data APIMsgData `json:"msg_data"`
}

// APIMsgData is
type APIMsgData interface{}