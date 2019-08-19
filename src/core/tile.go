package core

import "log"

// Tile is
type Tile struct {
	Type      TT
	X, Y      int
	Power     *Power
	Wormhole  *Wormhole
	Meteor    *Meteor
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
func NewTileHolder(x, y int) *Tile {
	return &Tile{
		Type: TileHolder,
		X:    x,
		Y:    y,
	}
}

// NewTilePower is
func NewTilePower(power *Power) *Tile {
	return &Tile{
		Type:  TilePower,
		X:     power.X,
		Y:     power.Y,
		Power: power,
	}
}

// NewTileWormhole is
func NewTileWormhole(wormhole *Wormhole) *Tile {
	return &Tile{
		Type:     TileWormhole,
		X:        wormhole.X,
		Y:        wormhole.Y,
		Wormhole: wormhole,
	}
}

// NewTileMeteor is
func NewTileMeteor(meteor *Meteor) *Tile {
	return &Tile{
		Type:   TileMeteor,
		X:      meteor.X,
		Y:      meteor.Y,
		Meteor: meteor,
	}
}

// NewTileTunnel is
func NewTileTunnel(tunnel *Tunnel) *Tile {
	return &Tile{
		Type:   TileTunnel,
		X:      tunnel.X,
		Y:      tunnel.Y,
		Tunnel: tunnel,
	}
}

// BecomeHolder is
func (t *Tile) BecomeHolder() {
	if t.Type != TilePower {
		log.Printf("cannot become from %v to to %v", t.Type.String(), TileHolder.String())
		return
	}
	t.Type = TileHolder
}

// AddPlayer is
func (t *Tile) AddPlayer(player *Player) {
	if !t.HasPlayer {
		t.Player = player
		t.HasPlayer = true
		return
	}
	if t.Multiple {
		t.Players = append(t.Players, player)
		return
	}
	t.Players = make([]*Player, 0)
	t.Players = append(t.Players, t.Player, player)
	t.Player = nil
	t.Multiple = true
}

// RemovePlayer is
func (t *Tile) RemovePlayer(player *Player) {
	if !t.HasPlayer {
		return
	}
	if !t.Multiple {
		t.Player = nil
		t.HasPlayer = false
		return
	}

	pi := -1
	for i, p := range t.Players {
		if p == player {
			pi = i
			break
		}
	}
	if pi == -1 {
		return
	}

	l := len(t.Players)
	if l == 2 {
		var left *Player
		for _, p := range t.Players {
			if p != player {
				left = p
				break
			}
		}
		t.Player = left
		t.Players = nil
		t.Multiple = false
		return
	}
	t.Players[l-1], t.Players[pi] = t.Players[pi], t.Players[l-1]
	t.Players = t.Players[:l-1]
}

// GetPlayer is
func (t *Tile) GetPlayer() *Player {
	if !t.HasPlayer {
		return nil
	}
	if t.Multiple {
		return t.Players[0]
	}
	return t.Player
}

// // Enter is
// func (t *Tile) Enter(player *Player, powerForce TeamForce) {
// 	switch t.Type {
// 	case TileHolder:
// 		t.enterHolder(player, powerForce)
// 	case TilePower:
// 		t.enterPower(player)
// 	case TileWormhole:
// 		t.enterWormhome(player, powerForce)
// 	case TileMeteor:
// 		// stay still and do nothing
// 	case TileTunnel:
// 		t.enterTunnel(player, powerForce)
// 	}
// }

// func (t *Tile) enterPower(player *Player) {
// }

// func (t *Tile) enterHolder(player *Player, powerForce TeamForce) {
// }

// func (t *Tile) enterWormhome(player *Player, powerForce TeamForce) {
// }

// func (t *Tile) enterTunnel(player *Player, powerForce TeamForce) {
// }
