package service

import (
	"context"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type FCWWNEntryService interface {
	GenericService[entity.FCWWNEntry]
	Customers(context.Context) ([]any, error)
}

type fcWWNEntryService struct {
	GenericService[entity.FCWWNEntry]
}

func NewFCWWNEntryService(p repository.FCWWNEntryRepository) FCWWNEntryService {
	return &fcWWNEntryService{
		GenericService: NewGenericService(p)}
}

func (s fcWWNEntryService) Customers(ctx context.Context) ([]any, error) {
	result, err := s.Collection().Distinct(context.Background(), "customer", bson.M{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s fcWWNEntryService) Find(ctx context.Context, filter Filter, opt SortOption) ([]entity.FCWWNEntry, error) {
	customer, ok := filter["customer"]
	if ok {
		if customer == entity.GLOBAL_CUSTOMER {
			delete(filter, "customer")
		}
	}

	return s.GenericService.Find(ctx, filter, opt)
}
