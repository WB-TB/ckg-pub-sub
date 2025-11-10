package pubsub

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub/v2"
)

type Receiver interface {
	Consume(ctx context.Context, messages []*pubsub.Message) (map[string]bool, error)
}

func (c *Client) StartConsumer(ctx context.Context, receiver Receiver) {
	c.Receiver = receiver

	// Create a context that can be cancelled
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	slog.Info("Starting message consumer...")

	// Main consumer loop
	for {
		select {
		case <-consumerCtx.Done():
			slog.Info("Consumer context cancelled, shutting down...")
			return
		case <-signalChan:
			slog.Info("Received termination signal, shutting down...")
			return
		default:
			// Pull messages from PubSub
			messages, err := c.PullMessages(consumerCtx, c.Config.Consumer.MaxMessagesPerPull)
			if err != nil {
				slog.Error("Error pulling messages", "error", err)
				// Wait before retrying
				time.Sleep(c.Config.Consumer.SleepTimeBetweenPulls)
				continue
			}

			if len(messages) == 0 {
				// No messages received, wait before pulling again
				time.Sleep(c.Config.Consumer.SleepTimeBetweenPulls)
				continue
			}

			slog.Debug("Received messages", "count", len(messages))

			if c.Receiver != nil {
				go c.Receiver.Consume(ctx, messages)
			}
		}
	}
}
