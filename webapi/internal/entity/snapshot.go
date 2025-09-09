package entity

import (
	"fmt"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/db"
)

type Snapshot struct {
	cdb.BaseModel `bson:",inline"`
	SnapshotID    string `bson:"snapshot_id"`
}

func (s Snapshot) CollectionName() string {
	return fmt.Sprintf("fc_wwn_entries_%s", s.SnapshotID)
}

func Snapshots() *cdb.CRUD[Snapshot] {
	return cdb.NewCRUD[Snapshot](db.Database(), "snapshots")
}

func GetSnapshot(s Snapshot) *cdb.CRUD[FCWWNEntry] {
	return cdb.NewCRUD[FCWWNEntry](db.Database(), s.CollectionName())
}
