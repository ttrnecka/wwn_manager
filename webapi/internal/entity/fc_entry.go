package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FCWWNEntry struct {
	cdb.BaseModel      `bson:",inline"`
	Customer           string             `bson:"customer"`
	WWN                string             `bson:"wwn"`
	Zones              []string           `bson:"zones"`
	Aliases            []string           `bson:"aliases"`
	Hostname           string             `bson:"hostname,omitempty"`
	LoadedHostname     string             `bson:"loaded_hostname"`
	Type               string             `bson:"type,omitempty"`
	TypeRule           primitive.ObjectID `bson:"type_rule,omitempty"`
	HostNameRule       primitive.ObjectID `bson:"hostname_rule,omitempty"`
	NeedsReconcile     bool               `bson:"needs_reconcile"`
	IsPrimaryCustomer  bool               `bson:"is_primary_customer"`
	DuplicateCustomers []string           `bson:"duplicate_customers"`
}

func FCWWNEntries() *cdb.CRUD[FCWWNEntry] {
	return cdb.NewCRUD[FCWWNEntry](db.Database(), "fc_wwn_entries")
}
