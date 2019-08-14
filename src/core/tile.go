package core

import "log"

// Tile is
type Tile struct {
	Type      TT
	Point     int
	Wormhole  *Wormhole
	Tunnel    *Tunnel
	HasPlayer bool
	Player    *Player
	Multiple  bool
	Players   []*Player
}

// TT is tile type
type TT int

func (t TT) String() string {
	s := "unknown"
	switch t {
	case TileHolder:
		s = "holder"
	case TilePower:
		s = "power"
	case TileWormhole:
		s = "wormhole"
	case TileMeteor:
		s = "meteor"
	case TileTunnel:
		s = "tunnel"
	}
	return s
}

const (
	// TileHolder is holder type (container)
	TileHolder = TT(0x00)

	// TilePower is power type (container)
	TilePower = TT(0x01)

	// TileWormhole is wormhole type (container)
	TileWormhole = TT(0x02)

	// TileMeteor is meteor type
	TileMeteor = TT(0x03)

	// TileTunnel is tunnel type
	TileTunnel = TT(0x04)
)

// NewTileHolder is
func NewTileHolder() *Tile {
	return &Tile{
		Type: TileHolder,
	}
}

// NewTilePower is
func NewTilePower(point int) *Tile {
	return &Tile{
		Type:  TilePower,
		Point: point,
	}
}

// NewTileWormhole is
func NewTileWormhole(wormhole *Wormhole) *Tile {
	return &Tile{
		Type:     TileWormhole,
		Wormhole: wormhole,
	}
}

// NewTileMeteor is
func NewTileMeteor() *Tile {
	return &Tile{
		Type: TileMeteor,
	}
}

// NewTileTunnel is
func NewTileTunnel(tunnel *Tunnel) *Tile {
	return &Tile{
		Type:   TileTunnel,
		Tunnel: tunnel,
	}
}

// BecomeHolder is
func (t *Tile) BecomeHolder() {
	if t.Type != TilePower {
		log.Println("cannot become from %v to to %v", t.Type.String(), TileHolder.String())
		return
	}
	t.Type = TileHolder
	t.Point = 0
}
