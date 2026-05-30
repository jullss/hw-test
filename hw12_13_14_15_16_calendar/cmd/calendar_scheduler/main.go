package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sqlstorage "github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage/sql"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/queue"
	"go.yaml.in/yaml/v3"
)

type SchedulerConfig struct {
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
	Storage struct {
		DBURL string `yaml:"dbUrl"`
	} `yaml:"storage"`
	RabbitMQ struct {
		URL   string `yaml:"url"`
		Queue string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
	Scheduler struct {
		ScanInterval string `yaml:"scan_interval"`
	} `yaml:"scheduler"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	var config SchedulerConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		f.Close()
		log.Fatalf("failed to decode config: %v", err)
	}
	f.Close()

	logg := logger.New(config.Logger.Level)
	logg.Info("scheduler is starting...")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbRepo := sqlstorage.New(nil)
	if err := dbRepo.Connect(ctx, config.Storage.DBURL); err != nil {
		logg.Error("scheduler failed to connect to db", "err", err)
		os.Exit(1)
	}
	defer dbRepo.Close(context.Background())

	rabbitClient, err := queue.NewRabbitClient(config.RabbitMQ.URL, config.RabbitMQ.Queue)
	if err != nil {
		logg.Error("scheduler failed to connect to rabbitmq", "err", err)
		os.Exit(1)
	}
	defer rabbitClient.Close()

	var publisher queue.Publisher = rabbitClient

	scanDuration, err := time.ParseDuration(config.Scheduler.ScanInterval)
	if err != nil {
		scanDuration = 5 * time.Second
	}

	ticker := time.NewTicker(scanDuration)
	defer ticker.Stop()

	logg.Info("scheduler is running...")

	for {
		select {
		case <-ticker.C:
			logg.Info("scheduler tick: scanning for notifications and cleaning old events")

			oneYearAgo := time.Now().AddDate(-1, 0, 0)

			if err := dbRepo.CleanOldEvents(ctx, oneYearAgo); err != nil {
				logg.Error("failed to clean old events", "err", err)
			} else {
				logg.Info("old events cleaned up successfully")
			}

			eventsToNotify, err := dbRepo.GetEventsToNotify(ctx)
			if err != nil {
				logg.Error("failed to fetch events to notify", "err", err)
				continue
			}

			for _, ev := range eventsToNotify {
				notification := queue.Notification{
					EventID: ev.ID,
					Title:   ev.Title,
					Date:    ev.StartTime,
					UserID:  ev.UserID,
				}

				if err := publisher.Publish(ctx, notification); err != nil {
					logg.Error("failed to publish notification", "event_id", ev.ID, "err", err)
				} else {
					logg.Info("notification pushed to rabbitmq", "event_id", ev.ID)
				}
			}

		case <-ctx.Done():
			logg.Info("scheduler is stopping gracefully...")
			return
		}
	}
}
