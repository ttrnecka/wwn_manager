package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	WWNSetSAN                  int = 1
	WWNSetManual               int = 2
	WWNSetAuto                 int = 3
	FCWWNEntriesCollectionName     = "fc_wwn_entries"
)

type DuplicateCustomer struct {
	Customer string `bson:"customer" json:"customer"`
	WWNSet   int    `bson:"wwn_set" json:"wwn_set"`
	Hostname string `bson:"hostname" json:"hostname"`
}

type FCWWNEntry struct {
	cdb.BaseModel      `bson:",inline"`
	Customer           string               `bson:"customer"`
	WWN                string               `bson:"wwn"`
	Zones              []string             `bson:"zones"`
	Aliases            []string             `bson:"aliases"`
	Hostname           string               `bson:"hostname,omitempty"`
	LoadedHostname     string               `bson:"loaded_hostname"`
	IsCSVLoad          bool                 `bson:"is_csv_load"`
	WWNSet             int                  `bson:"wwn_set"`
	Type               string               `bson:"type,omitempty"`
	TypeRule           primitive.ObjectID   `bson:"type_rule,omitempty"`
	HostNameRule       primitive.ObjectID   `bson:"hostname_rule,omitempty"`
	ReconcileRules     []primitive.ObjectID `bson:"reconcile_rules,omitempty"`
	NeedsReconcile     bool                 `bson:"needs_reconcile"`
	IsPrimaryCustomer  bool                 `bson:"is_primary_customer"`
	DuplicateCustomers []DuplicateCustomer  `bson:"duplicate_customers"`
	IgnoreLoaded       bool                 `bson:"ignore_loaded"`
	IgnoreEntry        bool                 `bson:"ignore_entry"`
	// DuplicateCustomers2 []map[string]string `bson:"duplicate_customers2"`
}

func FCWWNEntries(db *mongo.Database) *cdb.CRUD[FCWWNEntry] {
	return cdb.NewCRUD[FCWWNEntry](db, FCWWNEntriesCollectionName)
}
