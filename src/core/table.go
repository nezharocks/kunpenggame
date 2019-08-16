package core

// Table is
type Table struct {
	Tiles        [][]*Tile
	Map          *Map
	Teams        []*Team
	TeamsPlayers [][]*Player
}

// NewTable is
func NewTable(m *Map, teams []*Team, teamsPlayers [][]*Player) *Table {
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

	// init players
	tl := len(teamsPlayers)
	for ti := 0; ti < tl; ti++ {
		players := teamsPlayers[ti]
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
		Tiles:        tiles,
		Map:          m,
		Teams:        teams,
		TeamsPlayers: teamsPlayers,
	}
}
