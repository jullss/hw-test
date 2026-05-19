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
	internalhttp "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
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

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

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

	addr := net.JoinHostPort(config.Listen.Host, config.Listen.Port)
	if err := server.Start(ctx, addr); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
