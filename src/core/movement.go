package core

import (
	"fmt"
	"log"
)

// Movement is
type Movement struct {
	Dx     int
	Dy     int
	Action string
}

const (
	up    = "up"
	down  = "down"
	left  = "left"
	right = "right"
	stay  = ""
)

var (
	// Stay stay still
	Stay = Movement{0, 0, stay}

	// MoveUp moves up a tile
	MoveUp = Movement{0, -1, up}

	// MoveDown moves down a tile
	MoveDown = Movement{0, 1, down}

	// MoveLeft moves left a tile
	MoveLeft = Movement{-1, 0, left}

	// MoveRight moves right a tile
	MoveRight = Movement{1, 0, right}
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

func (m Movement) String() string {
	return fmt.Sprintf("%v - %v, %v", m.Action, m.Dx, m.Dy)
}
