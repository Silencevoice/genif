package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type TestEntity struct {
	ID    string `bson:"_id"`
	Value string `bson:"value"`
}

func TestMongoStore_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get existing entity", func(mt *mtest.T) {
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

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
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

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
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}
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
}

func TestMongoStore_GetMultipleByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get multiple existing entities", func(mt *mtest.T) {
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}
		stringObjectID1 := primitive.NewObjectID().Hex()
		stringObjectID2 := primitive.NewObjectID().Hex()
		// Simular respuesta para Find
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.test", mtest.FirstBatch,
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
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

		// Ejecutar el método con un ID no válido
		ids := []string{"invalid-id", "1"}
		result, err := store.GetMultipleByID(context.Background(), ids)

		// Validar resultados
		assert.Nil(mt, result)
		assert.Error(mt, err)
		assert.Equal(mt, "invalid ID format", err.Error())
	})

	mt.Run("No documents found", func(mt *mtest.T) {
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}
		stringObjectID1 := primitive.NewObjectID().Hex()

		// Simular respuesta vacía para Find
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.test", mtest.FirstBatch))

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
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		stringObjectID := primitive.NewObjectID().Hex()
		entity := TestEntity{ID: stringObjectID, Value: "test-value"}
		result, err := store.Insert(context.Background(), entity.ID, &entity)
		assert.NoError(mt, err)
		assert.Equal(mt, &entity, result)
	})

	mt.Run("Insert duplicate key", func(mt *mtest.T) {
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

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
}

func TestMongoStore_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Delete existing entity", func(mt *mtest.T) {
		stringObjectID := primitive.NewObjectID().Hex()
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}
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
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "errmsg", Value: "not found"},
		})

		err := store.Delete(context.Background(), stringObjectID)
		assert.Error(mt, err)
	})
}

func TestMongoStore_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Update existing entity", func(mt *mtest.T) {
		stringObjectID := primitive.NewObjectID().Hex()
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

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
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "errmsg", Value: "not found"},
		})

		entity := TestEntity{ID: stringObjectID, Value: "value"}
		err := store.Update(context.Background(), entity.ID, &entity)
		assert.Error(mt, err)
	})
}

func TestMongoStore_GetAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get all entities", func(mt *mtest.T) {
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.test", mtest.FirstBatch,
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

func TestMongoStore_Search(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Search with matching results", func(mt *mtest.T) {
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

		// Simular respuesta para Find con resultados coincidentes
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.test", mtest.FirstBatch,
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
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

		// Simular respuesta vacía para Find
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.test", mtest.FirstBatch))

		// Ejecutar el método
		filter := bson.M{"value": "non-existent"}
		result, err := store.ExecuteQuery(context.Background(), filter)

		// Validar resultados
		assert.NoError(mt, err)
		assert.Empty(mt, result)
	})

	mt.Run("Search with Find error", func(mt *mtest.T) {
		collection := mt.Coll
		store := &MongoStore[TestEntity]{collection: collection}

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
