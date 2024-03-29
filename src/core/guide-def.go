package core

// Guide is
type Guide struct {
	Map *Map

	MapVision *MapVision

	// T is a vector of all vertex's type
	T []int

	// O is a vector of all vertex's index which locates tile objects in Map
	O []int

	// D is a matrix of shortest distance.
	D [][]int

	// P is a matrix for shortest path.
	P [][]int
}

// NewGuide creates a guide from a Map which is parsed from a JSON bytes.
func NewGuide(m *Map) *Guide {
	n := m.Width * m.Height
	D := createMatrix(n)
	P := createMatrix(n)
	// create T and O, update all objects' vertex index
	T, O := initTnO(m)
	updateWormholeExits(m)
	updateTunnelExits(m, T, O)
	updateGraph(m, T, O, D)
	floyd(D, P)
	mv := NewMapVision(m.Width, m.Height, m.Vision)
	return &Guide{
		Map:       m,
		MapVision: mv,
		T:         T,
		O:         O,
		D:         D,
		P:         P,
	}
}

// ShortestPath is
func (g *Guide) ShortestPath(i, j int) []int {
	path := make([]int, 0, 10)
	k := i
	for {
		k = g.P[k][j]
		path = append(path, k)
		if k == j {
			break
		}
	}
	return path
}

// ShortestDistance is
func (g *Guide) ShortestDistance(i, j int) int {
	return g.D[i][j]
}

// Target is
type Target struct {
	Vertex int
	// player ID
	Subject *Player

	// flag see if it is a player, by default, it is false (power)
	Player bool
	Point  int
	Dist   int
	Value  float64
}
