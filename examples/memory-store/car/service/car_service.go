package service

import (
	"context"
	"exmemory/car/model"
	"exmemory/car/repository"
)

type CarServiceImpl struct {
	repo *repository.CarRepository
}

func (s CarServiceImpl) AddCar(ctx context.Context, car model.Car) (string, error) {
	id := car.Id
	_, err := s.repo.Insert(ctx, id, &car)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s CarServiceImpl) FindCarById(ctx context.Context, id string) (*model.Car, error) {
	car, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (s CarServiceImpl) FindCarsByModel(ctx context.Context, carModel string) ([]*model.Car, error) {
	filterFunc := func(ctx context.Context, data map[string]model.Car) ([]*model.Car, error) {
		var result []*model.Car
		for _, car := range data {
			if car.Model == carModel {
				result = append(result, &car)
			}
		}
		return result, nil
	}

	return s.repo.ExecuteQuery(ctx, filterFunc)
}

func NewCarService(carRepo *repository.CarRepository) CarService {
	return &CarServiceImpl{
		repo: carRepo,
	}
}
