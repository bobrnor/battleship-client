package main

import (
	"math/rand"
	"time"

	"git.nulana.com/bobrnor/battleship-client/game"
)

func main() {
	rand.Seed(time.Now().Unix())

	for {
		game := game.NewGame()
		<-game.Play()
		<-time.After(5 * time.Second)
	}
}
