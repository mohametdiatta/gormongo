package gormongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IModel interface {
	Init(ctx context.Context, db *mongo.Database, collection string)
	FindAll(filter interface{}, result interface{}, opts ...*options.FindOptions) error // ✅
	FindOne(filter interface{}, result interface{}) error
	FindById(id string, result interface{}) error
	Create(document interface{}) (*mongo.InsertOneResult, error)
	Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	UpdateById(id string, update interface{}) (*mongo.UpdateResult, error)
	Delete(filter interface{}) (*mongo.DeleteResult, error)
	DeleteById(id string) (*mongo.DeleteResult, error)
	Count(filter interface{}) (int64, error)
}

type Registry struct {
	models map[string]IModel
	ctx    context.Context
	db     *mongo.Database
}

func NewRegistry(ctx context.Context, db *mongo.Database) *Registry {
	return &Registry{
		models: make(map[string]IModel),
		ctx:    ctx,
		db:     db,
	}
}

func (r *Registry) Register(name string, model IModel, collection string) {
	model.Init(r.ctx, r.db, collection)
	r.models[name] = model
}

func (r *Registry) Get(name string) (IModel, error) {
	model, exists := r.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s not found", name)
	}
	return model, nil
}
