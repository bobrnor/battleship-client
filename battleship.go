package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/packer/common/uuid"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	request := map[string]interface{}{
		"client_id": uuid.TimeOrderedUUID(),
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		zap.S().Errorw("Can't marshal request",
			"err", err,
		)
		return
	}

	httpClient := http.Client{}
	r, err := http.NewRequest("POST", "http://0.0.0.0:8000/auth", bytes.NewReader(requestJSON))
	if err != nil {
		zap.S().Errorw("Can't create request",
			"err", err,
		)
		return
	}

	response, err := httpClient.Do(r)
	if err != nil {
		zap.S().Errorw("Can't read response",
			"err", err,
		)
		return
	}
	defer response.Body.Close()

	type authResponseData struct {
		SessionID string `json:"session_id"`
	}
	var data authResponseData
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		zap.S().Errorw("Can't decode response",
			"err", err,
		)
		return
	}

	zap.S().Infow("Response read",
		"response", data,
	)
}
