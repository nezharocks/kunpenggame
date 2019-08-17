package core

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

const (
	// TeamNum is
	TeamNum = 2

	// DefaultLegNum is
	DefaultLegNum = 2

	// DefaultLegModeNum is
	DefaultLegModeNum = 2

	// DefaultRoundNum is
	DefaultRoundNum = 2

	// DefaultPlayerNum is
	DefaultPlayerNum = 4

	// DefaultPlayerLives is
	DefaultPlayerLives = 8

	// ActionTimeout is
	ActionTimeout = 800
)

// DefaultTeamForces is teams' force settings
var DefaultTeamForces = [TeamNum]TeamForce{BeatForce, ThinkForce}

// DefaultBattleModes is battle's mode settings
var DefaultBattleModes = [DefaultLegNum][DefaultLegModeNum]BattleMode{
	{FireMode, WaterMode},
	{WaterMode, FireMode},
}

// TeamPlaceHolderIndices holds the index matrix of place holders per team per leg
var TeamPlaceHolderIndices = [DefaultLegNum][TeamNum]int{
	{0, 1},
	{1, 0},
}

// Judge is
type Judge struct {
	TeamSeq       int
	Width, Height int
	Vision        int
	LegNum        int
	RoundNum      int
	PlayerNum     int
	PlayerLives   int
	TeamForces    [TeamNum]TeamForce
	BattleModes   [][DefaultLegModeNum]BattleMode
	MapData       string
	Map           *Map
	PlayerOrders  []int
}

// NewFirstGame is
func NewFirstGame(mapData string, vision, width, height int) *Judge {
	j := &Judge{
		TeamSeq:     0,
		Width:       width,
		Height:      height,
		Vision:      vision,
		LegNum:      DefaultLegNum,
		RoundNum:    DefaultRoundNum,
		PlayerNum:   DefaultPlayerNum,
		PlayerLives: DefaultPlayerLives,
		TeamForces:  DefaultTeamForces,
		BattleModes: DefaultBattleModes[:],
		MapData:     mapData,
	}
	return j
}

// Init is
func (j *Judge) Init() error {
	// init players' initial location orders at each leg's beginning
	rand.Seed(time.Now().Unix())
	j.PlayerOrders = rand.Perm(j.PlayerNum)

	// init map
	m, err := NewMapFromString(j.MapData)
	if err != nil {
		log.Println(err)
		return err
	}
	err = m.Init(j.Vision, j.Width, j.Height)
	if err != nil {
		log.Println(err)
		return err
	}
	j.Map = m

	return nil
}

// NewTeamID is
func (j *Judge) NewTeamID() int {
	j.TeamSeq++
	return j.TeamSeq
}

// NewBattle is
func (j *Judge) NewBattle(teamBattle1, teamBattle2 TeamBattle) *JudgeBattle {
	b := &JudgeBattle{
		Judge: j,
	}
	b.TeamBattles[0] = teamBattle1
	b.TeamBattles[1] = teamBattle2
	b.init()
	return b
}

// JudgeBattle is
type JudgeBattle struct {
	Judge        *Judge
	TeamBattles  [TeamNum]TeamBattle
	Teams        [TeamNum]*Team
	TeamsPlayers [TeamNum][]int
	Map          *Map
	Legs         []*JudgeBattleLeg
}

// Run is
func (b *JudgeBattle) Run() {
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
func (b *JudgeBattle) init() {
	b.Map = b.Judge.Map
	// player := b.Map.TeamPlaceHolders[1][0]
	// player.X = 1
	// player.Y = 8

	b.initTeams()
	b.initTeamPlayers()
	b.initLegs()
}

func (b *JudgeBattle) initTeams() {
	for i := 0; i < TeamNum; i++ {
		b.Teams[i] = &Team{
			ID:   b.TeamBattles[i].GetTeamID(),
			Name: b.TeamBattles[i].GetTeamName(),
		}
	}
}

func (b *JudgeBattle) initTeamPlayers() {
	id := 0
	n := b.Judge.PlayerNum

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

func (b *JudgeBattle) newLeg(index int) *JudgeBattleLeg {
	leg := &JudgeBattleLeg{
		Battle: b,
		Index:  index,
	}

	// init teams
	for i := 0; i < TeamNum; i++ {
		force := b.Judge.TeamForces[i]
		leg.Teams[i] = &Team{
			ID:         b.TeamBattles[i].GetTeamID(),
			Players:    b.TeamsPlayers[i],
			Force:      force.String(),
			RemainLife: b.Judge.PlayerLives - b.Judge.PlayerNum,
		}
	}

	// init teams' players
	orders := b.Judge.PlayerOrders
	playerNum := b.Judge.PlayerNum
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
	leg.PlayersMap = make(map[TeamBattle][]*Player, b.Judge.PlayerNum)
	for t := 0; t < TeamNum; t++ {
		tb := b.TeamBattles[t]
		leg.TeamMap[tb] = leg.Teams[t]
		leg.PlayersMap[tb] = leg.TeamsPlayers[t]
	}

	leg.Table = NewTable(b.Map, b.TeamBattles, leg.TeamMap, leg.PlayersMap)
	return leg
}

func (b *JudgeBattle) newPlayer(teamID, playerID int, placeholder *PlaceHolder) *Player {
	player := &Player{
		Team: teamID,
		ID:   playerID,
		X:    placeholder.X,
		Y:    placeholder.Y,
	}
	return player
}

func (b *JudgeBattle) initLegs() {
	n := b.Judge.LegNum
	b.Legs = make([]*JudgeBattleLeg, n, n)
	for i := 0; i < n; i++ {
		b.Legs[i] = b.newLeg(i)
	}
}

// JudgeBattleLeg is
type JudgeBattleLeg struct {
	Battle       *JudgeBattle `json:"-"`
	Index        int
	Teams        [TeamNum]*Team
	TeamsPlayers [TeamNum][]*Player
	TeamMap      map[TeamBattle]*Team
	PlayersMap   map[TeamBattle][]*Player
	Table        *Table `json:"-"`
}

// JSON is
func (l *JudgeBattleLeg) JSON() string {
	bytes, err := json.MarshalIndent(l, "", "    ")
	// bytes, err := json.Marshal(l)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// getPowerForce is
// index is mode index in the same leg
func (l *JudgeBattleLeg) getPowerForce(index int) TeamForce {
	mode := l.Battle.Judge.BattleModes[l.Index][index]
	return mode.PowerForce()
}

// getTeamBattles is
func (l *JudgeBattleLeg) getTeamBattles(powerForce TeamForce) (escapee, hunter TeamBattle) {
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
func (l *JudgeBattleLeg) Run() {
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
	fullRound := l.Battle.Judge.RoundNum
	semiRound := l.Battle.Judge.RoundNum / 2

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
			if l.Action(&action) {
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
			if l.Action(&action) {
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
func (l *JudgeBattleLeg) Round(index int, tb TeamBattle, powerForce TeamForce) *Round {
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
	visions := make([]*Vision, 0, l.Battle.Judge.PlayerNum)
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
func (l *JudgeBattleLeg) Action(action *Action) bool {
	return false
}
