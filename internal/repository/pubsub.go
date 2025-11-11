package repository

import (
	"context"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type PubSub interface {
	GetIncomingIDs(messageIDs []string) ([]string, error)
	SaveNewIncoming(incoming models.IncomingMessageStatusTB) error
	UpdateIncoming(messageID string, processedAt *string) error
	DeleteIncomingMessage(dateExpired string)

	GetOutgoingIDs(messageIDs []string) ([]string, error)
	GetLastOutgoingTimestamp() (string, error)
	SaveNewOutgoing(outgoing models.OutgoingMessageSkriningTB) error
	UpdateOutgoing(messageID string, processedAt *string) error
}

type PubSubRepository struct {
	Configurations *config.Configurations
	Context        context.Context
	Connnection    *mongo.Client
}

func NewPubSubRepository(ctx context.Context, config *config.Configurations, conn *mongo.Client) *PubSubRepository {
	return &PubSubRepository{
		Configurations: config,
		Context:        ctx,
		Connnection:    conn,
	}
}

func (r *PubSubRepository) GetIncomingIDs(messageIDs []string) ([]string, error) {
	return nil, nil
}

func (r *PubSubRepository) SaveNewIncoming(incoming models.IncomingMessageStatusTB) error {
	return nil
}

func (r *PubSubRepository) UpdateIncoming(messageID string, processedAt *string) error {
	return nil
}

func (r *PubSubRepository) DeleteIncomingMessage(dateExpired string) {

}

func (r *PubSubRepository) GetOutgoingIDs(messageIDs []string) ([]string, error) {
	return nil, nil
}

func (r *PubSubRepository) GetLastOutgoingTimestamp() (string, error) {
	return "", nil
}

func (r *PubSubRepository) SaveNewOutgoing(outgoing models.OutgoingMessageSkriningTB) error {
	return nil
}

func (r *PubSubRepository) UpdateOutgoing(messageID string, processedAt *string) error {
	return nil
}
