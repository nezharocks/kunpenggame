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
	Width    int         `json:"width"`
	Height   int         `json:"height"`
	Vision   int         `json:"vision"`
	Meteor   []*Meteor   `json:"meteor"`
	Tunnel   []*Tunnel   `json:"tunnel"`
	Wormhole []*Wormhole `json:"wormhole"`
}

// Team is ,omitempty
type Team struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
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
