package core

import (
	"reflect"
	"testing"
	"time"
)

func teamsStartLegExample() []*Team {
	return []*Team{
		&Team{
			ID:      200,
			Name:    "ai",
			Players: []int{5, 6, 7, 8},
			Force:   "think",
		},
		&Team{
			ID:      100,
			Name:    "daolaji",
			Players: []int{1, 2, 3, 4},
			Force:   "beat",
		},
	}
}

func teamsLegExample() []*Team {
	return []*Team{
		&Team{
			ID:      100,
			Name:    "daolaji",
			Players: []int{1, 2, 3, 4},
			Force:   "beat",
		},
		&Team{
			ID:      200,
			Name:    "ai",
			Players: []int{5, 6, 7, 8},
			Force:   "think",
		},
	}
}

func mapExample() *Map {
	m := &Map{
		Width:  15,
		Height: 15,
		Vision: 3,
		Meteor: []*Meteor{
			&Meteor{1, 1},
			&Meteor{1, 4},
		},
		Tunnel: []*Tunnel{
			&Tunnel{
				X:         3,
				Y:         1,
				Direction: "down",
			},
			&Tunnel{
				X:         3,
				Y:         4,
				Direction: "down",
			},
		},
		Wormhole: []*Wormhole{
			&Wormhole{
				Name: "a",
				X:    4,
				Y:    1,
			},
			&Wormhole{
				Name: "b",
				X:    4,
				Y:    4,
			},
		},
	}
	return m
}

func TestNewBattle(t *testing.T) {
	battleTime := time.Now()
	type args struct {
		teamID     int
		battleTime time.Time
	}
	tests := []struct {
		name string
		args args
		want *Battle
	}{
		{
			name: "NewBattle - ok",
			args: args{
				teamID:     100,
				battleTime: battleTime,
			},
			want: &Battle{
				TeamID: 100,
				Time:   battleTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBattle(tt.args.teamID, tt.args.battleTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBattle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBattle_NewLeg(t *testing.T) {
	battleTime := time.Now()
	type fields struct {
		TeamID  int
		Time    time.Time
		Teams   []*Team
		Legs    []*Leg
		Current *Leg
	}
	type args struct {
		legStart *LegStart
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Leg
	}{
		{
			name: "NewLeg - ok",
			fields: fields{
				TeamID: 100,
				Time:   battleTime,
			},
			args: args{
				legStart: &LegStart{
					Map:   mapExample(),
					Teams: teamsStartLegExample(),
				},
			},
			want: &Leg{
				Index: 0,
				Map:   mapExample(),
				Teams: teamsLegExample(),
				Players: map[int]*Player{
					1: &Player{Team: 100, ID: 1},
					2: &Player{Team: 100, ID: 2},
					3: &Player{Team: 100, ID: 3},
					4: &Player{Team: 100, ID: 4},
					5: &Player{Team: 200, ID: 5},
					6: &Player{Team: 200, ID: 6},
					7: &Player{Team: 200, ID: 7},
					8: &Player{Team: 200, ID: 8},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Battle{
				TeamID:  tt.fields.TeamID,
				Time:    tt.fields.Time,
				Teams:   tt.fields.Teams,
				Legs:    tt.fields.Legs,
				Current: tt.fields.Current,
			}
			got := b.NewLeg(tt.args.legStart)
			// fmt.Printf("%v\n", got.JSON())
			// fmt.Printf("%v\n", tt.want.JSON())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Battle.NewLeg() = %v, want %v", got, tt.want)
			}
		})
	}
}
