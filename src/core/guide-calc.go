package core

import (
	"math/rand"
	"sort"
)

func (g *Guide) calcOneStepEatTargets(tss ...[]*Target) []*Target {
	var selected []*Target
	oneSteps := make([]*Target, 0)
	for _, ts := range tss {
		for _, t := range ts {
			if t.Dist == 1 {
				oneSteps = append(oneSteps, t)
			}
		}
	}
	if len(oneSteps) == 0 {
		return selected
	}
	selected = make([]*Target, 0)
	if len(oneSteps) > 1 {
		sort.Slice(oneSteps, func(i, j int) bool {
			return oneSteps[i].Point > oneSteps[j].Point
		})
		cache := make(map[int]bool, len(oneSteps)*2)
		for _, t := range oneSteps {
			picked, ok := cache[t.Vertex]
			if ok && picked {
				continue
			}
			selected = append(selected, t)
			cache[t.Vertex] = true
		}
	} else {
		selected = oneSteps
	}
	return selected
}

// func (g *Guide) calcOneStepEscapeTargets(ourPowerTargets, ourRivalTargets, rivalPowerTargets []*Target) []*Target {
// 	var selected []*Target
// 	oneSteps := make([]*Target, 0)

// 	for _, ts := range tss {
// 		for _, t := range ts {
// 			if t.Dist == 1 {
// 				oneSteps = append(oneSteps, t)
// 			}
// 		}
// 	}
// 	if len(oneSteps) == 0 {
// 		return selected
// 	}
// 	selected = make([]*Target, 0)
// 	if len(oneSteps) > 1 {
// 		sort.Slice(oneSteps, func(i, j int) bool {
// 			return oneSteps[i].Point > oneSteps[j].Point
// 		})
// 		cache := make(map[int]bool, len(oneSteps)*2)
// 		for _, t := range oneSteps {
// 			picked, ok := cache[t.Vertex]
// 			if ok && picked {
// 				continue
// 			}
// 			selected = append(selected, t)
// 			cache[t.Vertex] = true
// 		}
// 	} else {
// 		selected = oneSteps
// 	}
// 	return selected
// }

func filterTargets(ts []*Target, excludes []*Target) []*Target {
	filtered := make([]*Target, 0)
	for _, t := range ts {
		if t == nil {
			continue
		}

	loop:
		for _, e := range excludes {
			if e == nil {
				continue loop
			}
			if t.Vertex != e.Vertex {
				filtered = append(filtered, t)
			}
		}
	}
	return filtered
}

func filterPlayers(ps []*Player, excludes []*Target) []*Player {
	filtered := make([]*Player, 0)
	for _, p := range ps {
		if p == nil {
			continue
		}
	loop:
		for _, e := range excludes {
			if e == nil || e.Subject == nil { // todo
				continue loop
			}

			if p.ID != e.Subject.ID {
				filtered = append(filtered, p)
			}
		}
	}
	return filtered
}

func (g *Guide) calcEat(roundID int, ourPlayers []*Player, rivalPlayers []*Player, powers []*Power) []*PlayerAction {
	targets := make([]*Target, 0)
	ourPowerTargets := g.calcPowerTargets(ourPlayers, powers)
	ourRivalTargets := g.calcPlayerTargets(ourPlayers, rivalPlayers)
	oneSteps := g.calcOneStepEatTargets(ourPowerTargets, ourRivalTargets)
	if len(oneSteps) != 0 {
		targets = append(targets, oneSteps...)
		ourPowerTargets = filterTargets(ourPowerTargets, oneSteps)
		ourRivalTargets = filterTargets(ourRivalTargets, oneSteps)
		ourPlayers = filterPlayers(ourPlayers, oneSteps)
		rivalPlayers = filterPlayers(rivalPlayers, oneSteps)
	}
	// hunt power
	selected, ourPlayers := g.selectPowerTargets(ourPlayers, ourPowerTargets)
	targets = append(targets, selected...)
	if len(ourPlayers) > 1 && len(rivalPlayers) > 0 { // hunt prey
		ourRivalTargets := g.calcPlayerTargets(ourPlayers, rivalPlayers)
		selected = g.selectPlayerTargets(ourPlayers, ourRivalTargets)
		targets = append(targets, selected...)
	} else {
		ourVisionTargets := g.calcVisionTargets(ourPlayers)
		selected = g.selectVisionTargets(ourPlayers, ourVisionTargets)
		targets = append(targets, selected...)
	}

	return g.generateActions(targets)
}

func (g *Guide) calcEscape(roundID int, ourPlayers []*Player, rivalPlayers []*Player, powers []*Power) []*PlayerAction {
	targets := make([]*Target, 0)
	var actions []*PlayerAction
	ourPowerTargets := g.calcPowerTargets(ourPlayers, powers)
	// ourRivalTargets := g.calcPlayerTargets(ourPlayers, rivalPlayers)
	rivalPowerTargets := g.calcPowerTargets(rivalPlayers, powers)
	ourPowerTargets = g.calcEscapeePowerTargets(ourPowerTargets, rivalPowerTargets)
	selected, ourPlayers := g.selectEscapeePowerTargets(ourPlayers, ourPowerTargets)

	if len(selected) != 0 {
		targets = append(targets, selected...)
	}

	// escape
	if len(rivalPlayers) <= 1 {
		ourVisionTargets := g.calcVisionTargets(ourPlayers)
		selected = g.selectVisionTargets(ourPlayers, ourVisionTargets)
		targets = append(targets, selected...)
		return g.generateActions(targets)
	}
	actions = g.generateActions(targets)
	for _, p := range ourPlayers {
		actions = append(actions, &PlayerAction{
			Team:   p.TeamID,
			Player: p.ID,
			Move:   randomMove(),
		})
	}

	return actions
}
func randomMove() []string {
	m := allMoves[rand.Intn(len(allMoves))]
	if m == stay {
		return []string{}
	}
	return []string{m}
}

func (g *Guide) calcEscapeePowerTargets(ourPowerTargets, rivalPowerTargets []*Target) []*Target {
	targets := make([]*Target, 0)
	for _, ot := range ourPowerTargets {
		picked := true
	loop_rival:
		for _, rt := range rivalPowerTargets {
			if ot.Vertex == rt.Vertex {
				continue loop_rival
			}
			if ot.Dist >= rt.Dist {
				picked = false
			}
		}
		if picked {
			targets = append(targets, ot)
		}
	}

	return targets
}

func (g *Guide) selectEscapeePowerTargets(players []*Player, targets []*Target) ([]*Target, []*Player) {
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

func (g *Guide) visit(ourPlayers []*Player) {
	for _, p := range ourPlayers {
		g.MapVision.Visit(p.X, p.Y)
	}
}

// Calc is
func (g *Guide) Calc(roundID int, hunting bool, ourPlayers []*Player, rivalPlayers []*Player, powers []*Power) *Action {
	var actions []*PlayerAction
	g.visit(ourPlayers)

	halfRound := DefaultRoundNum / 2
	preRound := 5
	if hunting {
		if roundID < halfRound-preRound {
			actions = g.calcEat(roundID, ourPlayers, rivalPlayers, powers) // leg 1-1
		} else if roundID < halfRound {
			actions = g.calcEscape(roundID, ourPlayers, rivalPlayers, powers) // leg 1-2
		} else {
			actions = g.calcEat(roundID, ourPlayers, rivalPlayers, powers) // leg 2-3
		}
	} else {
		if roundID < halfRound-preRound {
			actions = g.calcEscape(roundID, ourPlayers, rivalPlayers, powers) // leg 2-1
		} else if roundID < halfRound {
			actions = g.calcEscape(roundID, ourPlayers, rivalPlayers, powers) // leg 2-2
		} else {
			actions = g.calcEscape(roundID, ourPlayers, rivalPlayers, powers) // leg 1-3
		}
	}
	return &Action{
		ID:      roundID,
		Actions: actions,
	}
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
				Dist:    dist,
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
		dist := rivalTotalDist / pl
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

func (g *Guide) shortestAreaPointDistance(pv int, area *MapArea) (vertex, dist int) {
	for _, av := range area.Points {
		d := g.ShortestDistance(pv, av)
		if d != InfDist {
			vertex = av
			dist = d
			return
		}
	}
	return area.V, InfDist
}

func (g *Guide) calcVisionTargets(players []*Player) []*Target {
	areas := g.MapVision.BlindAreas()
	targets := make([]*Target, 0)
	for _, p := range players {
		for _, a := range areas {
			vertex, dist := g.shortestAreaPointDistance(p.V, a)
			value := float64(a.Count) / float64(dist)
			target := &Target{
				Vertex:  vertex,
				Subject: p,
				Point:   0,
				Dist:    dist,
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

func (g *Guide) selectVisionTargets(players []*Player, targets []*Target) []*Target {
	selected := make([]*Target, 0)
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
		}
	}
	return selected
}

func (g *Guide) generateActions(targets []*Target) []*PlayerAction {
	l := len(targets)
	w, h := g.Map.Width, g.Map.Height
	actions := make([]*PlayerAction, l, l)
	for k, t := range targets {
		i := t.Subject.V
		path := g.ShortestPath(i, t.Vertex)
		j := path[0]
		dir := CalcMove(w, h, i, j)
		move := []string{}
		if dir != stay {
			move = []string{dir}
		}
		actions[k] = &PlayerAction{
			Team:   t.Subject.TeamID,
			Player: t.Subject.ID,
			Move:   move,
		}
	}
	return actions
}
