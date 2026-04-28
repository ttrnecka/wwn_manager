package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	database *mongo.Database
}

// connector is a func type so tests can inject a fake.
type connector func(ctx context.Context, uri string) (*mongo.Database, error)

// mongoConnector is the real implementation used in production.
func mongoConnector(ctx context.Context, uri string) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connection failed: %w", err)
	}
	return client.Database("wwn_identity"), nil
}

// connect is the internal, testable core.
func connect(uri string, conn connector) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := conn(ctx, uri)
	if err != nil {
		return nil, err
	}

	return &DB{database: database}, nil
}

// Connect is the public entrypoint — signature unchanged for callers.
func Connect() (*DB, error) {
	uri, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		uri = "mongodb://localhost:27017/"
	}

	return connect(uri, mongoConnector)
}
