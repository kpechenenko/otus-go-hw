package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/cfg"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/gateway"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/handler"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/repository"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/service"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/logger"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/middleware"
	desc "github.com/kpechenenko/hw12_13_14_15_calendar/calendar/pkg/api/calendar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CalendarApp struct {
	calendarRepo    repository.Repository
	calendarService service.Service
	gwService       *gateway.Service
	gwServer        *http.Server
	grpcServer      *grpc.Server
	config          *cfg.Config
}

func (c *CalendarApp) Start() {
	logger.Info("starting calendar")
	c.startGRPCServer()
	c.startGatewayServer()
}

func (c *CalendarApp) initCalendarRepository() {
	if c.calendarRepo != nil {
		logger.Debug("repository already initialized")
		return
	}
	logger.Debug("creating repository")
	switch {
	case c.config.Repository.UseInMemory:
		logger.Debug("using in memory repository")
		c.calendarRepo = repository.NewInMemory()
	case len(c.config.Repository.DataSourceName) > 0:
		logger.Debug("using database repository")
		var pool *pgxpool.Pool
		var err error
		if pool, err = pgxpool.Connect(context.Background(), c.config.Repository.DataSourceName); err != nil {
			logger.Fatalf("fail to connect to database: %v", err)
		}
		defer pool.Close()
		c.calendarRepo = repository.NewPg(pool)
	default:
		logger.Fatal("there is no repository configured in config")
	}
}

func (c *CalendarApp) initCalendarService() {
	if c.calendarService != nil {
		logger.Debug("calendar service already initialized")
		return
	}
	logger.Debug("creating calendar service")
	if c.calendarRepo == nil {
		c.initCalendarRepository()
	}
	c.calendarService = service.New(c.calendarRepo)
}

func (c *CalendarApp) initGatewayService() {
	if c.gwService != nil {
		logger.Debug("gateway service already initialized")
		return
	}
	logger.Debug("creating gateway service")
	if c.calendarService == nil {
		c.initCalendarService()
	}
	c.gwService = gateway.New(c.calendarService)
}

func (c *CalendarApp) Stop(ctx context.Context) error {
	logger.Info("stopping calendar")
	var err1, err2 error
	if c.gwServer != nil {
		err1 = c.gwServer.Shutdown(ctx)
	}
	if c.grpcServer != nil {
		err2 = c.gwServer.Shutdown(ctx)
	}
	return errors.Join(err1, err2)
}

func New(config *cfg.Config) *CalendarApp {
	return &CalendarApp{config: config}
}

func (c *CalendarApp) startGRPCServer() {
	c.initGatewayService()

	lis, err := net.Listen("tcp", c.config.GRPCServer.Address)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.GrpcGatewayLoggingRequest),
	)
	desc.RegisterCalendarServer(grpcServer, c.gwService)

	go func() {
		logger.Infof("grpc server listening at %v", lis.Addr())
		if err = grpcServer.Serve(lis); err != nil {
			logger.Fatalf("grpc server failed to serve: %v", err)
		}
	}()
}

func (c *CalendarApp) startGatewayServer() {
	conn, err := grpc.Dial(
		c.config.GRPCServer.Address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf("failed to dial grpc: %v", err)
	}
	gwMux := runtime.NewServeMux()
	if err = desc.RegisterCalendarHandler(context.Background(), gwMux, conn); err != nil {
		logger.Fatalf("failed to register gateway mux: %v", err)
	}

	h := handler.New(c.calendarService)
	if err = gwMux.HandlePath(http.MethodGet, fmt.Sprintf("/%s", handler.PingPath), h.PingWithParams); err != nil {
		logger.Errorf("fail to register ping handler to gateway mux: %v", err)
	}
	loggedMux := middleware.NewRequestLogger(gwMux)

	gwServer := &http.Server{
		Addr:              c.config.HTTPServer.Address,
		Handler:           loggedMux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		logger.Infof("gateway server listening at %s", gwServer.Addr)
		if err = gwServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("gateway server failed to serve: %v", err)
		}
	}()
}

// todo
// дописать в make файлинструкции по хагрузке бафа и остальных бинарников
// signup в мейне мб обработать по другому и перезапускать с новым конфигом
