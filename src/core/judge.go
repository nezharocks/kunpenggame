package core

import (
	"log"
)

// Judge is
type Judge struct {
	teamSeq       int
	Width, Height int
	Vision        int
	LegNum        int
	RoundNum      int
	PlayerNum     int
	FirstMode     ForceMode
	MapData       string
}

// JudgeBattle is
type JudgeBattle struct {
	Judge        *Judge
	TeamID1      int
	TeamID2      int
	TeamsPlayers [][]int
	Map          *Map
	CurrentLeg   int
	CurrentRound int
}

var legModeSequences = []ForceMode{BeatMode, ThinkMode}

// NewTeamID is
func (j *Judge) NewTeamID() int {
	j.teamSeq++
	return j.teamSeq
}

// NewBattle is
func (j *Judge) NewBattle(teamID1, teamID2 int) *JudgeBattle {
	b := &JudgeBattle{
		Judge:   j,
		TeamID1: teamID1,
		TeamID2: teamID2,
	}
	b.init()
	return b
}

// Init is
func (b *JudgeBattle) init() {
	b.TeamsPlayers = generateTeamPlayers(b.Judge.PlayerNum)
	b.initMap()
	b.initPlayerLocations()

}

// func (b *JudgeBattle) newLegTeam(teamID, n int, force string) *Team {
// 	team := &Team{
// 		ID:      teamID,
// 		Force:   force,
// 	}

// 	for i := 0; i < n; i++ {
// 		team.Players[i] = g.NewPlayerID()

// 	}
// 	return team
// }

func (b *JudgeBattle) initMap() {
	m, err := NewMapFromString(b.Judge.MapData)
	if err != nil {
		log.Println(err)
		return
	}
	m.Vision = b.Judge.Vision
	b.Map = m
}

func (b *JudgeBattle) initPlayerLocations() {
	mode := legModeSequences[b.CurrentLeg]

}

func generateTeamPlayers(n int) [][]int {
	id := 0
	teams := make([][]int, 2, 2)

	// generate players for the first team
	players := make([]int, n, n)
	for i := 0; i < n; i++ {
		id++
		players[i] = id
	}
	teams[0] = players

	// generate players for the second team
	players = make([]int, n, n)
	for i := 0; i < n; i++ {
		id++
		players[i] = id
	}
	teams[1] = players
	return teams
}
