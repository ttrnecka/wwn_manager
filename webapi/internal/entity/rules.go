package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type RuleType string

// zone, alias rules are used to decet host within zone or alias name
// wwn_match rule is used to match exact WWN to host
// range_rule is to mark WWN match as certain type - array, backup, HBA
const (
	ZoneRule                     RuleType = "zone"                            // host rule
	AliasRule                    RuleType = "alias"                           // host rule
	WWNHostMapRule               RuleType = "wwn_host_map"                    // host & reconcile rule
	WWNCustomerMapRule           RuleType = "wwn_customer_map"                // reconcile rule
	IgnoreLoaded                 RuleType = "ignore_loaded"                   // reconcile rule
	WWNArrayRangeRule            RuleType = "wwn_range_array"                 // range rule
	WWNBackupRangeRule           RuleType = "wwn_range_backup"                // range rule
	WWNHostRangeRule             RuleType = "wwn_range_host"                  // range rule
	WWNOtherRangeRule            RuleType = "wwn_range_other"                 // range rule
	DefaultReconcileRulePrimary  RuleType = "default_reconcile_rule_primary"  //reconcile rule
	DefaultReconcileRuleOverride RuleType = "default_reconcile_rule_override" //reconcile rule
	DefaultReconcileRuleIgnore   RuleType = "default_reconcile_rule_ignore"   //reconcile rule
)

var RangeRules []RuleType = []RuleType{
	WWNArrayRangeRule, WWNBackupRangeRule, WWNHostRangeRule, WWNOtherRangeRule,
}

var HostRules []RuleType = []RuleType{
	ZoneRule, AliasRule, WWNHostMapRule,
}

var ReconcileRules []RuleType = []RuleType{
	WWNCustomerMapRule, IgnoreLoaded,
}

type Rule struct {
	cdb.BaseModel `bson:",inline"`
	Customer      string   `bson:"customer"`
	Type          RuleType `bson:"type"`
	Regex         string   `bson:"regex"`
	Group         int      `bson:"group"`
	Order         int      `bson:"order"`
	Comment       string   `bson:"comment"`
}

func Rules(db *mongo.Database) *cdb.CRUD[Rule] {
	return cdb.NewCRUD[Rule](db, "rules")
}
