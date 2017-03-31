package main

import (
	"fmt"

	"bobrnor.com/battleship/grid"
)

func main() {

	for i := 0; i < 1000000; i++ {
		_, err := grid.Generate()
		if err != nil {
			fmt.Printf("%+v\n", err.Error())
			return
		}
		// g.Print()
	}
}
