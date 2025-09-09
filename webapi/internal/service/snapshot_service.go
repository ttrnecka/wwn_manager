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
