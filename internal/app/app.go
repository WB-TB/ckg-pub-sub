package app

import (
	"context"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	database "pubsub-ckg-tb/internal/db"

	pubsubInternal "pubsub-ckg-tb/internal/pubsub"

	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Configurations *config.Configurations
	Database       *mongo.Client
	PubSub         *pubsubInternal.Client
}

func InitApp() (*App, error) {
	// Load env variables dari file .env
	cfg := config.GetConfig()

	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(cfg.App.LogLevel)); err != nil {
		logLevel = slog.LevelInfo // Default jika parsing gagal
		slog.Warn("Gagal parsing log level dari config, menggunakan default 'info'", "error", err)
	}

	slog.SetLogLoggerLevel(logLevel)

	// Initialize database connection
	dbConn := database.GetConnection(cfg.Database)

	// Initialize PubSub client
	ctx := context.Background()
	pubsubClient, err := pubsubInternal.NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	defer pubsubClient.Close()

	return &App{
		Configurations: cfg,
		Database:       dbConn,
		PubSub:         pubsubClient,
	}, nil
}

func (a *App) RunPubSubConsumer(receiver pubsubInternal.Receiver) {
	ctx := context.Background()

	// Ensure topic and subscription exist
	if !a.PubSub.EnsureTopicExists(ctx) {
		slog.Error("Topic tidak ditemukan")
		return
	}

	if !a.PubSub.EnsureSubscriptionExists(ctx) {
		slog.Error("Subscription tidak ditemukan")
		return
	}

	// Start consuming messages in a loop
	a.PubSub.StartConsumer(ctx, receiver)
}

func (a *App) RunPubSubProducer(transmitter pubsubInternal.Transmitter) {
	ctx := context.Background()

	// Start procuce messages one time
	a.PubSub.StartProducer(ctx, transmitter)
}

func (a *App) Close() {
	slog.Info("Closing application resources...")
	database.CloseConnection()
	a.PubSub.Close()
}
