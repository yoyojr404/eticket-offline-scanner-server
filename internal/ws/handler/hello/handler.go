package hello

import (
	"context"
	"encoding/json"
	"fmt"
)

type HelloHandler struct {
}

func (h *HelloHandler) Handle(ctx context.Context, msg json.RawMessage) ([]byte, error) {
	var payload Payload
	if err := json.Unmarshal(msg, &payload); err != nil {
		return nil, err
	}

	return json.Marshal(fmt.Sprintf("hello: %s", payload.Message))
}
