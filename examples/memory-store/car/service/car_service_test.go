package service

import (
	"context"
	"testing"

	"github.com/Silencevoice/go-store/examples/memory-store/car/model"
	"github.com/Silencevoice/go-store/examples/memory-store/car/repository"
	"github.com/stretchr/testify/assert"
)

func Test_AddCar(t *testing.T) {
	ctx := context.Background()

	store := repository.NewCarRepository()
	car := model.Car{Id: "1", Make: "Toyota", Model: "Corolla", Year: 2020}
	inserted, err := store.Insert(ctx, car.Id, &car)
	assert.NoError(t, err)
	assert.Equal(t, car, *inserted)
}
