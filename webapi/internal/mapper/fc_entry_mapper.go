package mapper

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToFCEntryEntity(p dto.FCEntryDTO) entity.FCEntry {
	entry := entity.FCEntry{
		Customer: p.Customer,
		WWN:      p.WWN,
		Zone:     p.Zone,
		Alias:    p.Alias,
		Hostname: p.Hostname,
	}
	entry.ID, _ = primitive.ObjectIDFromHex(p.ID)
	return entry
}

func ToFCEntryDTO(p entity.FCEntry) dto.FCEntryDTO {
	return dto.FCEntryDTO{
		ID:       p.ID.Hex(),
		Customer: p.Customer,
		WWN:      p.WWN,
		Zone:     p.Zone,
		Alias:    p.Alias,
		Hostname: p.Hostname,
	}
}
