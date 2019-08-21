package core

// V is vertex type which store vertex information
// as of now, V store the distance of a edge
type V uint16

// MaxV is
const MaxV = 1<<16 - 1
const (
	vHolder   = byte(0)
	vMeteor   = byte(1)
	vTunnel   = byte(2)
	vWormhole = byte(3)
)

func initTiles(m *Map) [][]byte {
	w, h := m.Width, m.Height

	// init tiles with type
	t := make([][]byte, w, w)
	for i := 0; i < w; i++ {
		t[i] = make([]byte, h, h)
	}

	// set meteors
	for _, o := range m.Meteors {
		t[o.X][o.Y] = vMeteor
	}

	// set tunnels
	for _, o := range m.Tunnels {
		t[o.X][o.Y] = vTunnel
	}

	// set wormholes
	for _, o := range m.Wormholes {
		t[o.X][o.Y] = vWormhole
	}
	return t
}

func initGraphByTiles(m *Map, t [][]byte, g [][]V) {
	w, h := m.Width, m.Height

	// init tiles with type

	for i := 0; i < w; i++ {
		t[i] = make([]byte, h, h)
	}

	// set meteors
	for _, o := range m.Meteors {
		t[o.X][o.Y] = vMeteor
	}

	// set tunnels
	for _, o := range m.Tunnels {
		t[o.X][o.Y] = vTunnel
	}

	// set wormholes
	for _, o := range m.Wormholes {
		t[o.X][o.Y] = vWormhole
	}
}

func initGraph(m *Map) {
	w, h := m.Width, m.Height
	n := w * h
	t := initTiles(m)

	// init a graph
	g := make([][]V, n, n)
	for i := 0; i < n; i++ {
		g[i] = make([]V, n, n)
		for j := 0; j < n; j++ {
			if i != j {
				g[i][j] = MaxV
			}
		}
	}

	initGraphByTiles(m, t, g)
	for i := 0; i < n; i++ {
		// x, y := i/w, i%w
		// if x == y {

		// 	g[i] = MaxV
		// }
	}

	// for _, o := range m.Meteors {
	// 	i := o.X*w + o.Y
	// 	g[i] = MaxV
	// }
}
