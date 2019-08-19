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
		Meteors: []*Meteor{
			&Meteor{1, 1},
			&Meteor{1, 4},
		},
		Tunnels: []*Tunnel{
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
		Wormholes: []*Wormhole{
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

func battleWithStartedLeg(teamID int, battleTime time.Time) *Battle {
	battle := NewBattle(teamID, battleTime)
	legStart := &LegStart{
		Map:   mapExample(),
		Teams: teamsStartLegExample(),
	}
	battle.StartLeg(legStart)
	return battle
}

func legStartExample() *LegStart {
	return &LegStart{
		Map:   mapExample(),
		Teams: teamsStartLegExample(),
	}
}

func legEndExample() *LegEnd {
	return &LegEnd{
		Teams: []*Team{
			&Team{
				ID:    100,
				Point: 800,
			},
			&Team{
				ID:    200,
				Point: 400,
			},
		},
	}
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

func TestBattle_StartLeg(t *testing.T) {
	battleTime := time.Now()
	sortedTeams := teamsLegExample()
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
			name: "StartLeg - ok",
			fields: fields{
				TeamID: 100,
				Time:   battleTime,
			},
			args: args{
				legStart: legStartExample(),
			},
			want: &Leg{
				Index: 0,
				Map:   mapExample(),
				Teams: sortedTeams,
				TeamMap: map[int]*Team{
					sortedTeams[0].ID: sortedTeams[0],
					sortedTeams[1].ID: sortedTeams[1],
				},
				Players: map[int]*Player{
					1: &Player{TeamID: 100, ID: 1},
					2: &Player{TeamID: 100, ID: 2},
					3: &Player{TeamID: 100, ID: 3},
					4: &Player{TeamID: 100, ID: 4},
					5: &Player{TeamID: 200, ID: 5},
					6: &Player{TeamID: 200, ID: 6},
					7: &Player{TeamID: 200, ID: 7},
					8: &Player{TeamID: 200, ID: 8},
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
			got := b.StartLeg(tt.args.legStart)
			// fmt.Printf("%v\n", got.JSON())
			// fmt.Printf("%v\n", tt.want.JSON())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Battle.StartLeg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBattle_EndLeg(t *testing.T) {
	battleTime := time.Now()
	sortedTeams := teamsLegExample()
	battle := battleWithStartedLeg(100, battleTime)
	sortedTeams[0].Point = 800
	sortedTeams[1].Point = 400
	type fields struct {
		TeamID  int
		Time    time.Time
		Teams   []*Team
		Legs    []*Leg
		Current *Leg
	}
	type args struct {
		legEnd *LegEnd
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Leg
		wantErr bool
	}{
		{
			name: "EndLeg - ok",
			fields: fields{
				TeamID:  100,
				Time:    battleTime,
				Teams:   battle.Teams,
				Legs:    battle.Legs,
				Current: battle.Current,
			},
			args: args{
				legEnd: legEndExample(),
			},
			want: &Leg{
				Index: 0,
				Map:   mapExample(),
				Teams: sortedTeams,
				TeamMap: map[int]*Team{
					sortedTeams[0].ID: sortedTeams[0],
					sortedTeams[1].ID: sortedTeams[1],
				},
				Players: map[int]*Player{
					1: &Player{TeamID: 100, ID: 1},
					2: &Player{TeamID: 100, ID: 2},
					3: &Player{TeamID: 100, ID: 3},
					4: &Player{TeamID: 100, ID: 4},
					5: &Player{TeamID: 200, ID: 5},
					6: &Player{TeamID: 200, ID: 6},
					7: &Player{TeamID: 200, ID: 7},
					8: &Player{TeamID: 200, ID: 8},
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
			if err := b.EndLeg(tt.args.legEnd); (err != nil) != tt.wantErr {
				t.Errorf("Battle.EndLeg() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(b.Current, tt.want) {
				t.Errorf("Battle.EndLeg() = %v, want %v", b.Current, tt.want)
			}
		})
	}
}
