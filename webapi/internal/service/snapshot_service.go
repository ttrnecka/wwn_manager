package service

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
)

type SnapshotService interface {
	GenericService[entity.Snapshot]
}

type snapshotService struct {
	GenericService[entity.Snapshot]
}

func NewSnapshotService(p repository.SnapshotRepository) SnapshotService {
	return &snapshotService{
		GenericService: NewGenericService(p)}
}

func (s snapshotService) MakeSnapshot() (*entity.Snapshot, error) {
	//TODO
	// make function that coppies whole fc_wwn_entries to new collection
	// sets up indices on it, then creates snapshot record referincing it and returns the snapshot reference
	return nil, nil
}
