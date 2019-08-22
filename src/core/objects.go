package core

import "fmt"

// Power is
type Power struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Point int `json:"point"`
	V     int `json:"-"`
}

func (o Power) String() string {
	return fmt.Sprintf("power-%v@tile(%v,%v)", o.Point, o.X, o.Y)
}

// Meteor is
type Meteor struct {
	X int `json:"x"`
	Y int `json:"y"`
	V int `json:"-"`
}

func (o Meteor) String() string {
	return fmt.Sprintf("meteor@tile(%v,%v)", o.X, o.Y)
}

// Wormhole is
type Wormhole struct {
	X          int       `json:"x"`
	Y          int       `json:"y"`
	Name       string    `json:"name"`
	Exit       *Wormhole `json:"-"`
	ExitVertex int       `json:"-"`
	V          int       `json:"-"`
}

func (o Wormhole) String() string {
	entry := fmt.Sprintf("entry(%v,%v)", o.X, o.Y)
	exit := "exit()"
	if o.Exit != nil {
		exit = fmt.Sprintf("exit(%v,%v - %v)", o.Exit.X, o.Exit.Y, o.ExitVertex)
	}
	return fmt.Sprintf("wormhole-%v@%v->%v", o.Name, entry, exit)
}

// PlaceHolder is
type PlaceHolder struct {
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Player *Player `json:"-"`
}

// Team is
type Team struct {
	ID         int       `json:"id"`
	Name       string    `json:"-"`
	Players    []int     `json:"players,omitempty"`
	Force      string    `json:"force,omitempty"`
	Point      int       `json:"point"`
	RemainLife int       `json:"remain_life"`
	SleeperQ   []Sleeper `json:"-"`
}

// EnqueueSleeper is
func (t *Team) EnqueueSleeper(player *Player) {
	if t.SleeperQ == nil {
		t.SleeperQ = make([]Sleeper, 0)
	}
	t.SleeperQ = append(t.SleeperQ, Sleeper{SleepRound, player})
}

// DequeueSleepers is
func (t *Team) DequeueSleepers() []*Player {
	if len(t.SleeperQ) == 0 {
		return nil
	}
	l := -1
	var players []*Player
	for i, s := range t.SleeperQ {
		if s.Round == 0 {
			if players == nil {
				players = make([]*Player, 0)
			}
			players = append(players, s.Player)
		} else {
			l = i
			break
		}
	}
	if l != -1 && l != 0 {
		t.SleeperQ = t.SleeperQ[:l]
	}
	return players
}

// SleepLess is
func (t *Team) SleepLess() {
	for _, s := range t.SleeperQ {
		s.Round--
	}
}

// Sleeper is
type Sleeper struct {
	Round  int
	Player *Player
}

// Player is one of the member of a team which joins the game as a game player
type Player struct {
	TeamID int   `json:"team"`
	ID     int   `json:"id"`
	Point  int   `json:"score"`
	Sleep  int   `json:"sleep"`
	X      int   `json:"x"`
	Y      int   `json:"y"`
	Dead   bool  `json:"-"`
	Team   *Team `json:"-"`
	V      int   `json:"-"`
}

// String is
func (p Player) String() string {
	return fmt.Sprintf("%v-team-%v-player-%v@%v,%v", p.Team.Force, p.TeamID, p.ID, p.X, p.Y)
}

// IsDead is
func (p Player) IsDead() bool {
	return p.Dead
}

// IsAsleep is
func (p Player) IsAsleep() bool {
	return p.Sleep != 0
}
