package dto

type FCEntryDTO struct {
	ID             string `json:"id"`
	Customer       string `json:"customer" validate:"required"`
	WWN            string `json:"wwn" validate:"required"`
	Zone           string `json:"zone"`
	Alias          string `json:"alias"`
	Hostname       string `json:"hostname"`
	LoadedHostname string `json:"loaded_hostname"`
	Type           string `json:"type"`
	TypeRule       string `json:"type_rule"`
	HostNameRule   string `json:"hostname_rule"`
	NeedsReconcile bool   `json:"needs_reconcile"`
}
