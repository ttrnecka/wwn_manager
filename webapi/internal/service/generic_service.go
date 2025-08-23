package service

import (
	"context"
	"fmt"

	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DependencyDeleteFunc func(ctx context.Context, parentID primitive.ObjectID) error

type GenericService[T any] interface {
	All(context.Context) ([]T, error)
	Get(context.Context, string) (*T, error)
	DeleteAll(context.Context) error
	Delete(context.Context, string) error
	Update(context.Context, primitive.ObjectID, *T) (primitive.ObjectID, error)
	RegisterDependencies(...DependencyDeleteFunc)
	Collection() *mongo.Collection
	Find(context.Context, Filter, SortOption) ([]T, error)
	InsertAll(context.Context, []T) error
}

type Filter = map[string]interface{}
type SortOption = map[string]string

type genericService[T any] struct {
	MainRepo    repository.GenericRepository[T]
	dependecies []DependencyDeleteFunc
}

func NewGenericService[T any](r repository.GenericRepository[T]) GenericService[T] {
	return &genericService[T]{
		MainRepo: r,
	}
}

func (s *genericService[T]) All(ctx context.Context) ([]T, error) {
	return s.MainRepo.All(ctx)
}

func (s *genericService[T]) Find(ctx context.Context, filter Filter, opt SortOption) ([]T, error) {
	var mopts []*options.FindOptions
	for k, v := range opt {
		order := 1
		if v == "desc" {
			order = -1
		}
		mopts = append(mopts, options.Find().SetSort(bson.M{k: order}))
	}
	// explicint casting - both are map[string]interface{} but different types and Find ignores anything that is not bson.M
	bson := bson.M(filter)
	return s.MainRepo.Find(ctx, bson, mopts...)
}

func (s *genericService[T]) Get(ctx context.Context, id string) (*T, error) {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.MainRepo.GetByID(ctx, idp)
}

func (s *genericService[T]) Delete(ctx context.Context, id string) error {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = s.DeleteDependencies(ctx, idp)
	if err != nil {
		return err
	}
	return s.MainRepo.HardDeleteByID(ctx, idp)
}

func (s *genericService[T]) DeleteAll(ctx context.Context) error {
	return s.MainRepo.HardDelete(ctx, bson.M{})
}

func (s *genericService[T]) Update(ctx context.Context, id primitive.ObjectID, item *T) (primitive.ObjectID, error) {
	if id.IsZero() {
		return s.MainRepo.Create(ctx, item)
	}
	return id, s.MainRepo.UpdateByID(ctx, id, item)
}

func (s *genericService[T]) DeleteDependencies(ctx context.Context, parentID primitive.ObjectID) error {
	for _, fn := range s.dependecies {
		if err := fn(ctx, parentID); err != nil {
			return fmt.Errorf("failed to delete dependency: %w", err)
		}
	}
	return nil
}

func (s *genericService[T]) InsertAll(ctx context.Context, items []T) error {
	return s.MainRepo.InsertAll(ctx, items)
}

func (s *genericService[T]) RegisterDependencies(fn ...DependencyDeleteFunc) {
	s.dependecies = append(s.dependecies, fn...)
}

func (s *genericService[T]) Collection() *mongo.Collection {
	return s.MainRepo.GetCollection()
}
