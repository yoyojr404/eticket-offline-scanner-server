package handler

import (
	"context"
	"encoding/json"
)

type Handler interface {
	// executing handler
	Handle(ctx context.Context, msg json.RawMessage) ([]byte, error)
}
