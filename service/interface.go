package service

import (
	"Project/CarDealearship/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Cars interface {
	GetByID(ctx *gofr.Context, id string) (models.Car, error)
	GetByBrand(ctx *gofr.Context, brand string, isEngine bool) ([]models.Car, error)
	Create(ctx *gofr.Context, car *models.Car) (models.Car, error)
	Delete(ctx *gofr.Context, id string) error
	Update(ctx *gofr.Context, id string, car *models.Car) (models.Car, error)
}
