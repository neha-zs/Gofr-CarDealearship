package handlers

import (
	"Project/CarDealearship/models"
	"Project/CarDealearship/service"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"strconv"
)

type handler struct {
	service service.Cars
}

// nolint:revive // need not be exported
// New factory function
func New(c service.Cars) handler {
	return handler{service: c}
}

type response struct {
	Customers []models.Car
}

// GetByID function is the delivery function to get a car by its id
func (c handler) GetByID(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	resp, err := c.service.GetByID(ctx, id)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetByBrand is a handler function to get all the cars with the given brand name.
func (c handler) GetByBrand(ctx *gofr.Context) (interface{}, error) {
	brand := ctx.Param("brand")
	isEngine := ctx.Param("isEngine")

	isEng, err := strconv.ParseBool(isEngine)
	if err != nil {
		return nil, err
	}

	resp, err := c.service.GetByBrand(ctx, brand, isEng)
	if err != nil {
		return nil, err
	}

	r := response{Customers: resp}

	return r, nil
}

// Create is the delivery function to create a model of a car
func (c handler) Create(ctx *gofr.Context) (interface{}, error) {
	var car models.Car
	if err := ctx.Bind(&car); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := c.service.Create(ctx, &car)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update is a handler function to update a car record in database
func (c handler) Update(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	var car models.Car
	if err := ctx.Bind(&car); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	res, err := c.service.Update(ctx, id, &car)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Delete is a handler function to delete a car record from database.
func (c handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	if err := c.service.Delete(ctx, id); err != nil {
		return nil, err
	}

	return "Deleted successfully", nil
}
