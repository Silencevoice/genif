package mongo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type TestEntity struct {
	ID    string `bson:"_id"`
	Value string `bson:"value"`
}

func TestMongoStore_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get existing entity", func(mt *mtest.T) {

		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		// Configurar la respuesta simulada para FindOne
		stringObjectID := primitive.NewObjectID().Hex()
		entity := TestEntity{ID: stringObjectID, Value: "test-value"}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.test", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: entity.ID},
			{Key: "value", Value: entity.Value},
		}), mtest.CreateCursorResponse(int64(0), "foo.bar", mtest.NextBatch))

		// Ejecutar el método
		result, err := store.GetByID(context.Background(), stringObjectID)

		// Validar resultados
		assert.NoError(mt, err)
		assert.Equal(mt, &entity, result)
	})

	mt.Run("Get non-existent entity", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		// Configurar la respuesta simulada para FindOne con error
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "errmsg", Value: "not found"},
		})

		// Ejecutar el método
		result, err := store.GetByID(context.Background(), "non-existent")

		// Validar resultados
		assert.Nil(mt, result)
		assert.Error(mt, err)
	})

	mt.Run("Fail FindOne", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		stringObjectID := primitive.NewObjectID().Hex()

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    0,                        // Código de error para "no documents"
			Message: "no documents in result", // Mensaje descriptivo
		}))

		// Ejecutar el método
		result, err := store.GetByID(context.Background(), stringObjectID)

		// Validar resultados
		assert.Nil(mt, result)
		assert.Error(mt, err)
	})

	mt.Run("Incorrect ID format", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		// Ejecutar el método
		result, err := store.GetByID(context.Background(), "1")

		// Validar resultados
		assert.Nil(mt, result)
		assert.Error(mt, err)
		assert.ErrorContains(t, err, "invalid ID format")
	})
}

func TestMongoStore_GetMultipleByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get multiple existing entities", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		stringObjectID1 := primitive.NewObjectID().Hex()
		stringObjectID2 := primitive.NewObjectID().Hex()

		// Simular respuesta para Find
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch,
			bson.D{{Key: "_id", Value: stringObjectID1}, {Key: "value", Value: "value-1"}},
			bson.D{{Key: "_id", Value: stringObjectID2}, {Key: "value", Value: "value-2"}},
		))

		// Ejecutar el método
		ids := []string{stringObjectID1, stringObjectID2}
		result, err := store.GetMultipleByID(context.Background(), ids)

		// Validar resultados
		assert.NoError(mt, err)
		assert.Len(mt, result, 2)
		assert.Equal(mt, "value-1", result[0].Value)
		assert.Equal(mt, "value-2", result[1].Value)
	})

	mt.Run("Invalid ID format", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		// Ejecutar el método con un ID no válido
		ids := []string{"invalid-id", "1"}
		result, err := store.GetMultipleByID(context.Background(), ids)

		// Validar resultados
		assert.Nil(mt, result)
		assert.Error(mt, err)
		assert.Equal(mt, "invalid ID format", err.Error())
	})

	mt.Run("No documents found", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		stringObjectID1 := primitive.NewObjectID().Hex()

		// Simular respuesta vacía para Find
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch))

		// Ejecutar el método
		ids := []string{stringObjectID1}
		result, err := store.GetMultipleByID(context.Background(), ids)

		// Validar resultados
		assert.NoError(mt, err)
		assert.Empty(mt, result)
	})

}

func TestMongoStore_Insert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Insert successful", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		stringObjectID := primitive.NewObjectID().Hex()
		entity := TestEntity{ID: stringObjectID, Value: "test-value"}
		result, err := store.Insert(context.Background(), entity.ID, &entity)
		assert.NoError(mt, err)
		assert.Equal(mt, &entity, result)
	})

	mt.Run("Insert duplicate key", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "errmsg", Value: "already existing key"},
			{Key: "code", Value: 11000},
		})

		stringObjectID := primitive.NewObjectID().Hex()
		entity := TestEntity{ID: stringObjectID, Value: "test-value"}
		result, err := store.Insert(context.Background(), entity.ID, &entity)
		assert.Nil(mt, result)
		assert.Error(mt, err)
		assert.Contains(mt, err.Error(), "already existing key")
	})

	mt.Run("Invalid ID format", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		entity := TestEntity{ID: "1", Value: "test-value"}
		result, err := store.Insert(context.Background(), entity.ID, &entity)
		assert.Nil(mt, result)
		assert.Error(mt, err)
		assert.Contains(mt, err.Error(), "invalid ID format")
	})
}

func TestMongoStore_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Delete existing entity", func(mt *mtest.T) {
		stringObjectID := primitive.NewObjectID().Hex()
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		resp := bson.D{
			{Key: "n", Value: 1},
			{Key: "nDeleted", Value: 1},
		}
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(resp...),
		)
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := store.Delete(context.Background(), stringObjectID)
		assert.NoError(mt, err)
	})

	mt.Run("Delete non-existent entity", func(mt *mtest.T) {
		stringObjectID := primitive.NewObjectID().Hex()
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "errmsg", Value: "not found"},
		})

		err := store.Delete(context.Background(), stringObjectID)
		assert.Error(mt, err)
	})

	mt.Run("Invalid ID format", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		err := store.Delete(context.Background(), "1")
		assert.Error(mt, err)
		assert.ErrorContains(mt, err, "invalid ID format")
	})

}

func TestMongoStore_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Update existing entity", func(mt *mtest.T) {
		stringObjectID := primitive.NewObjectID().Hex()
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		resp := bson.D{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 1},
		}
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(resp...),
		)

		entity := TestEntity{ID: stringObjectID, Value: "updated-value"}
		err := store.Update(context.Background(), entity.ID, &entity)
		assert.NoError(mt, err)
	})

	mt.Run("Update non-existent entity", func(mt *mtest.T) {
		stringObjectID := primitive.NewObjectID().Hex()
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "errmsg", Value: "not found"},
		})

		entity := TestEntity{ID: stringObjectID, Value: "value"}
		err := store.Update(context.Background(), entity.ID, &entity)
		assert.Error(mt, err)
	})

	mt.Run("Invalid ID format", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		err := store.Update(context.Background(), "1", &TestEntity{})
		assert.Error(mt, err)
		assert.ErrorContains(mt, err, "invalid ID format")
	})
}

func TestMongoStore_GetAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get all entities", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch,
			bson.D{{Key: "_id", Value: "1"}, {Key: "value", Value: "value-1"}},
			bson.D{{Key: "_id", Value: "2"}, {Key: "value", Value: "value-2"}},
		))

		result, err := store.GetAll(context.Background())
		assert.NoError(mt, err)
		assert.Len(mt, result, 2)
		assert.Equal(mt, "value-1", result[0].Value)
		assert.Equal(mt, "value-2", result[1].Value)
	})
}

func TestMongoStore_ExecuteQuery(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Search with matching results", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		// Simular respuesta para Find con resultados coincidentes
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch,
			bson.D{{Key: "_id", Value: "1"}, {Key: "value", Value: "value-1"}},
			bson.D{{Key: "_id", Value: "2"}, {Key: "value", Value: "value-2"}},
		))

		// Ejecutar el método
		filter := bson.M{"value": bson.M{"$regex": "^value"}}
		result, err := store.ExecuteQuery(context.Background(), filter)

		// Validar resultados
		assert.NoError(mt, err)
		assert.Len(mt, result, 2)
		assert.Equal(mt, "value-1", result[0].Value)
		assert.Equal(mt, "value-2", result[1].Value)
	})

	mt.Run("Search with no matching results", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		// Simular respuesta vacía para Find
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch))

		// Ejecutar el método
		filter := bson.M{"value": "non-existent"}
		result, err := store.ExecuteQuery(context.Background(), filter)

		// Validar resultados
		assert.NoError(mt, err)
		assert.Empty(mt, result)
	})

	mt.Run("Search with Find error", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")

		// Simular error en Find
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "find error",
		}))

		// Ejecutar el método
		filter := bson.M{"value": bson.M{"$regex": "^value"}}
		result, err := store.ExecuteQuery(context.Background(), filter)

		// Validar resultados
		assert.Nil(mt, result)
		assert.Error(mt, err)
		assert.Contains(mt, err.Error(), "find error")
	})
}

func TestExecuteUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Update with no errors", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		update := func(ctx context.Context, collection *mongo.Collection) (int, error) {
			return 1, nil
		}
		num, err := store.ExecuteUpdate(context.Background(), update)

		assert.NoError(t, err)
		assert.Equal(t, 1, num)
	})

	mt.Run("Update with errors", func(mt *mtest.T) {
		store := NewMongoStore[TestEntity](mt.DB, "foo.bar")
		update := func(ctx context.Context, collection *mongo.Collection) (int, error) {
			return 0, errors.New("update error")
		}
		num, err := store.ExecuteUpdate(context.Background(), update)

		assert.Equal(t, 0, num)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "update error")
	})
}
