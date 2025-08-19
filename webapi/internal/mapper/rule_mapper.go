package mapper

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToRuleEntity(p dto.RuleDTO) entity.Rule {
	entry := entity.Rule{
		Customer: p.Customer,
		Type:     p.Type,
		Regex:    p.Regex,
		Order:    p.Order,
		Comment:  p.Comment,
	}
	entry.ID, _ = primitive.ObjectIDFromHex(p.ID)
	return entry
}

func ToRuleDTO(p entity.Rule) dto.RuleDTO {
	return dto.RuleDTO{
		ID:       p.ID.Hex(),
		Customer: p.Customer,
		Type:     p.Type,
		Regex:    p.Regex,
		Order:    p.Order,
		Comment:  p.Comment,
	}
}
