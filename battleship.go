package main

import (
	"sync"

	"math/rand"
	"time"

	"git.nulana.com/bobrnor/battleship-client/game"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	rand.Seed(time.Now().Unix())

	game := game.NewGame()
	game.Play()

	//c := game.NewClient()
	//c.Auth()
	//c.SearchRoom()
	//c.CreateBattlefield()
	//c.StartGame()
	//c.Longpoll()
	//c.Turn()
	//c.Longpoll()

	//if c.Error() != nil {
	//	log.Printf("Can't %+v", c.Error())
	//	wg.Done()
	//	return
	//}

	wg.Wait()
}
