package main

import (
	"math/rand"
	"time"

	"git.nulana.com/bobrnor/battleship-client/game"
)

func main() {
	rand.Seed(time.Now().Unix())

	game := game.NewGame()
	<-game.Play()
}
