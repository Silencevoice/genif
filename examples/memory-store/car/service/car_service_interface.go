package service

import (
	"context"

	"github.com/Silencevoice/go-store/examples/memory-store/car/model"
)

type CarService interface {
	AddCar(ctx context.Context, car model.Car) (string, error)
	FindCarById(ctx context.Context, id string) (*model.Car, error)
	FindCarsByModel(ctx context.Context, model string) ([]*model.Car, error)
}
