package core

// Table is
type Table struct {
	Tiles [][]*Tile
	Leg   *GameBattleLeg
	Map   *Map
}

// NewTable is
func NewTable(leg *GameBattleLeg) *Table {
	m := leg.Battle.Map

	// init the 2x2 matrix
	tiles := make([][]*Tile, m.Width, m.Width)
	for i := 0; i < m.Width; i++ {
		tiles[i] = make([]*Tile, m.Height, m.Height)
	}

	// init wormholes
	for _, wormhole := range m.Wormholes {
		tiles[wormhole.X][wormhole.Y] = NewTileWormhole(wormhole)
	}

	// init meteors
	for _, meteor := range m.Meteors {
		tiles[meteor.X][meteor.Y] = NewTileMeteor()
	}

	// init tunnels
	for _, tunnel := range m.Tunnels {
		tiles[tunnel.X][tunnel.Y] = NewTileTunnel(tunnel)
	}

	// init powers
	for _, power := range m.Powers {
		tiles[power.X][power.Y] = NewTilePower(power.Point)
	}

	// init players
	for _, players := range leg.TeamPlayers {
		for _, player := range players {
			tile := NewTileHolder()
			tile.AddPlayer(player)
			tiles[player.X][player.Y] = tile
		}
	}

	// init blank holders left
	for i := 0; i < m.Width; i++ {
	rows:
		for j := 0; j < m.Height; j++ {
			if tiles[i][j] != nil {
				continue rows
			}
			tiles[i][j] = NewTileHolder()
		}
	}

	return &Table{
		Tiles: tiles,
		Leg:   leg,
		Map:   m,
	}
}

// GetVisiblePowers is
func GetVisiblePowers(powers []*Power, visions []*Vision) []*Power {
	ret := make([]*Power, 0, len(powers))
	for _, p := range powers {
		// fmt.Printf("%+v\n", *p)
		if p == nil {
			continue
		}
	loop_vision:
		for _, v := range visions {
			visible := v.InVision(p.X, p.Y)
			// fmt.Printf("%+v\t%+v\t%+v\n", *p, *v, visible)
			if visible {
				ret = append(ret, p)
				break loop_vision
			}
		}
	}
	return ret
}

// GetVisiblePlayers is
func GetVisiblePlayers(players []*Player, visions []*Vision) []*Player {
	ret := make([]*Player, 0, len(players))
	for _, p := range players {
		if p.IsAsleep() {
			continue
		}
	loop_vision:
		for _, v := range visions {
			if v.InVision(p.X, p.Y) {
				ret = append(ret, p)
				break loop_vision
			}
		}
	}
	return ret
}

// Move is
func (t *Table) Move(player *Player, powerForce TeamForce, move Movement) {
	// todo
}
