package core

import (
	"fmt"
	"sort"
)

// Guide is
type Guide struct {
	Map *Map

	// T is a vector of all vertex's type
	T []int

	// O is a vector of all vertex's index which locates tile objects in Map
	O []int

	// D is a matrix of shortest distance.
	D [][]int

	// P is a matrix for shortest path.
	P [][]int
}

// NewGuide creates a guide from a Map which is parsed from a JSON bytes.
func NewGuide(m *Map) *Guide {
	n := m.Width * m.Height
	D := createMatrix(n)
	P := createMatrix(n)
	// create T and O, update all objects' vertex index
	T, O := initTnO(m)
	updateWormholeExits(m)
	updateTunnelExits(m, T, O)
	updateGraph(m, T, O, D)
	floyd(D, P)
	return &Guide{
		Map: m,
		T:   T,
		O:   O,
		D:   D,
		P:   P,
	}
}

// ShortestPath is
func (g *Guide) ShortestPath(i, j int) []int {
	path := make([]int, 0, 10)
	k := i
	for {
		k = g.P[k][j]
		path = append(path, k)
		if k == j {
			break
		}
	}
	return path
}

// ShortestDistance is
func (g *Guide) ShortestDistance(i, j int) int {
	return g.D[i][j]
}

func (g *Guide) calcPowerTargets(players []*Player, powers []*Power) []*Target {
	var targets []*Target
	pl, po := len(players), len(powers)
	n := pl * po
	if n == 0 {
		return targets
	}
	targets = make([]*Target, 0, n)
	for _, player := range players {
		for _, power := range powers {
			dist := g.ShortestDistance(player.V, power.V)
			value := float64(power.Point) / float64(dist)
			target := &Target{
				Vertex:  power.V,
				Subject: player,
				Point:   power.Point,
				Dist:    float64(dist),
				Value:   value,
			}
			targets = append(targets, target)
		}
	}
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Value > targets[j].Value
	})
	return targets
}

func (g *Guide) calcPlayerTargets(players []*Player, rivals []*Player) []*Target {
	var targets []*Target
	pl, po := len(players), len(rivals)
	n := pl * po
	if n == 0 {
		return targets
	}
	targets = make([]*Target, 0, po)
	for _, rival := range rivals {
		rivalTotalDist := 0
		for _, player := range players {
			rivalTotalDist += g.ShortestDistance(player.V, rival.V)
		}
		dist := float64(rivalTotalDist) / float64(pl)
		point := rival.Point + Bounty
		value := float64(point) / float64(dist)
		target := &Target{
			Vertex:  rival.V,
			Player:  true,
			Subject: nil, // for all players
			Point:   point,
			Dist:    dist,
			Value:   value,
		}
		targets = append(targets, target)
	}
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Value > targets[j].Value
	})
	return targets
}

func (g *Guide) selectPowerTargets(players []*Player, targets []*Target) ([]*Target, []*Player) {
	selected := make([]*Target, 0)
	idle := make([]*Player, 0)
	for _, p := range players {
		var picked *Target
	loop_targets:
		for i, t := range targets {
			if t == nil {
				continue loop_targets
			}
			if picked == nil {
				if p == t.Subject {
					picked = t
					targets[i] = nil
				}
			} else {
				if picked.Vertex == t.Vertex {
					targets[i] = nil
				}
			}
		}
		if picked != nil {
			selected = append(selected, picked)
		} else {
			idle = append(idle, p)
		}
	}
	return selected, idle
}

func (g *Guide) selectPlayerTargets(players []*Player, targets []*Target) []*Target {
	selected := make([]*Target, 0)
	topRivalTarget := targets[0]
	for _, p := range players {
		t := *topRivalTarget
		t.Subject = p
		selected = append(selected, &t)
	}
	return selected
}

func (g *Guide) generateActions(targets []*Target) []*PlayerAction {
	l := len(targets)
	actions := make([]*PlayerAction, l, l)
	for k, t := range targets {
		i := t.Subject.V
		path := g.ShortestPath(i, t.Vertex)
		j := path[0]
		actions[k] = &PlayerAction{
			Team:   t.Subject.TeamID,
			Player: t.Subject.ID,
			Move:   g.calcMove(i, j),
		}
	}
	return actions
}

// Neighbor is
type Neighbor struct {
	V    int
	Move string
}

func (n Neighbor) String() string {
	return fmt.Sprintf("%v to %d", n.Move, n.V)
}

func getNeighborMovements(i, w, h int) []*Neighbor {
	x, y := i%w, i/w
	nbs := make([]*Neighbor, 0, 4)
	if x-1 >= 0 {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x - 1,
			Move: left,
		})
	}
	if x+1 < w {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x + 1,
			Move: right,
		})
	}
	if y-1 >= 0 {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x - w,
			Move: up,
		})
	}
	if y+1 < h {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x + w,
			Move: down,
		})
	}
	return nbs
}

func (g *Guide) calcMove(src, dest int) []string {
	w, h := g.Map.Width, g.Map.Height
	nbs := getNeighborMovements(src, w, h)
	m := stay
	for _, nb := range nbs {
		fmt.Println(nb.String())
		if nb.V == dest {
			m = nb.Move
			break
		}
	}
	if m == stay {
		fmt.Println("horrible to stay still, need to fix")
		return []string{}
	}
	return []string{m}
}

// Calc is
func (g *Guide) Calc(roundID int, hunting bool, ourPlayers []*Player, rivalPlayers []*Player, powers []*Power) *Action {
	powerTargets := g.calcPowerTargets(ourPlayers, powers)
	selected, idle := g.selectPowerTargets(ourPlayers, powerTargets)
	var selectPlayerTargets []*Target
	if len(idle) > 0 && len(rivalPlayers) > 0 {
		playerTargets := g.calcPlayerTargets(idle, rivalPlayers)
		selectPlayerTargets = g.selectPlayerTargets(idle, playerTargets)
		selected = append(selected, selectPlayerTargets...)
	} else if len(idle) > 0 && len(rivalPlayers) == 0 {
		// todo: calculate and select vision targets
	}
	action := &Action{
		ID: roundID,
	}
	action.Actions = g.generateActions(selected)
	fmt.Println(action.Message().String())
	return action
}

// Target is
type Target struct {
	Vertex int
	// player ID
	Subject *Player

	// flag see if it is a player, by default, it is false (power)
	Player bool
	Point  int
	Dist   float64
	Value  float64
}
