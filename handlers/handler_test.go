package handlers

import (
	"Project/CarDealearship/models"
	"Project/CarDealearship/service"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

// TestGetByID to test the handler GetByID
func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCars(ctrl)
	s := New(mockService)
	app := gofr.New()

	id1 := uuid.New()
	id3 := uuid.New()

	testCar := models.Car{
		ID: id1, Engine: models.Engine{EngineID: id1, Displacement: 500, Cylinders: 2, Range: 200},
		Name: "Model 3", Year: 2018, Brand: "Tesla", FuelType: "petrol"}

	testCases := []struct {
		desc string
		id   uuid.UUID
		resp models.Car
		err  error
		mock []*gomock.Call
	}{
		{
			desc: "success case",
			id:   id1,
			err:  nil,
			resp: testCar,
			mock: []*gomock.Call{mockService.EXPECT().GetByID(gomock.Any(), id1.String()).
				Return(testCar, nil)}},
		{
			desc: "not found",
			id:   id3,
			resp: models.Car{},
			err:  errors.EntityNotFound{Entity: "Car", ID: id3.String()},
			mock: []*gomock.Call{mockService.EXPECT().GetByID(gomock.Any(), id3.String()).
				Return(models.Car{}, errors.EntityNotFound{Entity: "Car", ID: id3.String()})},
		},
		{
			desc: "not found",
			id:   uuid.Nil,
			resp: models.Car{},
			err:  errors.InvalidParam{Param: []string{"id"}},
			mock: []*gomock.Call{
				mockService.EXPECT().GetByID(gomock.Any(), uuid.Nil.String()).
					Return(models.Car{}, errors.InvalidParam{Param: []string{"id"}})},
		},
	}

	for i, tc := range testCases {
		r := httptest.NewRequest("GET", "/car/{id}"+tc.id.String(), nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id.String(),
		})

		resp, err := s.GetByID(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
