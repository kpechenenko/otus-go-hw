package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/app"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/cfg"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/logger"
)

var (
	configFile  string
	calendarApp *app.CalendarApp
	exit        = make(chan bool)
)

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar/config.yaml", "path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	go listenSignals()
	go start()

	<-exit
	logger.Infof("see you soon...")
}

func start() {
	var appConfig *cfg.Config
	var err error
	if appConfig, err = cfg.ReadConfigFromYamlFile(configFile); err != nil {
		logger.Fatalf("read config from yaml file: %v", err)
	}
	if err = appConfig.Check(); err != nil {
		logger.Fatalf("invalid сonfig content: %v", err)
	}
	logger.SetLevel(appConfig.Logger.Level)
	calendarApp = app.New(appConfig)
	calendarApp.Start()
}

func listenSignals() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-signalChan

		switch s {
		case syscall.SIGHUP:
			logger.Info("SIGHUP received")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			err := calendarApp.Stop(ctx)
			if err != nil {
				logger.Fatalf("can't stop calendar with SIGHUP: %v", err)
			}
			cancel()

			start()
		case syscall.SIGINT, syscall.SIGTERM:
			logger.Info("SIGTERM or SIGINT received")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			err := calendarApp.Stop(ctx)
			if err != nil {
				logger.Fatalf("can't stop calendar with SIGTERM or SIGINT: %v", err)
			}
			cancel()

			exit <- true
		}
	}
}
