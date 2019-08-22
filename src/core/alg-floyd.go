package core

import (
	"log"
	"strings"
)

// ExitTrap means a tunnel has no exit. because the exit is a meteor,
// or out of the borders of the map.
const ExitTrap = -1

// InfDist means there is no path between vi and vj
const InfDist = 10000

const (
	vHolder     = 0
	vMeteor     = 1
	vTunnel     = 2
	vWormhole   = 3
	vBirthPlace = 4
	vPower      = 5
)

// init tiles with type and Object index
func initTnO(m *Map) ([]int, []int) {
	w, h := m.Width, m.Height
	n := w * h
	T := make([]int, n, n) // by default, all tile are holders
	O := make([]int, n, n) // by default, objects index in object list in map

	for i, o := range m.Meteors {
		v := o.Y*w + o.X
		T[v] = vMeteor
		O[v] = i
		o.V = v
	}

	for i, o := range m.Tunnels {
		v := o.Y*w + o.X
		T[v] = vTunnel
		O[v] = i
		o.V = v
	}

	for i, o := range m.Wormholes {
		v := o.Y*w + o.X
		T[v] = vWormhole
		O[v] = i
		o.V = v
	}

	for i, o := range m.Powers {
		v := o.Y*w + o.X
		T[v] = vPower
		O[v] = i
		o.V = v
	}

	for i, o := range m.PlaceHolders {
		v := o.Y*w + o.X
		T[v] = vBirthPlace
		O[v] = i
		o.V = v
	}
	return T, O
}

func moveVertex(x, y int, w, h int, move *Movement) int {
	mx, my := x+move.Dx, y+move.Dy
	if mx < 0 || mx >= w {
		return ExitTrap
	}
	if my < 0 || my >= h {
		return ExitTrap
	}
	return my*w + mx
}

func calcTunnelExit(m *Map, T []int, O []int, o *Tunnel) int {
	w, h := m.Width, m.Height
	v := o.Y*w + o.X
	move := NewMovementFromTunnel(o.Direction)
	ev := moveVertex(o.X, o.Y, w, h, move)
	if ev == ExitTrap || ev == v {
		return ExitTrap
	}
	t := T[ev]
	switch t {
	case vHolder, vBirthPlace, vPower:
		return ev
	case vMeteor:
		return ExitTrap
	case vTunnel:
		return calcTunnelExit(m, T, O, m.Tunnels[O[ev]])
	case vWormhole:
		return m.Wormholes[O[ev]].ExitVertex
	default:
		return ev
	}
}

func updateTunnelExits(m *Map, T []int, O []int) {
	for _, o := range m.Tunnels {
		o.ExitVertex = calcTunnelExit(m, T, O, o)
	}
}

func updateWormholeExits(m *Map) {
	w := m.Width
	cache := make(map[string]*Wormhole, 20)
	for _, o := range m.Wormholes {
		// hit the pair
		name := strings.ToLower(o.Name)
		paired, ok := cache[name]
		if !ok {
			cache[name] = o
			continue
		}
		delete(cache, name)

		// couple the pair
		o.Exit = paired
		o.ExitVertex = paired.Y*w + paired.X
		paired.Exit = o
		paired.ExitVertex = o.Y*w + o.X
	}

	// update exceptional wormholes which are not paired
	for _, v := range cache {
		v.Exit = nil
		v.ExitVertex = -1
	}
}

func getNeighbors(i, w, h int) []int {
	x, y := i%w, i/w
	nbs := make([]int, 0, 4)
	if x-1 >= 0 {
		nbs = append(nbs, y*w+x-1)
	}
	if x+1 < w {
		nbs = append(nbs, y*w+x+1)
	}
	if y-1 >= 0 {
		nbs = append(nbs, y*w+x-w)
	}
	if y+1 < h {
		nbs = append(nbs, y*w+x+w)
	}
	return nbs
}

func updateGraphWeights(m *Map, T []int, O []int, G [][]int) {
	w, h := m.Width, m.Height
	n := w * h
	for i := 0; i < n; i++ {
		t := T[i]
		// meteor and tunnel have no outbound edges
		if t == vMeteor {
			continue
		}
		nbs := getNeighbors(i, w, h)

		// calculate edges from i to its neighbors
		switch t {
		case vHolder, vBirthPlace, vPower:
			// handle its neighbors
			for _, nb := range nbs {
				nbt := T[nb] // neighbor type
				switch nbt {
				case vHolder, vBirthPlace, vPower:
					G[i][nb] = 1
				case vMeteor:
					// no inbound edges, let it be
				case vWormhole:
					G[i][nb] = 1
				case vTunnel:
					G[i][nb] = 1
				}
			}
		case vWormhole:
			// handle its neighbors
			for _, nb := range nbs {
				nbt := T[nb] // neighbor type
				switch nbt {
				case vHolder, vBirthPlace, vPower:
					G[i][nb] = 1
				case vMeteor:
					// no inbound edges, let it be
				case vWormhole:
					G[i][nb] = 1
				case vTunnel:
					G[i][nb] = 1
				}
			}
			// handle its exit
			o := m.Wormholes[O[i]]
			exit := o.ExitVertex
			if exit == ExitTrap {
				log.Printf("%v has no exit at all, map data is illegal", o.String())
			} else {
				G[i][exit] = 0
				G[exit][i] = 0
				// x1, y1 := i%w, i/w
				// x2, y2 := exit%w, exit/w
				// fmt.Printf("wormhole %v,%v --> %v,%v\n", x1, y1, x2, y2)
			}

		case vTunnel:
			// handle its exit
			o := m.Tunnels[O[i]]
			exit := o.ExitVertex
			if exit == ExitTrap {
				log.Printf("%v has no exit at all, map data is illegal", o.String())
			} else {
				G[i][exit] = 0
				// x1, y1 := i%w, i/w
				// x2, y2 := exit%w, exit/w
				// fmt.Printf("tunnel %v,%v --> %v,%v\n", x1, y1, x2, y2)
			}
		} // end switch
	} // end for
}

func updateGraph(m *Map, T []int, O []int, D [][]int) {
	w, h := m.Width, m.Height
	n := w * h

	// init a graph
	for i := 0; i < n; i++ {
	to:
		for j := 0; j < n; j++ {
			if i == j {
				D[i][j] = 0
				continue to
			}
			D[i][j] = InfDist
		}
	}
	updateGraphWeights(m, T, O, D)
}

func createGraph(m *Map) ([][]int, []int, []int) {
	w, h := m.Width, m.Height
	n := w * h
	T, O := initTnO(m)

	// init a graph
	G := make([][]int, n, n)
	for i := 0; i < n; i++ {
		G[i] = make([]int, n, n)
	to:
		for j := 0; j < n; j++ {
			if i == j {
				continue to
			}
			G[i][j] = InfDist
		}
	}
	updateTunnelExits(m, T, O)
	updateGraphWeights(m, T, O, G)
	return G, T, O
}

func createMatrix(n int) [][]int {
	M := make([][]int, n, n)
	for i := 0; i < n; i++ {
		M[i] = make([]int, n, n)
	}
	return M
}

func floyd(G [][]int, P [][]int) {
	n := len(G)
	// init path matrix
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			P[i][j] = j
		}
	}

	// exhaustion
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if G[i][j] > G[i][k]+G[k][j] {
					G[i][j] = G[i][k] + G[k][j]
					P[i][j] = P[i][k]
				}
			}
		}
	}
}

func floydPath(P [][]int, i, j int) []int {
	path := make([]int, 0)
	k := i
	for {
		k = P[k][j]
		path = append(path, k)
		if k == j {
			break
		}
	}
	return path
}
