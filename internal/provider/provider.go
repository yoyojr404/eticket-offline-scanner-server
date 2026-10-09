package provider

import (
	"github.com/yoyojr404/internal/background"
	"github.com/yoyojr404/internal/preflight"
	"github.com/yoyojr404/internal/ticket_categories"
	"github.com/yoyojr404/internal/tickets"
	"github.com/yoyojr404/internal/ws"
)

// HTTPHandlerProvider holds all HTTP controllers/handlers required by the REST API server.
type HTTPHandlerProvider struct {
	CategoryHandler *ticket_categories.Handler
	TicketHandler   *tickets.Handler
	WSHandler       *ws.Handler
}

// ProvideHTTPHandlers injects the required services into HTTP handlers.
func ProvideHTTPHandlers(p *AppProvider) *HTTPHandlerProvider {
	return &HTTPHandlerProvider{
		CategoryHandler: ticket_categories.NewHandler(p.CategoryService),
		TicketHandler:   tickets.NewHandler(p.TicketService),
		WSHandler:       ws.NewHandler(p.Upgrader, p.WSHub),
	}
}

type PreflightProvider struct {
	Runner *preflight.Runner
}

func ProvidePreflight(p *AppProvider) *PreflightProvider {
	return &PreflightProvider{
		Runner: preflight.NewRunner(p.PreflightSeedService),
	}
}

type BackroundProvider struct {
	Runner *background.RunnerGroup
}

func ProvideBackround(p *AppProvider) *BackroundProvider {
	return &BackroundProvider{
		Runner: background.NewRunnerGroup(),
	}
}
