package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/handler"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/repository"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/service"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/cfg"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/logger"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/middleware"
)

type CalendarApp struct {
	repo   repository.Repository
	srv    service.Service
	server *http.Server
	config *cfg.Config
}

func (c *CalendarApp) Start() {
	logger.Info("starting calendar")
	c.startHTTP()
}

func (c *CalendarApp) initRepository() {
	logger.Debug("creating repository")
	switch {
	case c.config.Repository.UseInMemory:
		logger.Debug("using in memory repository")
		c.repo = repository.NewInMemory()
	case len(c.config.Repository.DataSourceName) > 0:
		logger.Debug("using database repository")
		var pool *pgxpool.Pool
		var err error
		if pool, err = pgxpool.Connect(context.Background(), c.config.Repository.DataSourceName); err != nil {
			logger.Fatalf("fail to connect to database: %v", err)
		}
		defer pool.Close()
		c.repo = repository.NewPg(pool)
	default:
		logger.Fatal("there is no repository configured in config")
	}
}

func (c *CalendarApp) initService() {
	logger.Debug("creating service")
	if c.repo == nil {
		logger.Fatal("repository must be initialized")
	}
	c.srv = service.New(c.repo)
}

func (c *CalendarApp) startHTTP() {
	c.initRepository()
	c.initService()

	logger.Debug("creating handlers")
	if c.srv == nil {
		logger.Fatal("service must be initialized")
	}
	h := handler.New(c.srv)

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("GET /%s/", handler.PingPath), h.Ping)

	logger.Debug("adding request logger")
	loggedMux := middleware.NewRequestLogger(mux)

	c.server = &http.Server{
		Addr:              c.config.Server.Address,
		Handler:           loggedMux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go c.listenHTTP()
}

func (c *CalendarApp) listenHTTP() {
	if c.server == nil {
		logger.Fatalf("http server has not been initialized")
	}
	if err := c.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("http listening error address=%q: %v", c.config.Server.Address, err)
	}
}

func (c *CalendarApp) Stop(ctx context.Context) error {
	if c.server == nil {
		return nil
	}
	return c.server.Shutdown(ctx)
}

func New(config *cfg.Config) *CalendarApp {
	return &CalendarApp{config: config}
}
