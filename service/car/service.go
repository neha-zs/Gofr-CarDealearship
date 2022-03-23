package car

import (
	"Project/CarDealearship/models"
	"Project/CarDealearship/stores"
	"reflect"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/google/uuid"
)

type service struct {
	carStore    stores.Car
	engineStore stores.Engine
}

// nolint:revive // need not be exported
// New factory function
func New(c stores.Car, e stores.Engine) service {
	return service{carStore: c, engineStore: e}
}

// GetByID function is the service function to get a car by its id
func (service service) GetByID(ctx *gofr.Context, id string) (models.Car, error) {
	if id == uuid.Nil.String() {
		return models.Car{}, errors.InvalidParam{Param: []string{id}}
	}

	c, err := service.carStore.GetCarByID(ctx, id)
	if err != nil {
		return models.Car{}, err
	}

	engine, err := service.engineStore.EngineGetByID(ctx, id)
	if err != nil {
		return models.Car{}, err
	}

	c.Engine = engine

	return c, nil
}

// GetByBrand is a service layer function to get all the cars with the given brand name.
func (service service) GetByBrand(ctx *gofr.Context, brand string, isEngine bool) ([]models.Car, error) {
	res, err := service.carStore.GetCarsByBrand(ctx, brand)

	if isEngine {
		for i := 0; i < len(res); i++ {
			res[i].Engine, err = service.engineStore.EngineGetByID(ctx, res[i].ID.String())
		}
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create is the service layer function to create a model of a car
func (service service) Create(ctx *gofr.Context, car *models.Car) (models.Car, error) {
	c := validateCreateCar(car)

	if reflect.DeepEqual(c, models.Car{}) {
		return models.Car{}, errors.InvalidParam{}
	}

	engine, err := service.engineStore.EngineCreate(ctx, &c.Engine)
	if err != nil {
		return models.Car{}, err
	}

	c.Engine = engine
	c.ID = c.Engine.EngineID

	c, err = service.carStore.CreateCar(ctx, &c)
	if err != nil {
		return models.Car{}, err
	}

	return c, nil
}

// Update is a service layer function to update a car record in database
func (service service) Update(ctx *gofr.Context, id string, car *models.Car) (models.Car, error) {
	c, err := service.carStore.UpdateCar(ctx, id, car)
	if err != nil {
		return models.Car{}, err
	}

	engine, err := service.engineStore.EngineUpdate(ctx, id, &car.Engine)
	if err != nil {
		return models.Car{}, err
	}

	c.ID = uuid.MustParse(id)
	c.Engine = engine

	return c, nil
}

// Delete to service layer function to delete the car from database
func (service service) Delete(ctx *gofr.Context, id string) error {
	if id == uuid.Nil.String() {
		return errors.EntityNotFound{ID: id}
	}

	err := service.carStore.DeleteCar(ctx, id)
	if err != nil {
		return err
	}

	err = service.engineStore.EngineDelete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// checkBrand check if the validity of brand
func checkBrand(car *models.Car) *models.Car {
	brands := [5]string{"Tesla", "Porsche", "Ferrari", "Mercedes", "BMW"}
	flag := false

	for i := range brands {
		if brands[i] == car.Brand {
			flag = true
			break
		}
	}

	if !flag || car.Brand == "" {
		*car = models.Car{}
		return car
	}

	return car
}

// checkFuel check if the validity of fuel
func checkFuel(car *models.Car) *models.Car {
	fuelType := [3]string{"Petrol", "Diesel", "Electric"}
	flag := false

	for i := range fuelType {
		if fuelType[i] == car.FuelType {
			flag = true
			break
		}
	}

	if !flag || car.FuelType == "" {
		*car = models.Car{}
		return car
	}

	return car
}

// checkAge check if the manufacture year of the car
func checkAge(car *models.Car) *models.Car {
	y := time.Now().Year()
	if car.Year > y || car.Year < 1900 {
		*car = models.Car{}
	}

	return car
}

// validateCreateCar checks the validity of the car
func validateCreateCar(car *models.Car) models.Car {
	car = checkAge(car)
	car = checkBrand(car)
	car = checkFuel(car)

	return *car
}
