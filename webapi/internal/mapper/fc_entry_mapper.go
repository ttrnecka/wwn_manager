package mapper

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToFCWWNEntryEntity(p dto.FCWWNEntryDTO) entity.FCWWNEntry {
	entry := entity.FCWWNEntry{
		Customer:       p.Customer,
		WWN:            p.WWN,
		Zones:          p.Zones,
		Aliases:        p.Aliases,
		Hostname:       p.Hostname,
		LoadedHostname: p.LoadedHostname,
		Type:           p.Type,
		NeedsReconcile: p.NeedsReconcile,
	}
	entry.ID, _ = primitive.ObjectIDFromHex(p.ID)
	entry.TypeRule, _ = primitive.ObjectIDFromHex(p.TypeRule)
	entry.HostNameRule, _ = primitive.ObjectIDFromHex(p.HostNameRule)
	return entry
}

func ToFCWWNEntryDTO(p entity.FCWWNEntry) dto.FCWWNEntryDTO {
	return dto.FCWWNEntryDTO{
		ID:             p.ID.Hex(),
		Customer:       p.Customer,
		WWN:            p.WWN,
		Zones:          p.Zones,
		Aliases:        p.Aliases,
		Hostname:       p.Hostname,
		LoadedHostname: p.LoadedHostname,
		Type:           p.Type,
		TypeRule:       p.TypeRule.Hex(),
		HostNameRule:   p.HostNameRule.Hex(),
		NeedsReconcile: p.NeedsReconcile,
	}
}
