package app

import (
	"context"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	database "pubsub-ckg-tb/internal/db"
	"pubsub-ckg-tb/internal/db/connection"

	pubsubInternal "pubsub-ckg-tb/internal/pubsub"
)

type App struct {
	Configurations *config.Configurations
	Context        context.Context
	Database       connection.DatabaseConnection
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

	ctx := context.Background()

	// Initialize database connection
	dbConn := database.GetConnection(&cfg.Database)
	dbConn.Connect(ctx)

	// Initialize PubSub client
	pubsubClient, err := pubsubInternal.NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Configurations: cfg,
		Context:        ctx,
		Database:       dbConn,
		PubSub:         pubsubClient,
	}, nil
}

func (a *App) RunPubSubConsumer(receiver pubsubInternal.Receiver) {
	// Ensure topic and subscription exist
	if !a.PubSub.EnsureTopicExists(a.Context) {
		slog.Error("Topic tidak ditemukan")
		return
	}

	if !a.PubSub.EnsureSubscriptionExists(a.Context) {
		slog.Error("Subscription tidak ditemukan")
		return
	}

	// Start consuming messages in a loop
	a.PubSub.StartConsumer(a.Context, receiver)
}

func (a *App) RunPubSubProducer(transmitter pubsubInternal.Transmitter, watchMode bool) {
	// Start procuce messages one time
	a.PubSub.StartProducer(a.Context, transmitter, watchMode)
}

func (a *App) Close() {
	slog.Info("Closing application resources...")
	database.CloseConnection(a.Context)
	a.PubSub.Close()
}
