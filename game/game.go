package game

import (
	"math/rand"
	"sync"

	"log"

	"git.nulana.com/bobrnor/battleship-client/grid"
	json "git.nulana.com/bobrnor/json.git"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"
)

type Game struct {
	sync.Mutex

	UID     string
	RoomUID string
	Grid    *grid.Grid

	longpollClient *LongpollClient
	jsonClient     json.Client
	err            error
}

type ResponseError struct {
	Code uint64 `json:"code"`
	Msg  string `json:"msg"`
}

const (
	authPath   = "http://0.0.0.0:8000/auth"
	searchPath = "http://0.0.0.0:8000/game/search"
	startPath  = "http://0.0.0.0:8000/game/start"
	turnPath   = "http://0.0.0.0:8000/game/turn"

	searchResultType = "search_result"
	gameType         = "game"
	opponentTurnType = "opponent_turn"
)

func NewGame() *Game {
	game := &Game{
		UID:        uuid.TimeOrderedUUID(),
		jsonClient: json.Client{},
	}
	return game
}

func (g *Game) Play() {
	log.Printf("Starting game `%s`", g.UID)
	g.auth()
	g.searchGame()
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
	log.Printf("Longpoll message received %+v", message)

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
	default:
		log.Printf("Unknown longpoll message type %s", messageType)
	}
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
	coord := map[string]interface{}{
		"x": uint(rand.Intn(10)),
		"y": uint(rand.Intn(10)),
	}
	log.Printf("My turn: %+v", coord)
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

	if result == "hit" {
		g.turn()
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
		log.Printf("Opponent missed, my turn")
		g.turn()
	} else {
		log.Printf("Opponent hit")
	}
}

func (g *Game) checkTurn(x, y uint) bool {
	if g.err != nil {
		return false
	}

	pos := y*10 + x
	if pos > 99 {
		return false
	}

	bytePos := pos / 8
	byte := g.Grid[bytePos]

	bitPos := pos % 8
	if byte&(1<<bitPos) != 0 {
		return true
	}

	return false
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

	log.Printf("Response read %+v", response)
}
