package repository

import (
	"genif/car/model"
	"genif/store/memory"
)

type CarRepository struct {
	memory.MemStore[model.Car]
}

func NewCarRepository() *CarRepository {
	return &CarRepository{
		MemStore: *memory.NewMemStore[model.Car](),
	}
}
