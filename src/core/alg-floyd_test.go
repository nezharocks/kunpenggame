package core

import (
	"fmt"
	"log"
	"testing"
	"time"
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

func Test_initTnO(t *testing.T) {
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
			name: "initTnO - print map 1",
			args: args{
				m: makeMockMap(Map3),
			},
			want:  nil,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			T, O := initTnO(tt.args.m)
			m := tt.args.m
			printMap(m, T, O, m.Width, m.Height)
			// printObjectIndexes(O, m.Width, m.Height)
			updateTunnelExits(m, T, O)
			printMapObjects(m)
			// if !reflect.DeepEqual(T, tt.want) {
			// 	t.Errorf("initTnO() T = %v, want %v", T, tt.want)
			// }
			// if !reflect.DeepEqual(O, tt.want1) {
			// 	t.Errorf("initTnO() O = %v, want %v", O, tt.want1)
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
		want  [][]int
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
			start := time.Now()
			G, T, O := createGraph(m)
			//printGraph(G)
			n := m.Width * m.Height
			P := createMatrix(n)
			floyd(G, P)
			fmt.Printf("%v\n", time.Since(start))
			printMap(m, T, O, m.Width, m.Height)
			printPath(G, P, 0, 0)
			printPath(G, P, 360, 360)
			printPath(G, P, 0, 360)
			printPath(G, P, 60, 86)
			printPath(G, P, 0, 4)
			printPath(G, P, 0, 385)
			printPath(G, P, 14, 246)
			printPath(G, P, 250, 150)

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
func printPath(G [][]int, P [][]int, i, j int) {
	path := floydPath(P, i, j)
	fmt.Printf("%v: %v\n", G[i][j], path)
}
