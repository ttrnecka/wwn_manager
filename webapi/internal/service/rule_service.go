package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
)

type RuleService interface {
	GenericService[entity.Rule]
	CreateReconcileRules(context.Context, *entity.FCWWNEntry, dto.EntryReconcileDTO) error
	ExportRules(context.Context, *os.File) error
	BackupRules(context.Context) error
}

type ruleService struct {
	GenericService[entity.Rule]
}

func NewRuleService(p repository.RuleRepository) RuleService {
	return &ruleService{
		GenericService: NewGenericService(p)}
}

func (s ruleService) CreateReconcileRules(ctx context.Context, entry *entity.FCWWNEntry, reconcileData dto.EntryReconcileDTO) error {

	rules := make([]*entity.Rule, 0)

	if reconcileData.PrimaryHostname != "" {
		// if entry decode host is primary
		if entry.Hostname == reconcileData.PrimaryHostname {
			rules = append(rules, &entity.Rule{
				Customer: entry.Customer,
				Type:     entity.IgnoreLoaded,
				Regex:    fmt.Sprintf("%s,%s", entry.LoadedHostname, entry.WWN),
				Group:    0,
				Order:    0,
				Comment:  fmt.Sprintf("%s hostname reconciliation", reconcileData.PrimaryHostname),
			})
		}

		// if entry loaded host is primary
		if entry.LoadedHostname == reconcileData.PrimaryHostname {
			rules = append(rules, &entity.Rule{
				Customer: entry.Customer,
				Type:     entity.WWNHostMapRule,
				Regex:    entry.WWN,
				Group:    0,
				Order:    0,
				Comment:  reconcileData.PrimaryHostname,
			})
		}
	}

	// if primary customer is set and different than current entry customer
	if reconcileData.PrimaryCustomer != "" {
		rules = append(rules, &entity.Rule{
			Customer: entity.GLOBAL_CUSTOMER,
			Type:     entity.WWNCustomerMapRule,
			Regex:    entry.WWN,
			Group:    0,
			Order:    0,
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

func (s ruleService) ExportRules(ctx context.Context, file *os.File) error {
	items, err := s.All(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rules: %v", err)
	}
	writer := csv.NewWriter(file)
	for _, item := range items {
		itemDTO := mapper.ToRuleDTO(item)
		err := writer.Write([]string{strconv.Itoa(itemDTO.Order), itemDTO.Customer, itemDTO.Regex, strconv.Itoa(itemDTO.Group), string(itemDTO.Type), itemDTO.Comment})
		if err != nil {
			return fmt.Errorf("failed to write csv file: %v", err)
		}
	}
	writer.Flush()
	return nil
}

func (s ruleService) BackupRules(ctx context.Context) error {
	exportDir := "backup"

	if err := ensureDir(exportDir); err != nil {
		return fmt.Errorf("cannot create %s folder: %v", exportDir, err)
	}

	file, err := createTimestampedFile(exportDir, "rules", "csv")
	if err != nil {
		return fmt.Errorf("cannot create rules file: %v", err)
	}
	defer file.Close()

	err = s.ExportRules(ctx, file)
	if err != nil {
		return fmt.Errorf("cannot backup rules: %v", err)
	}
	return nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0750)
}

func createTimestampedFile(dir, baseName, ext string) (*os.File, error) {
	timestamp := time.Now().Format("20060102_150405") // e.g., 20251023_142355
	fileName := fmt.Sprintf("%s_%s.%s", baseName, timestamp, ext)
	fullPath := filepath.Join(dir, fileName)

	fullPath = filepath.Clean(fullPath)

	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
