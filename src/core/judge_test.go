package core

import (
	"reflect"
	"testing"
)

func mockJudge() *Judge {
	return &Judge{
		Width:     20,
		Height:    20,
		Vision:    4,
		LegNum:    2,
		RoundNum:  300,
		PlayerNum: 4,
		MapData:   Map1,
	}
}

func mockMap() *Map {
	m, _ := NewMapFromString(Map1)
	m.Vision = 4
	return m
}

func mockTeamBattles() (tb1, tb2 TeamBattle) {
	tb1 = NewTeamImpl("guest team")
	tb1.SetTeamID(100)
	tb2 = NewTeamImpl("ai team")
	tb2.SetTeamID(200)
	return
}

func TestJudge_NewBattle(t *testing.T) {
	j := mockJudge()
	tb1, tb2 := mockTeamBattles()
	type fields struct {
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
	}
	type args struct {
		teamBattle1, teamBattle2 TeamBattle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *JudgeBattle
	}{
		{
			name: "NewBattle - ok",
			fields: fields{
				TeamSeq:   j.TeamSeq,
				Width:     j.Width,
				Height:    j.Height,
				Vision:    j.Vision,
				LegNum:    j.LegNum,
				RoundNum:  j.RoundNum,
				PlayerNum: j.PlayerNum,
				MapData:   j.MapData,
			},
			args: args{
				teamBattle1: tb1,
				teamBattle2: tb2,
			},
			want: &JudgeBattle{
				Judge:       j,
				TeamBattles: [2]TeamBattle{tb1, tb2},
				TeamsPlayers: [2][]int{
					[]int{1, 2, 3, 4},
					[]int{5, 6, 7, 8},
				},
				Map: mockMap(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Judge{
				TeamSeq:     tt.fields.TeamSeq,
				Width:       tt.fields.Width,
				Height:      tt.fields.Height,
				Vision:      tt.fields.Vision,
				LegNum:      tt.fields.LegNum,
				RoundNum:    tt.fields.RoundNum,
				PlayerNum:   tt.fields.PlayerNum,
				PlayerLives: tt.fields.PlayerLives,
				MapData:     tt.fields.MapData,
			}
			if got := j.NewBattle(tt.args.teamBattle1, tt.args.teamBattle2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Judge.NewBattle() = %v, want %v", got, tt.want)
			}
		})
	}
}
