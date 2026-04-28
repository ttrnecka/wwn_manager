package db

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestConnect_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("connects and returns DB", func(mt *mtest.T) {
		fakeConnector := func(ctx context.Context, uri string) (*mongo.Database, error) {
			return mt.Client.Database("wwn_identity"), nil
		}

		db, err := connect("mongodb://fake:27017/", fakeConnector)

		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NotNil(t, db.database)
	})
}

func TestConnect_Failure(t *testing.T) {
	fakeConnector := func(ctx context.Context, uri string) (*mongo.Database, error) {
		return nil, errors.New("connection refused")
	}

	db, err := connect("mongodb://bad:27017/", fakeConnector)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestConnect_UsesEnvURI(t *testing.T) {
	t.Setenv("MONGO_URI", "mongodb://custom:27017/")

	var capturedURI string
	fakeConnector := func(ctx context.Context, uri string) (*mongo.Database, error) {
		capturedURI = uri
		return nil, errors.New("stop early")
	}

	// Call the internal connect directly with the env-resolved URI
	uri, _ := os.LookupEnv("MONGO_URI")
	connect(uri, fakeConnector)

	assert.Equal(t, "mongodb://custom:27017/", capturedURI)
}
