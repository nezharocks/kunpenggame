package core

import (
	"encoding/json"
	"fmt"
)

// Map is
type Map struct {
	Width        int            `json:"width"`
	Height       int            `json:"height"`
	Vision       int            `json:"vision"`
	Meteors      []*Meteor      `json:"meteor"`
	Tunnels      []*Tunnel      `json:"tunnel"`
	Wormholes    []*Wormhole    `json:"wormhole"`
	Powers       []*Power       `json:"power"`
	PlaceHolders []*PlaceHolder `json:"holder"`
	Mat          [][]*Tile      `json:"-"`
}

// NewMapFromString is
func NewMapFromString(data string) (m *Map, err error) {
	x, y := 0, 0
	m = &Map{
		Meteors:      make([]*Meteor, 0),
		Tunnels:      make([]*Tunnel, 0),
		Wormholes:    make([]*Wormhole, 0),
		Powers:       make([]*Power, 0),
		PlaceHolders: make([]*PlaceHolder, 0),
	}
loop:
	for _, c := range data {
		switch c {
		case '.':
		case '…':
			x = x + 2
		case '\n':
			if m.Width == 0 {
				m.Width = x
			} else {
				if m.Width != x {
					return nil, fmt.Errorf("the %dth row has different width (%v) from other width (%v)", y, x, m.Width)
				}
			}

			x = 0
			y++
			continue loop
		case '#':
			m.Meteors = append(m.Meteors, &Meteor{x, y})
		case '^', 'v', '<', '>':
			tunnel := NewTunnelFromChar(c, x, y)
			m.Tunnels = append(m.Tunnels, tunnel)
		case '1', '2', '3', '4', '5':
			m.Powers = append(m.Powers, &Power{x, y, int(c - 48)})
		case 'O', 'X':
			m.PlaceHolders = append(m.PlaceHolders, &PlaceHolder{string(c), x, y})
		default:
			if (c > 'a' && c < 'z') || (c > 'A' || c < 'Z') {
				m.Wormholes = append(m.Wormholes, &Wormhole{string(c), x, y})
			} else {
				fmt.Printf("char %v is not supported at (%v,%v)", c, x, y)
			}
		}
		x++
	}
	m.Height = y + 1

	if m.Width != m.Height {
		return nil, fmt.Errorf("map's width %v and height %v are different", m.Width, m.Height)
	}

	return m, nil
}

// JSON is
func (m *Map) JSON() string {
	// bytes, err := json.MarshalIndent(m, "", "    ")
	bytes, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// Matrix is
func (m *Map) Matrix() [][]*Tile {
	m.Mat = make([][]*Tile, m.Width, m.Width)
	for i := 0; i < m.Width; i++ {
		m.Mat[i] = make([]*Tile, m.Height, m.Height)
	}

	// init wormholes
	for _, wormhole := range m.Wormholes {
		m.Mat[wormhole.X][wormhole.Y] = NewTileWormhole(wormhole)
	}

	// init meteors
	for _, meteor := range m.Meteors {
		m.Mat[meteor.X][meteor.Y] = NewTileMeteor()
	}

	// init tunnels
	for _, tunnel := range m.Tunnels {
		m.Mat[tunnel.X][tunnel.Y] = NewTileTunnel(tunnel)
	}

	// init blank holders left
	for i := 0; i < m.Width; i++ {
	rows:
		for j := 0; j < m.Height; j++ {
			if m.Mat[i][j] != nil {
				continue rows
			}
			m.Mat[i][j] = NewTileHolder()
		}
	}

	return m.Mat
}
