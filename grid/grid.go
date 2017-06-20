package grid

import "fmt"

// Grid ...
type Grid [13]uint8

func (g *Grid) Set(x, y uint) {
	pos := 10*y + x
	byteIdx := pos / 8
	bitIdx := pos % 8
	g[byteIdx] = g[byteIdx] | (1 << bitIdx)
}

func (g *Grid) Get(x, y uint) bool {
	pos := 10*y + x
	byteIdx := pos / 8
	bitIdx := pos % 8
	return g[byteIdx]&(1<<bitIdx) != 0
}

// Print ...
func (g *Grid) Print() {
	for x := uint(0); x < 10; x++ {
		for y := uint(0); y < 10; y++ {
			if g.Get(x, y) {
				fmt.Printf(" ◉ ")
			} else {
				fmt.Printf(" ○ ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n")
}
