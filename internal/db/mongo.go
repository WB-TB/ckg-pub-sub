package database

import (
	"context"
	"fmt"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	"strings"
	"sync"

	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbConn *mongoDB.Client
	lockDB sync.Mutex
)

func GetConnection(config config.DatabaseConfig) *mongoDB.Client {
	if dbConn == nil {
		lockDB.Lock()
		defer lockDB.Unlock()
		dbConn = connectDB(config)
	}
	return dbConn
}

func NewConnectionDB(config config.DatabaseConfig) (*mongoDB.Client, error) {
	dbConn = connectDB(config)
	return dbConn, nil
}

func connectDB(config config.DatabaseConfig) *mongoDB.Client {
	ctx := context.Background()

	driver := config.Driver
	if driver == "" {
		driver = "mongodb"
	}
	// Build connection string with authentication
	authSource := config.Database
	if authSource == "" {
		authSource = "admin"
	}

	var host string
	if strings.Contains(config.Host, ",") || config.Port == 0 {
		host = config.Host
	} else {
		host = fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	var connectionString string
	if config.Username != "" && config.Password != "" {
		connectionString = fmt.Sprintf("%s://%s:%s@%s/",
			driver,
			config.Username,
			config.Password,
			host,
		)
	} else {
		connectionString = fmt.Sprintf("%s://%s/",
			driver,
			host,
		)
	}

	if config.Attributes != "" {
		connectionString += "?" + config.Attributes
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
	client, err := mongoDB.Connect(ctx, clientOpts)
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err)
		return nil
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error("Failed to ping MongoDB", "error", err)
		return nil
	}

	slog.Info("Successfully connected to MongoDB")
	return client
}

func CloseConnection() {
	if dbConn != nil {
		if err := dbConn.Disconnect(context.Background()); err != nil {
			slog.Error("got an error while disconnecting database server", "error", err)
		} else {
			slog.Info("Successfully disconnected from MongoDB")
		}
	}
}
