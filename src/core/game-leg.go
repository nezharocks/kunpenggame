package core

import (
	"encoding/json"
	"log"
	"time"
)

// GameBattleLeg is
type GameBattleLeg struct {
	Battle       *GameBattle `json:"-"`
	Index        int
	Teams        [TeamNum]*Team
	TeamsPlayers [TeamNum][]*Player
	TeamMap      map[TeamBattle]*Team
	PlayersMap   map[TeamBattle][]*Player
	Table        *Table `json:"-"`
}

// JSON is
func (l *GameBattleLeg) JSON() string {
	bytes, err := json.MarshalIndent(l, "", "    ")
	// bytes, err := json.Marshal(l)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// getPowerForce is
// index is mode index in the same leg
func (l *GameBattleLeg) getPowerForce(index int) TeamForce {
	mode := l.Battle.Config.BattleModes[l.Index][index]
	return mode.PowerForce()
}

// getTeamBattles is
func (l *GameBattleLeg) getTeamBattles(powerForce TeamForce) (escapee, hunter TeamBattle) {
	for i, team := range l.Teams {
		if powerForce.Equal(team.Force) {
			hunter = l.Battle.TeamBattles[i]
		} else {
			escapee = l.Battle.TeamBattles[i]
		}
	}
	return
}

// Run is
func (l *GameBattleLeg) Run() {
	var (
		err     error
		escapee TeamBattle
		hunter  TeamBattle
		round   *Round
	)

	powerForce := l.getPowerForce(0) // first half in a leg
	escapee, hunter = l.getTeamBattles(powerForce)

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
	for i := 0; i < fullRound; i++ {
		if i == semiRound {
			powerForce := l.getPowerForce(1) // second half in a leg
			escapee, hunter = l.getTeamBattles(powerForce)
		}
		// escapee action
		round = l.Round(i, escapee, powerForce)
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
			if l.Action(&action, i, escapee, powerForce) {
				break loop_rounds
			}
		case <-time.After(time.Millisecond * ActionTimeout):
			log.Printf("team %v timeout at the %vth round of the %vth leg\n", escapee.GetTeamID(), i, l.Index)
		}

		// hunter action
		round = l.Round(i, hunter, powerForce)
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
			if l.Action(&action, i, hunter, powerForce) {
				break loop_rounds
			}
		case <-time.After(time.Millisecond * ActionTimeout):
			log.Printf("team %v timeout at the %vth round of the %vth leg\n", hunter.GetTeamID(), i, l.Index)
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
}

// Round is
func (l *GameBattleLeg) Round(index int, tb TeamBattle, powerForce TeamForce) *Round {
	r := &Round{
		ID:   index,
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
	activePlayers := l.Table.TeamActivePlayers(tb)
	visions := make([]*Vision, 0, l.Battle.Config.PlayerNum)
	for _, p := range activePlayers {
		v := l.Battle.Map.GetVision(p.X, p.Y)
		visions = append(visions, v)
		// fmt.Printf("%v\n", *v)
	}

	// calculate powers and players (including rival's and mime) in the visions
	r.Powers = l.Table.GetVisiblePowers(visions)
	currentPlayers := l.Table.TeamAlivePlayers(tb)
	rival := l.Table.GetRival(tb)
	rivalPlayers := l.Table.GetVisiblePlayers(l.PlayersMap[rival], visions)
	r.Players = append(currentPlayers, rivalPlayers...)
	return r
}

// Action apply team's movements to the battle and return if the battle should be over
func (l *GameBattleLeg) Action(action *Action, index int, tb TeamBattle, powerForce TeamForce) bool {
	for _, a := range action.Actions {
		a.Movement = NewMovement(a.Move)
		// fmt.Println(a.Player)
	}

	// l.TeamsPlayers
	return false
}
