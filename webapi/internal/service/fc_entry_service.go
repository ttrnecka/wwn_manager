package service

import (
	"context"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FCWWNEntryService interface {
	GenericService[entity.FCWWNEntry]
	Customers(context.Context) ([]any, error)
	FlagDuplicateWWNs(context.Context) error
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

func (s fcWWNEntryService) FlagDuplicateWWNs(ctx context.Context) error {
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", "$wwn"},
			{"customers", bson.D{{"$addToSet", "$customer"}}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$match", bson.D{
			{"count", bson.D{{"$gt", 1}}},
		}}},
		{{"$lookup", bson.D{
			{"from", "fc_wwn_entries"},
			{"localField", "_id"},
			{"foreignField", "wwn"},
			{"as", "docs"},
		}}},
		{{"$unwind", "$docs"}},
		{{"$replaceRoot", bson.D{
			{"newRoot", bson.D{
				{"$mergeObjects", bson.A{"$docs", bson.D{
					{"duplicate_customers", "$customers"},
				}}},
			}},
		}}},
		{{"$merge", bson.D{
			{"into", "fc_wwn_entries"},
			{"on", "_id"},
			{"whenMatched", "merge"},
			{"whenNotMatched", "discard"},
		}}},
	}
	_, err := s.Collection().Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	return nil
}
