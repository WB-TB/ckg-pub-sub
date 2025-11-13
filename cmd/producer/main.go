package main

import (
	"log/slog"
	"os"
	"pubsub-ckg-tb/internal/app"
	"pubsub-ckg-tb/internal/app/ckg"
)

func main() {
	app, err := app.InitApp()
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// str, _ := json.MarshalIndent(app.Configurations, "", "  ")
	// log.Printf("Configuration: %s", str)

	// Ensure proper cleanup when the application exits
	defer app.Close()

	slog.Info("Application initialized successfully")
	watchMode := false
	if app.Database.GetName() == "MongoDB" {
		watchMode = true
	}

	app.RunPubSubProducer(ckg.NewCkgTransmitter(
		app.Context,
		app.Configurations,
		app.Database,
		app.PubSub,
	), watchMode)
}
