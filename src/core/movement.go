package core

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

// NewMovement is
func NewMovement(m []string) *Movement {
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
		return &Stay
	}
}
