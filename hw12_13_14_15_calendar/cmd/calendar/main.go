package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/handler"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/logger"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/middleware"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/repository"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/service"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := NewConfigFomFile(configFile)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("create config from file: %v", err))
		os.Exit(1) //nolint:gocritic
	}

	logger.SetLevel(cfg.Logger.Level)

	var db *sql.DB
	var repo repository.EventRepository
	logger.Infof("use in memory repo: %s", cfg.Storage.UseInMemory)
	if cfg.Storage.UseInMemory {
		repo = repository.NewInMemoryEventRepository()
	} else {
		if db, err = sql.Open("pgx", cfg.Storage.DataSourceName); err != nil {
			logger.Fatalf("open conn to db: %v", err)
		}
		if err = db.Ping(); err != nil {
			logger.Fatalf("ping db: %v", err)
		}
		repo = repository.NewPgEventRepository(db)
	}
	srv := service.NewEventService(repo)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorw("create listener", "error", err, "addr", addr)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if err = listener.Close(); err != nil {
			logger.Errorf("stop http server: %v", err)
		} else {
			logger.Info("close listener")
		}
		if db != nil {
			err = db.Close()
			if err != nil {
				logger.Errorf("close db: %v", err)
			} else {
				logger.Info("close db")
			}
		}
	}()

	mux := http.NewServeMux()

	logger.Info("creating handlers")
	h := handler.NewHandler(srv)
	mux.HandleFunc("/", h.HelloWorld)

	logger.Info("add request logger")
	requestLogger := middleware.NewRequestLogger(mux)

	logger.Info("starting server")
	if err = http.Serve(listener, requestLogger); err != nil {
		logger.Fatalf("start server: %v", err)
	}

}
