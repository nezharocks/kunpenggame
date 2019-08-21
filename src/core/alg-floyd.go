package core

import (
	"log"
)

// V is vertex type which store vertex information
// as of now, V store the distance of a edge
type V uint16

// ExitTrap means a tunnel has no exit. because the exit is a meteor,
// or out of the borders of the map.
const ExitTrap = -1

// InfinitDist means there is no path between vi and vj
const InfinitDist = V(10000)

const (
	vHolder   = byte(0)
	vMeteor   = byte(1)
	vTunnel   = byte(2)
	vWormhole = byte(3)
)

// init tiles with type and index
func initTileObjects(m *Map) ([]byte, []int) {
	w, h := m.Width, m.Height
	n := w * h
	T := make([]byte, n, n) // by default, all tile are holders
	O := make([]int, n, n)  // by default, objects index in object list in map

	for i, o := range m.Meteors {
		v := o.Y*w + o.X
		T[v] = vMeteor
		O[v] = i
	}

	for i, o := range m.Tunnels {
		v := o.Y*w + o.X
		T[v] = vTunnel
		O[v] = i
	}

	for i, o := range m.Wormholes {
		v := o.Y*w + o.X
		T[v] = vWormhole
		O[v] = i
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

func calcTunnelExit(m *Map, T []byte, O []int, o *Tunnel) int {
	w, h := m.Width, m.Height
	v := o.Y*w + o.X
	move := NewMovementFromTunnel(o.Direction)
	ev := moveVertex(o.X, o.Y, w, h, move)
	if ev == ExitTrap || ev == v {
		return ExitTrap
	}
	t := T[ev]
	switch t {
	case vHolder:
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

func updateTunnelExits(m *Map, T []byte, O []int) {
	for _, o := range m.Tunnels {
		o.ExitVertex = calcTunnelExit(m, T, O, o)
	}
}

func updateGraphWeights(m *Map, T []byte, O []int, G [][]V) {
	w, h := m.Width, m.Height
	n := w * h
	for i := 0; i < n; i++ {
		t := T[i]
		// meteor and tunnel have no outbound edges
		if t == vMeteor || t == vTunnel {
			continue
		}
		x, y := i%w, i/w
		// get neighbors
		nbs := make([]V, 0)
		if x-1 >= 0 {
			nbs = append(nbs, V(y*w+x-1))
		}
		if x+1 < w {
			nbs = append(nbs, V(y*w+x+1))
		}
		if y-1 >= 0 {
			nbs = append(nbs, V(y*w+x-w))
		}
		if y+1 < h {
			nbs = append(nbs, V(y*w+x+w))
		}

		// calculate edges from i to its neighbors
		switch t {
		case vHolder:
			for _, nb := range nbs {
				nbt := T[nb] // neighbor type
				switch nbt {
				case vHolder:
					G[i][nb] = 1
				case vMeteor:
					// no inbound edges, let it be
				case vWormhole:
					G[i][nb] = 1 // todo
					// handle the edge with its exit
					o := m.Wormholes[O[nb]]
					exit := o.ExitVertex
					if exit == ExitTrap {
						log.Printf("%v has no exit at all, map data is illegal", o.String())
					} else {
						G[i][exit] = 1
						// todo: record the tunnel pair
					}
				case vTunnel:
					G[i][nb] = 1 // todo
					// handle the edge with its exit
					o := m.Tunnels[O[nb]]
					exit := o.ExitVertex
					if exit == ExitTrap {
						log.Printf("%v has no exit at all, map data has a trap from point %v", o.String(), i)
					} else if i == exit {
						log.Printf("%v has an exit which is back to point %v. bypass the edge", o.String(), i)
					} else {
						G[i][exit] = 1
						// todo: record the tunnel pair
					}
				}
			}
		case vWormhole:
			// handle the edge from current wormhole to its exit
			o := m.Wormholes[O[i]]
			exit := o.ExitVertex
			if exit == ExitTrap {
				log.Printf("%v has no exit at all, map data is illegal", o.String())
			} else {
				G[i][exit] = 0
				G[exit][i] = 0
			}
			// handle neighbors' edges
			for _, nb := range nbs {
				nbt := T[nb] // neighbor type
				switch nbt {
				case vHolder:
					G[i][nb] = 1
				case vMeteor:
					// no inbound edges, let it be
				case vWormhole:
					G[i][nb] = 1
					// handle the edge with its exit
					o := m.Wormholes[O[nb]]
					exit := o.ExitVertex
					if exit == ExitTrap {
						log.Printf("%v has no exit at all, map data is illegal", o.String())
					} else if i == exit {
						log.Printf("%v has an exit which is back to point %v. the wormhole pair are neighbors, bypass the edge", o.String(), i)
					} else {
						G[i][exit] = 1
						// todo: record the tunnel pair
					}
				case vTunnel:
					G[i][nb] = 1 // todo
					// handle the edge with its exit
					o := m.Tunnels[O[nb]]
					exit := o.ExitVertex
					if exit == ExitTrap {
						log.Printf("%v has no exit at all, map data has a trap from point %v", o.String(), i)
					} else if i == exit {
						log.Printf("%v has an exit which is back to point %v. bypass the edge", o.String(), i)
					} else {
						G[i][exit] = 1
						// todo: record the tunnel pair
					}
				}
			}
		}
	}
}

func createGraph(m *Map) ([][]V, []byte, []int) {
	w, h := m.Width, m.Height
	n := w * h
	T, O := initTileObjects(m)

	// init a graph
	G := make([][]V, n, n)
	for i := 0; i < n; i++ {
		G[i] = make([]V, n, n)
	to:
		for j := 0; j < n; j++ {
			if i == j {
				continue to
			}
			G[i][j] = InfinitDist
		}
	}
	updateTunnelExits(m, T, O)
	updateGraphWeights(m, T, O, G)
	return G, T, O
}

func createMatrix(n int) [][]V {
	M := make([][]V, n, n)
	for i := 0; i < n; i++ {
		M[i] = make([]V, n, n)
	}
	return M
}
func floyd(G [][]V, P [][]V) {
	n := len(G)
	// init path matrix
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			P[i][j] = V(j)
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
func floydPath(P [][]V, i, j V) []V {
	path := make([]V, 10)
	k := i
	for {
		path = append(path, k)
		if k == j {
			break
		}
		k = P[k][j]
	}
	return path
}
