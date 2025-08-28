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
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$wwn"},
			{Key: "customers", Value: bson.D{{Key: "$addToSet", Value: bson.D{
				{Key: "customer", Value: "$customer"},
				{Key: "wwn_set", Value: "$wwn_set"},
				{Key: "hostname", Value: "$hostname"},
			}}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "fc_wwn_entries"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "wwn"},
			{Key: "as", Value: "docs"},
		}}},
		{{Key: "$unwind", Value: "$docs"}},
		{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: bson.D{
				{Key: "$mergeObjects", Value: bson.A{"$docs", bson.D{
					{Key: "duplicate_customers", Value: "$customers"},
				}}},
			}},
		}}},
		{{Key: "$merge", Value: bson.D{
			{Key: "into", Value: "fc_wwn_entries"},
			{Key: "on", Value: "_id"},
			{Key: "whenMatched", Value: "merge"},
			{Key: "whenNotMatched", Value: "discard"},
		}}},
	}
	_, err := s.Collection().Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	return nil
}
