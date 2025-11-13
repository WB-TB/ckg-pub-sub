package ckg

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/db/mongo"
	"pubsub-ckg-tb/internal/models"
	"pubsub-ckg-tb/internal/pubsub"
	"pubsub-ckg-tb/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type CkgTransmitter struct {
	Configurations *config.Configurations
	Database       connection.DatabaseConnection
	PubSub         *pubsub.Client
	PubSubRepo     repository.PubSub
	CkgRepo        repository.CKGTB
}

func NewCkgTransmitter(ctx context.Context, config *config.Configurations, db connection.DatabaseConnection, pubsub *pubsub.Client) *CkgTransmitter {
	pubsubRepo := repository.NewPubSubRepository(ctx, config, db)
	ckgRepo := repository.NewCKGTBRepository(ctx, config, db)

	return &CkgTransmitter{
		Configurations: config,
		Database:       db,
		PubSub:         pubsub,
		PubSubRepo:     pubsubRepo,
		CkgRepo:        ckgRepo,
	}
}

func (t *CkgTransmitter) Watch(ctx context.Context) {
	// Get database and collection names from config
	// databaseName := t.Configurations.Database.Database
	// collectionName := t.Configurations.CKG.TableSkrining

	// Only support MongoDB for now
	if t.Database.GetDriver() != "mongodb" {
		slog.Error("Only MongoDB is supported for change streams")
		return
	}

	// mongoConn := t.Database.GetMongoConnection()
	// db := mongoConn.Database(databaseName)
	// collection := db.Collection(collectionName)
	mongoDB := t.Database.(*mongo.MongoDBConnection)
	// collection := mongoDB.GetCollection(collectionName)

	// // Create change stream options
	// changeStreamOptions := options.ChangeStream()
	// changeStreamOptions.SetFullDocument(options.UpdateLookup)

	// // Create change stream
	// changeStream, err := collection.Watch(ctx, []bson.M{}, changeStreamOptions)
	changeStream := mongoDB.Watch(ctx, t.Configurations.CKG.TableSkrining)
	if changeStream == nil {
		return
	}
	defer changeStream.Close(ctx)

	slog.Info("Memulai watch untuk perubahan pada collection", "database", t.Configurations.Database.Database, "collection", t.Configurations.CKG.TableSkrining)

	// Listen for changes
	for {
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, stopping watch...")
			return
		default:
			if changeStream.Next(ctx) {
				var changeDoc bson.M
				if err := changeStream.Decode(&changeDoc); err != nil {
					slog.Error("Gagal decode change stream document", "error", err)
					continue
				}

				// Process the change
				if err := t.processChange(ctx, changeDoc); err != nil {
					slog.Error("Gagal memproses perubahan", "error", err)
				}
			} else {
				// Check for errors
				if err := changeStream.Err(); err != nil {
					slog.Error("Error pada change stream", "error", err)

					// Check if error is due to client disconnection
					if err.Error() == "client is disconnected" {
						slog.Info("MongoDB client terputus, mencoba reconnect...")

						// Try to reconnect and recreate change stream
						select {
						case <-ctx.Done():
							return
						default:
							// Close existing change stream
							// changeStream.Close(ctx)

							// Check MongoDB connection health
							if err := t.Database.Ping(ctx); err != nil {
								slog.Error("MongoDB connection is not healthy, waiting before retry...", "error", err)
								time.Sleep(5 * time.Second)
								continue
							}

							// Create new change stream
							// changeStream, err = collection.Watch(ctx, []bson.M{}, changeStreamOptions)
							changeStream, err = mongoDB.RestartWatch(ctx, t.Configurations.CKG.TableSkrining, changeStream)
							if err != nil {
								// Wait before retrying
								time.Sleep(5 * time.Second)
								continue
							}

							slog.Info("Berhasil membuat ulang change stream")
							continue
						}
					} else {
						// For other errors, try to recreate the change stream
						select {
						case <-ctx.Done():
							return
						default:
							slog.Info("Mencoba membuat ulang change stream...")
							// changeStream.Close(ctx)
							// changeStream, err = collection.Watch(ctx, []bson.M{}, changeStreamOptions)
							changeStream, err = mongoDB.RestartWatch(ctx, t.Configurations.CKG.TableSkrining, changeStream)
							if err != nil {
								// Wait before retrying
								time.Sleep(5 * time.Second)
								continue
							}
							continue
						}
					}
				}
			}
		}
	}
}

func (t *CkgTransmitter) processChange(ctx context.Context, changeDoc bson.M) error {
	// Extract operation type
	operation, ok := changeDoc["operationType"].(string)
	if !ok {
		return fmt.Errorf("operationType tidak ditemukan dalam change document")
	}

	// Extract document ID
	documentKey, ok := changeDoc["documentKey"].(bson.M)
	if !ok {
		str, _ := bson.MarshalExtJSON(changeDoc, false, false)
		log.Printf("changeDoc: %s", str)
		return fmt.Errorf("documentKey tidak ditemukan dalam change document")
	}

	id, ok := documentKey["_id"]
	if !ok {
		return fmt.Errorf("_id tidak ditemukan dalam documentKey")
	}

	slog.Info("Mendeteksi perubahan data", "operation", operation, "id", id)

	// Get the full document based on operation type
	var fullDocument bson.M
	switch operation {
	case "insert", "replace", "update":
		if doc, ok := changeDoc["fullDocument"].(bson.M); ok {
			fullDocument = doc
		} else {
			return fmt.Errorf("fullDocument tidak ditemukan untuk operation: %s", operation)
		}
	default:
		slog.Debug("Operation type tidak didukung", "operation", operation)
		return nil
	}

	// Convert bson.M to bytes for unmarshaling
	docBytes, err := bson.Marshal(fullDocument)
	if err != nil {
		return fmt.Errorf("gagal marshal document: %v", err)
	}

	skriningResult, err := t.CkgRepo.GetOnePendingTbSkrining(t.Configurations.CKG.TableSkrining, docBytes)
	if err != nil || skriningResult == nil {
		return err
	}

	// Bukan terduga TB abaikan saja
	if skriningResult.TerdugaTb == nil || *skriningResult.TerdugaTb != "Ya" {
		return nil
	}

	// Prepare data for PubSub
	output := []*models.SkriningCKGResult{skriningResult} // array dengan 1 entry
	// pubsubObjectWrapper := models.PubSubObjectWrapper[*models.SkriningCKGResult]{
	// 	Data: output,
	// }
	pubsubObjectWrapper := models.NewPubSubProducerWrapper(output)

	jsonStr, err := pubsubObjectWrapper.ToJSON()
	if err != nil {
		return err
	}

	attributes := t.Configurations.Producer.MessageAttributes
	if attributes == nil {
		attributes = make(map[string]string)
	}
	attributes["environment"] = t.Configurations.App.Environment
	attributes["timestamp"] = time.Now().Format(time.RFC3339)
	attributes["operation_type"] = operation

	// Send data via PubSub
	slog.Debug("Publish Message", "message", jsonStr, "attributes", attributes)
	t.PubSub.PublishMessage(ctx, []byte(jsonStr), attributes)

	return nil
}

func (t *CkgTransmitter) Produce(ctx context.Context) error {
	output, err := t.Prepare(ctx)
	if err != nil {
		slog.Warn("Gagal menjalankan producer", "error", err)
		return err
	}

	if len(output) > 0 {
		batchSize := t.Configurations.Producer.BatchSize
		totalItems := len(output)
		count := 0

		// Pecah data ke dalam batch dan kirim satu-satu
		for i := 0; i < totalItems; i += batchSize {
			end := i + batchSize
			if end > totalItems {
				end = totalItems
			}

			batch := output[i:end]

			// pubsubObjectWrapper := models.PubSubObjectWrapper[*models.SkriningCKGResult]{
			// 	Data: batch,
			// }
			pubsubObjectWrapper := models.NewPubSubProducerWrapper(batch)

			jsonStr, err := pubsubObjectWrapper.ToJSON()
			if err != nil {
				return err
			}

			attributes := t.Configurations.Producer.MessageAttributes
			if attributes == nil {
				attributes = map[string]string{}
			}
			attributes["environment"] = t.Configurations.App.Environment
			attributes["timestamp"] = time.Now().Format(time.RFC3339)

			slog.Debug("JSON: " + jsonStr)
			// Kirim data via PubSub
			msgID, err := t.PubSub.PublishMessage(ctx, []byte(jsonStr), attributes)
			if err == nil {
				// Simpan log outgoing message
				outgoing := models.OutgoingMessageSkriningTB{
					ID:        msgID,
					CreatedAt: time.Now().Format(time.RFC1123),
				}
				t.PubSubRepo.SaveOutgoing(outgoing)
			}

			// tiap 10 langkah istirahat bentar, serius santai
			count++
			if count%10 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

func (t *CkgTransmitter) Prepare(ctx context.Context) ([]*models.SkriningCKGResult, error) {
	start := ""
	end := ""
	limit := int64(0)
	output := make([]*models.SkriningCKGResult, 0)

	// Get last timestamp from outgoing table
	if start == "" || start == "last" {
		start, _ = t.PubSubRepo.GetLastOutgoingTimestamp()
	}

	// If no last timestamp, use default start time
	now := time.Now()
	if start == "" {
		defaultTime := now.Add(-48 * time.Hour) // Default to 6 hours ago
		start = defaultTime.Format(time.RFC3339)
	}

	if end == "" {
		end = now.Format(time.RFC3339)
	}

	if limit == 0 {
		limit = int64(t.Configurations.Producer.BatchSize)
	}

	// Get status pasien data from database
	pending, err := t.CkgRepo.GetPendingTbSkrining(start, end, limit)
	if err == nil {
		for _, skrining := range pending {
			output = append(output, &skrining)
		}
	}

	return output, nil
}
