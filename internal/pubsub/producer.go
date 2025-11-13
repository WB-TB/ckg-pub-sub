package pubsub

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Transmitter interface {
	Watch(ctx context.Context)
	Produce(ctx context.Context) error
}

func (c *Client) StartProducer(ctx context.Context, transmitter Transmitter, watchMode bool) {
	c.Transmitter = transmitter

	// Create a context that can be cancelled
	producerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	if c.Transmitter != nil {
		// Start the transmitter in a goroutine
		go func() {
			slog.Info("Starting message producer...")
			if watchMode {
				c.Transmitter.Watch(producerCtx)
			} else {
				c.Transmitter.Produce(producerCtx)
			}
		}()

		// Wait for shutdown signal
		select {
		case <-signalChan:
			slog.Info("Received shutdown signal, cancelling context...")
			cancel()
		case <-producerCtx.Done():
			slog.Info("Context cancelled, shutting down...")
		}

		slog.Info("Producer stopped gracefully")
	}
}
