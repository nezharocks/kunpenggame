package core

import (
	"log"
)

// GameBattle is
type GameBattle struct {
	Config       *GameConfig
	TeamBattles  [TeamNum]TeamBattle
	Teams        [TeamNum]*Team
	TeamsPlayers [TeamNum][]int
	Map          *Map
	Legs         []*GameBattleLeg
}

// Run is
func (b *GameBattle) Run() {
	// run legs to battle
	for _, leg := range b.Legs {
		leg.Run()
		for i := 0; i < TeamNum; i++ {
			b.Teams[i].Point += leg.Teams[i].Point
		}
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
}

// Init is
func (b *GameBattle) init() {
	b.Map = b.Config.Map
	// player := b.Map.TeamPlaceHolders[1][0]
	// player.X = 1
	// player.Y = 8

	b.initTeams()
	b.initTeamPlayers()
	b.initLegs()
}

func (b *GameBattle) initTeams() {
	for i := 0; i < TeamNum; i++ {
		b.Teams[i] = &Team{
			ID:   b.TeamBattles[i].GetTeamID(),
			Name: b.TeamBattles[i].GetTeamName(),
		}
	}
}

func (b *GameBattle) initTeamPlayers() {
	id := 0
	n := b.Config.PlayerNum

	// generate players for the first team
	players := make([]int, n, n)
	for i := 0; i < n; i++ {
		players[i] = id
		id++
	}
	b.TeamsPlayers[0] = players

	// generate players for the second team
	players = make([]int, n, n)
	for i := 0; i < n; i++ {
		players[i] = id
		id++
	}
	b.TeamsPlayers[1] = players
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
			Players:    b.TeamsPlayers[i],
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
		leg.TeamsPlayers[t] = make([]*Player, playerNum, playerNum)
		for i := 0; i < playerNum; i++ {
			playerIndex := orders[i]
			holder := holders[playerIndex%holderNum]
			teamID := b.TeamBattles[t].GetTeamID()
			playerID := b.TeamsPlayers[t][playerIndex]
			leg.TeamsPlayers[t][playerIndex] = b.newPlayer(teamID, playerID, holder)
		}
	}

	leg.TeamMap = make(map[TeamBattle]*Team, TeamNum)
	leg.PlayersMap = make(map[TeamBattle][]*Player, b.Config.PlayerNum)
	for t := 0; t < TeamNum; t++ {
		tb := b.TeamBattles[t]
		leg.TeamMap[tb] = leg.Teams[t]
		leg.PlayersMap[tb] = leg.TeamsPlayers[t]
	}

	leg.Table = NewTable(b.Map, b.TeamBattles, leg.TeamMap, leg.PlayersMap)
	return leg
}

func (b *GameBattle) newPlayer(teamID, playerID int, placeholder *PlaceHolder) *Player {
	player := &Player{
		Team: teamID,
		ID:   playerID,
		X:    placeholder.X,
		Y:    placeholder.Y,
	}
	return player
}

func (b *GameBattle) initLegs() {
	n := b.Config.LegNum
	b.Legs = make([]*GameBattleLeg, n, n)
	for i := 0; i < n; i++ {
		b.Legs[i] = b.newLeg(i)
	}
}
