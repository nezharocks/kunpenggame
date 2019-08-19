package core

import "fmt"

// Tunnel is
type Tunnel struct {
	Direction string `json:"direction"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
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
	return &Tunnel{dir, x, y}
}

func (o Tunnel) String() string {
	return fmt.Sprintf("tunnel-%v@%v,%v", o.Direction, o.X, o.Y)
}
