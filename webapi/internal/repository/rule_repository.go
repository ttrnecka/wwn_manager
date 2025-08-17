package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
)

type RuleRepository interface {
	GenericRepository[entity.Rule]
}

func NewRuleRepository(db *cdb.CRUD[entity.Rule]) RuleRepository {
	return NewGenericRepository(db)
}
