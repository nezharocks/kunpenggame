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
		FirstMode: BeatMode,
		MapData:   map1,
	}
}

func mockMap() *Map {
	m, _ := NewMapFromString(map1)
	m.Vision = 4
	return m
}

func TestJudge_NewBattle(t *testing.T) {
	j := mockJudge()
	type fields struct {
		teamSeq   int
		Width     int
		Height    int
		Vision    int
		LegNum    int
		RoundNum  int
		PlayerNum int
		FirstMode ForceMode
		MapData   string
	}
	type args struct {
		teamID1 int
		teamID2 int
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
				teamSeq:   j.teamSeq,
				Width:     j.Width,
				Height:    j.Height,
				Vision:    j.Vision,
				LegNum:    j.LegNum,
				RoundNum:  j.RoundNum,
				PlayerNum: j.PlayerNum,
				FirstMode: j.FirstMode,
				MapData:   j.MapData,
			},
			args: args{
				teamID1: 100,
				teamID2: 200,
			},
			want: &JudgeBattle{
				Judge:   j,
				TeamID1: 100,
				TeamID2: 200,
				TeamsPlayers: [][]int{
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
				teamSeq:   tt.fields.teamSeq,
				Width:     tt.fields.Width,
				Height:    tt.fields.Height,
				Vision:    tt.fields.Vision,
				LegNum:    tt.fields.LegNum,
				RoundNum:  tt.fields.RoundNum,
				PlayerNum: tt.fields.PlayerNum,
				FirstMode: tt.fields.FirstMode,
				MapData:   tt.fields.MapData,
			}
			if got := j.NewBattle(tt.args.teamID1, tt.args.teamID2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Judge.NewBattle() = %v, want %v", got, tt.want)
			}
		})
	}
}
