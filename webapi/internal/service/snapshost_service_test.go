package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// startMongoContainer spins up MongoDB in Docker for testing
func startMongoContainer(t *testing.T) (testcontainers.Container, string) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0", // or "mongo:latest"
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := mongoC.Host(ctx)
	require.NoError(t, err)

	port, err := mongoC.MappedPort(ctx, "27017")
	require.NoError(t, err)

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	return mongoC, uri
}

func TestMakeSnapshot(t *testing.T) {
	ctx := context.Background()

	// --- Start Mongo container
	mongoC, uri := startMongoContainer(t)
	defer func() {
		_ = mongoC.Terminate(ctx)
	}()

	// --- Connect Mongo client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("testdb_snapshot")
	snapshots := entity.Snapshots(db)
	entries := entity.FCWWNEntries(db)

	// --- Prepare repositories
	snapRepo := repository.NewSnapshotRepository(snapshots)
	entryRepo := repository.NewFCWWNEntryRepository(entries)

	// --- Prepare service
	svc := NewSnapshotService(snapRepo, entryRepo)

	// --- Insert some sample data
	entriesColl := entryRepo.GetCollection()
	_, err = entriesColl.InsertMany(ctx, []interface{}{
		bson.M{"_id": "1", "name": "entry1"},
		bson.M{"_id": "2", "name": "entry2"},
	})
	require.NoError(t, err)

	// --- Call MakeSnapshot
	snap, err := svc.MakeSnapshot(ctx)
	require.NoError(t, err)
	require.NotNil(t, snap)

	// --- Verify snapshot collection exists
	targetColl := db.Collection(snap.EntryCollectionName())
	count, err := targetColl.CountDocuments(ctx, bson.D{})
	require.NoError(t, err)
	require.Equal(t, int64(2), count)

	// --- Verify content matches
	cur, err := targetColl.Find(ctx, bson.D{})
	require.NoError(t, err)
	defer cur.Close(ctx)

	var docs []bson.M
	err = cur.All(ctx, &docs)
	require.NoError(t, err)
	require.Len(t, docs, 2)
}
