package service

import (
	"context"
	"testing"

	"github.com/Silencevoice/go-store/examples/memory-store/car/model"
	"github.com/Silencevoice/go-store/examples/memory-store/car/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddCar(t *testing.T) {
	ctx := context.Background()

	t.Run("Add Car successful", func(t *testing.T) {
		store := repository.NewCarRepository()
		srv := NewCarService(store)

		car := model.Car{Id: "1", Make: "Toyota", Model: "Corolla", Year: 2020}
		inserted, err := srv.AddCar(ctx, car)
		assert.NoError(t, err)
		assert.Equal(t, car.Id, inserted)
	})

	t.Run("Add car repo error", func(t *testing.T) {
		store := repository.NewCarRepository()
		srv := NewCarService(store)

		car := model.Car{Id: "1", Make: "Toyota", Model: "Corolla", Year: 2020}
		inserted, err := srv.AddCar(ctx, car)
		assert.NoError(t, err)
		assert.Equal(t, car.Id, inserted)

		_, err = srv.AddCar(ctx, car)
		assert.ErrorContains(t, err, "already existing key")
	})
}

func TestFindCarById(t *testing.T) {
	ctx := context.Background()

	t.Run("GetByID successful", func(t *testing.T) {
		store := repository.NewCarRepository()
		srv := NewCarService(store)

		car := model.Car{Id: "1", Make: "Toyota", Model: "Corolla", Year: 2020}
		inserted, err := srv.AddCar(ctx, car)
		assert.NoError(t, err)
		assert.Equal(t, car.Id, inserted)

		found, err := srv.FindCarById(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, car, *found)
	})

	t.Run("GetByID error", func(t *testing.T) {
		store := repository.NewCarRepository()
		srv := NewCarService(store)

		found, err := srv.FindCarById(ctx, "1")

		assert.Nil(t, found)
		assert.ErrorContains(t, err, "entity not found")

	})
}

func TestFindCarsByModel(t *testing.T) {
	ctx := context.Background()
	store := repository.NewCarRepository()
	srv := NewCarService(store)

	car1 := model.Car{Id: "1", Make: "Toyota", Model: "Corolla", Year: 2020}
	inserted1, err := srv.AddCar(ctx, car1)
	assert.NoError(t, err)
	assert.Equal(t, car1.Id, inserted1)
	car2 := model.Car{Id: "2", Make: "Toyota", Model: "Yaris", Year: 2020}
	inserted2, err := srv.AddCar(ctx, car2)
	assert.NoError(t, err)
	assert.Equal(t, car2.Id, inserted2)

	found, err := srv.FindCarsByModel(ctx, "Corolla")
	require.NoError(t, err)
	assert.Len(t, found, 1)
	assert.Equal(t, car1, *found[0])
}
