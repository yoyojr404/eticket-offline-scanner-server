package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/yoyojr404/internal/provider"
	"github.com/yoyojr404/internal/router"
	"github.com/yoyojr404/pkg/config"
	"github.com/yoyojr404/pkg/database"
	"github.com/yoyojr404/pkg/middleware"
)

// App manages the lifecycle, dependencies, and servers of the application.
type App struct {
	Config     *config.Config
	DB         *sqlx.DB
	Provider   *provider.AppProvider
	Preflight  *provider.PreflightProvider
	Background *provider.BackroundProvider
	Router     *gin.Engine
	HTTPServer *http.Server
}

// NewApp wires the application using the Factory and Builder patterns.
func NewApp() (*App, error) {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Database Connection via dedicated database package
	db, err := database.InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: database initialization failed: %w", err)
	}

	// 3. Resolve Delivery-Agnostic Dependencies (Core Services & Repositories)
	appProvider := provider.NewAppProvider(cfg, db)

	// 4. Resolve HTTP Handlers via Provider
	httpHandlers := provider.ProvideHTTPHandlers(appProvider)
	go appProvider.WSHub.Run()

	// 5. Resolve Preflight Provider
	preflightProvider := provider.ProvidePreflight(appProvider)

	backgroundProvider := provider.ProvideBackround(appProvider)
	backgroundProvider.Runner.Add(appProvider.ScanDeviceWorker)

	// 6. Resolve Route Middlewares via Middleware Provider
	mws := middleware.NewProvider(cfg)

	// 7. Produce Feature Route Modules via RouteFactory with injected Middlewares
	routeFactory := router.NewRouteFactory(httpHandlers, mws)

	// 8. Assemble Application via AppBuilder
	return NewAppBuilder().
		WithConfig(cfg).
		WithDatabase(db).
		WithProvider(appProvider).
		WithPreflight(preflightProvider).
		WithBackground(backgroundProvider).
		WithRouteModules(routeFactory.CreateModules()...).
		Build()
}

// Run starts the HTTP server and handles graceful shutdown upon OS signals.
func (a *App) Run() error {
	// 1. Execute preflight tasks (seeding, schema setup, etc.) before opening any port
	if a.Preflight != nil {
		log.Println("[INFO] Running preflight tasks...")
		if err := a.Preflight.Runner.Run(context.Background()); err != nil {
			return fmt.Errorf("preflight failed: %w", err)
		}
		log.Println("[INFO] Preflight completed successfully.")
	}

	// 2. Execute background tasks (workers, etc.) before opening any port
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	if a.Background != nil {
		log.Println("[INFO] Running background tasks...")
		if err := a.Background.Runner.RunAll(workerCtx); err != nil {
			return fmt.Errorf("background failed: %w", err)
		}
		log.Println("[INFO] Background spawning successfully.")
	}

	serverErrCh := make(chan error, 1)

	go func() {
		log.Printf("[INFO] Server listening on %s", a.HTTPServer.Addr)
		if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	// Wait for OS interrupt signal or server launch failure
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		return fmt.Errorf("server failure: %w", err)
	case sig := <-quit:
		log.Printf("[INFO] Received signal %v. Initiating graceful shutdown...", sig)
	}

	// Graceful shutdown with 5-second deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Server forced to shutdown: %v", err)
	}

	// Close database connection pool
	if a.DB != nil {
		if err := database.Close(a.DB); err != nil {
			log.Printf("[ERROR] Error closing database: %v", err)
		}
	}

	workerCancel()
	a.Background.Runner.ShutdownAll(ctx)

	log.Println("[INFO] Application shutdown complete.")
	return nil
}
