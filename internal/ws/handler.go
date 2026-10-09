package ws

import (
	"context"
	"log"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Handler handles HTTP requests for ticket categories.
// It is strictly focused on request processing, invoking services, and rendering responses.
type Handler struct {
	// service  Service
	upgrader *websocket.Upgrader
	hub      *Hub
}

// NewHandler creates a new Handler with injected Service (DI).
func NewHandler(upgrader *websocket.Upgrader, hub *Hub) *Handler {
	return &Handler{
		// service:  service,
		upgrader: upgrader,
		hub:      hub,
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(apperror.NewInternal("Error upgrading to WebSocket", err))
		return
	}

	client := NewClient(h.hub, conn)

	h.hub.register <- client

	log.Println("Client connected via WebSocket")

	// Detach from the HTTP request context: it cancels as soon as this
	// handler returns, but the WS goroutines must live for the full
	// connection lifetime.
	wsCtx := context.WithoutCancel(c.Request.Context())

	go client.ReadPump()
	go client.WritePump()
	go h.hub.HandleMessage(wsCtx, client)
}
