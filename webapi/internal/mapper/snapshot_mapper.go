package mapper

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToSnapshotEntity(s dto.SnapshotDTO) entity.Snapshot {
	snapshot := entity.Snapshot{
		SnapshotID: s.SnapshotID,
		Comment:    s.Comment,
	}
	snapshot.ID, _ = primitive.ObjectIDFromHex(s.ID)
	return snapshot
}

func ToSnapshotDTO(s entity.Snapshot) dto.SnapshotDTO {
	return dto.SnapshotDTO{
		ID:         s.ID.Hex(),
		SnapshotID: s.SnapshotID,
		Comment:    s.Comment,
	}
}
