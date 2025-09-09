package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
)

type SnapshotRepository interface {
	GenericRepository[entity.Snapshot]
}

func NewSnapshotRepository(db *cdb.CRUD[entity.Snapshot]) SnapshotRepository {
	return NewGenericRepository(db)
}
