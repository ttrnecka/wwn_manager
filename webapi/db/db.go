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
	client   *mongo.Client
	database *mongo.Database
}

func Connect() (*DB, error) {
	uri, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		uri = "mongodb://localhost:27017/"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db := DB{}
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connection failed: %w", err)
	}
	db.database = client.Database("wwn_identity")
	db.client = client
	return &db, nil
}
