package ws

import (
	"github.com/gin-gonic/gin"
)

// Route manages the route registration for the ticket_categories feature.
type Route struct {
	handler *Handler
}

// NewRoute creates a new Route instance with injected handler and middleware.
func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// RegisterRoutes mounts all category routes onto the supplied Gin router group.
func (r *Route) RegisterRoutes(rg *gin.RouterGroup) {
	wss := rg.Group("/ws")
	{
		wss.GET("", r.handler.HandleWebSocket)

	}
}
