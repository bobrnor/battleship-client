package grid

import "fmt"

// Grid ...
type Grid [10][10]bool

// Print ...
func (g *Grid) Print() {
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if g[x][y] {
				fmt.Printf(" * ")
			} else {
				fmt.Printf(" . ")
			}
		}
		fmt.Printf("\n")
	}
}
