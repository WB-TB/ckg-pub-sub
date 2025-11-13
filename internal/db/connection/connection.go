package connection

import (
	"context"
	"pubsub-ckg-tb/internal/db/dbtypes"
)

type DatabaseConnection interface {
	// Common methods that both databases should implement
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Ping(ctx context.Context) error

	// Get the underlying connection for specific operations
	GetConnection() any
	GetDriver() string
	GetName() string

	Find(ctx context.Context, table string, column []string, filter dbtypes.M, sort map[string]int, limit int64, skip int64) (any, error)
	FindOne(ctx context.Context, result any, table string, column []string, filter dbtypes.M, sort map[string]int) error
	InsertOne(ctx context.Context, table string, data any) (any, error)
	UpdateOne(ctx context.Context, table string, filter dbtypes.M, data any) (int64, error)
	DeleteOne(ctx context.Context, table string, filter dbtypes.M) (any, error)
}
