package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FCEntry struct {
	cdb.BaseModel  `bson:",inline"`
	Customer       string             `bson:"customer"`
	WWN            string             `bson:"wwn"`
	Zone           string             `bson:"zone"`
	Alias          string             `bson:"alias"`
	Hostname       string             `bson:"hostname,omitempty"`
	LoadedHostname string             `bson:"loaded_hostname"`
	Type           string             `bson:"type,omitempty"`
	TypeRule       primitive.ObjectID `bson:"type_rule,omitempty"`
	HostNameRule   primitive.ObjectID `bson:"hostname_rule,omitempty"`
	NeedsReconcile bool               `bson:"needs_reconcile"`
}

func FCEntries() *cdb.CRUD[FCEntry] {
	return cdb.NewCRUD[FCEntry](db.Database(), "fc_entries")
}

type FCWWNEntry struct {
	cdb.BaseModel  `bson:",inline"`
	Customer       string             `bson:"customer"`
	WWN            string             `bson:"wwn"`
	Zones          []string           `bson:"zones"`
	Aliases        []string           `bson:"aliases"`
	Hostname       string             `bson:"hostname,omitempty"`
	LoadedHostname string             `bson:"loaded_hostname"`
	Type           string             `bson:"type,omitempty"`
	TypeRule       primitive.ObjectID `bson:"type_rule,omitempty"`
	HostNameRule   primitive.ObjectID `bson:"hostname_rule,omitempty"`
	NeedsReconcile bool               `bson:"needs_reconcile"`
}

func FCWWNEntries() *cdb.CRUD[FCWWNEntry] {
	return cdb.NewCRUD[FCWWNEntry](db.Database(), "fc_wwn_entries")
}
