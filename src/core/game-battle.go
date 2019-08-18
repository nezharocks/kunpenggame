package core

import (
	"log"
)

// GameBattle is
type GameBattle struct {
	Config      *GameConfig
	TeamBattles [TeamNum]TeamBattle
	TeamIDs     [TeamNum]int
	Teams       [TeamNum]*Team
	TeamPlayers [TeamNum][]int
	Legs        []*GameBattleLeg
	Map         *Map
}

// Run is
func (b *GameBattle) Run() {
	// run legs to battle
	for _, leg := range b.Legs {
		leg.Run()
	}

	// send game over commands
	for _, teamBattle := range b.TeamBattles {
		if err := teamBattle.GameOver(&GameOver{}); err != nil {
			log.Println(err)
		}
	}

	// print battle results
	for i := 0; i < TeamNum; i++ {
		team := b.Teams[i]
		log.Printf("team \t%v\t%v\n", team.ID, team.Point)
	}

	// todo: save the battle results to DB
}

// Init is
func (b *GameBattle) Init() {
	b.Map = b.Config.Map
	// player := b.Map.TeamPlaceHolders[1][0]
	// player.X = 1
	// player.Y = 8

	// init teams and players
	pn, pid := b.Config.PlayerNum, 0
	for i := 0; i < TeamNum; i++ {
		b.TeamIDs[i] = b.TeamBattles[i].GetTeamID()
		b.Teams[i] = &Team{
			ID:    b.TeamBattles[i].GetTeamID(),
			Name:  b.TeamBattles[i].GetTeamName(),
			Force: b.Config.TeamForces[i].String(),
		}
		b.TeamPlayers[i] = make([]int, pn, pn)
		for j := 0; j < pn; j++ {
			b.TeamPlayers[i][j] = pid
			pid++
		}
	}
	b.newLegs()
}

func (b *GameBattle) newLegs() {
	n := b.Config.LegNum
	b.Legs = make([]*GameBattleLeg, n, n)
	for i := 0; i < n; i++ {
		b.Legs[i] = b.newLeg(i)
	}
}

func (b *GameBattle) newLeg(index int) *GameBattleLeg {
	leg := &GameBattleLeg{
		Battle: b,
		Index:  index,
	}

	// init teams
	for i := 0; i < TeamNum; i++ {
		force := b.Config.TeamForces[i]
		leg.Teams[i] = &Team{
			ID:         b.TeamBattles[i].GetTeamID(),
			Players:    b.TeamPlayers[i],
			Force:      force.String(),
			RemainLife: b.Config.PlayerLives - b.Config.PlayerNum,
		}
	}

	// init teams' players
	orders := b.Config.PlayerOrders
	playerNum := b.Config.PlayerNum
	for t := 0; t < TeamNum; t++ {
		phIndex := TeamPlaceHolderIndices[index][t]
		holders := b.Map.TeamPlaceHolders[phIndex]
		holderNum := len(holders)
		leg.TeamPlayers[t] = make([]*Player, playerNum, playerNum)
		for i := 0; i < playerNum; i++ {
			playerIndex := orders[i]
			holder := holders[playerIndex%holderNum]
			teamID := b.TeamBattles[t].GetTeamID()
			playerID := b.TeamPlayers[t][playerIndex]
			player := &Player{Team: teamID, ID: playerID, X: holder.X, Y: holder.Y}
			leg.TeamPlayers[t][playerIndex] = player
			leg.IDPlayers[playerID] = player
		}
	}
	leg.Table = NewTable(leg)
	return leg
}

// TeamIndex id
func (b *GameBattle) TeamIndex(teamID int) int {
	for i, id := range b.TeamIDs {
		if id == teamID {
			return i
		}
	}
	return -1
}

// RivalID is
func (b *GameBattle) RivalID(teamID int) (id int, index int) {
	for i, id := range b.TeamIDs {
		if id != teamID {
			return id, i
		}
	}
	return -1, -1
}

// GetPowerForce is
func (b *GameBattle) GetPowerForce(legIndex, partIndex int) TeamForce {
	return b.Config.BattleModes[legIndex][partIndex].PowerForce()
}

// GetEscapeeHunter is
func (b *GameBattle) GetEscapeeHunter(powerForce TeamForce) (escapee, hunter TeamBattle) {
	for i, team := range b.Teams {
		if powerForce.Equal(team.Force) {
			hunter = b.TeamBattles[i]
		} else {
			escapee = b.TeamBattles[i]
		}
	}
	return
}
