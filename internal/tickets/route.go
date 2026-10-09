package tickets

import (
	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	BasicAuth gin.HandlerFunc
	RateLimit gin.HandlerFunc
}

type Route struct {
	handler *Handler
	mws     Middlewares
}

func NewRoute(handler *Handler, mws Middlewares) *Route {
	return &Route{
		handler: handler,
		mws:     mws,
	}
}

func (r *Route) RegisterRoutes(rg *gin.RouterGroup) {
	tickets := rg.Group("/tickets")
	{
		tickets.GET("", r.handler.GetAll)
		tickets.GET("/:id", r.mws.RateLimit, r.handler.GetByID)
	}
}
