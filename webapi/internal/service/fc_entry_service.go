package service

import (
	"context"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s fcWWNEntryService) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]entity.FCWWNEntry, error) {
	bson, ok := filter.(bson.M)
	if ok {
		customer, ok := bson["customer"]
		if ok {
			if customer == entity.GLOBAL_CUSTOMER {
				delete(bson, "customer")
			}
		}
	}
	return s.GenericService.Find(ctx, filter, opts...)
}
