package core

import (
	"fmt"
	"log"
	"testing"
)

func makeMockMap(mapData string) *Map {
	m, err := NewMapFromString(mapData)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = m.Init(defaultVision, defaultWidth, defaultHeight)
	if err != nil {
		log.Println(err)
		return nil
	}

	return m
}

func Test_initTileObjects(t *testing.T) {
	type args struct {
		m *Map
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 []int
	}{
		{
			name: "initTileObjects - print map 1",
			args: args{
				m: makeMockMap(Map3),
			},
			want:  nil,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			T, O := initTileObjects(tt.args.m)
			m := tt.args.m
			printTypes(T, m.Width, m.Height)
			printObjectIndexes(O, m.Width, m.Height)
			updateTunnelExits(m, T, O)
			printMapObjects(m)
			// if !reflect.DeepEqual(T, tt.want) {
			// 	t.Errorf("initTileObjects() T = %v, want %v", T, tt.want)
			// }
			// if !reflect.DeepEqual(O, tt.want1) {
			// 	t.Errorf("initTileObjects() O = %v, want %v", O, tt.want1)
			// }
		})
	}
}

func Test_createGraph(t *testing.T) {
	type args struct {
		m *Map
	}
	tests := []struct {
		name  string
		args  args
		want  [][]V
		want1 []byte
		want2 []int
	}{
		{
			name: "createGraph - print map 1",
			args: args{
				m: makeMockMap(Map1),
			},
			want:  nil,
			want1: nil,
			want2: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.args.m
			G, T, _ := createGraph(m)
			//printGraph(G)
			printTypes(T, m.Width, m.Height)

			n := m.Width * m.Height
			P := createMatrix(n)
			floyd(G, P)
			path := floydPath(P, 60, 86)
			fmt.Printf("%v\n", path)

			// if !reflect.DeepEqual(G, tt.want) {
			// 	t.Errorf("createGraph() G = %v, want %v", G, tt.want)
			// }
			// if !reflect.DeepEqual(T, tt.want1) {
			// 	t.Errorf("createGraph() T = %v, want %v", T, tt.want1)
			// }
			// if !reflect.DeepEqual(O, tt.want2) {
			// 	t.Errorf("createGraph() O = %v, want %v", O, tt.want2)
			// }
		})
	}
}
