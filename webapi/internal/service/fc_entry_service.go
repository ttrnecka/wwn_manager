package service

import (
	"context"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FCWWNEntryService interface {
	GenericService[entity.FCWWNEntry]
	Customers(context.Context) ([]any, error)
	FlagDuplicateWWNs(context.Context, Filter) error
	GetUniqueRules(context.Context) ([]string, error)
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

func (s fcWWNEntryService) FindWithSoftDeleted(ctx context.Context, filter Filter, opt SortOption) ([]entity.FCWWNEntry, error) {
	customer, ok := filter["customer"]
	if ok {
		if customer == entity.GLOBAL_CUSTOMER {
			delete(filter, "customer")
		}
	}

	return s.GenericService.FindWithSoftDeleted(ctx, filter, opt)
}

func (s fcWWNEntryService) FlagDuplicateWWNs(ctx context.Context, filter Filter) error {
	filters := bson.D{}

	filters = append(filters, bson.E{
		Key:   "count",
		Value: bson.D{{Key: "$gt", Value: 1}},
	})

	wwns, ok := filter["wwn"]
	if ok {
		filters = append(filters, bson.E{
			Key:   "_id",
			Value: bson.D{{Key: "$in", Value: wwns}},
		})
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$wwn"},
			{Key: "customers", Value: bson.D{{Key: "$addToSet", Value: bson.D{
				{Key: "customer", Value: "$customer"},
				{Key: "wwn_set", Value: "$wwn_set"},
				{Key: "hostname", Value: "$hostname"},
			}}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$match", Value: filters}},
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

// GetUniqueRules pulls unique ObjectIDs from type_rule, hostname_rule, and reconcile_rules
func (s fcWWNEntryService) GetUniqueRules(ctx context.Context) ([]string, error) {
	pipeline := mongo.Pipeline{
		{
			{"$project", bson.D{
				{"allRules", bson.D{
					{"$setUnion", bson.A{
						bson.A{bson.D{{"$ifNull", bson.A{"$type_rule", nil}}}},
						bson.A{bson.D{{"$ifNull", bson.A{"$hostname_rule", nil}}}},
						bson.D{{"$ifNull", bson.A{"$reconcile_rules", bson.A{}}}},
					}},
				}},
			}},
		},
		{
			{"$unwind", "$allRules"},
		},
		{
			{"$group", bson.D{
				{"_id", nil},
				{"uniqueRules", bson.D{{"$addToSet", "$allRules"}}},
			}},
		},
		{
			{"$project", bson.D{
				{"_id", 0},
				{"uniqueRules", 1},
			}},
		},
	}

	cursor, err := s.Collection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		UniqueRules []primitive.ObjectID `bson:"uniqueRules"`
	}

	var uniqueRules []string

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	for _, r := range result.UniqueRules {
		uniqueRules = append(uniqueRules, r.Hex())
	}

	return uniqueRules, nil
}
