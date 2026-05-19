package gormongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

	ctx            context.Context
	db             *mongo.Database
	collectionName string
}

func (m *Model) Init(ctx context.Context, db *mongo.Database, collection string) {
	m.ctx = ctx
	m.db = db
	m.collectionName = collection
}

func (m *Model) collection() *mongo.Collection {
	return m.db.Collection(m.collectionName)
}

func (m *Model) FindAll(filter interface{}, result interface{}, opts ...*options.FindOptions) error {
	if filter == nil {
		filter = bson.D{}
	}

	cursor, err := m.collection().Find(m.ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(m.ctx)

	if err = cursor.All(m.ctx, result); err != nil {
		return err
	}

	return nil
}

func (m *Model) FindOne(filter interface{}, result interface{}) error {
	if filter == nil {
		filter = bson.D{}
	}

	err := m.collection().FindOne(m.ctx, filter).Decode(result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("document not found")
		}
		return err
	}

	return nil
}

func (m *Model) FindById(id string, result interface{}) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	return m.FindOne(bson.M{"_id": objectId}, result)
}

func (m *Model) Create(document interface{}) (*mongo.InsertOneResult, error) {
	now := time.Now()
	switch d := document.(type) {
	case map[string]interface{}:
		d["created_at"] = now
		d["updated_at"] = now
	}

	result, err := m.collection().InsertOne(m.ctx, document)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Model) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	updateWithTime := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Merge avec l'update existant
	if u, ok := update.(bson.M); ok {
		if set, ok := u["$set"].(bson.M); ok {
			set["updated_at"] = time.Now()
			updateWithTime = u
		}
	}

	result, err := m.collection().UpdateOne(m.ctx, filter, updateWithTime)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Model) UpdateById(id string, update interface{}) (*mongo.UpdateResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	return m.Update(bson.M{"_id": objectId}, update)
}

func (m *Model) Delete(filter interface{}) (*mongo.DeleteResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	result, err := m.collection().DeleteOne(m.ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Model) DeleteById(id string) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	return m.Delete(bson.M{"_id": objectId})
}

func (m *Model) Count(filter interface{}) (int64, error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := m.collection().CountDocuments(m.ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
