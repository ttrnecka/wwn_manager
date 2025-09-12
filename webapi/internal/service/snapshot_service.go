package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SnapshotService interface {
	GenericService[entity.Snapshot]
	MakeSnapshot(context.Context, string) (*entity.Snapshot, error)
	GetSnapshotEntries(context.Context, entity.Snapshot) ([]entity.FCWWNEntry, error)
	GetEntryService(entity.Snapshot) FCWWNEntryService
}

type snapshotService struct {
	GenericService[entity.Snapshot]
	EntryService GenericService[entity.FCWWNEntry]
}

func NewSnapshotService(p repository.SnapshotRepository, e repository.FCWWNEntryRepository) SnapshotService {
	return &snapshotService{
		GenericService: NewGenericService(p),
		EntryService:   NewGenericService(e),
	}
}

func (s snapshotService) MakeSnapshot(ctx context.Context, comment string) (*entity.Snapshot, error) {
	snapshot := entity.Snapshot{
		SnapshotID: time.Now().Unix(),
		Comment:    comment,
	}
	snapshotID, err := s.Update(ctx, entity.NilObjectID(), &snapshot)
	if err != nil {
		return nil, err
	}
	snapshot.ID = snapshotID

	sourceColl := s.EntryService.Collection()
	targetColl := s.EntryService.Collection().Database().Collection(snapshot.EntryCollectionName())

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{}}},        // match everything
		{{Key: "$out", Value: targetColl.Name()}}, // write into snapshot collection
	}

	cursor, err := s.EntryService.Collection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("snapshot failed: %w", err)
	}
	defer cursor.Close(ctx)

	indexes, err := sourceColl.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer indexes.Close(ctx)

	var idxDocs []bson.M
	if err := indexes.All(ctx, &idxDocs); err != nil {
		return nil, fmt.Errorf("failed to decode indexes: %w", err)
	}

	for _, idx := range idxDocs {
		// Skip the default _id index (MongoDB creates it automatically)
		if name, ok := idx["name"].(string); ok && name == "_id_" {
			continue
		}

		keyDoc, ok := idx["key"].(bson.M)
		if !ok {
			continue
		}

		keys := bson.D{}
		for field, order := range keyDoc {
			keys = append(keys, bson.E{Key: field, Value: order})
		}

		opts := options.Index()

		if unique, ok := idx["unique"].(bool); ok {
			opts.SetUnique(unique)
		}
		if sparse, ok := idx["sparse"].(bool); ok {
			opts.SetSparse(sparse)
		}
		if name, ok := idx["name"].(string); ok {
			opts.SetName(name)
		}
		model := mongo.IndexModel{
			Keys:    keys,
			Options: opts,
		}

		_, err := targetColl.Indexes().CreateOne(ctx, model)
		if err != nil {
			return nil, fmt.Errorf("failed to create index %v: %w", idx, err)
		}
	}

	return &snapshot, nil
}

func (s snapshotService) GetSnapshotEntries(ctx context.Context, snapshot entity.Snapshot) ([]entity.FCWWNEntry, error) {
	db := s.EntryService.Collection().Database()
	entryCrud := snapshot.GetEntries(db)
	entryRepo := repository.NewFCWWNEntryRepository(entryCrud)
	entrySvc := NewFCWWNEntryService(entryRepo)

	entries, err := entrySvc.FindWithSoftDeleted(ctx, Filter{}, SortOption{"wwn": "asc"})
	return entries, err
}

func (s snapshotService) GetEntryService(snapshot entity.Snapshot) FCWWNEntryService {
	db := s.EntryService.Collection().Database()
	entryCrud := snapshot.GetEntries(db)
	entryRepo := repository.NewFCWWNEntryRepository(entryCrud)
	entrySvc := NewFCWWNEntryService(entryRepo)
	return entrySvc
}
