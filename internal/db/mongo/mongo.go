package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/db/dbtypes"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDBConnection implements DatabaseConnection for MongoDB
type MongoDBConnection struct {
	conn        *mongo.Client
	config      *config.DatabaseConfig
	collections map[string]*mongo.Collection
}

func NewDBConnection(config *config.DatabaseConfig) connection.DatabaseConnection {
	return &MongoDBConnection{
		config:      config,
		collections: make(map[string]*mongo.Collection),
	}
}

func (p *MongoDBConnection) GetConnection() any {
	return p.conn
}

func (p *MongoDBConnection) GetDriver() string {
	return p.config.Driver
}

func (p *MongoDBConnection) GetName() string {
	return "MongoDB"
}

func (p *MongoDBConnection) Connect(ctx context.Context) error {
	driver := p.config.Driver

	// Build connection string with authentication
	authSource := p.config.Database
	if authSource == "" {
		authSource = "admin"
	}

	var host string
	if strings.Contains(p.config.Host, ",") || p.config.Port == 0 {
		host = p.config.Host
	} else {
		host = fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)
	}

	var connectionString string
	if p.config.Username != "" && p.config.Password != "" {
		connectionString = fmt.Sprintf("%s://%s:%s@%s/",
			driver,
			p.config.Username,
			p.config.Password,
			host,
		)
	} else {
		connectionString = fmt.Sprintf("%s://%s/",
			driver,
			host,
		)
	}

	if p.config.Attributes != "" {
		connectionString += "?" + p.config.Attributes
	}

	slog.Debug("Attempting to connect to MongoDB with", "URI", connectionString)

	// Create client options
	clientOpts := options.Client().
		SetRetryWrites(true).
		SetRetryReads(true).
		SetMinPoolSize(5).
		SetMaxConnecting(100).
		ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err)
		return nil
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error("Failed to ping MongoDB", "error", err)
		return err
	}

	slog.Info("Successfully connected to MongoDB")
	return nil
}

func (m *MongoDBConnection) Close(ctx context.Context) error {
	if m.conn != nil {
		err := m.conn.Disconnect(ctx)
		if err == nil {
			slog.Info("Successfully disconnected from " + m.GetName())
		}
		return err
	}
	return nil
}

func (m *MongoDBConnection) Ping(ctx context.Context) error {
	if m.conn != nil {
		return m.conn.Ping(ctx, nil)
	}
	return nil
}

func (m *MongoDBConnection) Watch(ctx context.Context, table string) *mongo.ChangeStream {
	collection := m.GetCollection(table)

	// Create change stream options
	changeStreamOptions := options.ChangeStream()
	changeStreamOptions.SetFullDocument(options.UpdateLookup)

	// Create change stream
	changeStream, err := collection.Watch(ctx, []bson.M{}, changeStreamOptions)
	if err != nil {
		slog.Error("Gagal membuat change stream", "error", err, "collection", table)
		return nil
	}
	defer changeStream.Close(ctx)

	return changeStream
}

func (m *MongoDBConnection) RestartWatch(ctx context.Context, table string, changeStream *mongo.ChangeStream) (*mongo.ChangeStream, error) {
	changeStream.Close(ctx)

	collection := m.GetCollection(table)

	// Create change stream options
	changeStreamOptions := options.ChangeStream()
	changeStreamOptions.SetFullDocument(options.UpdateLookup)

	changeStream, err := collection.Watch(ctx, []bson.M{}, changeStreamOptions)
	if err != nil {
		slog.Error("Gagal membuat ulang change stream", "error", err)
	}
	return changeStream, err
}

func (m *MongoDBConnection) Find(ctx context.Context, table string, column []string, filter dbtypes.M, sort map[string]int, limit int64, skip int64) (any, error) {
	collection := m.GetCollection(table)
	if collection == nil {
		return nil, fmt.Errorf("collection %s not found", table)
	}

	// Create projection if columns are specified
	projection := bson.M{}
	if len(column) > 0 {
		// Create projection
		for _, col := range column {
			projection[col] = 1
		}
	}

	// Build find options
	findOptions := options.Find()
	if len(projection) > 0 {
		findOptions.SetProjection(projection)
	}

	if len(sort) > 0 {
		sortBson := bson.M{}
		for field, order := range sort {
			if order == 1 {
				sortBson[field] = 1
			} else {
				sortBson[field] = -1
			}
		}
		findOptions.SetSort(sortBson)
	}

	if limit > 0 {
		findOptions.SetLimit(limit)
	}

	if skip > 0 {
		findOptions.SetSkip(skip)
	}

	mfilter := bson.M{}
	copyToBsonMap(filter, &mfilter)
	cursor, err := collection.Find(ctx, mfilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *MongoDBConnection) FindOne(ctx context.Context, result any, table string, column []string, filter dbtypes.M, sort map[string]int) error {
	collection := m.GetCollection(table)
	if collection == nil {
		return fmt.Errorf("collection %s not found", table)
	}

	// Create projection if columns are specified
	projection := bson.M{}
	if len(column) > 0 {
		// Create projection
		for _, col := range column {
			projection[col] = 1
		}
	}

	// Build find options
	findOptions := options.FindOne()
	if len(projection) > 0 {
		findOptions.SetProjection(projection)
	}

	if len(sort) > 0 {
		sortBson := bson.M{}
		for field, order := range sort {
			if order == 1 {
				sortBson[field] = 1
			} else {
				sortBson[field] = -1
			}
		}
		findOptions.SetSort(sortBson)
	}

	mfilter := bson.M{}
	copyToBsonMap(filter, &mfilter)

	err := collection.FindOne(ctx, mfilter, findOptions).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDBConnection) InsertOne(ctx context.Context, table string, data any) (any, error) {
	collection := m.GetCollection(table)
	if collection == nil {
		return nil, fmt.Errorf("collection %s not found", table)
	}

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *MongoDBConnection) UpdateOne(ctx context.Context, table string, filter dbtypes.M, data any) (int64, error) {
	collection := m.GetCollection(table)
	if collection == nil {
		return 0, fmt.Errorf("collection %s not found", table)
	}

	mfilter := bson.M{}
	copyToBsonMap(filter, &mfilter)

	result, err := collection.UpdateOne(ctx, mfilter, bson.M{"$set": data})
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

func (m *MongoDBConnection) DeleteOne(ctx context.Context, table string, filter dbtypes.M) (any, error) {
	collection := m.GetCollection(table)
	if collection == nil {
		return nil, fmt.Errorf("collection %s not found", table)
	}

	mfilter := bson.M{}
	copyToBsonMap(filter, &mfilter)
	result, err := collection.DeleteOne(ctx, mfilter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoDBConnection) GetCollection(collectionName string) *mongo.Collection {
	if collection, ok := m.collections[collectionName]; ok {
		return collection
	}
	collection := m.conn.Database(config.GetConfig().Database.Database).Collection(collectionName)
	return collection
}

func copyToBsonMap(src dbtypes.M, dst *bson.M) any {
	for k, v := range src {
		vType := reflect.TypeOf(v)
		switch vType.Kind() {
		case reflect.Map:
			m := bson.M{}
			copyToBsonMap(v.(dbtypes.M), &m)
			(*dst)[k] = m
		case reflect.Slice:
			arr := v.([]any)
			for i, item := range arr {
				iType := reflect.TypeOf(item)
				if iType.Kind() == reflect.Map {
					m := bson.M{}
					copyToBsonMap(item.(dbtypes.M), &m)
					arr[i] = m
				} else {
					arr[i] = item
				}
			}
			(*dst)[k] = arr
		default:
			(*dst)[k] = v
		}
	}

	return dst
}
