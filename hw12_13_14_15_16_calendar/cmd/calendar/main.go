package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/pb"
	internalgrpc "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"google.golang.org/grpc"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	var repo storage.Repository
	var dbRepo *sqlstorage.Storage
	if config.Storage.Type == "sql" {
		dbRepo = sqlstorage.New(nil)

		if err := dbRepo.Connect(ctx, config.Storage.DBURL); err != nil {
			logg.Error("failed to connect to db", "err", err)
			os.Exit(1)
		}

		if err := dbRepo.Migrate(ctx); err != nil {
			logg.Error("failed to run migrations", "err", err)
			os.Exit(1)
		}

		repo = dbRepo
		logg.Info("using sql storage")
	} else {
		repo = memorystorage.New()
		logg.Info("using memory storage")
	}

	calendar := app.New(logg, repo)

	server := internalhttp.NewServer(logg, calendar)

	grpcImpl := internalgrpc.NewServer(logg, calendar)
	grpcLoggerInterceptor := func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		if err != nil {
			logg.Error("gRPC fail", "method", info.FullMethod, "duration", duration, "err", err)
		} else {
			logg.Info("gRPC success", "method", info.FullMethod, "duration", duration)
		}

		return resp, err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcLoggerInterceptor),
	)
	pb.RegisterCalendarServiceServer(grpcServer, grpcImpl)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("grpc server stopping")
		grpcServer.GracefulStop()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if dbRepo != nil {
			logg.Info("closing database connection...")
			if err := dbRepo.Close(context.Background()); err != nil {
				logg.Error("failed to close db", "err", err)
			}
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		grpcAddr := net.JoinHostPort(config.Listen.Host, config.Listen.GRPCPort)

		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logg.Error("failed to listen for grpc", "err", err)
			return
		}

		logg.Info("grpc server starting", "addr", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			logg.Error("grpc server failed", "err", err)
		}
	}()

	addr := net.JoinHostPort(config.Listen.Host, config.Listen.Port)
	if err := server.Start(ctx, addr); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
