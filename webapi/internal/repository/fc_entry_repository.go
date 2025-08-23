package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
)

type FCWWNEntryRepository interface {
	GenericRepository[entity.FCWWNEntry]
}

func NewFCWWNEntryRepository(db *cdb.CRUD[entity.FCWWNEntry]) FCWWNEntryRepository {
	return NewGenericRepository(db)
}
