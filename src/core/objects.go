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

// Wormhole is
type Wormhole struct {
	Name string    `json:"name"`
	X    int       `json:"x"`
	Y    int       `json:"y"`
	Exit *Wormhole `json:"-"`
}

// PlaceHolder is
type PlaceHolder struct {
	Name   string  `json:"name"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Player *Player `json:"-"`
}

// Team is
type Team struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Players    []int  `json:"players,omitempty"`
	Force      string `json:"force,omitempty"`
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
