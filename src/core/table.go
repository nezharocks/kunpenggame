package core

// Table is
type Table struct {
	Tiles       [][]*Tile
	Map         *Map
	TeamBattles [TeamNum]TeamBattle
	TeamMap     map[TeamBattle]*Team
	PlayersMap  map[TeamBattle][]*Player
}

// NewTable is
func NewTable(m *Map, teamBattles [TeamNum]TeamBattle, teamMap map[TeamBattle]*Team, playersMap map[TeamBattle][]*Player) *Table {
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
	for _, players := range playersMap {
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
		Tiles:       tiles,
		Map:         m,
		TeamMap:     teamMap,
		PlayersMap:  playersMap,
		TeamBattles: teamBattles,
	}
}

// TeamActivePlayers is
func (t *Table) TeamActivePlayers(tb TeamBattle) []*Player {
	players := t.PlayersMap[tb]
	activePlayers := make([]*Player, 0, len(players))
	for _, p := range players {
		if !p.IsAsleep() {
			activePlayers = append(activePlayers, p)
		}
	}
	return activePlayers
}

// TeamAlivePlayers is
func (t *Table) TeamAlivePlayers(tb TeamBattle) []*Player {
	players := t.PlayersMap[tb]
	activePlayers := make([]*Player, 0, len(players))
	for _, p := range players {
		if !p.IsDead() {
			activePlayers = append(activePlayers, p)
		}
	}
	return activePlayers
}

// GetRival is
func (t *Table) GetRival(tb TeamBattle) TeamBattle {
	for _, v := range t.TeamBattles {
		if v != tb {
			return v
		}
	}
	return nil
}

// GetVisiblePowers is
func (t *Table) GetVisiblePowers(visions []*Vision) []*Power {
	ret := make([]*Power, 0, len(t.Map.Powers))
	for _, p := range t.Map.Powers {
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
func (t *Table) GetVisiblePlayers(players []*Player, visions []*Vision) []*Player {
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
