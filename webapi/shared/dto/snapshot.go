package dto

type SnapshotDTO struct {
	ID         string `json:"id"`
	SnapshotID int64  `json:"snapshot_id"`
	Comment    string `json:"comment"`
}
