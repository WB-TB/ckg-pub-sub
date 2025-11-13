package db

import (
	"context"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/db/mongo"
	"pubsub-ckg-tb/internal/db/sql"
	"sync"
)

var (
	dbConn connection.DatabaseConnection
	lockDB sync.Mutex
)

func GetConnection(config *config.DatabaseConfig) connection.DatabaseConnection {
	if dbConn == nil {
		lockDB.Lock()
		defer lockDB.Unlock()
		dbConn = connectDB(config)
	}
	return dbConn
}

func connectDB(config *config.DatabaseConfig) connection.DatabaseConnection {
	driver := config.Driver
	if driver == "postgres" || driver == "mysql" {
		return sql.NewDBConnection(config)
	} else {
		return mongo.NewDBConnection(config)
	}
}

func CloseConnection(ctx context.Context) {
	if dbConn != nil {
		if err := dbConn.Close(ctx); err != nil {
			slog.Error("got an error while disconnecting database server", "error", err)
		}
	}
}
