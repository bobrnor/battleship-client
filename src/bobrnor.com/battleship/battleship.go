package main

import (
	"fmt"

	"bobrnor.com/battleship/grid"
)

func main() {

	g, err := grid.Generate()
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		return
	}
	g.Print()
}
