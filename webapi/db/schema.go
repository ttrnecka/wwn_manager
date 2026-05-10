package db

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type schema struct {
	logger *zerolog.Logger
	db     *DB
}

func (s *schema) EnsureEntryCollection() error {
	ctx := context.Background()
	col := s.db.database.Collection("fc_wwn_entries")

	// 1. Create compound index for customer + wwn
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "customer", Value: 1},
			{Key: "wwn", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "wwn", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "createdAt", Value: -1},
		},
	})

	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	s.logger.Info().Msgf("Indexes ensured for collection: fc_wwn_entries")
	return nil
}

func (s *schema) EnsureUserCollection() error {
	ctx := context.Background()
	col := s.db.database.Collection("users")

	// 1. Create compound index for email + status
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
			{Key: "status", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "createdAt", Value: -1},
		},
	})

	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	s.logger.Info().Msgf("Indexes ensured for collection: users")
	// 2. Insert default user if not exists
	defaultUser := bson.M{
		"username":  "admin",
		"email":     "admin@poc.com",
		"status":    "active",
		"password":  string(hashedPw),
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}

	filter := bson.M{"username": defaultUser["username"]}
	update := bson.M{"$setOnInsert": defaultUser}

	_, err = col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("insert default user: %w", err)
	}
	s.logger.Info().Msgf("Default user inserted")
	return nil
}
