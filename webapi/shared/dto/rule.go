package dto

import "github.com/ttrnecka/wwn_identity/webapi/internal/entity"

type RuleDTO struct {
	ID       string          `json:"id"`
	Customer string          `json:"customer" validate:"required"`
	Type     entity.RuleType `json:"type" validate:"required"`
	Regex    string          `json:"regex" validate:"required"`
	Order    int             `json:"order" validate:"required"`
}
