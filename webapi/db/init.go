package db

import (
	"fmt"

	"github.com/rs/zerolog"
	logging "github.com/ttrnecka/agent_poc/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger zerolog.Logger
var dB *DB

func Database() *mongo.Database {
	return dB.database
}

func init() {
	logger = logging.SetupLogger("db")
}

func Init() error {
	d, err := Connect()
	if err != nil {
		err = fmt.Errorf("initializing db failed: %w", err)
		logger.Error().Err(err).Msg("")
		return err
	}
	dB = d

	// Ensure all indexes before starting application logic
	if err := EnsureUserCollection(); err != nil {
		err = fmt.Errorf("setting user collection failed: %w", err)
		logger.Error().Err(err).Msg("")
		return err
	}

	if err := EnsureEntryCollection(); err != nil {
		err = fmt.Errorf("setting fc_wwn_entry collection failed: %w", err)
		logger.Error().Err(err).Msg("")
		return err
	}
	return nil
}
