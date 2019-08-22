package core

import (
	"fmt"
	"testing"
)

func TestMapVision_Visit(t *testing.T) {
	type fields struct {
		MapVision *MapVision
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "MapVision.Visit - visit center point",
			fields: fields{
				MapVision: NewMapVision(20, 20, 3),
			},
			args: args{10, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields.MapVision
			// m.Visit(tt.args.x, tt.args.y)
			areas := m.Areas
			// m.Visit(areas[1].X, areas[1].Y)
			// m.Visit(areas[5].X, areas[5].Y)
			m.Visit(areas[1].X, areas[1].Y)
			m.Visit(areas[6].X, areas[6].Y)
			m.Visit(areas[8].X, areas[8].Y)
			m.Visit(5, 5)
			fmt.Println(m.String())

			blind := m.BlindAreas()
			for _, b := range blind {
				fmt.Println(b.String())
			}
		})
	}
}
