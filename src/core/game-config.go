package core

import (
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

// GameConfig is
type GameConfig struct {
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
func NewFirstGame(mapData string, vision, width, height int) *GameConfig {
	g := &GameConfig{
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
	return g
}

// Init is
func (g *GameConfig) Init() error {
	// init players' initial location orders at each leg's beginning
	rand.Seed(time.Now().Unix())
	g.PlayerOrders = rand.Perm(g.PlayerNum)

	// init map
	m, err := NewMapFromString(g.MapData)
	if err != nil {
		log.Println(err)
		return err
	}
	err = m.Init(g.Vision, g.Width, g.Height)
	if err != nil {
		log.Println(err)
		return err
	}
	g.Map = m

	return nil
}

// NewTeamID is
func (g *GameConfig) NewTeamID() int {
	g.TeamSeq++
	return g.TeamSeq
}

// NewBattle is
func (g *GameConfig) NewBattle(teamBattle1, teamBattle2 TeamBattle) *GameBattle {
	b := &GameBattle{
		Config: g,
	}
	b.TeamBattles[0] = teamBattle1
	b.TeamBattles[1] = teamBattle2
	b.init()
	return b
}
