package core

import (
	"fmt"
)

// Neighbor is
type Neighbor struct {
	V    int
	Move string
}

func (n Neighbor) String() string {
	return fmt.Sprintf("%v to %d", n.Move, n.V)
}

// GetNeighborMovements is
func GetNeighborMovements(i, w, h int) []*Neighbor {
	x, y := i%w, i/w
	nbs := make([]*Neighbor, 0, 4)
	if x-1 >= 0 {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x - 1,
			Move: left,
		})
	}
	if x+1 < w {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x + 1,
			Move: right,
		})
	}
	if y-1 >= 0 {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x - w,
			Move: up,
		})
	}
	if y+1 < h {
		nbs = append(nbs, &Neighbor{
			V:    y*w + x + w,
			Move: down,
		})
	}
	return nbs
}

// CalcMove is
func CalcMove(w, h int, src, dest int) string {
	nbs := GetNeighborMovements(src, w, h)
	m := stay
	for _, nb := range nbs {
		if nb.V == dest {
			m = nb.Move
			break
		}
	}
	return m
}
