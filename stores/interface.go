package stores

import (
	"Project/CarDealearship/models"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Car interface {
	GetCarByID(ctx *gofr.Context, id string) (models.Car, error)
	GetCarsByBrand(ctx *gofr.Context, brand string) ([]models.Car, error)
	CreateCar(ctx *gofr.Context, car *models.Car) (models.Car, error)
	DeleteCar(ctx *gofr.Context, id string) error
	UpdateCar(ctx *gofr.Context, id string, car *models.Car) (models.Car, error)
}

type Engine interface {
	EngineGetByID(ctx *gofr.Context, id string) (models.Engine, error)
	EngineCreate(ctx *gofr.Context, engine *models.Engine) (models.Engine, error)
	EngineDelete(ctx *gofr.Context, id string) error
	EngineUpdate(ctx *gofr.Context, id string, engine *models.Engine) (models.Engine, error)
}
