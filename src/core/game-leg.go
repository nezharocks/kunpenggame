package core

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"
)

// GameBattleLeg is
type GameBattleLeg struct {
	Battle      *GameBattle `json:"-"`
	Index       int
	Teams       [TeamNum]*Team
	TeamPlayers [TeamNum][]*Player
	IDPlayers   map[int]*Player
	Table       *Table `json:"-"`
}

// JSON is
func (l *GameBattleLeg) JSON() string {
	bytes, _ := json.MarshalIndent(l, "", "    ")
	return string(bytes)
}

// Run is
func (l *GameBattleLeg) Run() {
	var (
		err     error
		escapee TeamBattle
		hunter  TeamBattle
		round   *Round
	)
	partIndex := 0 // the first part of the leg
	powerForce := l.Battle.GetPowerForce(l.Index, partIndex)
	escapee, hunter = l.Battle.GetEscapeeHunter(powerForce)

	// send leg starts
	legStart := &LegStart{
		Map:   l.Battle.Map,
		Teams: l.Teams[:], // todo
	}
	log.Printf("%+v\n", legStart.Message())
	err = escapee.LegStart(legStart)
	if err != nil {
		log.Println(err)
	}
	err = hunter.LegStart(legStart)
	if err != nil {
		log.Println(err)
	}
	fullRound := l.Battle.Config.RoundNum
	semiRound := l.Battle.Config.RoundNum / 2

	// handle rounds and actions in the leg
loop_rounds:
	for rid := 0; rid < fullRound; rid++ {
		if rid == semiRound {
			partIndex++ // the second part of the leg
			powerForce := l.Battle.GetPowerForce(l.Index, partIndex)
			escapee, hunter = l.Battle.GetEscapeeHunter(powerForce)
		}
		// escapee action
		round = l.Round(rid, escapee.GetTeamID(), powerForce)
		if debugRound {
			log.Printf("%+v\n", round.Message())
		}
		err = escapee.Round(round)
		if err != nil {
			log.Println(err)
		}

		select {
		case action := <-escapee.GetActionCh():
			if debugRound {
				log.Printf("%+v\n", action.Message())
			}
			if l.Action(&action, rid, escapee.GetTeamID(), powerForce) {
				break loop_rounds
			}
		case <-time.After(time.Millisecond * ActionTimeout):
			log.Printf("team %v timeout at the %vth round of the %vth leg\n", escapee.GetTeamID(), rid, l.Index)
		}

		// hunter action
		round = l.Round(rid, hunter.GetTeamID(), powerForce)
		if debugRound {
			log.Printf("%+v\n", round.Message())
		}
		err = hunter.Round(round)
		if err != nil {
			log.Println(err)
		}

		select {
		case action := <-hunter.GetActionCh():
			if debugRound {
				log.Printf("%+v\n", action.Message())
			}
			if l.Action(&action, rid, hunter.GetTeamID(), powerForce) {
				break loop_rounds
			}
		case <-time.After(time.Millisecond * ActionTimeout):
			log.Printf("team %v timeout at the %vth round of the %vth leg\n", hunter.GetTeamID(), rid, l.Index)
		}
	}

	// send leg ends
	legEnd := &LegEnd{
		Teams: l.Teams[:],
	}
	err = escapee.LegEnd(legEnd)
	if err != nil {
		log.Println(err)
	}
	err = hunter.LegEnd(legEnd)
	if err != nil {
		log.Println(err)
	}

	// update teams's point for the battle
	for i := 0; i < TeamNum; i++ {
		l.Battle.Teams[i].Point += l.Teams[i].Point
	}
}

// Round is
func (l *GameBattleLeg) Round(roundID, teamID int, powerForce TeamForce) *Round {
	r := &Round{
		ID:   roundID,
		Mode: powerForce.String(),
	}
	// calculate teams
	var roundTeams [TeamNum]*Team
	for i, t := range l.Teams {
		nt := *t
		nt.Players = nil
		nt.Force = ""
		roundTeams[i] = &nt
	}
	r.Teams = roundTeams[:]

	// calculate the visions of the active players
	activePlayers := l.ActivePlayers(teamID)
	visions := make([]*Vision, 0, l.Battle.Config.PlayerNum)
	for _, p := range activePlayers {
		v := l.Battle.Map.GetVision(p.X, p.Y)
		visions = append(visions, v)
		// fmt.Printf("%v\n", *v)
	}

	// calculate powers and players (including rival's and mime) in the visions
	r.Powers = GetVisiblePowers(l.Battle.Map.Powers, visions)
	currentPlayers := l.AlivePlayers(teamID)
	_, rivalIndex := l.Battle.RivalID(teamID)
	visibleRivalPlayers := GetVisiblePlayers(l.TeamPlayers[rivalIndex], visions)
	r.Players = append(currentPlayers, visibleRivalPlayers...)
	return r
}

// Action apply team's movements to the battle and return if the battle should be over
func (l *GameBattleLeg) Action(action *Action, roundID, teamID int, powerForce TeamForce) bool {
	actions := action.Actions
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Player < actions[j].Player
	})
	for _, a := range actions {
		move := NewMovement(a.Move)
		player, ok := l.IDPlayers[a.Player]
		if !ok {
			continue
		}
		l.Table.Move(player, powerForce, move)
		fmt.Println(a.Player)
	}

	// l.TeamPlayers
	return false
}

// ActivePlayers is
func (l *GameBattleLeg) ActivePlayers(teamID int) []*Player {
	teamIndex := l.Battle.TeamIndex(teamID)
	players := l.TeamPlayers[teamIndex]
	activePlayers := make([]*Player, 0, len(players))
	for _, p := range players {
		if !p.IsAsleep() {
			activePlayers = append(activePlayers, p)
		}
	}
	return activePlayers
}

// AlivePlayers is
func (l *GameBattleLeg) AlivePlayers(teamID int) []*Player {
	teamIndex := l.Battle.TeamIndex(teamID)
	players := l.TeamPlayers[teamIndex]
	activePlayers := make([]*Player, 0, len(players))
	for _, p := range players {
		if !p.IsDead() {
			activePlayers = append(activePlayers, p)
		}
	}
	return activePlayers
}
