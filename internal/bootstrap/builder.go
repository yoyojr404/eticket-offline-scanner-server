package bootstrap

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/yoyojr404/internal/provider"
	"github.com/yoyojr404/internal/router"
	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/config"
	"github.com/yoyojr404/pkg/database"
	"github.com/yoyojr404/pkg/middleware"
	"github.com/yoyojr404/pkg/response"
)

// AppBuilder provides a fluent builder interface for assembling the App instance.
type AppBuilder struct {
	cfg          *config.Config
	db           *sqlx.DB
	provider     *provider.AppProvider
	preflight    *provider.PreflightProvider
	background   *provider.BackroundProvider
	middlewares  []gin.HandlerFunc
	apiPrefix    string
	routeModules []router.Module
	err          error
}

// NewAppBuilder initializes a new AppBuilder with sensible defaults.
func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		apiPrefix: "/api/v1",
	}
}

// WithConfig sets the application configuration.
func (b *AppBuilder) WithConfig(cfg *config.Config) *AppBuilder {
	b.cfg = cfg
	return b
}

// WithDatabase sets an existing database connection.
func (b *AppBuilder) WithDatabase(db *sqlx.DB) *AppBuilder {
	b.db = db
	return b
}

// WithProvider sets the delivery-agnostic AppProvider.
func (b *AppBuilder) WithProvider(p *provider.AppProvider) *AppBuilder {
	b.provider = p
	return b
}

// WithMiddlewares appends custom global middlewares to the router.
func (b *AppBuilder) WithMiddlewares(mws ...gin.HandlerFunc) *AppBuilder {
	b.middlewares = append(b.middlewares, mws...)
	return b
}

// WithAPIPrefix sets the API versioning group prefix (defaults to "/api/v1").
func (b *AppBuilder) WithAPIPrefix(prefix string) *AppBuilder {
	b.apiPrefix = prefix
	return b
}

// WithRouteModules registers route modules created by the RouteFactory.
func (b *AppBuilder) WithRouteModules(modules ...router.Module) *AppBuilder {
	b.routeModules = append(b.routeModules, modules...)
	return b
}

// WithPreflight sets the PreflightProvider whose Runner will execute before the server starts.
func (b *AppBuilder) WithPreflight(p *provider.PreflightProvider) *AppBuilder {
	b.preflight = p
	return b
}

// WithBackground sets the BackroundProvider whose Runner will execute before the server starts.
func (b *AppBuilder) WithBackground(p *provider.BackroundProvider) *AppBuilder {
	b.background = p
	return b
}

// Build compiles all dependencies, wires router and modules, and returns the assembled App.
func (b *AppBuilder) Build() (*App, error) {
	if b.err != nil {
		return nil, b.err
	}

	// 1. Resolve Config
	if b.cfg == nil {
		b.cfg = config.Load()
	}
	gin.SetMode(b.cfg.GinMode)

	// 2. Resolve Database
	if b.db == nil && b.cfg.DBDriver != "memory" {
		dbConn, err := database.InitDB(b.cfg)
		if err != nil {
			return nil, fmt.Errorf("builder: database initialization failed: %w", err)
		}
		b.db = dbConn
	}

	// 3. Resolve AppProvider
	if b.provider == nil {
		b.provider = provider.NewAppProvider(b.cfg, b.db)
	}

	// 4. Build Gin Router & Middlewares
	r := gin.New()
	r.Use(middleware.NewBlacklistMap(b.provider.IPWhitelist).Blacklist())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	for _, mw := range b.middlewares {
		r.Use(mw)
	}

	// Standard response envelope for unhandled routes
	r.NoRoute(func(c *gin.Context) {
		response.Error(c, apperror.NewNotFound("Resource or endpoint not found", nil))
	})
	r.NoMethod(func(c *gin.Context) {
		response.Error(c, apperror.New(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil))
	})

	// Health Check Route
	r.GET("/health", func(c *gin.Context) {
		response.OK(c, "Service is healthy", gin.H{
			"status": "UP",
			"time":   time.Now().UTC(),
		})
	})

	// 5. Mount Route Modules to API Group
	apiGroup := r.Group(b.apiPrefix)
	for _, module := range b.routeModules {
		module.RegisterRoutes(apiGroup)
	}

	// 6. Build HTTP Server
	srv := &http.Server{
		Addr:         ":" + b.cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		Config:     b.cfg,
		DB:         b.db,
		Provider:   b.provider,
		Preflight:  b.preflight,
		Background: b.background,
		Router:     r,
		HTTPServer: srv,
	}, nil
}
