package main

import (
	"Project/CarDealearship/handlers"
	car2 "Project/CarDealearship/service/car"
	"Project/CarDealearship/stores/car"
	"Project/CarDealearship/stores/engine"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()
	k.Server.ValidateHeaders = false

	st := car.New()
	engin := engine.New()
	svc := car2.New(st, engin)
	h := handlers.New(svc)

	k.GET("/car/{id}", h.GetByID)
	k.GET("/cars", h.GetByBrand)
	k.POST("/car", h.Create)
	k.PUT("/car/{id}", h.Update)

	k.Start()

}
