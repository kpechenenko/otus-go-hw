package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/handler"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/logger"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/middleware"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/repository"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/service"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	var cfg *Config
	var err error

	if cfg, err = readConfigFromYaml(configFile); err != nil {
		logger.Fatalf("read config from yaml file: %v", err)
	}
	if err = cfg.check(); err != nil {
		logger.Fatalf("invalid сonfig content: %v", err)
	}
	logger.SetLevel(cfg.Logger.Level)
	appCtx := context.Background()

	var repo repository.Repository
	if cfg.Repository.UseInMemory {
		logger.Debug("using in memory repository")
		repo = repository.NewInMemory()
	} else {
		logger.Debug("using database repository")
		var pool *pgxpool.Pool
		if pool, err = pgxpool.Connect(appCtx, cfg.Repository.DataSourceName); err != nil {
			logger.Fatalf("fail to connect to database: %v", err)
		}
		defer pool.Close()
		repo = repository.NewPg(pool)
	}
	srv := service.New(repo)

	logger.Debug("creating handlers")
	h := handler.New(srv)

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("GET /%s/", handler.PingPath), h.Ping)

	logger.Debug("adding  request logger")
	loggedMux := middleware.NewRequestLogger(mux)

	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: loggedMux,
	}

	ctx, stop := signal.NotifyContext(appCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	go func() {
		<-ctx.Done()
		logger.Debug("receive stop signal")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = server.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("failed to shutdown server: %v", err)
		}
		logger.Info("graceful shutdown complete")
	}()

	logger.Infof("starting http server at %s", cfg.Server.Address)
	if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("listen and serve error: %v", err)
	}
	logger.Info("stopped")
}
