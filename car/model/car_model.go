package model

type Car struct {
	Id           string  `json:"id" bson:"_id"`                    // Identificador del coche
	Make         string  `json:"make" bson:"make"`                 // Marca del coche (e.g., Toyota, Ford)
	Model        string  `json:"model" bson:"model"`               // Modelo del coche (e.g., Corolla, Mustang)
	Year         int     `json:"year" bson:"year"`                 // Año de fabricación
	Color        string  `json:"color" bson:"color"`               // Color del coche
	Mileage      float64 `json:"mileage" bson:"mileage"`           // Kilometraje del coche
	EngineSize   float64 `json:"engine_size" bson:"engine_size"`   // Tamaño del motor en litros (e.g., 2.0)
	FuelType     string  `json:"fuel_type" bson:"fuel_type"`       // Tipo de combustible (e.g., Gasoline, Diesel, Electric)
	Transmission string  `json:"transmission" bson:"transmission"` // Tipo de transmisión (e.g., Manual, Automatic)
	Doors        int     `json:"doors" bson:"doors"`               // Número de puertas
	Price        float64 `json:"price" bson:"price"`               // Precio del coche
	VIN          string  `json:"vin" bson:"vin"`                   // Número de identificación del vehículo
	IsUsed       bool    `json:"is_used" bson:"is_used"`           // Si el coche es usado (true) o nuevo (false)
}
