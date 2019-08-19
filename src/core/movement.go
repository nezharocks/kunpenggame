package core

import "log"

// Movement is
type Movement struct {
	Dx int
	Dy int
}

const (
	up    = "up"
	down  = "down"
	left  = "left"
	right = "right"
)

var (
	// Stay stay still
	Stay = Movement{0, 0}

	// MoveUp moves up a tile
	MoveUp = Movement{0, -1}

	// MoveDown moves down a tile
	MoveDown = Movement{0, 1}

	// MoveLeft moves left a tile
	MoveLeft = Movement{-1, -0}

	// MoveRight moves right a tile
	MoveRight = Movement{1, 0}
)

// NewMovementFromAction is
func NewMovementFromAction(m []string) *Movement {
	if len(m) == 0 {
		return &Stay
	}
	move := m[0]
	switch move {
	case up:
		return &MoveUp
	case down:
		return &MoveDown
	case left:
		return &MoveLeft
	case right:
		return &MoveRight
	default:
		log.Printf("wrong action move %q", move)
		return &Stay
	}
}

// NewMovementFromTunnel is
func NewMovementFromTunnel(dir string) *Movement {
	switch dir {
	case dirUp:
		return &MoveUp
	case dirDown:
		return &MoveDown
	case dirLeft:
		return &MoveLeft
	case dirRight:
		return &MoveRight
	default:
		log.Printf("wrong tunnel direction %q", dir)
		return nil
	}
}
