package repository

import (
	"context"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/models"
)

type PubSub interface {
	GetIncomingIDs(messageIDs []string) ([]string, error)
	SaveNewIncoming(incoming models.IncomingMessageStatusTB) error
	UpdateIncoming(messageID string, processedAt *string) error
	DeleteIncomingMessage(dateExpired string)

	GetOutgoingIDs(messageIDs []string) ([]string, error)
	GetLastOutgoingTimestamp() (string, error)
	SaveOutgoing(outgoing models.OutgoingMessageSkriningTB) error
}

type PubSubRepository struct {
	Configurations *config.Configurations
	Context        context.Context
	Connnection    connection.DatabaseConnection
}

func NewPubSubRepository(ctx context.Context, config *config.Configurations, conn connection.DatabaseConnection) *PubSubRepository {
	return &PubSubRepository{
		Configurations: config,
		Context:        ctx,
		Connnection:    conn,
	}
}

func (r *PubSubRepository) GetIncomingIDs(messageIDs []string) ([]string, error) {
	filter := map[string]any{
		"id": map[string]any{
			"$in": messageIDs,
		},
	}
	ids, err := r.Connnection.Find(r.Context, r.Configurations.CKG.TableIncoming, []string{"id"}, filter, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, entry := range ids.([]any) {
		id := entry.(map[string]any)["id"]
		result = append(result, id.(string))
	}

	return result, nil
}

func (r *PubSubRepository) SaveNewIncoming(incoming models.IncomingMessageStatusTB) error {
	_, err := r.Connnection.InsertOne(r.Context, r.Configurations.CKG.TableIncoming, incoming)
	if err != nil {
		return err
	}

	return nil
}

func (r *PubSubRepository) UpdateIncoming(messageID string, processedAt *string) error {
	filter := map[string]any{
		"id": messageID,
	}
	update := map[string]any{
		"processed_at": processedAt,
	}
	_, err := r.Connnection.UpdateOne(r.Context, r.Configurations.CKG.TableIncoming, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *PubSubRepository) DeleteIncomingMessage(dateExpired string) {
	filter := map[string]any{
		"received_at": map[string]any{
			"$lt": dateExpired,
		},
	}
	r.Connnection.DeleteOne(r.Context, r.Configurations.CKG.TableIncoming, filter)
}

func (r *PubSubRepository) GetOutgoingIDs(messageIDs []string) ([]string, error) {
	filter := map[string]any{
		"id": map[string]any{
			"$in": messageIDs,
		},
	}
	ids, err := r.Connnection.Find(r.Context, r.Configurations.CKG.TableOutgoing, []string{"id"}, filter, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, entry := range ids.([]any) {
		id := entry.(map[string]any)["id"]
		result = append(result, id.(string))
	}
	return result, nil
}

func (r *PubSubRepository) GetLastOutgoingTimestamp() (string, error) {
	sort := map[string]int{
		"created_at": -1,
	}
	var outgoing models.OutgoingMessageSkriningTB
	err := r.Connnection.FindOne(r.Context, &outgoing, r.Configurations.CKG.TableOutgoing, nil, nil, sort)
	if err != nil {
		return "", err
	}
	return outgoing.CreatedAt, nil
}

func (r *PubSubRepository) SaveOutgoing(outgoing models.OutgoingMessageSkriningTB) error {
	var out models.OutgoingMessageSkriningTB
	filter := map[string]any{
		"id": outgoing.ID,
	}
	err := r.Connnection.FindOne(r.Context, &out, r.Configurations.CKG.TableOutgoing, nil, filter, nil)
	if err == nil && out.ID != "" {
		filter := map[string]any{
			"id": out.ID,
		}
		update := map[string]any{
			"updated_at": out.UpdatedAt,
		}
		r.Connnection.UpdateOne(r.Context, r.Configurations.CKG.TableOutgoing, filter, update)
	} else {
		r.Connnection.InsertOne(r.Context, r.Configurations.CKG.TableOutgoing, outgoing)
	}

	return nil
}
