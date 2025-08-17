package service

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
)

type RuleService interface {
	GenericService[entity.Rule]
}

type ruleService struct {
	GenericService[entity.Rule]
}

func NewRuleService(p repository.RuleRepository) RuleService {
	return &ruleService{
		GenericService: NewGenericService(p)}
}
