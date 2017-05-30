package grid

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	sizes = []int{4, 3, 3, 2, 2, 2, 1, 1, 1, 1}
)

// generator ...
type generator struct {
	ships []*ship
	sizes []int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Generate() (Grid, error) {
	g := generator{}

	g.randomizeSizes()
	if err := g.generateArrangement(); err != nil {
		return Grid{}, err
	}

	return g.grid(), nil
}

func (g *generator) generateArrangement() error {
	g.initArrangement()
	if err := g.findValidArrangement(); err != nil {
		return err
	}
	return nil
}

func (g *generator) randomizeSizes() {
	g.sizes = []int{}
	randIndexes := rand.Perm(10)
	for _, index := range randIndexes {
		g.sizes = append(g.sizes, sizes[index])
	}
}

func (g *generator) initArrangement() {
	g.ships = []*ship{}
	for i := 0; i < 10; i++ {
		sizeIndex := g.sizes[i]
		size := sizes[sizeIndex]
		g.ships = append(g.ships, NewShip(size))
	}
}

func (g *generator) findValidArrangement() error {
	for i := 0; i < 10; {
		if g.isSubArrangementValid(i) {
			i++
		} else {
			for !g.findNextPosition(i) {
				if i > 0 {
					i--
				} else {
					return fmt.Errorf("Can't find valid arrangement:\n%+v", g)
				}
			}
		}
	}
	return nil
}

func (g *generator) findNextPosition(shipIndex int) bool {
	ship := g.ships[shipIndex]
	if !ship.FindNextPosition() {
		return false
	}
	g.ships[shipIndex] = ship
	return true
}

func (g *generator) isSubArrangementValid(lastShipIndex int) bool {
	ship := g.ships[lastShipIndex]

	if !ship.IsValid() {
		return false
	}

	for i := lastShipIndex - 1; i >= 0; i-- {
		s := g.ships[i]
		if s.Intersect(ship) {
			return false
		}
	}

	return true
}

func (g *generator) grid() Grid {
	grid := Grid{}
	for _, ship := range g.ships {
		for x := ship.Left; x <= ship.Right; x++ {
			for y := ship.Top; y <= ship.Bottom; y++ {
				grid[x][y] = true
			}
		}
	}
	return grid
}
