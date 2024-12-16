package service

import (
	"context"
	"exmemory/car/model"
)

type CarService interface {
	AddCar(ctx context.Context, car model.Car) (string, error)
	FindCarById(ctx context.Context, id string) (*model.Car, error)
	FindCarsByModel(ctx context.Context, model string) ([]*model.Car, error)
}
