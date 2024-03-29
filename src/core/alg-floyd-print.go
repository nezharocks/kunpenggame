package core

import (
	"fmt"
	"math"
)

func printGraph(G [][]int) {
	w := len(G)
	// for i := 0; i < w; i++ {
	// 	fmt.Printf("%-3d ", i)
	// }
	fmt.Println()
	for j := 0; j < w; j++ {
		for i := 0; i < w; i++ {
			v := "  "
			if G[i][j] != math.MaxUint16 {
				v = fmt.Sprintf("%-2d", G[i][j])
			}
			fmt.Print(v)
		}
		fmt.Println()
		fmt.Println()
	}
	fmt.Println()
}

func printObjectIndexes(O []int, w, h int) {
	// column header
	fmt.Print("    ")
	for i := 0; i < w; i++ {
		fmt.Printf("%-3d ", i)
	}
	fmt.Println()

	// column line
	fmt.Print("  |-")
	for i := 0; i < w; i++ {
		fmt.Printf("----")
	}
	fmt.Println()

	for y := 0; y < h; y++ {
		fmt.Printf("%-2d| ", y)
		for x := 0; x < w; x++ {
			i := y*w + x
			fmt.Printf("%-3d ", O[i])
		}
		fmt.Println()
	}
	fmt.Println()
}

func printMap(m *Map, T []int, O []int, w, h int) {
	// column header
	fmt.Print("    ")
	for i := 0; i < w; i++ {
		fmt.Printf("%-2d ", i)
	}
	fmt.Println()

	// column line
	fmt.Print("  |-")
	for i := 0; i < w; i++ {
		fmt.Printf("---")
	}
	fmt.Println()

	for y := 0; y < h; y++ {
		fmt.Printf("%-2d| ", y)
		for x := 0; x < w; x++ {
			i := y*w + x
			t := T[i]
			switch t {
			case vHolder:
				fmt.Print(".  ")
			case vBirthPlace:
				o := m.PlaceHolders[O[i]]
				fmt.Printf("%v  ", o.Name)
			case vPower:
				o := m.Powers[O[i]]
				fmt.Printf("%v  ", o.Point)
			case vMeteor:
				fmt.Print("#  ")
			case vWormhole:
				o := m.Wormholes[O[i]]
				fmt.Printf("%v  ", o.Name)
			case vTunnel:
				o := m.Tunnels[O[i]]
				fmt.Printf("%v  ", o.Char())
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func printMapObjects(m *Map) {
	fmt.Println("meteors\t==============================")
	for i, o := range m.Meteors {
		fmt.Printf("%-2d  %v\n", i, o.String())
	}

	fmt.Println("tunnels\t==============================")
	for i, o := range m.Tunnels {
		fmt.Printf("%-2d  %v\n", i, o.String())
	}

	fmt.Println("wormholes\t==============================")
	for i, o := range m.Wormholes {
		fmt.Printf("%-2d  %v\n", i, o.String())
	}

	fmt.Println("powers\t==============================")
	total := 0
	for i, o := range m.Powers {
		fmt.Printf("%-2d  %v\n", i, o.String())
		total += o.Point
	}
	fmt.Println("total points", total)
}
