package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/db"
)

type RuleType string

// zone, alias rules are used to decet host within zone or alias name
// wwn_match rule is used to match exact WWN to host
// range_rule is to mark WWN match as certain type - array, backup, HBA
const (
	ZoneRule           RuleType = "zone"
	AliasRule          RuleType = "alias"
	WWNMapRule         RuleType = "wwn_map"
	WWNArrayRangeRule  RuleType = "wwn_range_array"
	WWNBackupRangeRule RuleType = "wwn_range_backup"
	WWNHostRangeRule   RuleType = "wwn_range_host"
	WWNOtherRangeRule  RuleType = "wwn_range_other"
)

type Rule struct {
	cdb.BaseModel `bson:",inline"`
	Customer      string   `bson:"customer"`
	Type          RuleType `bson:"type"`
	Regex         string   `bson:"regex"`
	Group         int      `bson:"group"`
	Order         int      `bson:"order"`
	Comment       string   `bson:"comment"`
}

func Rules() *cdb.CRUD[Rule] {
	return cdb.NewCRUD[Rule](db.Database(), "rules")
}
