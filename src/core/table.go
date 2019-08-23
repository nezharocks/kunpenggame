package core

import "log"

// Table is
type Table struct {
	Tiles         [][]*Tile
	Leg           *GameBattleLeg
	Map           *Map
	tunnelHistory []*Tile
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
		tiles[meteor.X][meteor.Y] = NewTileMeteor(meteor)
	}

	// init tunnels
	for _, tunnel := range m.Tunnels {
		tiles[tunnel.X][tunnel.Y] = NewTileTunnel(tunnel)
	}

	// init powers
	for _, power := range m.Powers {
		tiles[power.X][power.Y] = NewTilePower(power)
	}

	// init players
	for _, players := range leg.TeamPlayers {
		for _, player := range players {
			tile := NewTileHolder(player.X, player.Y)
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
			tiles[i][j] = NewTileHolder(i, j)
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
func (t *Table) Move(player *Player, move *Movement, powerForce TeamForce) {
	x, y, still := t.nextCoordinates(player.X, player.Y, move)
	if still {
		return
	}
	t.EnterTile(player, t.Tiles[x][y], powerForce)
}

// EnterTile is
func (t *Table) EnterTile(player *Player, tile *Tile, powerForce TeamForce) {
	switch tile.Type {
	case TileHolder:
		t.enterHolder(player, tile, powerForce)
	case TilePower:
		t.enterPower(player, tile)
	case TileWormhole:
		t.enterWormhole(player, tile, powerForce)
	case TileMeteor:
		log.Printf("%v hits %v, then stay still", player.String(), tile.Meteor.String())
	case TileTunnel:
		t.enterTunnel(player, tile, powerForce)
	}
	t.tunnelHistory = nil
}

func (t *Table) enterPower(player *Player, tileIn *Tile) {
	tileOut := t.Tiles[player.X][player.Y]
	tileOut.RemovePlayer(player)
	point := tileIn.Power.Point
	player.Point += point
	player.Team.Point += point
	player.X = tileIn.X
	player.Y = tileIn.Y
	tileIn.AddPlayer(player)
	tileIn.BecomeHolder()
}

func (t *Table) enterHolder(player *Player, tileIn *Tile, powerForce TeamForce) {
	tileOut := t.Tiles[player.X][player.Y]
	tileOut.RemovePlayer(player)
	if !tileIn.HasPlayer {
		player.X = tileIn.X
		player.Y = tileIn.Y
		tileIn.AddPlayer(player)
		return
	}
	tp := tileIn.GetPlayer()
	tileForce := tp.Team.Force
	playerForce := player.Team.Force
	if tileForce == playerForce {
		player.X = tileIn.X
		player.Y = tileIn.Y
		tileIn.AddPlayer(player)
		return
	}

	// the player suicides
	if tileForce == powerForce.String() {
		t.eatEscapee(tp, player)
		return
	}

	// the player eats prey (enemy players in the tile)
	if tileIn.Multiple {
		for _, escapee := range tileIn.Players {
			t.eatEscapee(player, escapee)
		}
	} else {
		t.eatEscapee(player, tp)
	}
	player.X = tileIn.X
	player.Y = tileIn.Y
	tileIn.AddPlayer(player)
}

func (t *Table) eatEscapee(hunter, escapee *Player) {
	// calculate earned point
	hunter.Point += escapee.Point
	hunter.Team.Point += escapee.Point + Bounty

	// make escapee lose its life
	escapee.Point = 0
	escapee.Sleep = 1
	if escapee.Team.RemainLife <= 0 {
		escapee.Dead = true
		return
	}
	escapee.Team.EnqueueSleeper(escapee)
}

func (t *Table) enterWormhole(player *Player, tile *Tile, powerForce TeamForce) {
	if tile.Wormhole == nil || tile.Wormhole.Exit == nil {
		log.Println("state illegal - tile or tile's wormhole has no wormhole")
		return
	}
	exit := tile.Wormhole.Exit
	t.enterHolder(player, t.Tiles[exit.X][exit.Y], powerForce)
}

func (t *Table) enterTunnel(player *Player, tile *Tile, powerForce TeamForce) {
	if tile.Tunnel == nil {
		log.Println("state illegal - tile has no tunnel")
		log.Printf("%v skip moving to the empty tunnel", player.String())
		return
	}
	move := NewMovementFromTunnel(tile.Tunnel.Direction)
	if move == nil {
		log.Printf("%v get null movement from %v\n", player.String(), tile.Tunnel.String())
		log.Printf("%v skip moving to %v", player.String(), tile.Tunnel.String())
		return
	}

	x, y, still := t.nextCoordinates(tile.X, tile.Y, move)
	if still {
		log.Printf("%v get wrong movement from %v\n", player.String(), tile.Tunnel.String())
		log.Printf("%v skip moving to %v", player.String(), tile.Tunnel.String())
		return
	}
	t.EnterTile(player, t.Tiles[x][y], powerForce)
}

func (t *Table) nextCoordinates(tx, ty int, move *Movement) (x, y int, still bool) {
	x, y = tx+move.Dx, ty+move.Dy
	w, h := t.Map.Width-1, t.Map.Height-1
	if x < 0 {
		x = 0
	} else if x > w {
		x = w
	}
	if y < 0 {
		y = 0
	} else if y > h {
		y = h
	}
	still = tx == x && ty == y
	return x, y, still
}
