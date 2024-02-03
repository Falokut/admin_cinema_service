package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/admin_cinema_service/internal/config"
	"github.com/Falokut/admin_cinema_service/internal/repository"
	"github.com/Falokut/admin_cinema_service/internal/service"
	admin_cinema_service "github.com/Falokut/admin_cinema_service/pkg/admin_cinema_service/v1/protos"
	jaegerTracer "github.com/Falokut/admin_cinema_service/pkg/jaeger"
	"github.com/Falokut/admin_cinema_service/pkg/metrics"
	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logging.NewEntry(logging.ConsoleOutput)
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(logLevel)

	tracer, closer, err := jaegerTracer.InitJaeger(cfg.JaegerConfig)
	if err != nil {
		logger.Errorf("Shutting down, error while creating tracer %v", err)
		return
	}
	logger.Info("Jaeger connected")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(cfg.PrometheusConfig.Name)
	if err != nil {
		logger.Errorf("Shutting down, error while creating metrics %v", err)
		return
	}

	shutdown := make(chan error, 1)
	go func() {
		logger.Info("Metrics server running")
		if err := metrics.RunMetricServer(cfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Errorf("Shutting down, error while running metrics server %v", err)
			shutdown <- err
		}
	}()

	moviesConn, err := getGrpcConnection(cfg.MoviesService.Addr, cfg.MoviesService.ConnectionConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the movies service is not established %v", err)
		return
	}
	defer moviesConn.Close()

	moviesService := service.NewMoviesService(moviesConn)

	cinemaDB, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database not established %s", err.Error())
		return
	}
	defer cinemaDB.Close()

	go func() {
		healthckeckManager := healthcheck.NewHealthManager(logger.Logger, []healthcheck.HealthcheckResource{cinemaDB}, cfg.HealthcheckPort, nil)
		if err := healthckeckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Errorf("Shutting down, error while running healthcheck endpoint %s", err.Error())
			shutdown <- err
		}
	}()
	cinemaRepo := repository.NewCinemaRepository(logger.Logger, cinemaDB)
	repository := service.NewCinemaRepositoryWrapper(logger.Logger, cinemaRepo, moviesService)
	logger.Info("Service initializing")
	service := service.NewCinemaService(logger.Logger, repository)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	go func() {
		if err := s.Run(getListenServerConfig(cfg), metric, nil, nil); err != nil {
			logger.Errorf("Shutting down, error while running server %s", err.Error())
			shutdown <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case <-quit:
		break
	case <-shutdown:
		break
	}
	s.Shutdown()
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:        cfg.Listen.Mode,
		Host:        cfg.Listen.Host,
		Port:        cfg.Listen.Port,
		ServiceDesc: &admin_cinema_service.CinemaServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(admin_cinema_service.CinemaServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}

			return admin_cinema_service.RegisterCinemaServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
	}
}

func getGrpcConnection(addr string, connConfig config.ConnectionSecureConfig) (*grpc.ClientConn, error) {
	creds, err := connConfig.GetGrpcTransportCredentials()
	if err != nil {
		return nil, err
	}

	return grpc.Dial(addr, creds,
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
}
