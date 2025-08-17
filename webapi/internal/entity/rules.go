package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/db"
)

type RuleType string

const (
	ZoneRule  RuleType = "zone"
	AliasRule RuleType = "alias"
)

type Rule struct {
	cdb.BaseModel `bson:",inline"`
	Customer      string   `bson:"customer"`
	Type          RuleType `bson:"type"`
	Regex         string   `bson:"regex"`
	Order         int      `bson:"order"`
}

func Rules() *cdb.CRUD[Rule] {
	return cdb.NewCRUD[Rule](db.Database(), "rules")
}
