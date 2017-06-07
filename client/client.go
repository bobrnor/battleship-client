package client

import (
	"git.nulana.com/bobrnor/json.git"

	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	UID     string
	RoomUID string

	seq    uint64
	seqSet bool

	jsonClient json.Client
	err        error
}

const (
	authPath    = "http://0.0.0.0:8000/auth"
	searchPath  = "http://0.0.0.0:8000/game/search"
	confirmPath = "http://0.0.0.0:8000/game/confirm"
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
