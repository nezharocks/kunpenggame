package core

import (
	"fmt"
	"math/rand"
	"time"
)

// BattleAlg is
type BattleAlg struct {
	TeamID  int
	Time    time.Time
	Teams   []*Team
	Legs    []*LegAlg
	Current *LegAlg
}

// NewBattleAlg is
func NewBattleAlg(teamID int, battleTime time.Time) *BattleAlg {
	rand.Seed(time.Now().Unix())
	return &BattleAlg{
		TeamID: teamID,
		Time:   battleTime,
	}
}

// LegAlg is
type LegAlg struct {
	Index   int
	Map     *Map
	Guide   *Guide
	Teams   []*Team
	TeamMap map[int]*Team
	Players map[int]*Player
	Rounds  []*LegRound
	Current *LegRound
}

// StartLeg is
func (b *BattleAlg) StartLeg(legStart *LegStart) *LegAlg {
	guide := NewGuide(legStart.Map)
	leg := &LegAlg{
		Index:   len(b.Legs),
		Map:     legStart.Map,
		Guide:   guide,
		Teams:   make([]*Team, 0, 2),
		TeamMap: make(map[int]*Team, 2),
	}

	// ensure the current (guest) team is the first item in leg.Teams slice
	otherTeams := make([]*Team, 0)
	for _, team := range legStart.Teams {
		leg.TeamMap[team.ID] = team
		if team.ID == b.TeamID {
			leg.Teams = append(leg.Teams, team)
		} else {
			otherTeams = append(otherTeams, team)
		}
	}
	leg.Teams = append(leg.Teams, otherTeams...)

	// on the coming first leg, update battle's teams' basic information
	// including id, name and players
	if leg.Index == 0 {
		teamCount := len(leg.Teams)
		b.Teams = make([]*Team, teamCount, teamCount)
		for i := 0; i < teamCount; i++ {
			t := leg.Teams[i]
			b.Teams[i] = &Team{
				ID:      t.ID,
				Players: t.Players,
				Force:   t.Force,
			}
		}
	}

	// // init leg players' map
	leg.Players = make(map[int]*Player, 8)
	for _, t := range leg.Teams {
		for _, p := range t.Players {
			leg.Players[p] = &Player{
				TeamID: t.ID,
				ID:     p,
			}
		}
	}

	b.Legs = append(b.Legs, leg)
	b.Current = leg
	return leg
}

// EndLeg is
func (b *BattleAlg) EndLeg(legEnd *LegEnd) error {
	if b.Current == nil {
		return fmt.Errorf("battle state illegal - the battle has no any started legs")
	}
	for _, t1 := range legEnd.Teams {
		t2, ok := b.Current.TeamMap[t1.ID]
		if !ok {
			continue
		}
		t2.Point = t1.Point
	}
	return nil
}

// AddRound is
func (b *BattleAlg) AddRound(round *Round) error {
	if b.Current == nil {
		return fmt.Errorf("battle state illegal - the battle has no any started legs")
	}
	leg := b.Current
	legRound := &LegRound{Round: round}
	leg.Rounds = append(leg.Rounds, legRound)
	leg.Current = legRound

	// update points and remain lives each team
	for _, t1 := range round.Teams {
		t2, ok := leg.TeamMap[t1.ID]
		if !ok {
			continue
		}
		t2.Point = t1.Point
		t2.RemainLife = t1.RemainLife
	}

	// update players' points, sleep and location in current view
	for _, p1 := range round.Players {
		p2, ok := leg.Players[p1.ID]
		if !ok {
			continue
		}
		p2.Point = p1.Point
		p2.Sleep = p1.Sleep
		p2.X = p1.X
		p2.Y = p1.Y
	}
	return nil
}

// CalcAction is
// 0.1 - first version which can work
func (b *BattleAlg) CalcAction() (*Action, error) {
	leg := b.Current
	round := leg.Current.Round
	hunting := leg.Teams[0].Force == round.Mode
	// if hunting {
	// 	fmt.Println("hunt round", round.ID)
	// } else {
	// 	fmt.Println("escape round", round.ID)
	// }

	w := leg.Map.Width
	for _, p := range round.Powers {
		p.V = p.Y*w + p.X
	}
	for _, p := range round.Players {
		p.V = p.Y*w + p.X
	}

	ourPlayers := make([]*Player, 0, DefaultPlayerNum)
	rivalPlayers := make([]*Player, 0, DefaultPlayerNum)
	for _, p := range round.Players {
		if p.TeamID == b.TeamID {
			ourPlayers = append(ourPlayers, p)
		} else {
			rivalPlayers = append(rivalPlayers, p)
		}
		// todo check it out
		// if p.TeamID == b.TeamID && p.Sleep == 0 {
		// 	ourPlayers = append(ourPlayers, p)
		// }
	}
	action := leg.Guide.Calc(round.ID, hunting, ourPlayers, rivalPlayers, round.Powers)
	return action, nil
}
