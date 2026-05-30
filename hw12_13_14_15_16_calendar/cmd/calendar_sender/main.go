package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/queue"
	"go.yaml.in/yaml/v3"
)

type SenderConfig struct {
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
	RabbitMQ struct {
		URL   string `yaml:"url"`
		Queue string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/sender_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	var config SenderConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		f.Close()
		log.Fatalf("failed to decode config: %v", err)
	}
	f.Close()

	logg := logger.New(config.Logger.Level)
	logg.Info("sender is starting...")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rabbitClient, err := queue.NewRabbitClient(config.RabbitMQ.URL, config.RabbitMQ.Queue)
	if err != nil {
		logg.Error("failed to connect to rabbitmq", "err", err)
		os.Exit(1)
	}
	defer rabbitClient.Close()

	var consumer queue.Consumer = rabbitClient
	msgs, err := consumer.Consume(ctx)
	if err != nil {
		logg.Error("failed to start consuming", "err", err)
		os.Exit(1)
	}

	logg.Info("sender is running and waiting for notifications...")

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				logg.Info("rabbitmq channel closed, stopping sender")
				return
			}

			logg.Info("NOTIFICATION SENT",
				"event_id", msg.EventID,
				"title", msg.Title,
				"user_id", msg.UserID,
				"date", msg.Date.Format("2006-01-02 15:04:05"),
			)
		case <-ctx.Done():
			logg.Info("sender is stopping gracefully...")
			return
		}
	}
}
