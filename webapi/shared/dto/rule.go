package dto

import "github.com/ttrnecka/wwn_identity/webapi/internal/entity"

type RuleDTO struct {
	ID       string          `json:"id"`
	Customer string          `json:"customer" validate:"required"` // customer __GLOBAL__ has special meaning for global rules, mostly used for WWN Range rules
	Type     entity.RuleType `json:"type" validate:"required"`
	Regex    string          `json:"regex" validate:"required"`
	Group    int             `json:"group"`
	Order    int             `json:"order" validate:"required"`
	Comment  string          `json:"comment"`
}
