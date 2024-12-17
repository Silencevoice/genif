package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestEntity struct {
	ID    string
	Value string
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()

	t.Run("Insert entity", func(t *testing.T) {
		entity := TestEntity{ID: "1", Value: "test-value"}
		inserted, err := store.Insert(ctx, entity.ID, &entity)
		require.NoError(t, err)
		assert.Equal(t, entity, *inserted)
	})

	t.Run("Insert duplicate entity", func(t *testing.T) {
		entity := TestEntity{ID: "1", Value: "test-value-duplicate"}
		_, err := store.Insert(ctx, entity.ID, &entity)
		assert.Error(t, err)
		assert.Equal(t, "already existing key", err.Error())
	})
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()
	store.Insert(ctx, "1", &TestEntity{ID: "1", Value: "test-value"})

	t.Run("Get existing entity", func(t *testing.T) {
		entity, err := store.GetByID(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "test-value", entity.Value)
	})

	t.Run("Get non-existent entity", func(t *testing.T) {
		_, err := store.GetByID(ctx, "non-existent")
		assert.Error(t, err)
		assert.Equal(t, "entity not found", err.Error())
	})
}

func TestGetMultipleByID(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()
	store.Insert(ctx, "1", &TestEntity{ID: "1", Value: "value-1"})
	store.Insert(ctx, "2", &TestEntity{ID: "2", Value: "value-2"})

	t.Run("Get multiple existing entities", func(t *testing.T) {
		entities, err := store.GetMultipleByID(ctx, []string{"1", "2"})
		require.NoError(t, err)
		assert.Len(t, entities, 2)
	})

	t.Run("Get with some non-existent IDs", func(t *testing.T) {
		_, err := store.GetMultipleByID(ctx, []string{"1", "non-existent"})
		assert.Error(t, err)
		assert.Equal(t, "entity not found", err.Error())
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()
	store.Insert(ctx, "1", &TestEntity{ID: "1", Value: "test-value"})

	t.Run("Update existing entity", func(t *testing.T) {
		updatedEntity := TestEntity{ID: "1", Value: "updated-value"}
		err := store.Update(ctx, "1", &updatedEntity)
		require.NoError(t, err)

		entity, err := store.GetByID(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "updated-value", entity.Value)
	})

	t.Run("Update non-existent entity", func(t *testing.T) {
		entity := TestEntity{ID: "non-existent", Value: "value"}
		err := store.Update(ctx, "non-existent", &entity)
		assert.Error(t, err)
		assert.Equal(t, "entity not found", err.Error())
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()
	store.Insert(ctx, "1", &TestEntity{ID: "1", Value: "test-value"})

	t.Run("Delete existing entity", func(t *testing.T) {
		err := store.Delete(ctx, "1")
		require.NoError(t, err)

		_, err = store.GetByID(ctx, "1")
		assert.Error(t, err)
		assert.Equal(t, "entity not found", err.Error())
	})

	t.Run("Delete non-existent entity", func(t *testing.T) {
		err := store.Delete(ctx, "non-existent")
		assert.Error(t, err)
		assert.Equal(t, "entity not found", err.Error())
	})
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()
	store.Insert(ctx, "1", &TestEntity{ID: "1", Value: "value-1"})
	store.Insert(ctx, "2", &TestEntity{ID: "2", Value: "value-2"})

	t.Run("Get all entities", func(t *testing.T) {
		all, err := store.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, all, 2)

		expected := map[string]string{
			"1": "value-1",
			"2": "value-2",
		}

		for _, entity := range all {
			assert.Equal(t, expected[entity.ID], entity.Value)
		}
	})
}

func TestExecuteQuery(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()
	store.Insert(ctx, "1", &TestEntity{ID: "1", Value: "value-1"})
	store.Insert(ctx, "2", &TestEntity{ID: "2", Value: "value-2"})

	t.Run("Search entities", func(t *testing.T) {
		result, err := store.ExecuteQuery(ctx, func(ctx context.Context, data map[string]TestEntity) ([]*TestEntity, error) {
			var entities []*TestEntity
			for _, entity := range data {
				if entity.Value == "value-2" {
					entities = append(entities, &entity)
				}
			}
			return entities, nil
		})

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "value-2", result[0].Value)
	})
}

func TestExecuteUpdate(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()
	store.Insert(ctx, "1", &TestEntity{ID: "1", Value: "value-1"})
	store.Insert(ctx, "2", &TestEntity{ID: "2", Value: "value-2"})

	t.Run("Update entities OK", func(t *testing.T) {
		num, err := store.ExecuteUpdate(ctx, func(ctx context.Context, data map[string]TestEntity) (int, error) {
			ent, ok := data["1"]
			if !ok {
				return 0, errors.New("not found entity")
			}
			ent.Value = "updated-value-1"
			data["1"] = ent
			return 1, nil
		})

		require.NoError(t, err)
		assert.Equal(t, 1, num)

		ent, err := store.GetByID(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "updated-value-1", ent.Value)
	})
}
