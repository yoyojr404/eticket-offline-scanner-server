package ticket_categories

import (
	"github.com/gin-gonic/gin"
)

// Route manages the route registration for the ticket_categories feature.
type Route struct {
	handler        *Handler
	authMiddleware gin.HandlerFunc
}

// NewRoute creates a new Route instance with injected handler and middleware.
func NewRoute(handler *Handler, authMiddleware gin.HandlerFunc) *Route {
	return &Route{
		handler:        handler,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes mounts all category routes onto the supplied Gin router group.
func (r *Route) RegisterRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/ticket-categories")
	{
		// Public read routes
		categories.GET("", r.handler.GetAll)
		categories.GET("/:id", r.handler.GetByID)

	}
}
