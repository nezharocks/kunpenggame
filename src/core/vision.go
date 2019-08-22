package core

import (
	"fmt"
	"strings"
)

// Vision is
type Vision struct {
	X1, Y1 int // the left and top point
	X2, Y2 int // the left and top point
	X, Y   int
}

func (v *Vision) String() string {
	return fmt.Sprintf("%-2d,%-2d %-2d,%-2d %-2d,%-2d\n", v.X1, v.Y1, v.X2, v.Y2, v.X, v.Y)
}

// GetVision is
func GetVision(w, h, v, x, y int) *Vision {
	rbx := w - 1
	rby := h - 1
	x1 := x - v
	if x1 < 0 {
		x1 = 0
	}
	x2 := x + v
	if x2 > rbx {
		x2 = rbx
	}
	y1 := y - v
	if y1 < 0 {
		y1 = 0
	}
	y2 := y + v
	if y2 > rby {
		y2 = rby
	}
	return &Vision{X1: x1, Y1: y1, X2: x2, Y2: y2, X: x, Y: y}
}

// InVision is
func (v Vision) InVision(x, y int) bool {
	if x < v.X1 || x > v.X2 {
		return false
	}
	if y < v.Y1 || y > v.Y2 {
		return false
	}
	return true
}

// MapArea is
type MapArea struct {
	X1, Y1 int // left top pixel
	X2, Y2 int // right bottom pixel
	X, Y   int // center pixel
	Count  int
}

func (a *MapArea) String() string {
	return fmt.Sprintf("%-2d,%-2d %-2d,%-2d %-2d,%-2d %v", a.X1, a.Y1, a.X2, a.Y2, a.X, a.Y, a.Count)
}

// MapVision is
type MapVision struct {
	Pixels [][]bool
	Width  int
	Height int
	Vision int
	Areas  []*MapArea
}

// NewMapVision is
func NewMapVision(w, h, v int) *MapVision {
	pixels := make([][]bool, w, w)
	for i := range pixels {
		pixels[i] = make([]bool, h, h)
	}
	m := &MapVision{
		Pixels: pixels,
		Width:  w,
		Height: h,
		Vision: v,
	}
	m.initAreas()
	return m
}

func (m *MapVision) initAreas() {
	w, h, v := m.Width, m.Height, m.Vision*2+1
	var (
		x1, y1 int
		x2, y2 int
		x, y   int
		i      int
		areas  [9]*MapArea
	)

	// first line
	y1 = 0

	i = 0
	x1 = 0
	x2, y2 = x1+v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	i = 1
	x1 = v
	x2, y2 = w-v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	i = 2
	x1 = w - v
	x2, y2 = w-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	// second line
	y1 = v
	y2 = h - v - 1

	i = 3
	x1 = 0
	x2 = x1 + v - 1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	i = 4
	x1 = v
	x2 = w - v - 1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	i = 5
	x1 = w - v
	x2 = w - 1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	// the third line
	y1 = h - v

	i = 6
	x1 = 0
	x2, y2 = x1+v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	i = 7
	x1 = v
	x2, y2 = w-v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	i = 8
	x1 = w - v
	x2, y2 = w-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	areas[i] = &MapArea{x1, y1, x2, y2, x, y, 0}

	m.Areas = areas[:]
}

func (m *MapVision) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Width: %v, Height: %v, Vision: %v\n", m.Width, m.Height, m.Vision))
	for i, a := range m.Areas {
		s := fmt.Sprintf("%v\t\t%v\n", i, a.String())
		sb.WriteString(s)
	}
	sb.WriteString("\n")
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			c := "."
			if m.Pixels[x][y] {
				c = " "
			}
			sb.WriteString(c)
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	return sb.String()
}

// Visit is
func (m *MapVision) Visit(x, y int) {
	v := GetVision(m.Width, m.Height, m.Vision, x, y)
	for i := v.X1; i <= v.X2; i++ {
		for j := v.Y1; j <= v.Y2; j++ {
			m.Pixels[i][j] = true
		}
	}
}

// BlindAreas is
func (m *MapVision) BlindAreas() []*MapArea {
	blind := make([]*MapArea, 0, 9)
	for _, a := range m.Areas {
		count := 0
		for x := a.X1; x <= a.X2; x++ {
			for y := a.Y1; y <= a.Y2; y++ {
				if !m.Pixels[x][y] {
					count++

				}
			}
		}
		if count == 0 {
			continue
		}
		b := *a
		b.Count = count
		blind = append(blind, &b)
	}

	return blind
}
