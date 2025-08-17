package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
)

type FCEntryRepository interface {
	GenericRepository[entity.FCEntry]
}

func NewFCEntryRepository(db *cdb.CRUD[entity.FCEntry]) FCEntryRepository {
	return NewGenericRepository(db)
}
