package router

import (
	"github.com/yoyojr404/internal/provider"
	"github.com/yoyojr404/internal/ticket_categories"
	"github.com/yoyojr404/internal/tickets"
	"github.com/yoyojr404/internal/ws"
	"github.com/yoyojr404/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// Module defines the contract that each feature route must implement.
type Module interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

// RouteFactory is responsible for producing all registered feature route modules.
// It maps the shared middleware provider to the specific middlewares required by each feature.
type RouteFactory struct {
	handlers *provider.HTTPHandlerProvider
	mws      *middleware.Provider
}

// NewRouteFactory creates a new RouteFactory with injected handlers and middleware provider.
func NewRouteFactory(handlers *provider.HTTPHandlerProvider, mws *middleware.Provider) *RouteFactory {
	return &RouteFactory{
		handlers: handlers,
		mws:      mws,
	}
}

// CreateModules instantiates and returns all route modules in the system.
func (f *RouteFactory) CreateModules() []Module {
	return []Module{
		// ticket_categories only requires basic auth
		ticket_categories.NewRoute(f.handlers.CategoryHandler, f.mws.BasicAuth),

		// tickets requires both basic auth and rate limiting
		tickets.NewRoute(f.handlers.TicketHandler, tickets.Middlewares{
			BasicAuth: f.mws.BasicAuth,
			RateLimit: f.mws.RateLimit,
		}),
		// tickets requires both basic auth and rate limiting
		ws.NewRoute(f.handlers.WSHandler),
	}
}
