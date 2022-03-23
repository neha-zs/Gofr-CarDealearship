package car

import (
	"Project/CarDealearship/models"
	"Project/CarDealearship/stores"

	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetByID to test the service GetByID
func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCar := stores.NewMockCar(ctrl)
	mockEngine := stores.NewMockEngine(ctrl)
	carService := New(mockCar, mockEngine)
	ctx := gofr.NewContext(nil, nil, gofr.New())

	id := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	testCases := []struct {
		desc     string
		id       uuid.UUID
		expected models.Car
	}{
		{
			desc: "success case",
			id:   id,
			expected: models.Car{ID: id, Name: "Model 3", Year: 2010, Brand: "Tesla", FuelType: "diesel",
				Engine: models.Engine{EngineID: id, Displacement: 200, Cylinders: 1, Range: 0}},
		},
		{
			desc:     "not found",
			id:       uuid.Nil,
			expected: models.Car{}},
		{
			desc:     "uuid not present",
			id:       uuid.New(),
			expected: models.Car{},
		},
		{
			desc:     "GetByID returns error",
			id:       id2,
			expected: models.Car{},
		},
		{
			desc:     "Cargetbyid  returns error in engine",
			id:       id3,
			expected: models.Car{},
		},
		{
			desc:     "EngineGetByID returns error",
			id:       id2,
			expected: models.Car{},
		},
	}

	for i, tc := range testCases {
		if tc.id == id2 {
			mockCar.EXPECT().GetCarByID(ctx, tc.id.String()).Return(tc.expected, errors.Error("err"))
		} else if tc.id == id3 {
			mockCar.EXPECT().GetCarByID(ctx, tc.id.String()).Return(tc.expected, nil)
			mockEngine.EXPECT().EngineGetByID(ctx, tc.id.String()).
				Return(tc.expected.Engine, errors.Error("err"))
		} else if tc.id != uuid.Nil {
			mockCar.EXPECT().GetCarByID(ctx, tc.id.String()).Return(tc.expected, nil)
			mockEngine.EXPECT().EngineGetByID(ctx, tc.id.String()).Return(tc.expected.Engine, nil)
		}

		res, _ := carService.GetByID(ctx, tc.id.String())

		assert.Equal(t, res, tc.expected,
			"%v [TEST%d]Failed. Got %v\tExpected %v\n", tc.desc, i+1, res, tc.expected)
	}
}

// TestGetByBrand to test the GetByBrand service
func TestGetByBrand(t *testing.T) {
	id := uuid.New()
	id2 := uuid.New()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockCar := stores.NewMockCar(ctrl)
	mockEngine := stores.NewMockEngine(ctrl)
	carService := New(mockCar, mockEngine)
	ctx := gofr.NewContext(nil, nil, gofr.New())

	testCases := []struct {
		desc        string
		id          uuid.UUID
		CarBrand    string
		CarEngine   bool
		expectedRes []models.Car
	}{
		{desc: "Success Case", id: id, CarBrand: "Tesla", CarEngine: true,
			expectedRes: []models.Car{{ID: id, Engine: models.Engine{EngineID: id, Displacement: 200, Cylinders: 4},
				Name: "Model 3", Year: 1990, Brand: "Tesla", FuelType: "Electric",
			}}},
		{
			desc:        "EngineGetByID returns error",
			id:          id2,
			CarBrand:    "BMW",
			CarEngine:   true,
			expectedRes: []models.Car(nil),
		},
	}

	mockCar.EXPECT().GetCarsByBrand(ctx, "BMW").
		Return([]models.Car(nil), errors.InvalidParam{})

	mockCar.EXPECT().GetCarsByBrand(ctx, "Tesla").
		Return([]models.Car{{ID: id, Engine: models.Engine{EngineID: id, Displacement: 200, Cylinders: 4},
			Name: "Model 3", Year: 1990, Brand: "Tesla", FuelType: "Electric"}}, nil)

	for i, tc := range testCases {
		if tc.CarBrand == "" || tc.CarBrand == "Tesla" {
			mockEngine.EXPECT().EngineGetByID(ctx, tc.id.String()).Return(tc.expectedRes[i].Engine, nil)
		}

		res, _ := carService.GetByBrand(ctx, tc.CarBrand, tc.CarEngine)

		assert.Equal(t, res, tc.expectedRes,
			" [TEST%d]Failed. Got %v\tExpected %v\n", i+1, res, tc.expectedRes)
	}
}

// TestCreate to test the Create service
func TestCreate(t *testing.T) {
	var (
		id = uuid.New()
		c1 = models.Car{ID: id, Name: "Model 3", Year: 2020, Brand: "Tesla", FuelType: "Diesel",
			Engine: models.Engine{EngineID: id, Displacement: 200, Cylinders: 6}}
		c2 = models.Car{ID: id, Name: "Model 5", Year: 2021, Brand: "BMW", FuelType: "Diesel",
			Engine: models.Engine{EngineID: id, Displacement: 400, Cylinders: 2}}
		c3 = models.Car{}
		c5 = models.Car{ID: id, Name: "Model 7", Year: 2020, Brand: "ABC", FuelType: "Diesel",
			Engine: models.Engine{EngineID: id, Displacement: 250, Cylinders: 4}}
		c4 = models.Car{ID: id, Name: "Mod 2", Year: 2020, Brand: "BMW", FuelType: "Diesel",
			Engine: models.Engine{EngineID: id, Displacement: 250, Cylinders: 3}}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCar := stores.NewMockCar(ctrl)
	mockEngine := stores.NewMockEngine(ctrl)
	carService := New(mockCar, mockEngine)
	ctx := gofr.NewContext(nil, nil, gofr.New())

	testCases := []struct {
		desc   string
		input  models.Car
		output models.Car
	}{
		{desc: "success case", input: c1, output: c1},
		{desc: "Create returns error", input: c2, output: c3},
		{desc: "EngineCreate returns error", input: c4, output: c3},
		{desc: "Brand not present", input: c5, output: c3},
	}

	mockCar.EXPECT().CreateCar(ctx, &c1).Return(c1, nil)
	mockEngine.EXPECT().EngineCreate(ctx, &c1.Engine).Return(c1.Engine, nil)

	mockCar.EXPECT().CreateCar(ctx, &c2).Return(c2, errors.InvalidParam{})
	mockEngine.EXPECT().EngineCreate(ctx, &c2.Engine).Return(c2.Engine, nil)

	mockEngine.EXPECT().EngineCreate(ctx, &c4.Engine).Return(c4.Engine, errors.InvalidParam{})

	for i := range testCases {
		res, _ := carService.Create(ctx, &testCases[i].input)
		assert.Equal(t, res, testCases[i].output,
			" [TEST%d]Failed. Got %v\tExpected %v\n", i+1, res, testCases[i].output)
	}
}

// TestUpdate to test the Update service
func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCar := stores.NewMockCar(ctrl)
	mockEngine := stores.NewMockEngine(ctrl)
	carService := New(mockCar, mockEngine)
	ctx := gofr.NewContext(nil, nil, gofr.New())
	var (
		id = uuid.New()
		c1 = models.Car{ID: id, Name: "Cayenne", Year: 2020, Brand: "Porsche", FuelType: "diesel",
			Engine: models.Engine{EngineID: id, Displacement: 100, Cylinders: 6, Range: 120}}
		c2 = models.Car{ID: id, Name: "Cayenne", Year: 2020, Brand: "Porsche", FuelType: "diesel",
			Engine: models.Engine{EngineID: id, Displacement: 100, Cylinders: 6, Range: 120}}
		c3 = models.Car{}
		c4 = models.Car{ID: id, Name: "Cayenne", Year: 2020, Brand: "Porsche", FuelType: "diesel",
			Engine: models.Engine{EngineID: id, Displacement: 100, Cylinders: 6, Range: 120}}
	)

	testCases := []struct {
		desc   string
		id     uuid.UUID
		input  models.Car
		output models.Car
	}{
		{desc: "success case", id: id, input: c1, output: c1},
		{desc: "Error in updateCar", id: id, input: c2, output: c3},
		{desc: "error in UpdateEngine", id: id, input: c4, output: c3},
	}

	mockCar.EXPECT().UpdateCar(ctx, c1.ID.String(), &c1).Return(c1, nil)
	mockEngine.EXPECT().EngineUpdate(ctx, c1.ID.String(), &c1.Engine).Return(c1.Engine, nil)

	mockCar.EXPECT().UpdateCar(ctx, c2.ID.String(), &c2).Return(c3, errors.InvalidParam{})

	mockCar.EXPECT().UpdateCar(ctx, c4.ID.String(), &c4).Return(c3, nil)
	mockEngine.EXPECT().EngineUpdate(ctx, c4.ID.String(), &c4.Engine).Return(c3.Engine, errors.InvalidParam{})

	for i, tc := range testCases {
		car, _ := carService.Update(ctx, tc.id.String(), &tc.input)
		assert.Equal(t, car, tc.output, "[TEST%d]Failed. Got %v\tExpected %v\n", i+1, car, tc.output)
	}
}

// TestDelete to test the Delete handler
func TestDelete(t *testing.T) {
	var (
		id  = uuid.New()
		id2 = uuid.New()
		id3 = uuid.New()
	)

	tests := []struct {
		desc   string
		id     uuid.UUID
		output models.Car
		err    error
	}{
		{"success case", id, models.Car{}, nil},
		{"Nil UUID", uuid.Nil, models.Car{}, errors.EntityNotFound{ID: uuid.Nil.String()}},
		{"error in Delete", id2, models.Car{}, errors.InvalidParam{}},
		{"error in DeleteEngine", id3, models.Car{}, errors.InvalidParam{}},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCar := stores.NewMockCar(ctrl)
	mockEngine := stores.NewMockEngine(ctrl)
	carService := New(mockCar, mockEngine)
	ctx := gofr.NewContext(nil, nil, gofr.New())

	mockCar.EXPECT().DeleteCar(ctx, id.String()).Return(nil)
	mockEngine.EXPECT().EngineDelete(ctx, id.String()).Return(nil)
	mockCar.EXPECT().DeleteCar(ctx, id2.String()).Return(errors.InvalidParam{})
	mockCar.EXPECT().DeleteCar(ctx, id3.String()).Return(nil)
	mockEngine.EXPECT().EngineDelete(ctx, id3.String()).Return(errors.InvalidParam{})

	for _, tc := range tests {
		err := carService.Delete(ctx, tc.id.String())
		assert.Equal(t, err, tc.err)
	}
}

// TestValidateCreate tests that the cal being created is valid or not.
func TestValidateCreate(t *testing.T) {
	var (
		id  = uuid.New()
		id2 = uuid.New()
		c1  = models.Car{ID: id, Name: "Model 3", Year: 2010, Brand: "Tesla", FuelType: "Diesel",
			Engine: models.Engine{EngineID: id, Displacement: 400, Cylinders: 6}}
		c2 = models.Car{ID: id2, Name: "Model 3", Year: 2010, Brand: "wer", FuelType: "Diesel",
			Engine: models.Engine{EngineID: id2, Displacement: 400, Cylinders: 6}}
		c3 = models.Car{}
		c4 = models.Car{ID: id, Name: "Model 3", Year: 2010, Brand: "Tesla", FuelType: "Solar",
			Engine: models.Engine{EngineID: id, Displacement: 400, Cylinders: 6}}
		c5 = models.Car{ID: id, Name: "Model 3", Year: 2033, Brand: "Tesla", FuelType: "Diesel",
			Engine: models.Engine{EngineID: id, Displacement: 400, Cylinders: 6}}
	)

	testCases := []struct {
		desc   string
		input  models.Car
		output models.Car
	}{
		{desc: "Success case", input: c1, output: c1},
		{desc: "Wrong brand", input: c2, output: c3},
		{desc: "Wrong Fuel Type", input: c4, output: c3},
		{desc: "Wrong year", input: c5, output: c3},
	}

	for i := range testCases {
		c := validateCreateCar(&testCases[i].input)

		assert.Equal(t, testCases[i].output, c,
			" [TEST%d]Failed. Got %v\tExpected %v\n", i+1, c, testCases[i].output)
	}
}
