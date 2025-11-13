package ckg

import (
	"context"
	"fmt"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/models"
	"pubsub-ckg-tb/internal/repository"
	"slices"

	"cloud.google.com/go/pubsub/v2"
)

type CkgReceiver struct {
	Configurations *config.Configurations
	Database       connection.DatabaseConnection
	PubSubRepo     repository.PubSub
	CkgRepo        repository.CKGTB
}

func NewCkgReceiver(ctx context.Context, config *config.Configurations, db connection.DatabaseConnection) *CkgReceiver {
	pubsubRepo := repository.NewPubSubRepository(ctx, config, db)
	ckgRepo := repository.NewCKGTBRepository(ctx, config, db)

	return &CkgReceiver{
		Configurations: config,
		Database:       db,
		PubSubRepo:     pubsubRepo,
		CkgRepo:        ckgRepo,
	}
}

func (r *CkgReceiver) Prepare(ctx context.Context, messages []*pubsub.Message) map[string][]any {
	validMessages := make(map[string][]any)

	// Extract message IDs
	messageIDs := make([]string, 0, len(messages))
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.ID)
	}

	// Periksa semua message ID lalu hanya ambil yang belum pernah diproses saja
	existingIDs, err := r.PubSubRepo.GetIncomingIDs(messageIDs)
	if err != nil {
		slog.Debug("Gagal mengambil daftar message ID existing", "error", err)
		existingIDs = []string{}
	}

	// Process semua message satu-satu
	for _, msg := range messages {
		// Acknowledge terlebih dahulu
		msg.Ack()

		// Skip jika message ID sudah diproses sebelumnya
		if slices.Contains(existingIDs, msg.ID) {
			slog.Debug("Skip message", "id", msg.ID)
			continue
		}

		// Parse message data
		dataStr := string(msg.Data)
		pubsubObjectWrapper := models.NewPubSubConsumerWrapper[*models.StatusPasien]()
		err := pubsubObjectWrapper.FromJSON(dataStr)
		if err != nil {
			slog.Debug("Gagal parsing", "id", msg.ID, "error", err)
			continue
		}

		// Hanya pedulikan Object CKG yang valid
		if !pubsubObjectWrapper.IsCKGObject() {
			slog.Debug("Abaikan message non-CKG", "id", msg.ID)
			continue
		}

		// Simpan incomming message agar tidak diproses berulang kali
		incoming := models.IncomingMessageStatusTB{
			ID:          msg.ID,
			Data:        &dataStr,
			ReceivedAt:  msg.PublishTime.String(),
			ProcessedAt: nil,
		}
		if err := r.PubSubRepo.SaveNewIncoming(incoming); err != nil {
			slog.Info("Gagal menyimpan incoming message", "id", msg.ID, "error", err)
		}

		// register ke validMessages
		validMessages[msg.ID] = []any{incoming, msg, pubsubObjectWrapper.Data}
	}

	return validMessages
}

func (r *CkgReceiver) Consume(ctx context.Context, messages []*pubsub.Message) (map[string]bool, error) {
	results := make(map[string]bool)

	// Filter message hanya yang belum diproses saja
	validMessages := r.Prepare(ctx, messages)

	// Process each valid message
	for msgID, data := range validMessages {
		// incoming := data[0].(*models.IncomingMessageStatusTB)
		msg := data[1].(*pubsub.Message)
		statusPasien := data[2].([]models.StatusPasien)

		// Process the message
		err := r.Process(ctx, statusPasien, msg)
		if err != nil {
			slog.Info("Saat memproses message", "id", msgID, "error", err)
			results[msgID] = false
			continue
		}

		results[msgID] = true
	}

	return results, nil
}

func (r *CkgReceiver) Process(ctx context.Context, statusPasien []models.StatusPasien, msg *pubsub.Message) error {
	slog.Debug(fmt.Sprintf("Received valid CKG SkriningCKG object [%s].\n Data: %s\n Attributes: %v", msg.ID, string(msg.Data), msg.Attributes))

	// Save to database
	_, err := r.CkgRepo.UpdateTbPatientStatus(statusPasien)
	r.PubSubRepo.UpdateIncoming(msg.ID, nil)

	return err
}
