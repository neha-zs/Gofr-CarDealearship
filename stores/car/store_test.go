package car

import (
	"context"
	"fmt"
	"testing"

	"Project/CarDealearship/models"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCarByID(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	defer db.Close()
	a := New()

	id1 := uuid.New()
	id2 := uuid.New()

	testCases := []struct {
		desc string
		id   string
		resp models.Car
		err  error
		mock interface{}
	}{
		{
			desc: "Success Case",
			id:   id1.String(),
			resp: models.Car{ID: id1, Engine: models.Engine{EngineID: id1, Displacement: 0, Cylinders: 0, Range: 0},
				Name: "Model 2", Year: 2000, Brand: "Tesla", FuelType: "Petrol"},
			err: nil,
			mock: mock.ExpectQuery("SELECT * FROM Car WHERE ID=?;").WithArgs(id1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "engine_id", "name", "year", "brand", "fuelType"}).
					AddRow(id1.String(), id1.String(), "Model 2", 2000, "Tesla", "Petrol")),
		},
		{
			desc: "ID not present",
			id:   id2.String(),
			resp: models.Car{},
			err:  errors.EntityNotFound{Entity: "Car", ID: id2.String()},
			mock: mock.ExpectQuery("SELECT * FROM Car WHERE ID=?;").WithArgs(id2).
				WillReturnError(errors.EntityNotFound{Entity: "Car", ID: id2.String()}),
		},
	}

	for i, tc := range testCases {

		resp, err := a.GetCarByID(ctx, tc.id)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

// TestGetCarsByBrand tests the datastore function GetCarsByBrand
func TestGetCarsByBrand(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	defer db.Close()
	a := New()

	var (
		id1 = uuid.New()
		id2 = uuid.New()
		id3 = uuid.New()

		car = models.Car{ID: id1, Name: "GenX", Year: 2015, Brand: "Tesla",
			FuelType: "electric", Engine: models.Engine{EngineID: id1}}

		car2 = models.Car{ID: id2, Name: "Model 3", Year: 2020, Brand: "Tesla",
			FuelType: "electric", Engine: models.Engine{EngineID: id2}}

		car3 = models.Car{ID: id3, Name: "Model 3", Year: 2020, Brand: "BMW",
			FuelType: "electric", Engine: models.Engine{EngineID: id3}}

		rows = sqlmock.NewRows([]string{"id", "engine_id", "name", "year", "brand", "fuel_type"}).
			AddRow(id1.String(), id1.String(), car.Name, car.Year, car.Brand, car.FuelType).
			AddRow(id2.String(), id2.String(), car2.Name, car2.Year, car2.Brand, car2.FuelType)

		rwbmw = sqlmock.NewRows([]string{"id", "engine_id", "name", "year", "brand"}).
			AddRow(id3.String(), id3.String(), car3.Name, car3.Year, car3.Brand)

		rowFerrari = sqlmock.NewRows([]string{"id", "engine_id", "name", "year", "brand"}).
				AddRow(id3.String(), id3.String(), car3.Name, car3.Year, "Ferrari").
				RowError(0, errors.Error("Row error"))

		rowPorsche = sqlmock.NewRows([]string{"id", "engine_id", "name", "year", "brand", "fuel_type"}).
				CloseError(fmt.Errorf("close error"))
	)

	testCases := []struct {
		desc   string
		brand  string
		output []models.Car
		err    error
	}{
		{desc: "Get all Tesla Car", brand: "Tesla", output: []models.Car{car, car2}, err: nil},
		{desc: "insufficient argument", brand: "BMW", output: nil,
			err: errors.Error("Scan Error")},
		{"Brand name not mentioned", "", nil, errors.MissingParam{}},
		{desc: "error in return rows", brand: "Ferrari", output: []models.Car{}, err: errors.Error("Row error")},
		{desc: "error in close row", brand: "Porsche", output: nil, err: nil},
	}

	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("Tesla").WillReturnRows(rows)
	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("BMW").WillReturnRows(rwbmw)
	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("").
		WillReturnError(errors.MissingParam{})
	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("Ferrari").
		WillReturnRows(rowFerrari)
	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("Porsche").WillReturnRows(rowPorsche)

	for i, tc := range testCases {
		car, err := a.GetCarsByBrand(ctx, tc.brand)

		assert.Equal(t, err, tc.err,
			"\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)

		assert.Equal(t, car, tc.output,
			"\n[TEST %v] Failed \n got %v\nGot \n Expected %v", i, car, tc.output)
	}
}

// TestCreateCar test the Create functionality of the datastore layer
func TestCreateCar(t *testing.T) {
	id := uuid.New()

	car := models.Car{ID: id, Name: "GenX", Year: 2015, Brand: "Tesla",
		FuelType: "electric", Engine: models.Engine{EngineID: id}}
	car2 := models.Car{ID: uuid.Nil, Name: "GenX", Year: 2015, Brand: "Tesla",
		FuelType: "electric", Engine: models.Engine{EngineID: id}}

	testCases := []struct {
		desc           string
		input          models.Car
		expectedOutput models.Car
		err            error
	}{
		{"Car created successfully", car, car, nil},
		{"failure", car2, models.Car{}, errors.Error("query error")},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()
	a := New()

	defer db.Close()

	mock.ExpectExec("INSERT INTO Car (id,engine_id,name,year,brand,fuel_type) VALUES(?,?,?,?,?,?)").
		WithArgs(car.ID, car.Engine.EngineID, car.Name, car.Year, car.Brand, car.FuelType).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO Car (id,engine_id,name,year,brand,fuel_type) VALUES(?,?,?,?,?,?)").
		WithArgs(uuid.Nil, car.Engine.EngineID, car.Name, car.Year, car.Brand, car.FuelType).
		WillReturnError(errors.Error("query error"))

	for i, tc := range testCases {
		res, err := a.CreateCar(ctx, &tc.input)

		if res != tc.expectedOutput {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, res, tc.expectedOutput)
		}

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestUpdateCar test the Update functionality of the datastore layer
func TestUpdateCar(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	car := models.Car{ID: id, Name: "BMW", Year: 2018, Brand: "Rolls-Royce", FuelType: "petrol"}
	updateFailed := errors.Error("Update Failed")

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	a := New()

	defer db.Close()

	mock.ExpectExec("UPDATE Car SET name=?,year=?,brand=?,fuel_type=? WHERE id=?").
		WithArgs(car.Name, car.Year, car.Brand, car.FuelType, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE Car SET name=?,year=?,brand=?,fuel_type=? WHERE id=?").
		WithArgs(car.Name, car.Year, car.Brand, car.FuelType, id).
		WillReturnError(errors.Error("Update Failed"))

	cases := []struct {
		desc  string
		input models.Car
		err   error
	}{
		{"success", car, nil},
		{"failure", car, updateFailed},
	}

	for i, tc := range cases {
		_, err := a.UpdateCar(ctx, tc.input.ID.String(), &tc.input)
		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestDeleteCar test the Delete functionality of the datastore layer
func TestDeleteCar(t *testing.T) {
	id1 := uuid.New()
	//deleteErr := errors.New("delete failed")

	testCases := []struct {
		desc         string
		id           uuid.UUID
		rowsEffected int
		err          error
	}{
		{"Success", id1, 1, nil},
		{"ID does not exists", uuid.Nil, 0, errors.EntityNotFound{}},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	a := New()

	defer db.Close()

	mock.ExpectExec("DELETE FROM Car WHERE ID=?").WithArgs(id1.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM Car WHERE ID=?").WithArgs(uuid.Nil).
		WillReturnError(errors.EntityNotFound{})

	for i, tc := range testCases {
		err := a.DeleteCar(ctx, tc.id.String())

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
