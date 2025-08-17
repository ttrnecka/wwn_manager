package dto

type FCEntryDTO struct {
	ID       string `json:"id"`
	Customer string `json:"customer" validate:"required"`
	WWN      string `json:"wwn" validate:"required"`
	Zone     string `json:"zone"`
	Alias    string `json:"alias"`
	Hostname string `json:"hostname"`
}
