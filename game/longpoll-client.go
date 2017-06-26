package game

import (
	"sync"

	"log"

	json "git.nulana.com/bobrnor/json.git"
)

type LongpollClientFunc func(map[string]interface{})

type LongpollClient struct {
	sync.RWMutex

	uid string
	fn  LongpollClientFunc

	seq    uint64
	seqSet bool

	jsonClient json.Client
}

const (
	longpollPath = "http://172.25.0.3:80/longpoll"
)

func NewLongpollClient(uid string, fn LongpollClientFunc) *LongpollClient {
	c := &LongpollClient{
		uid:        uid,
		fn:         fn,
		jsonClient: json.Client{},
	}
	go c.loop()
	return c
}

func (c *LongpollClient) loop() {
	for {
		var response struct {
			Type    string                 `json:"type"`
			Seq     uint64                 `json:"seq"`
			Content map[string]interface{} `json:"content"`
			Error   ResponseError          `json:"error"`
		}
		request := c.requestMessage()

		err := c.jsonClient.Post(longpollPath, request, &response)
		if err != nil {
			log.Fatalf("Can't make lp request %+v", err.Error())
		}

		//log.Printf("LP %+v", response)

		c.Lock()
		c.seq = response.Seq
		if !c.seqSet {
			c.seqSet = true
		}
		c.Unlock()

		c.fn(response.Content)
	}
}

func (c *LongpollClient) requestMessage() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	request := map[string]interface{}{}
	request["client_uid"] = c.uid
	if c.seqSet {
		request["seq"] = c.seq
	} else {
		request["reset"] = true
	}
	return request
}
