package repository

import (
	"github.com/Silencevoice/go-store/examples/memory-store/car/model"
	"github.com/Silencevoice/go-store/memory"
)

type CarRepository struct {
	memory.MemStore[model.Car]
}

func NewCarRepository() *CarRepository {
	return &CarRepository{
		MemStore: *memory.NewMemStore[model.Car](),
	}
}
