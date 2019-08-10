package core

import (
	"encoding/json"
	"time"
)

// Battle is
type Battle struct {
	TeamID  int
	Time    time.Time
	Teams   []*Team
	Legs    []*Leg
	Current *Leg
}

// NewBattle is
func NewBattle(teamID int, battleTime time.Time) *Battle {
	return &Battle{
		TeamID: teamID,
		Time:   battleTime,
	}
}

// Leg is
type Leg struct {
	Index   int
	Map     *Map
	Teams   []*Team
	Players map[int]*Player
	Rounds  []*LegRound
	Current *LegRound
}

// JSON is
func (b *Leg) JSON() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}

// LegRound is
type LegRound struct {
	Round  Round
	Action Action
}

// NewLeg is
func (b *Battle) NewLeg(legStart *LegStart) *Leg {
	leg := &Leg{
		Index: len(b.Legs),
		Map:   legStart.Map,
		Teams: make([]*Team, 0, 2),
	}

	// ensure the current (guest) team is the first item in leg.Teams slice
	otherTeams := make([]*Team, 0)
	for _, team := range legStart.Teams {
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
				Team: t.ID,
				ID:   p,
			}
		}
	}

	b.Legs = append(b.Legs, leg)
	b.Current = leg
	return leg
}
