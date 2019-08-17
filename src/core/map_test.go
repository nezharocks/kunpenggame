package core

import (
	"fmt"
	"reflect"
	"testing"
)

func createMap(mapData string) *Map {
	m, err := NewMapFromString(mapData)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = m.Init(defaultVision, defaultWidth, defaultHeight)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return m
}

func mockMap1() *Map {
	return createMap(Map1)
}

func TestNewMapFromStrTing(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantM   *Map
		wantErr bool
	}{
		{
			name: "TestNewMapFromString - Map1",
			args: args{
				data: Map1,
			},
			wantErr: false,
		},
		{
			name: "TestNewMapFromString - Map2",
			args: args{
				data: Map2,
			},
			wantErr: false,
		},
		{
			name: "TestNewMapFromString - Map3",
			args: args{
				data: Map3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := NewMapFromString(tt.args.data)
			fmt.Printf("%v\n\n", gotM.JSON())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMapFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotM, tt.wantM) {
			// 	t.Errorf("NewMapFromString() = %v, want %v", gotM, tt.wantM)
			// }
		})
	}
}

func TestMap_GetVision(t *testing.T) {
	type fields struct {
		Map *Map
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Vision
	}{
		{
			name:   "Map_GetVision - left top point",
			fields: fields{mockMap1()},
			args: args{
				x: 2, y: 2,
			},
			want: &Vision{
				X1: 0, Y1: 0, X2: 6, Y2: 6,
			},
		},
		{
			name:   "Map_GetVision - right top point",
			fields: fields{mockMap1()},
			args: args{
				x: 18, y: 1,
			},
			want: &Vision{
				X1: 14, Y1: 0, X2: 19, Y2: 5,
			},
		},
		{
			name:   "Map_GetVision - left bottom point",
			fields: fields{mockMap1()},
			args: args{
				x: 3, y: 17,
			},
			want: &Vision{
				X1: 0, Y1: 13, X2: 7, Y2: 19,
			},
		},
		{
			name:   "Map_GetVision - right bottom point",
			fields: fields{mockMap1()},
			args: args{
				x: 16, y: 17,
			},
			want: &Vision{
				X1: 12, Y1: 13, X2: 19, Y2: 19,
			},
		},
		{
			name:   "Map_GetVision - center point",
			fields: fields{mockMap1()},
			args: args{
				x: 10, y: 10,
			},
			want: &Vision{
				X1: 6, Y1: 6, X2: 16, Y2: 16,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields.Map
			got := m.GetVision(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.GetVision() = %v, want %v", got, tt.want)
			}
			fmt.Println(got.InVision(tt.args.x, tt.args.y))
		})
	}
}
