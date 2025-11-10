package pubsub

import (
	"context"
	"fmt"
	"log/slog"

	"pubsub-ckg-tb/internal/config"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	Client       *pubsub.Client
	ProjectID    string
	Topic        string
	Subscription string
	Config       *config.Configurations
	Receiver     Receiver
	Transmitter  Transmitter
}

func NewClient(ctx context.Context, cfg *config.Configurations) (*Client, error) {
	// Set credentials if file exists
	var opts []option.ClientOption
	if cfg.GoogleCloud.CredentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.GoogleCloud.CredentialsPath))
	}

	// Create pubsub client
	client, err := pubsub.NewClient(ctx, cfg.GoogleCloud.ProjectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %v", err)
	}
	defer client.Close()

	slog.Debug("PubSub client initialized for",
		"project", cfg.GoogleCloud.ProjectID,
		"topic", cfg.PubSub.Topic,
		"subscription", cfg.PubSub.Subscription)

	return &Client{
		Client:       client,
		ProjectID:    cfg.GoogleCloud.ProjectID,
		Topic:        cfg.PubSub.Topic,
		Subscription: cfg.PubSub.Subscription,
		Config:       cfg,
	}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

// EnsureTopicExists checks if the topic exists
func (c *Client) EnsureTopicExists(ctx context.Context) bool {
	return c.GetTopicInfo(ctx) != nil
}

// EnsureSubscriptionExists checks if the subscription exists
func (c *Client) EnsureSubscriptionExists(ctx context.Context) bool {
	return c.GetSubscriptionInfo(ctx) != nil
}

// GetTopicInfo returns information about the topic
func (c *Client) GetTopicInfo(ctx context.Context) *pubsubpb.Topic {
	req := &pubsubpb.ListTopicsRequest{
		Project: fmt.Sprintf("projects/%s", c.ProjectID),
	}
	it := c.Client.TopicAdminClient.ListTopics(ctx, req)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			continue
		}

		slog.Debug("Found topic", "name", topic.Name)
		if topic.Name == fmt.Sprintf("projects/%s/topics/%s", c.ProjectID, c.Topic) {
			return topic
		}
	}

	return nil
}

// GetSubscriptionInfo returns information about the subscription
func (c *Client) GetSubscriptionInfo(ctx context.Context) *pubsubpb.Subscription {
	exists := c.ListSubscriptions(ctx)

	for _, sub := range exists {
		if sub != nil && sub.Name == fmt.Sprintf("projects/%s/subscriptions/%s", c.ProjectID, c.Subscription) {
			return sub
		}
	}

	return nil
}

// PublishMessage publishes a message to the topic
func (c *Client) PublishMessage(ctx context.Context, data []byte, attributes map[string]string) (string, error) {
	publisher := c.Client.Publisher(c.Topic)
	result := publisher.Publish(ctx, &pubsub.Message{
		Data:       data,
		Attributes: attributes,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	msgID, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}

	slog.Debug("Published message with", "msgID", msgID)
	return msgID, nil
}

// PublishMessages publishes multiple messages to the topic
func (c *Client) PublishMessages(ctx context.Context, messages [][]byte, attributes map[string]string) ([]string, error) {
	results := []string{}

	var err error
	var msgID string
	i := 0
	for _, msg := range messages {
		msgID, err = c.PublishMessage(ctx, msg, attributes)
		if msgID != "" {
			results[i] = msgID
			i++
		}
	}

	return results, err
}

// PullMessages pulls messages from the subscription
func (c *Client) PullMessages(ctx context.Context, maxMessages int) ([]*pubsub.Message, error) {
	subscriber := c.Client.Subscriber(c.Subscription)
	subscriber.ReceiveSettings.MaxOutstandingMessages = maxMessages

	messages := make([]*pubsub.Message, 0)

	err := subscriber.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		messages = append(messages, msg)
		msg.Ack() // Acknowledge the message
	})

	if err != nil {
		return nil, fmt.Errorf("failed to pull messages: %v", err)
	}

	return messages, nil
}

// ListSubscriptions lists all subscriptions for this topic
func (c *Client) ListSubscriptions(ctx context.Context) []*pubsubpb.Subscription {
	var subs []*pubsubpb.Subscription
	req := &pubsubpb.ListSubscriptionsRequest{
		Project: fmt.Sprintf("projects/%s", c.ProjectID),
	}
	it := c.Client.SubscriptionAdminClient.ListSubscriptions(ctx, req)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}
		subs = append(subs, s)
	}
	return subs
}
