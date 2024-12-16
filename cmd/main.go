package main

import (
	"context"
	"fmt"
	"genif/car/model"
	"genif/car/repository"
	"genif/car/service"
	"os"
)

func main() {
	// Get Car Service
	srv := service.NewCarService(repository.NewCarRepository())

	// Insert a new Car
	car := model.Car{
		Id:           "123456789",
		Make:         "Ford",
		Model:        "Fiesta",
		Year:         2022,
		Color:        "Blue",
		Mileage:      56000,
		EngineSize:   1.4,
		FuelType:     "Gasoline",
		Transmission: "Manual",
		Doors:        3,
		Price:        14000,
		VIN:          "54895245-5658623",
		IsUsed:       true,
	}
	id, err := srv.AddCar(context.Background(), car)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding car: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Inserted id: %s\n", id)

	// Find the car
	found, err := srv.FindCarById(context.Background(), id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding car: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Found car %v\n", *found)

	// Find all cars by model
	res, err := srv.FindCarsByModel(context.Background(), "Fiesta")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding cars by model: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Found cars by model Fiesta: %v\n", res)

}
