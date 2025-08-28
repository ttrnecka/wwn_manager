package mapper

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToFCWWNEntryEntity(p dto.FCWWNEntryDTO) entity.FCWWNEntry {
	entry := entity.FCWWNEntry{
		Customer:           p.Customer,
		WWN:                p.WWN,
		Zones:              p.Zones,
		Aliases:            p.Aliases,
		Hostname:           p.Hostname,
		LoadedHostname:     p.LoadedHostname,
		IsCSVLoad:          p.IsCSVLoad,
		WWNSet:             p.WWNSet,
		Type:               p.Type,
		NeedsReconcile:     p.NeedsReconcile,
		IsPrimaryCustomer:  p.IsPrimaryCustomer,
		DuplicateCustomers: p.DuplicateCustomers,
		IgnoreLoaded:       p.IgnoreLoaded,
		IgnoreEntry:        p.IgnoreEntry,
	}
	entry.ID, _ = primitive.ObjectIDFromHex(p.ID)
	entry.TypeRule, _ = primitive.ObjectIDFromHex(p.TypeRule)
	entry.HostNameRule, _ = primitive.ObjectIDFromHex(p.HostNameRule)
	for _, r := range p.ReconcileRules {
		id, _ := primitive.ObjectIDFromHex(r)
		entry.ReconcileRules = append(entry.ReconcileRules, id)
	}
	return entry
}

func ToFCWWNEntryDTO(p entity.FCWWNEntry) dto.FCWWNEntryDTO {
	entry := dto.FCWWNEntryDTO{
		ID:                 p.ID.Hex(),
		Customer:           p.Customer,
		WWN:                p.WWN,
		Zones:              p.Zones,
		Aliases:            p.Aliases,
		Hostname:           p.Hostname,
		LoadedHostname:     p.LoadedHostname,
		IsCSVLoad:          p.IsCSVLoad,
		WWNSet:             p.WWNSet,
		Type:               p.Type,
		TypeRule:           p.TypeRule.Hex(),
		HostNameRule:       p.HostNameRule.Hex(),
		NeedsReconcile:     p.NeedsReconcile,
		IsPrimaryCustomer:  p.IsPrimaryCustomer,
		DuplicateCustomers: p.DuplicateCustomers,
		IgnoreLoaded:       p.IgnoreLoaded,
		IgnoreEntry:        p.IgnoreEntry,
	}

	for _, r := range p.ReconcileRules {
		entry.ReconcileRules = append(entry.ReconcileRules, r.Hex())
	}
	return entry
}
