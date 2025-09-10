package entity

import (
	"fmt"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type Snapshot struct {
	cdb.BaseModel `bson:",inline"`
	SnapshotID    int64 `bson:"snapshot_id"`
}

func (s Snapshot) EntryCollectionName() string {
	return fmt.Sprintf("%s_%d", FCWWNEntriesCollectionName, s.SnapshotID)
}

func Snapshots(db *mongo.Database) *cdb.CRUD[Snapshot] {
	return cdb.NewCRUD[Snapshot](db, "snapshots")
}

func (s Snapshot) GetEntries(db *mongo.Database) *cdb.CRUD[FCWWNEntry] {
	return cdb.NewCRUD[FCWWNEntry](db, s.EntryCollectionName())
}
