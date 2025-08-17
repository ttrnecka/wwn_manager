package service

import (
	"context"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type FCEntryService interface {
	GenericService[entity.FCEntry]
	Customers(context.Context) ([]any, error)
}

type fcEntryService struct {
	GenericService[entity.FCEntry]
}

func NewFCEntryService(p repository.FCEntryRepository) FCEntryService {
	return &fcEntryService{
		GenericService: NewGenericService(p)}
}

func (s fcEntryService) Customers(ctx context.Context) ([]any, error) {
	result, err := s.Collection().Distinct(context.Background(), "customer", bson.M{})
	if err != nil {
		return nil, err
	}
	return result, nil
}
