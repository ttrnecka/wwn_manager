package service

import (
	"context"
	"fmt"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
)

type RuleService interface {
	GenericService[entity.Rule]
	CreateReconcileRules(context.Context, *entity.FCWWNEntry, dto.EntryReconcileDTO) error
}

type ruleService struct {
	GenericService[entity.Rule]
}

func NewRuleService(p repository.RuleRepository) RuleService {
	return &ruleService{
		GenericService: NewGenericService(p)}
}

func (s ruleService) CreateReconcileRules(ctx context.Context, entry *entity.FCWWNEntry, reconcileData dto.EntryReconcileDTO) error {

	rules := make([]entity.Rule, 0)

	if reconcileData.PrimaryHostname != "" {
		// if entry decode host is primary
		if entry.Hostname == reconcileData.PrimaryHostname {
			rules = append(rules, entity.Rule{
				Customer: entry.Customer,
				Type:     entity.IgnoreLoaded,
				Regex:    entry.LoadedHostname,
				Group:    0,
				Order:    1,
				Comment:  fmt.Sprintf("%s hostname reconciliation", reconcileData.PrimaryHostname),
			})
		}

		// if entry loaded host is primary
		if entry.LoadedHostname == reconcileData.PrimaryHostname {
			rules = append(rules, entity.Rule{
				Customer: entry.Customer,
				Type:     entity.WWNHostMapRule,
				Regex:    entry.WWN,
				Group:    0,
				Order:    1,
				Comment:  reconcileData.PrimaryHostname,
			})
		}
	}

	// if primary customer is set and different than current entry customer
	if reconcileData.PrimaryCustomer != "" {
		rules = append(rules, entity.Rule{
			Customer: entity.GLOBAL_CUSTOMER,
			Type:     entity.WWNCustomerMapRule,
			Regex:    entry.WWN,
			Group:    0,
			Order:    1,
			Comment:  reconcileData.PrimaryCustomer,
		})
	}

	// there will always be just 1 or 2 so loop is fine
	for _, rule := range rules {
		err := s.DeleteMany(ctx, Filter{"customer": rule.Customer, "regex": rule.Regex, "type": rule.Type})
		if err != nil {
			return err
		}
	}

	if len(rules) == 0 {
		return nil
	}
	return s.InsertAll(ctx, rules)
}
