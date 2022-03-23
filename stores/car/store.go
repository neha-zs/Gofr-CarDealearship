package car

import (
	"Project/CarDealearship/models"
	"context"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type store struct{}

func New() store {
	return store{}
}

// GetCarByID function is the datastore layer function to get a car by its id
func (s store) GetCarByID(ctx *gofr.Context, Id string) (models.Car, error) {
	var c models.Car

	query := "SELECT * FROM Car WHERE ID=?;"
	err := ctx.DB().QueryRowContext(ctx, query, Id).
		Scan(&c.ID, &c.Engine.EngineID, &c.Name, &c.Year, &c.Brand, &c.FuelType)

	if err != nil {
		return models.Car{}, err
	}

	return c, nil
}

// GetCarsByBrand is a datastore layer function to get all the cars with the given brand name.
func (s store) GetCarsByBrand(ctx *gofr.Context, brand string) ([]models.Car, error) {
	var car []models.Car

	rows, err := ctx.DB().QueryContext(ctx, "select * from Car where brand=?;", brand)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var c models.Car

		err = rows.Scan(&c.ID, &c.Engine.EngineID, &c.Name, &c.Year, &c.Brand, &c.FuelType)
		if err != nil {
			return nil, errors.Error("Scan Error")
		}

		car = append(car, c)
	}

	err = rows.Err()
	if err != nil {
		return []models.Car{}, err
	}

	return car, nil
}

// CreateCar is the datastore layer function to create a model of a car
func (s store) CreateCar(ctx *gofr.Context, car *models.Car) (models.Car, error) {
	_, err := ctx.DB().ExecContext(ctx, "INSERT INTO Car (id,engine_id,name,year,brand,fuel_type) VALUES(?,?,?,?,?,?)",
		car.ID, car.Engine.EngineID, car.Name, car.Year, car.Brand, car.FuelType)
	if err != nil {
		return models.Car{}, err
	}

	return *car, nil
}

// DeleteCar to service layer function to delete the car from database
func (s store) DeleteCar(ctx *gofr.Context, id string) error {
	ctx.Context = context.TODO()
	_, err := ctx.DB().ExecContext(ctx, "DELETE FROM Car WHERE ID=?", id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCar is a datastore layer function to update a car record in database
func (s store) UpdateCar(ctx *gofr.Context, id string, car *models.Car) (models.Car, error) {
	_, err := ctx.DB().ExecContext(ctx, "UPDATE Car SET name=?,year=?,brand=?,fuel_type=? WHERE id=?",
		car.Name, car.Year, car.Brand, car.FuelType, id)
	if err != nil {
		return models.Car{}, err
	}

	return *car, nil
}
