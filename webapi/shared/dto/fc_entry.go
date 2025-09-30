package dto

import (
	"time"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
)

type FCWWNEntryDTO struct {
	ID                       string                     `json:"id"`
	Customer                 string                     `json:"customer" validate:"required"`
	WWN                      string                     `json:"wwn" validate:"required"`
	Zones                    []string                   `json:"zones"`
	Aliases                  []string                   `json:"aliases"`
	Hostname                 string                     `json:"hostname"`
	LoadedHostname           string                     `json:"loaded_hostname"`
	IsCSVLoad                bool                       `json:"is_csv_load"`
	WWNSet                   int                        `json:"wwn_set"`
	Type                     string                     `json:"type"`
	TypeRule                 string                     `json:"type_rule"`
	HostNameRule             string                     `json:"hostname_rule"`
	HostNameRuleType         entity.RuleType            `json:"hostname_rule_type,omitempty"`
	ReconcileRules           []string                   `json:"reconcile_rules"`
	DefaultReconcileMessages []entity.RuleType          `json:"default_reconcile_messages"`
	NeedsReconcile           bool                       `json:"needs_reconcile"`
	IsPrimaryCustomer        bool                       `json:"is_primary_customer"`
	DuplicateCustomers       []entity.DuplicateCustomer `json:"duplicate_customers,omitempty"`
	IgnoreLoaded             bool                       `json:"ignore_loaded"`
	IgnoreEntry              bool                       `json:"ignore_entry"`
	DeletedAt                *time.Time                 `json:"deleted_at,omitempty"`
}

type EntryReconcileDTO struct {
	PrimaryCustomer string `json:"primary_customer,omitempty"`
	PrimaryHostname string `json:"primary_hostname,omitempty"`
}
