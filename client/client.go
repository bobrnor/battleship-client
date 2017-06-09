package client

import (
	"git.nulana.com/bobrnor/battleship-client/grid"
	"git.nulana.com/bobrnor/json.git"

	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	UID         string
	RoomUID     string
	Battlefield grid.Grid

	seq    uint64
	seqSet bool

	jsonClient json.Client
	err        error
}

const (
	authPath     = "http://0.0.0.0:8000/auth"
	searchPath   = "http://0.0.0.0:8000/game/search"
	confirmPath  = "http://0.0.0.0:8000/game/confirm"
	startPath    = "http://0.0.0.0:8000/game/start"
	longpollPath = "http://0.0.0.0:8000/game/longpoll"
	turnPath     = "http://0.0.0.0:8000/game/turn"
)

func NewClient() *Client {
	return &Client{
		UID:        uuid.TimeOrderedUUID(),
		jsonClient: json.Client{},
	}
}

func (c *Client) Error() error {
	return c.err
}

func (c *Client) Auth() {
	request := map[string]interface{}{
		"client_uid": c.UID,
	}

	var response struct {
		Status int64 `json:"status"`
	}

	err := c.jsonClient.Post(authPath, request, &response)
	if err != nil {
		c.err = err
		return
	}

	zap.S().Infof("Response read %+v", response)

	if response.Status != 0 {
		c.err = errors.Errorf("Bad status %+v", response.Status)
	}
}

func (c *Client) SearchRoom() {
	if c.err != nil {
		return
	}

	request := map[string]interface{}{
		"client_uid": c.UID,
	}

	if c.seqSet {
		request["seq"] = c.seq
	} else {
		request["reset"] = true
	}

	var response struct {
		RoomUID string `json:"room_uid"`
		Seq     uint64 `json:"seq"`
		Status  int64  `json:"status"`
	}

	err := c.jsonClient.Post(searchPath, request, &response)
	if err != nil {
		c.err = err
		return
	}

	zap.S().Infof("Response read %+v", response)

	if response.Status != 0 {
		c.err = errors.Errorf("Bad status %+v", response.Status)
	} else if len(response.RoomUID) > 0 {
		c.RoomUID = response.RoomUID
		c.seq = response.Seq
		c.seqSet = true
	}
}

func (c *Client) ConfirmRoom() {
	if c.err != nil {
		return
	}

	request := map[string]interface{}{
		"client_uid": c.UID,
		"room_uid":   c.RoomUID,
	}

	var response struct {
		Status int64 `json:"status"`
	}

	err := c.jsonClient.Post(confirmPath, request, &response)
	if err != nil {
		c.err = err
		return
	}

	zap.S().Infof("Response read %+v", response)

	if response.Status != 0 {
		c.err = errors.Errorf("Bad status %+v", response.Status)
	}
}

func (c *Client) CreateBattlefield() {
	if c.err != nil {
		return
	}

	g, err := grid.Generate()
	if err != nil {
		c.err = err
		return
	}

	c.Battlefield = g
}

func (c *Client) StartGame() {
	if c.err != nil {
		return
	}

	c.seq = 0
	c.seqSet = false

	request := map[string]interface{}{
		"client_uid": c.UID,
		"room_uid":   c.RoomUID,
		"grid":       c.Battlefield,
	}

	var response struct {
		Status int64 `json:"status"`
	}

	err := c.jsonClient.Post(startPath, request, &response)
	if err != nil {
		c.err = err
		return
	}

	zap.S().Infof("Response read %+v", response)

	if response.Status != 0 {
		c.err = errors.Errorf("Bad status %+v", response.Status)
	}
}

func (c *Client) Longpoll() {
	if c.err != nil {
		return
	}

	request := map[string]interface{}{
		"client_uid": c.UID,
	}

	if c.seqSet {
		request["seq"] = c.seq
	} else {
		request["reset"] = true
	}

	response := map[string]interface{}{}
	err := c.jsonClient.Post(longpollPath, request, &response)
	if err != nil {
		c.err = err
		return
	}

	zap.S().Infof("Response read %+v", response)

	if status, ok := response["status"].(float64); !ok || status != 0 {
		c.err = errors.Errorf("Bad status %+v", status)
	}
}

func (c *Client) Turn() {
	if c.err != nil {
		return
	}

	request := map[string]interface{}{
		"client_uid": c.UID,
		"room_uid":   c.RoomUID,
		"X":          0,
		"Y":          1,
	}

	var response struct {
		Status int64 `json:"status"`
	}

	err := c.jsonClient.Post(turnPath, request, &response)
	if err != nil {
		c.err = err
		return
	}

	zap.S().Infof("Response read %+v", response)

	if response.Status != 0 {
		c.err = errors.Errorf("Bad status %+v", response.Status)
	}
}
