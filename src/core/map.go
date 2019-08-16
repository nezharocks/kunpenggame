package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Map is
type Map struct {
	Width            int                     `json:"width"`
	Height           int                     `json:"height"`
	Vision           int                     `json:"vision"`
	Meteors          []*Meteor               `json:"meteor"`
	Tunnels          []*Tunnel               `json:"tunnel"`
	Wormholes        []*Wormhole             `json:"wormhole"`
	Powers           []*Power                `json:"power"`
	TeamPlaceHolders [TeamNum][]*PlaceHolder `json:"-"`
}

// NewMapFromString is
func NewMapFromString(data string) (m *Map, err error) {
	x, y := 0, 0
	m = &Map{
		Meteors:   make([]*Meteor, 0),
		Tunnels:   make([]*Tunnel, 0),
		Wormholes: make([]*Wormhole, 0),
		Powers:    make([]*Power, 0),
	}
	oPlaceHolders := make([]*PlaceHolder, 0, DefaultPlayerNum)
	xPlaceHolders := make([]*PlaceHolder, 0, DefaultPlayerNum)
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
		case 'O':
			oPlaceHolders = append(oPlaceHolders, &PlaceHolder{x, y, nil})
		case 'X':
			xPlaceHolders = append(xPlaceHolders, &PlaceHolder{x, y, nil})
		default:
			if (c > 'a' && c < 'z') || (c > 'A' || c < 'Z') {
				m.Wormholes = append(m.Wormholes, &Wormhole{x, y, string(c), nil})
			} else {
				fmt.Printf("char %v is not supported at (%v,%v)", c, x, y)
			}
		}
		x++
	}
	m.Height = y + 1
	m.TeamPlaceHolders[0] = oPlaceHolders
	m.TeamPlaceHolders[1] = xPlaceHolders
	return m, nil
}

// Init is
func (m *Map) Init(vision, width, height int) error {
	m.Vision = vision

	// pair wormholes
	wormholeMap := make(map[string]*Wormhole, 10)
	for _, wormhole := range m.Wormholes {
		name := strings.ToLower(wormhole.Name)
		existed, ok := wormholeMap[name]
		if !ok {
			wormholeMap[name] = wormhole
			continue
		}
		existed.Exit = wormhole
		wormhole.Exit = existed
	}

	// check if the specified width and the parsed width are matched
	if width != m.Width {
		return fmt.Errorf("the given width %v is different from the width %v parsed from map data", width, m.Width)
	}

	// check if the specified height and the parsed height are matched
	if height != m.Height {
		return fmt.Errorf("the given height %v is different from the height %v parsed from map data", height, m.Height)
	}

	// check if the numbers of 'O' and 'X' place holders are th e same
	ol := len(m.TeamPlaceHolders[0])
	xl := len(m.TeamPlaceHolders[1])
	if ol != xl {
		return fmt.Errorf("the numbers of 'O' and 'X' place holder are different: %v:%v", ol, xl)
	}

	return nil
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
