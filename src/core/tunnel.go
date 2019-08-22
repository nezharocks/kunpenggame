package core

import "fmt"

// Tunnel is
type Tunnel struct {
	Direction  string `json:"direction"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	ExitVertex int    `json:"-"`
	V          int    `json:"-"`
}

const (
	dirUp    = "up"
	dirDown  = "down"
	dirLeft  = "left"
	dirRight = "right"
)

// NewTunnelFromChar is
func NewTunnelFromChar(c rune, x, y int) *Tunnel {
	dir := ""
	switch c {
	case '^':
		dir = dirUp
	case 'v':
		dir = dirDown
	case '<':
		dir = dirLeft
	case '>':
		dir = dirRight
	}
	return &Tunnel{dir, x, y, -1, -1}
}

func (o Tunnel) String() string {
	entry := fmt.Sprintf("entry(%v,%v)", o.X, o.Y)
	exit := fmt.Sprintf("exit(%v)", o.ExitVertex)
	return fmt.Sprintf("tunnel-%v@%v->%v", o.Direction, entry, exit)
}

// Char returns a arrow string
func (o Tunnel) Char() string {
	switch o.Direction {
	case dirUp:
		return "^"
	case dirDown:
		return "v"
	case dirLeft:
		return "<"
	case dirRight:
		return ">"
	default:
		return "?"
	}

}
