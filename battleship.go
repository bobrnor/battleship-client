package main

import (
	"git.nulana.com/bobrnor/battleship-client/client"
	"go.uber.org/zap"
)

var (
	c *client.Client
)

func main() {
	setupLogger()

	c := client.NewClient()
	c.Auth()
	for c.Error() == nil && len(c.RoomUID) == 0 {
		c.SearchRoom()
	}
	c.ConfirmRoom()

	if c.Error() != nil {
		zap.S().Errorf("Can't auth %+v", c.Error())
		return
	}
}

func setupLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
