package main

import (
	"log"

	"git.nulana.com/bobrnor/battleship-client/client"
)

var (
	c *client.Client
)

func main() {
	c := client.NewClient()
	c.Auth()
	for c.Error() == nil && len(c.RoomUID) == 0 {
		c.SearchRoom()
	}
	//c.ConfirmRoom()
	c.CreateBattlefield()
	c.StartGame()
	c.Longpoll()
	c.Turn()
	c.Longpoll()

	if c.Error() != nil {
		log.Printf("Can't %+v", c.Error())
		return
	}
}
