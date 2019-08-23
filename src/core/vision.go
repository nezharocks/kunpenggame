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
	Index   int
	X1, Y1  int // left top pixel
	X2, Y2  int // right bottom pixel
	X, Y, V int // center pixel
	Count   int
	Points  []int
}

func (a *MapArea) String() string {
	return fmt.Sprintf("%v\t\t%-2d,%-2d %-2d,%-2d %-2d,%-2d %v %v", a.Index, a.X1, a.Y1, a.X2, a.Y2, a.X, a.Y, a.V, a.Count)
}

// MapPoint is
type MapPoint struct {
	X, Y int
	W, H int
}

// NewMapPoint is
func NewMapPoint(x, y, w, h int) *MapPoint {
	return &MapPoint{x, y, w, h}
}

// Vertex returns the vertex index in the map tree
func (p *MapPoint) Vertex() int {
	return p.Y*p.W + p.X
}

// Up moves the point up
func (p *MapPoint) Up() *MapPoint {
	if p.Y-1 >= 0 {
		p.Y--
	}
	return p
}

// Down moves the point down
func (p *MapPoint) Down() *MapPoint {
	if p.Y+1 < p.H {
		p.Y++
	}
	return p
}

// Left moves the point left
func (p *MapPoint) Left() *MapPoint {
	if p.X-1 >= 0 {
		p.X--
	}
	return p
}

// Right moves the point Right
func (p *MapPoint) Right() *MapPoint {
	if p.X+1 < p.W {
		p.X++
	}
	return p
}

// UpdatePoints is
func (a *MapArea) UpdatePoints(w, h int) {
	p := NewMapPoint(a.X, a.Y, w, h)
	points := make([]int, 0, 4)
	switch a.Index {
	case 0:
		points = append(points, p.Vertex())
		points = append(points, p.Left().Vertex())
		points = append(points, p.Up().Vertex())
		points = append(points, p.Right().Vertex())
	case 1:
		points = append(points, p.Vertex())
		points = append(points, p.Up().Vertex())
		points = append(points, p.Left().Vertex())
		points = append(points, p.Right().Right().Vertex())
	case 2:
		points = append(points, p.Vertex())
		points = append(points, p.Right().Vertex())
		points = append(points, p.Up().Vertex())
		points = append(points, p.Left().Vertex())
	case 3:
		points = append(points, p.Vertex())
		points = append(points, p.Left().Vertex())
		points = append(points, p.Up().Vertex())
		points = append(points, p.Down().Down().Vertex())
	case 4:
		points = append(points, p.Vertex())
		points = append(points, p.Up().Vertex())
		points = append(points, p.Down().Down().Vertex())
		points = append(points, p.Up().Left().Vertex())
		points = append(points, p.Right().Right().Vertex())
	case 5:
		points = append(points, p.Vertex())
		points = append(points, p.Right().Vertex())
		points = append(points, p.Up().Vertex())
		points = append(points, p.Down().Down().Vertex())
	case 6:
		points = append(points, p.Vertex())
		points = append(points, p.Left().Vertex())
		points = append(points, p.Down().Vertex())
		points = append(points, p.Right().Vertex())
	case 7:
		points = append(points, p.Vertex())
		points = append(points, p.Down().Vertex())
		points = append(points, p.Left().Vertex())
		points = append(points, p.Right().Right().Vertex())
	case 8:
		points = append(points, p.Vertex())
		points = append(points, p.Right().Vertex())
		points = append(points, p.Down().Vertex())
		points = append(points, p.Left().Vertex())
	}
	a.Points = points
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
		x1, y1    int
		x2, y2    int
		x, y      int
		i, vertex int
		areas     [9]*MapArea
	)

	// first line
	y1 = 0

	i = 0
	x1 = 0
	x2, y2 = x1+v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	i = 1
	x1 = v
	x2, y2 = w-v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	i = 2
	x1 = w - v
	x2, y2 = w-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	// second line
	y1 = v
	y2 = h - v - 1

	i = 3
	x1 = 0
	x2 = x1 + v - 1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	i = 4
	x1 = v
	x2 = w - v - 1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	i = 5
	x1 = w - v
	x2 = w - 1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	// the third line
	y1 = h - v

	i = 6
	x1 = 0
	x2, y2 = x1+v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	i = 7
	x1 = v
	x2, y2 = w-v-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	i = 8
	x1 = w - v
	x2, y2 = w-1, y1+v-1
	x, y = (x2+x1)/2, (y2+y1)/2
	vertex = y*w + x
	areas[i] = &MapArea{i, x1, y1, x2, y2, x, y, vertex, 0, nil}

	// update points
	for _, a := range areas {
		a.UpdatePoints(w, h)
	}
	m.Areas = areas[:]
}

func (m *MapVision) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Width: %v, Height: %v, Vision: %v\n", m.Width, m.Height, m.Vision))
	for _, a := range m.Areas {
		s := fmt.Sprintf("%v\n", a.String())
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
			count = 1
		}
		a.Count = count
	}

	return m.Areas
}
