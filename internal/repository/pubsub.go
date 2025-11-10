package repository

import (
	"context"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type PubSub interface {
	GetIncomingIDs(ctx context.Context, messageIDs []string) ([]string, error)
	SaveNewIncoming(ctx context.Context, incoming models.IncomingMessageStatusTB) error
	UpdateIncoming(ctx context.Context, messageID string, processedAt *string) error
	DeleteIncomingMessage(ctx context.Context, dateExpired string)

	GetOutgoingIDs(ctx context.Context, messageIDs []string) ([]string, error)
	GetLastOutgoingTimestamp(ctx context.Context) (string, error)
	SaveNewOutgoing(ctx context.Context, outgoing models.OutgoingMessageSkriningTB) error
	UpdateOutgoing(ctx context.Context, messageID string, processedAt *string) error
}

type PubSubRepository struct {
	Configurations *config.Configurations
	Connnection    *mongo.Client
}

func NewPubSubRepository(config *config.Configurations, conn *mongo.Client) *PubSubRepository {
	return &PubSubRepository{
		Configurations: config,
		Connnection:    conn,
	}
}

func (r *PubSubRepository) GetIncomingIDs(ctx context.Context, messageIDs []string) ([]string, error) {
	return nil, nil
}

func (r *PubSubRepository) SaveNewIncoming(ctx context.Context, incoming models.IncomingMessageStatusTB) error {
	return nil
}

func (r *PubSubRepository) UpdateIncoming(ctx context.Context, messageID string, processedAt *string) error {
	return nil
}

func (r *PubSubRepository) DeleteIncomingMessage(ctx context.Context, dateExpired string) {

}

func (r *PubSubRepository) GetOutgoingIDs(ctx context.Context, messageIDs []string) ([]string, error) {
	return nil, nil
}

func (r *PubSubRepository) GetLastOutgoingTimestamp(ctx context.Context) (string, error) {
	return "", nil
}

func (r *PubSubRepository) SaveNewOutgoing(ctx context.Context, outgoing models.OutgoingMessageSkriningTB) error {
	return nil
}

func (r *PubSubRepository) UpdateOutgoing(ctx context.Context, messageID string, processedAt *string) error {
	return nil
}
