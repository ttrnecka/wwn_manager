package db

import (
	"fmt"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

var dB *DB

func Database() *mongo.Database {
	return dB.database
}

func Init(logger *zerolog.Logger) error {
	l := logger.With().Str("component", "db").Logger()
	logger = &l
	d, err := Connect()
	if err != nil {
		err = fmt.Errorf("initializing db failed: %w", err)
		logger.Error().Err(err).Msg("")
		return err
	}
	dB = d

	schema := &schema{
		logger: logger,
		db:     dB,
	}
	// Ensure all indexes before starting application logic
	if err := schema.EnsureUserCollection(); err != nil {
		err = fmt.Errorf("setting user collection failed: %w", err)
		logger.Error().Err(err).Msg("")
		return err
	}

	if err := schema.EnsureEntryCollection(); err != nil {
		err = fmt.Errorf("setting fc_wwn_entry collection failed: %w", err)
		logger.Error().Err(err).Msg("")
		return err
	}
	return nil
}
