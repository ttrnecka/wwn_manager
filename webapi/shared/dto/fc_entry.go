package dto

type FCWWNEntryDTO struct {
	ID                 string   `json:"id"`
	Customer           string   `json:"customer" validate:"required"`
	WWN                string   `json:"wwn" validate:"required"`
	Zones              []string `json:"zones"`
	Aliases            []string `json:"aliases"`
	Hostname           string   `json:"hostname"`
	LoadedHostname     string   `json:"loaded_hostname"`
	Type               string   `json:"type"`
	TypeRule           string   `json:"type_rule"`
	HostNameRule       string   `json:"hostname_rule"`
	DuplicateRule      string   `json:"duplicate_rule"`
	NeedsReconcile     bool     `json:"needs_reconcile"`
	IsPrimaryCustomer  bool     `json:"is_primary_customer"`
	DuplicateCustomers []string `json:"duplicate_customers,omitempty"`
}

type EntryReconcileDTO struct {
	PrimaryCustomer string `json:"primary_customer,omitempty"`
	PrimaryHostname string `json:"primary_hostname,omitempty"`
}
