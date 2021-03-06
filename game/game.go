package game

import (
	"math/rand"
	"sync"

	"log"

	"time"

	"git.nulana.com/bobrnor/battleship-grid.git"
	json "git.nulana.com/bobrnor/json.git"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"
)

type Game struct {
	sync.Mutex

	UID     string
	RoomUID string
	Grid    *grid.Grid

	turns *grid.Grid
	hits  *grid.Grid

	longpollClient *LongpollClient
	jsonClient     json.Client
	err            error

	gameOverChan chan struct{}
}

type ResponseError struct {
	Code uint64 `json:"code"`
	Msg  string `json:"msg"`
}

const (
	authPath   = "http://battleship_server:80/auth"
	searchPath = "http://battleship_server:80/game/search"
	startPath  = "http://battleship_server:80/game/start"
	turnPath   = "http://battleship_server:80/game/turn"

	searchResultType = "search_result"
	gameType         = "game"
	opponentTurnType = "opponent_turn"
	gameOverType     = "game_over"
)

func NewGame() *Game {
	game := &Game{
		UID:          uuid.TimeOrderedUUID(),
		jsonClient:   json.Client{},
		gameOverChan: make(chan struct{}, 1),
	}
	return game
}

func (g *Game) Play() <-chan struct{} {
	log.Printf("Starting game `%s`", g.UID)
	g.auth()
	g.searchGame()
	g.checkError()
	return g.gameOverChan
}

func (g *Game) stop() {
	select {
	case g.gameOverChan <- struct{}{}:
	default:
	}
}

func (g *Game) auth() {
	g.Lock()
	defer g.Unlock()

	if g.err != nil {
		return
	}

	request := g.authRequest()

	var response struct {
		Type      string        `json:"type"`
		ClientUID string        `json:"client_uid"`
		Error     ResponseError `json:"error"`
	}

	g.doRequest(authPath, request, &response)
	g.checkResponseError(response.Error)
	g.authDone()
}

func (g *Game) authRequest() interface{} {
	return map[string]interface{}{
		"client_uid": g.UID,
	}
}

func (g *Game) authDone() {
	if g.err != nil {
		return
	}

	g.longpollClient = NewLongpollClient(g.UID, g.LongpollMessageReceived)
}

func (g *Game) searchGame() {
	g.Lock()
	defer g.Unlock()

	if g.err != nil {
		return
	}

	request := g.searchRequest()

	var response struct {
		Type  string        `json:"type"`
		Error ResponseError `json:"error"`
	}

	g.doRequest(searchPath, request, &response)
	g.checkResponseError(response.Error)
}

func (g *Game) searchRequest() interface{} {
	return map[string]interface{}{
		"client_uid": g.UID,
	}
}

func (g *Game) LongpollMessageReceived(message map[string]interface{}) {
	//log.Printf("Longpoll message received %+v", message)

	if len(message) == 0 {
		return
	}

	messageType, ok := message["type"].(string)
	if !ok {
		log.Printf("Longpoll message has bad format %+v", message)
	}

	switch messageType {
	case searchResultType:
		g.searchDone(message)
	case gameType:
		g.startDone(message)
	case opponentTurnType:
		g.opponentTurn(message)
	case gameOverType:
		g.gameOver(message)
	default:
		log.Printf("Unknown longpoll message type %s", messageType)
	}

	g.checkError()
}

func (g *Game) searchDone(message map[string]interface{}) {
	g.Lock()
	defer g.Unlock()

	if g.err != nil {
		return
	}

	roomUID, ok := message["room_uid"].(string)
	if !ok {
		g.err = errors.Errorf("`search_result` message has bad format %+v", message)
		return
	}

	g.RoomUID = roomUID
	g.generateGrid()
	g.startGame()
}

func (g *Game) startDone(message map[string]interface{}) {
	g.Lock()
	defer g.Unlock()

	if g.err != nil {
		return
	}

	action, ok := message["action"].(string)
	if !ok {
		g.err = errors.Errorf("`game` message has bad format %+v", message)
		return
	}

	if action == "turn" {
		g.turn()
	}
}

func (g *Game) generateGrid() {
	gr, err := grid.Generate()
	if err != nil {
		g.err = err
		return
	}

	gr.Print()

	g.Grid = gr
	g.hits = &grid.Grid{}
	g.turns = &grid.Grid{}
}

func (g *Game) startGame() {
	request := g.startRequest()

	var response struct {
		Type  string        `json:"type"`
		Error ResponseError `json:"error"`
	}

	g.doRequest(startPath, request, &response)
	g.checkResponseError(response.Error)
}

func (g *Game) startRequest() interface{} {
	return map[string]interface{}{
		"client_uid": g.UID,
		"room_uid":   g.RoomUID,
		"grid":       g.Grid,
	}
}

func (g *Game) turn() {
	// throttling
	<-time.After(1 * time.Second)

	if g.turns.IsFull() {
		g.err = errors.Errorf("No turns left")
		g.stop()
		return
	}

	request := g.turnRequest()

	var response struct {
		Type   string        `json:"type"`
		Result string        `json:"result"`
		Error  ResponseError `json:"error"`
	}

	g.doRequest(turnPath, request, &response)
	g.checkResponseError(response.Error)
	g.afterTurn(response.Result)
}

func (g *Game) turnRequest() map[string]interface{} {
	var x, y uint
	for {
		x = uint(rand.Intn(10))
		y = uint(rand.Intn(10))
		if !g.turns.Get(x, y) {
			g.turns.Set(x, y)
			break
		}
	}

	coord := map[string]interface{}{
		"x": x,
		"y": y,
	}
	//log.Printf("My turn: %+v", coord)
	return map[string]interface{}{
		"client_uid": g.UID,
		"room_uid":   g.RoomUID,
		"coord":      coord,
	}
}

func (g *Game) afterTurn(result string) {
	if g.err != nil {
		return
	}

	switch result {
	case "hit":
		g.turn()
	case "win":
		g.stop()
	}
}

func (g *Game) opponentTurn(message map[string]interface{}) {
	g.Lock()
	defer g.Unlock()

	if g.err != nil {
		return
	}

	x, ok := message["x"].(float64)
	if !ok {
		g.err = errors.Errorf("`opponent_turn` message has bad format %+v", message)
		return
	}

	y, ok := message["y"].(float64)
	if !ok {
		g.err = errors.Errorf("`opponent_turn` message has bad format %+v", message)
		return
	}

	if !g.checkTurn(uint(x), uint(y)) {
		//log.Printf("Opponent missed, my turn")
		g.turn()
	} else {
		//log.Printf("Opponent hit")
	}

	g.Grid.PrintWithHitsOverlay(g.hits)
}

func (g *Game) checkTurn(x, y uint) bool {
	if g.err != nil {
		return false
	}

	g.hits.Set(x, y)
	return g.Grid.Get(x, y)
}

func (g *Game) gameOver(message map[string]interface{}) {
	g.Lock()
	defer g.Unlock()

	if g.err != nil {
		return
	}

	g.stop()
}

func (g *Game) checkResponseError(err ResponseError) {
	if g.err != nil {
		return
	}

	if err.Code != 0 {
		g.err = errors.Errorf("Bad status [%d] `%s`", err.Code, err.Msg)
		return
	}
}

func (g *Game) doRequest(path string, data interface{}, response interface{}) {
	if g.err != nil {
		return
	}

	err := g.jsonClient.Post(path, data, response)
	if err != nil {
		g.err = err
		return
	}

	//log.Printf("Response read %+v", response)
}

func (g *Game) checkError() {
	g.Lock()
	defer g.Unlock()

	if g.err != nil {
		log.Printf("Error: %+v", g.err.Error())
		g.stop()
	}
}
