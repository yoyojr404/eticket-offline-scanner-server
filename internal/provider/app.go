package provider

import (
	"net/http"

	workerscanDeviceservice "github.com/yoyojr404/internal/background/scandevice"
	"github.com/yoyojr404/internal/common/datastore"
	preflighrepo "github.com/yoyojr404/internal/preflight/repository"
	"github.com/yoyojr404/internal/preflight/repository/ticketcategories"
	preflightticket "github.com/yoyojr404/internal/preflight/repository/tickets"
	"github.com/yoyojr404/internal/preflight/repository/ticketverifications"
	prefligtseedservice "github.com/yoyojr404/internal/preflight/service/seed"
	"github.com/yoyojr404/internal/ticket_categories"
	"github.com/yoyojr404/internal/tickets"
	"github.com/yoyojr404/internal/ws"
	"github.com/yoyojr404/internal/ws/handler"
	"github.com/yoyojr404/internal/ws/handler/hello"
	"github.com/yoyojr404/internal/ws/handler/ticketverification"
	"github.com/yoyojr404/pkg/config"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

type AppProvider struct {
	Config *config.Config
	DB     *sqlx.DB

	CategoryRepo   ticket_categories.Repository
	TicketRepo     tickets.Repository
	SeedRepoBucket prefligtseedservice.ServiceRepoBucket

	CategoryService      ticket_categories.Service
	TicketService        tickets.Service
	PreflightSeedService prefligtseedservice.SeedService
	ScanDeviceWorker     workerscanDeviceservice.DeviceScan

	WSHub      *ws.Hub
	Upgrader   *websocket.Upgrader
	Dispatcher *handler.Dispatcher

	// IP Whitelist, this is used for any Server listens to connection
	IPWhitelist map[string]bool
}

func NewAppProvider(cfg *config.Config, db *sqlx.DB) *AppProvider {
	var catRepo ticket_categories.Repository
	var ticketRepo tickets.Repository

	if db == nil {
		panic("database connection is nil")
	}

	transactor := datastore.NewTransactor(db)

	verificationRepo := ticketverification.NewSQLRepository(db)
	verificationService := ticketverification.NewService(transactor, verificationRepo)
	dispatcher := handler.NewDispatcher()
	dispatcher.Register(handler.CmdHello, &hello.HelloHandler{})
	dispatcher.Register(handler.CmdVerifyTicket, ticketverification.NewTicketVerificationHandler(verificationService))

	catRepo = ticket_categories.NewSQLRepository(db)
	ticketRepo = tickets.NewSQLRepository(db)
	seedRepoBucket := prefligtseedservice.ServiceRepoBucket{
		DDLRepo:                   preflighrepo.NewDDLRepo(db),
		TicketCategoriesReader:    ticketcategories.NewCSVReaderRepository(cfg.TicketCategoriesPath),
		TicketCategoriesWriter:    ticketcategories.NewDBWriterRepository(db),
		TicketsReader:             preflightticket.NewCSVReaderRepository(cfg.TicketsPath),
		TicketsWriter:             preflightticket.NewDBWriterRepository(db),
		TicketVerificationsReader: ticketverifications.NewCSVReaderRepository(cfg.TicketVerificationsPath),
		TicketVerificationsWriter: ticketverifications.NewDBWriterRepository(db),
	}

	catService := ticket_categories.NewService(catRepo)
	ticketService := tickets.NewService(ticketRepo)
	seedService := prefligtseedservice.NewService(
		transactor,
		cfg.BatchSeedSize,
		seedRepoBucket,
	)
	ipWhiteList := make(map[string]bool)
	// !for development only
	scanDeviceService := workerscanDeviceservice.NewPortScanService(
		9000,
		9010,
		cfg.CAPath,
		cfg.ScannerKeyPath,
		cfg.ScannerCertPath,
		ipWhiteList,
	)

	return &AppProvider{
		Config:               cfg,
		DB:                   db,
		CategoryRepo:         catRepo,
		TicketRepo:           ticketRepo,
		SeedRepoBucket:       seedRepoBucket,
		PreflightSeedService: seedService,
		ScanDeviceWorker:     scanDeviceService,
		CategoryService:      catService,
		TicketService:        ticketService,
		WSHub:                ws.NewHub(dispatcher),
		Dispatcher:           dispatcher,
		IPWhitelist:          ipWhiteList,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
