package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore[T any] struct {
	collection *mongo.Collection
}

func NewMongoStore[T any](db *mongo.Database, collectionName string) *MongoStore[T] {
	return &MongoStore[T]{
		collection: db.Collection(collectionName),
	}
}

func (m *MongoStore[T]) GetByID(ctx context.Context, id string) (*T, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var result T
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("entity not found")
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoStore[T]) GetMultipleByID(ctx context.Context, ids []string) ([]*T, error) {
	objectIDs := []primitive.ObjectID{}
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New("invalid ID format")
		}
		objectIDs = append(objectIDs, objectID)
	}

	cursor, err := m.collection.Find(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*T
	for cursor.Next(ctx) {
		var entity T
		if err := cursor.Decode(&entity); err != nil {
			return nil, err
		}
		results = append(results, &entity)
	}

	return results, nil
}

func (m *MongoStore[T]) GetAll(ctx context.Context) ([]*T, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*T
	for cursor.Next(ctx) {
		var entity T
		if err := cursor.Decode(&entity); err != nil {
			return nil, err
		}
		results = append(results, &entity)
	}

	return results, nil
}

func (m *MongoStore[T]) Insert(ctx context.Context, id string, entity *T) (*T, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	_, err = m.collection.InsertOne(ctx, bson.M{
		"_id":  objectID,
		"data": entity,
	})
	if mongo.IsDuplicateKeyError(err) {
		return nil, errors.New("already existing key")
	} else if err != nil {
		return nil, err
	}

	return entity, nil
}

func (m *MongoStore[T]) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	res, err := m.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("entity not found")
	}

	return nil
}

func (m *MongoStore[T]) Update(ctx context.Context, id string, entity *T) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	res, err := m.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": entity})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("entity not found")
	}

	return nil
}

func (m *MongoStore[T]) ExecuteQuery(ctx context.Context, filter bson.M) ([]*T, error) {
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*T
	for cursor.Next(ctx) {
		var entity T
		if err := cursor.Decode(&entity); err != nil {
			return nil, err
		}
		results = append(results, &entity)
	}

	return results, nil
}

func (m *MongoStore[T]) ExecuteUpdate(ctx context.Context, f func(ctx context.Context, collection *mongo.Collection) (int, error)) (int, error) {
	return f(ctx, m.collection)
}
