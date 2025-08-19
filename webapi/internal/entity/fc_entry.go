package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/db"
)

type FCEntry struct {
	cdb.BaseModel `bson:",inline"`
	Customer      string `bson:"customer"`
	WWN           string `bson:"wwn"`
	Zone          string `bson:"zone"`
	Alias         string `bson:"alias"`
	Hostname      string `bson:"hostname,omitempty"`
	Type          string `bson:"type,omitempty"`
	TypeRule      string `bson:"type_rule,omitempty"`
	HostNameRule  string `bson:"hostname_rule,omitempty"`
}

func FCEntries() *cdb.CRUD[FCEntry] {
	return cdb.NewCRUD[FCEntry](db.Database(), "fc_entries")
}
